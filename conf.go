/**
 * A config parser for Golang. Support Int, Float, String and Array.
 *      e.g. config file:
 *          > StringItem: value
 *          > IntItem: 1000
 *          > FloatItem: 90.5
 *          >
 *          > [@IntArray]: 10 12 13
 *          > [@IntArray1@,]: 1, 2, 3, 4, 5
 *
 *  The rule to define an Array:
 *          1) [@ARRAY_KEY]
 *          2) [@ARRAY_KEY@ELEMENT_SEPARATOR]
 *      Default element separator is ' '.
 *      And it's possible to specify a customed separator using the latter way.
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2014/11/03 15:03:49
 */

package goconf

import (
    "os"
    "bufio"
    "strings"
    "io"
    "errors"
    "strconv"
)


const (
    _KV_SEP = ':'
    _NEWLINE = '\n'
    _SPACE_CHARS = " \t\n"
    _ARRAY_TAG = '@'
    _MAP_TAG = '$'
    _DEFAULT_SEP = ' '
    _COMPOSITE_LEFT = '['
    _COMPOSITE_RIGHT = ']'
    _COMMENT_TAG = '#'
)

type ParseState int


// ------- Item ------- //
type Item struct {
    rawKey  string
    val     string
    key     string
}

func (item *Item) Key() string {
    return item.key
}

func (item *Item) ToInt() (int64, error) {
    return strconv.ParseInt(item.val, 10, 64)
}

func (item *Item) ToString() string {
    return item.val
}

func (item *Item) ToFloat() (float64, error) {
    return strconv.ParseFloat(item.val, 64)
}

func (item *Item) ToIntArray() ([]int64, error) {
    eleStr, err := item.ToStringArray()
    if err != nil {
        return nil, err
    }

    values := make([]int64, len(eleStr))
    for idx, ele := range eleStr {
        ele = strings.Trim(ele, _SPACE_CHARS)
        val, err := strconv.ParseInt(ele, 10, 64)
        if err != nil {
            return nil, err
        }
        values[idx] = val
    }

    return values, nil
}

func (item *Item) ToFloatArray() ([]float64, error) {
    eleStr, err := item.ToStringArray()
    if err != nil {
        return nil, err
    }

    values := make([]float64, len(eleStr))
    for idx, ele := range eleStr {
        ele = strings.Trim(ele, _SPACE_CHARS)
        val, err := strconv.ParseFloat(ele, 64)
        if err != nil {
            return nil, err
        }
        values[idx] = val
    }

    return values, nil
}

func (item *Item) ToStringArray() ([]string, error) {
    // Make sure the item is a Array
    // Key: [$ArrayKey$;ArraySep] or [$ArrayKey]
    keyLen := len(item.rawKey)
    if item.rawKey[0] != _COMPOSITE_LEFT || item.rawKey[keyLen - 1] != _COMPOSITE_RIGHT ||
            item.rawKey[1] != _ARRAY_TAG || keyLen < 4 {
        errMsg := strings.Join([]string{"item is not a Array, key:", item.rawKey}, " ")
        return nil, errors.New(errMsg)
    }

    // Extract sep
    var sep byte
    sepIdx := strings.LastIndexAny(item.rawKey, string([]byte{_ARRAY_TAG}))
    if sepIdx == -1 {
        return nil, errors.New("not found Array Tag: @")
    }

    if sepIdx == 1 {
        sep = _DEFAULT_SEP
    } else if sepIdx != keyLen -3 {
        return nil, errors.New("Array Sep can only set to one char")
    } else {
        sep = item.rawKey[sepIdx + 1]
    }

    return strings.Split(item.val, string(sep)), nil
}


// ------- Conf ------- //
type Conf struct {
    in          io.Reader
    filePath    string
    items       map[string]*Item
    curState    ParseState
}

func New(filePath string) (*Conf, error) {
    conf := &Conf{}
    conf.filePath = filePath
    conf.items = make(map[string]*Item)

    if err := conf.openFile(); err != nil {
        return nil, err
    }

    return conf, nil
}

func (conf *Conf) Parse() error {
    buf := bufio.NewReader(conf.in)
    for {
        line, err := buf.ReadString(_NEWLINE)
        if len(line) == 0 && err == io.EOF {
            return nil
        } else if err != nil && err != io.EOF {
            return err
        }

        // Trim left space chars
        lineStr := strings.TrimLeft(line, _SPACE_CHARS)

        // Found an empty line
        if len(lineStr) == 0 {
            continue
        }

        if lineStr[len(lineStr) - 1] == _NEWLINE {
            lineStr = lineStr[:len(lineStr) - 1]
        }

        // Found an comment line
        if lineStr[0] == _COMMENT_TAG {
            continue
        }

        // Find 'Key : Value'
        parts := strings.SplitN(lineStr, string(_KV_SEP), 2)
        if len(parts) != 2 {
            return errors.New("the line need a ':', line: " + lineStr)
        }
        rawKey := strings.Trim(parts[0], _SPACE_CHARS)
        key := parseKey(rawKey)
        val := strings.Trim(parts[1], _SPACE_CHARS)
        if len(val) == 0 {
            return errors.New("an empty value")
        }

        conf.items[key] = &Item{rawKey, val, key}
    }
    return nil
}

func (conf *Conf) GetItem(key string) (*Item, error) {
    item, ok := conf.items[key]
    if !ok {
        return nil, errors.New("non-exist key: " + key)
    }
    return item, nil
}

func (conf *Conf) HasItem(key string) bool {
    _, ok := conf.items[key]
    return ok
}

func (conf *Conf) Items() []*Item {
    items := make([]*Item, len(conf.items))
    idx := 0
    for _, v := range conf.items {
        items[idx] = v
        idx++
    }

    return items
}

func (conf *Conf) GetInt(key string) (int64, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return -1, err
    }

    return item.ToInt()
}

func (conf *Conf) GetFloat(key string) (float64, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return -1, err
    }

    return item.ToFloat()
}

func (conf *Conf) GetString(key string) (string, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return "", err
    }

    return item.val, nil
}

func (conf *Conf) GetIntArray(key string) ([]int64, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return nil, err
    }

    return item.ToIntArray()
}

func (conf *Conf) GetFloatArray(key string) ([]float64, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return nil, err
    }

    return item.ToFloatArray()
}

func (conf *Conf) GetStringArray(key string) ([]string, error) {
    item, err := conf.GetItem(key)
    if err != nil {
        return nil, err
    }

    return item.ToStringArray()
}

func (conf *Conf) Close() error {
    f := conf.in.(*os.File)
    if err := f.Close(); err != nil {
        return err
    }

    return nil
}

func (conf *Conf) openFile() error {
    f, err := os.Open(conf.filePath)
    conf.in = f
    return err
}

func extraceString(buf []byte) string {
    str := string(buf[:len(buf) - 1])
    return strings.Trim(str, _SPACE_CHARS)
}

func isSpaceBuf(buf []byte) bool {
    for _, c := range buf {
        if c != ' ' && c != '\n' && c != '\t' {
            return false
        }
    }

    return true
}

func parseKey(key string) string {
    if len(key) < 4 {
        return key
    }

    if key[0] != _COMPOSITE_LEFT || key[len(key) - 1] != _COMPOSITE_RIGHT ||
            (key[1] != _ARRAY_TAG && key[1] != _MAP_TAG) {
        return key
    }

    c := key[len(key) - 3]
    if len(key) != 4 && (c == _ARRAY_TAG || c == _MAP_TAG) {
        return key[2:len(key) - 3]
    }

    return key[2:len(key) - 1]
}

