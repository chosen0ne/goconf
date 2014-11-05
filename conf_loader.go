/**
 * Use to load config into a Config Object
 * @author  chosen0ne(louzhenlin86@126.com)
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

func LoadDefault(
            configObjPtr interface{},
            configFile string,
            defaultConfigPtr interface{}) error {
    configObj := reflect.ValueOf(configObjPtr).Elem()
    defaultObj := reflect.ValueOf(defaultConfigPtr).Elem()

    // Make sure same type
    if configObj.Type() != defaultObj.Type() {
        return errors.New("type of configObjPtr and defaultConfigPtr must be same")
    }

    if err := Load(configObjPtr, configFile); err != nil {
        return err
    }

    fieldMetaInfo := configObj.Type()
    for i := 0; i < configObj.NumField(); i++ {
        fieldMeta := fieldMetaInfo.Field(i)
        srcField := configObj.FieldByName(fieldMeta.Name)

        // Not set by config
        if reflect.ValueOf(srcField.Interface()) == reflect.Zero(srcField.Type()) ||
                (srcField.Kind() == reflect.Slice && srcField.Len() == 0) {
            distField := defaultObj.FieldByName(fieldMeta.Name)
            srcField.Set(distField)
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

func LoadDefaultOrPanic(
            configObjPtr interface{},
            configFile string,
            defaultConfigPtr interface{}) {
    if err := LoadDefault(configObjPtr, configFile, defaultConfigPtr); err != nil {
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


