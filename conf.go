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
	"bufio"
	"fmt"
	"github.com/chosen0ne/goutils"
	"io"
	"os"
	"strings"
)

const (
	_KV_SEP      = ':'
	_NEWLINE     = '\n'
	_SPACE_CHARS = " \t\n"
	_GLOBAL      = "__global__"

	_DEFAULT_SEP   = ' '
	_SECTION_LEFT  = '['
	_SECTION_RIGHT = ']'
	_COMMENT_TAG   = '#'
)

var (
	elementSep byte
)

type _Section map[string]*Item

// ------- Conf ------- //
type Conf struct {
	filePath string
	sections map[string]_Section
	eleSep   byte
	cur      _Section // current section
}

func New(filePath string) *Conf {
	conf := &Conf{}
	conf.filePath = filePath
	conf.sections = make(map[string]_Section)
	conf.cur = make(map[string]*Item)
	conf.sections[_GLOBAL] = conf.cur

	return conf
}

func (conf *Conf) Parse() error {
	// Open config file
	f, err := os.Open(conf.filePath)
	if err != nil {
		return goutils.WrapErr(err)
	}

	defer f.Close()
	buf := bufio.NewReader(f)

	if err := conf._parse(buf); err != nil {
		return err
	}

	conf.cur = conf.sections[_GLOBAL]

	return nil
}

func (conf *Conf) _parse(buf *bufio.Reader) error {
	for {
		line, err := buf.ReadString(_NEWLINE)
		if len(line) == 0 && err == io.EOF {
			return nil
		} else if err != nil && err != io.EOF {
			return goutils.WrapErr(err)
		}

		// Trim space chars
		lineStr := strings.Trim(line, _SPACE_CHARS)

		// Found an empty line
		if len(lineStr) == 0 {
			continue
		}

		// Remove '\n'
		if lineStr[len(lineStr)-1] == _NEWLINE {
			lineStr = lineStr[:len(lineStr)-1]
		}

		// Found a comment line
		if lineStr[0] == _COMMENT_TAG {
			continue
		}

		if isSection(lineStr) {
			sectionName := strings.Trim(lineStr[1:len(lineStr)-1], _SPACE_CHARS)
			if _, ok := conf.sections[sectionName]; ok {
				return goutils.NewErr("section '%s' already exist", sectionName)
			}

			conf.cur = make(map[string]*Item)
			conf.sections[sectionName] = conf.cur
		} else {
			// Find 'Key : Value'
			parts := strings.SplitN(lineStr, string(_KV_SEP), 2)
			if len(parts) != 2 {
				return goutils.NewErr("need ':' in a line, line: %s", lineStr)
			}
			key := strings.Trim(parts[0], _SPACE_CHARS)
			val := strings.Trim(parts[1], _SPACE_CHARS)
			if len(val) == 0 {
				return goutils.NewErr("an empty value")
			}

			conf.cur[key] = &Item{key, val}
		}
	}

	return nil
}

func (conf *Conf) GetItem(key string) (*Item, error) {
	item, ok := conf.cur[key]
	if !ok {
		return nil, goutils.NewErr("non-exist item: %s", key)
	}
	return item, nil
}

func (conf *Conf) HasItem(key string) bool {
	_, ok := conf.cur[key]
	return ok
}

func (conf *Conf) Items() []*Item {
	items := make([]*Item, len(conf.cur))
	idx := 0
	for _, v := range conf.cur {
		items[idx] = v
		idx++
	}

	return items
}

func (conf *Conf) GetInt(key string) (int64, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return -1, goutils.WrapErr(err)
	}

	return item.ToInt()
}

func (conf *Conf) GetFloat(key string) (float64, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return -1, goutils.WrapErr(err)
	}

	return item.ToFloat()
}

func (conf *Conf) GetString(key string) (string, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return "", goutils.WrapErr(err)
	}

	return item.val, nil
}

func (conf *Conf) GetIntArray(key string) ([]int64, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return nil, goutils.WrapErr(err)
	}

	return item.ToIntArray()
}

func (conf *Conf) GetFloatArray(key string) ([]float64, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return nil, goutils.WrapErr(err)
	}

	return item.ToFloatArray()
}

func (conf *Conf) GetStringArray(key string) ([]string, error) {
	item, err := conf.GetItem(key)
	if err != nil {
		return nil, goutils.WrapErr(err)
	}

	return item.ToStringArray(), nil
}

func (conf *Conf) Section(name string) error {
	section, ok := conf.sections[name]
	if ok {
		conf.cur = section
		return nil
	}

	return goutils.NewErr("no section '%s'", name)
}

func (conf *Conf) HasSection(name string) bool {
	_, ok := conf.sections[name]
	return ok
}

func (conf *Conf) SetGlobalSection() {
	conf.cur = conf.sections[_GLOBAL]
}

func (conf *Conf) LoadSection(name string, configObj interface{}) error {
	section := conf.sections[name]
	if section == nil {
		return goutils.NewErr("no section named '%s'", name)
	}

	return nil
}

// SetElementSep: set the separator of elements in an array
func SetElementSep(sep byte) {
	elementSep = sep
}

func init() {
	elementSep = _DEFAULT_SEP
}

func isSection(line string) bool {
	if line[0] == _SECTION_LEFT && line[len(line)-1] == _SECTION_RIGHT {
		return true
	}

	return false
}
