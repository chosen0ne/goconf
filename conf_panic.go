/**
 * A panic-recover way to use Conf.
 * All the methods suffixed by BPY(By panic way) may throw exception,
 * so they need to recover.
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2014/11/04 21:50:38
 */

package goconf

func (conf *Conf) ParseOrPanic() {
    if err := conf.Parse(); err != nil {
        panic(err)
    }
}

func (conf *Conf) GetItemOrPanic(key string) *Item {
    item, err := conf.GetItem(key)
    if err != nil {
        panic(err)
    }
    return item
}

func (conf *Conf) ToInt(key string) int64 {
    val, err := conf.GetInt(key)
    if err != nil {
        panic(err)
    }
    return val
}

func (conf *Conf) ToFloat(key string) float64 {
    val, err := conf.GetFloat(key)
    if err != nil {
        panic(err)
    }
    return val
}

func (conf *Conf) ToString(key string) string {
    val, err := conf.GetString(key)
    if err != nil {
        panic(err)
    }
    return val
}

func (conf *Conf) ToIntArray(key string) []int64 {
    val, err := conf.GetIntArray(key)
    if err != nil {
        panic(err)
    }
    return val
}

func (conf *Conf) ToFloatArray(key string) []float64 {
    val, err := conf.GetFloatArray(key)
    if err != nil {
        panic(err)
    }
    return val
}

func (conf *Conf) ToStringArray(key string) []string {
    val, err := conf.GetStringArray(key)
    if err != nil {
        panic(err)
    }
    return val
}

