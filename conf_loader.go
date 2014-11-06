/**
 * Use to load config into a Config Object.
 * Make sure field name of the object is same as the item name in config file.
 * And the function 'Load' will read the config item by the type of each field,
 * then fill it.
 *
 *      e.g. config file:
 *          > StringItem: value
 *          > IntItem: 1000
 *          > FloatItem: 90.5
 *          >
 *          > [@IntArray]: 10 12 13
 *          > [@IntArray1@,]: 1, 2, 3, 4, 5
 *
 *      And the corresponding Config Struct is:
 *          type ConfigObj struct {
 *              StringItem  string      // field must be public, or it can't be set by reflection
 *              IntItem     int
 *              FloatItem   float32
 *              IntArray    []int64     // slice type of integer can only set int64, other int types aren't supported.
 *              IntArray1   []float64   // slice type of float can only set float64, float32 isn't supported.
 *          }
 *
 *          confObj := &ConfigObj{StringItem: "default value"} // default values can be set
 *          // Because of using reflect package, the error emitted by 'panic' must 'recover'
 *          defer func() {
 *              if err := recover(); err != nil {
 *                  // err recover
 *              }
 *          }
 *          LoadOrPanic(confObj, "config.conf")
 *
 * @author  chosen0ne(louzhenlin86@126.com
 * @date    2014/11/05 11:50:13
 */

package goconf

import (
    "reflect"
    "errors"
)

func Load(configObjPtr interface{}, configFile string) error {
    // Settable?
    configObj := reflect.ValueOf(configObjPtr).Elem()
    if !configObj.CanSet() {
        return errors.New("configObj must be settable")
    }

    // Create and Parse conf
    conf, err := New(configFile)
    if err != nil {
        return err
    }

    defer conf.Close()

    if err := conf.Parse(); err != nil {
        return err
    }

    // Load fields from conf
    value := reflect.ValueOf(configObjPtr).Elem()
    t := configObj.Type()
    for i := 0; i < value.NumField(); i++ {
        fieldValue := value.Field(i)
        fieldMeta := t.Field(i)
        if err := loadField(configObj, &fieldMeta, &fieldValue, conf); err != nil {
            return err
        }
    }

    return nil
}

// ------- Panic mode ------- //
func LoadOrPanic(configObjPtr interface{}, configFile string) {
    if err := Load(configObjPtr, configFile); err != nil {
        panic(err)
    }
}

func loadField(
            configObj interface{},
            fieldMeta *reflect.StructField,
            fieldValue *reflect.Value,
            conf *Conf) error {
    fieldName := fieldMeta.Name
    // Check field settable?
    if !fieldValue.CanSet() {
        return errors.New("field not settable, field: " + fieldName)
    }

    if !conf.HasItem(fieldName) {
        return nil
    }

    // Fetch value from conf, and load Config Object
    kind := fieldValue.Kind()
    if isInt(kind) {
        val, err := conf.GetInt(fieldName)
        if err != nil {
            return err
        }
        fieldValue.SetInt(val)
    } else if kind == reflect.Float32 || kind == reflect.Float64 {
        val, err := conf.GetFloat(fieldName)
        if err != nil {
            return err
        }
        fieldValue.SetFloat(val)
    } else if kind == reflect.String {
        val, err := conf.GetString(fieldName)
        if err != nil {
            return err
        }
        fieldValue.SetString(val)
    } else if kind == reflect.Slice {
        if err := loadSliceField(configObj, fieldMeta, fieldValue, conf); err != nil {
            return err
        }
    } else {
        return errors.New("not support type: " + kind.String())
    }

    return nil
}

func loadSliceField(
            configObj interface{},
            fieldMeta *reflect.StructField,
            fieldValue *reflect.Value,
            conf *Conf) error {
    fieldName := fieldMeta.Name
    eleValue := fieldMeta.Type.Elem()

    eleKind := eleValue.Kind()
    if isInt(eleKind) {
        vals, err := conf.GetIntArray(fieldName)
        if err != nil {
            return err
        }
        for _, val := range vals {
            fieldValue.Set(reflect.Append(*fieldValue, reflect.ValueOf(val)))
        }
    } else if eleKind == reflect.Float32 || eleKind == reflect.Float64 {
        vals, err := conf.GetFloatArray(fieldName)
        if err != nil {
            return err
        }
        for _, val := range vals {
            fieldValue.Set(reflect.Append(*fieldValue, reflect.ValueOf(val)))
        }
    } else if eleKind == reflect.String {
        vals, err := conf.GetStringArray(fieldName)
        if err != nil {
            return err
        }
        for _, val := range vals {
            fieldValue.Set(reflect.Append(*fieldValue, reflect.ValueOf(val)))
        }
    } else {
        return errors.New("not support element type for slice")
    }

    return nil
}

func isInt(k reflect.Kind) bool {
    if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 ||
            k == reflect.Int32 || k == reflect.Int64 || k == reflect.Uint ||
            k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 ||
            k == reflect.Uint64 {
        return true
    }

    return false
}

