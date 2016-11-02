/**
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2014/11/28 11:36:58
 */

package goconf

import (
	"github.com/chosen0ne/goutils"
	"strconv"
	"strings"
)

// ------- Item ------- //
type Item struct {
	key string
	val string
}

func (item *Item) Key() string {
	return item.key
}

func (item *Item) String() string {
	return item.key + "=>" + item.val
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
	eleStr := item.ToStringArray()

	values := make([]int64, len(eleStr))
	for idx, ele := range eleStr {
		ele = strings.Trim(ele, _SPACE_CHARS)
		val, err := strconv.ParseInt(ele, 10, 64)
		if err != nil {
			return nil, goutils.WrapErr(err)
		}
		values[idx] = val
	}

	return values, nil
}

func (item *Item) ToFloatArray() ([]float64, error) {
	eleStr := item.ToStringArray()

	values := make([]float64, len(eleStr))
	for idx, ele := range eleStr {
		ele = strings.Trim(ele, _SPACE_CHARS)
		val, err := strconv.ParseFloat(ele, 64)
		if err != nil {
			return nil, goutils.WrapErr(err)
		}
		values[idx] = val
	}

	return values, nil
}

func (item *Item) ToStringArray() []string {
	parts := strings.Split(item.val, string(elementSep))

	var eles []string
	for _, p := range parts {
		if p != "" {
			eles = append(eles, strings.Trim(p, _SPACE_CHARS))
		}
	}

	return eles
}
