/**
 * Unit test cases
 *
 * @author  chosen0ne(louzhenlin86@126.com)
 * @date    2014/11/04 15:06:47
 */

package goconf

import (
	"bufio"
	"bytes"
	"chosen0ne.com/utils"
	"reflect"
	"testing"
)

// ------- Tests for Item ------- //
func TestItemSingleValOk(t *testing.T) {
	var num int64 = 100
	expected := []interface{}{num, 1.101, "astring"}
	input := []string{"100", "1.101", "astring"}
	callMethod := []string{"ToInt", "ToFloat", "ToString"}
	returnCount := []int{2, 2, 1}

	for i, _ := range input {
		item := &Item{val: input[i]}
		value := reflect.ValueOf(item)
		method := value.MethodByName(callMethod[i])
		results := method.Call(make([]reflect.Value, 0))

		if len(results) != returnCount[i] {
			t.Errorf("count of return args of '%s' error", callMethod[i])
		}

		// check error
		if returnCount[i] > 1 {
			err := results[1].Interface()
			if err != nil && err.(error) != nil {
				t.Errorf("failed to test '%s', err: %s", callMethod[i], err)
			}
		}

		// check return value
		expectVal := expected[i]
		outputVal := results[0].Interface()
		if outputVal != expectVal {
			t.Error("output is not expected, output:", outputVal, ", expected:", expectVal)
		}
	}
}

func TestItemIntErr(t *testing.T) {
	item := &Item{val: "not a int: 1"}
	val, err := item.ToInt()
	if err == nil || val != 0 {
		t.Errorf("should return a err, val: %d, err: %s", val, err)
	}
}

func TestItemFloatErr(t *testing.T) {
	item := &Item{val: "not a float: 1.011"}
	val, err := item.ToFloat()
	if err == nil || val != 0 {
		t.Errorf("should return a err, val: %f, err: %s", val, err)
	}
}

func matchStringArray(output, expected []string) error {
	if len(output) != len(expected) {
		return utils.NewErr("length of expected and output is different output: %d, expected: %d",
			len(output), len(expected))
	}

	for idx, str := range output {
		if str != expected[idx] {
			return utils.NewErr("not expected output, output: %s, expected: %s", output, expected)
		}
	}

	return nil
}

// Test for Array use default separator ' '
func TestItemStringArrayOk1(t *testing.T) {
	item := &Item{"key1", "abc de fg h"}
	expected := []string{"abc", "de", "fg", "h"}

	strArray := item.ToStringArray()

	err := matchStringArray(strArray, expected)
	if err != nil {
		t.Errorf("not expected output, err: %s", err)
	}
}

func TestItemIntArrayOk(t *testing.T) {
	item := &Item{"IntArray", "12 23 44 55"}
	expected := []int64{12, 23, 44, 55}

	intArray, err := item.ToIntArray()
	if err != nil {
		t.Fatalf("failed to IntArray, err: %s", err)
	}

	if len(intArray) != len(expected) {
		t.Errorf("length of expected and output is different output: %d, expected: %d",
			len(intArray), len(expected))
	}

	for idx, v := range intArray {
		if v != expected[idx] {
			t.Errorf("not expected output, output: %s, expected: %s", intArray, expected)
		}
	}
}

func TestItemFloatArrayOk(t *testing.T) {
	item := &Item{"FloatArray", "1.1 1.2 12.33"}
	expected := []float64{1.1, 1.2, 12.33}

	floatArray, err := item.ToFloatArray()
	if err != nil {
		t.Fatalf("failed to FloatArray, err: %s", err)
	}

	if len(floatArray) != len(expected) {
		t.Errorf("length of expected and output is different output: %d, expected: %d",
			len(floatArray), len(expected))
	}

	for idx, v := range floatArray {
		if v != floatArray[idx] {
			t.Errorf("not expected output, output: %s, expected: %s", floatArray, expected)
		}
	}
}

// ------- Tests for Conf ------- //
func genConf(s string) (*Conf, *bufio.Reader) {
	buf := bytes.NewBufferString(s)
	return New(""), bufio.NewReader(buf)
}

func TestConfParseOk1(t *testing.T) {
	conf, buf := genConf("item1: value1\n\n\nitem2: value2")

	if err := conf._parse(buf); err != nil {
		t.Errorf("failed to parse, err: %s", err)
	}
}

func TestConfParseOk2(t *testing.T) {
	conf, buf := genConf("[@int@;]: a;b;c\n[@int]: 1 2 3")

	if err := conf._parse(buf); err != nil {
		t.Errorf("failed to parse, err: %s", err)
	}
}

// Partial Key, without value
func TestConfParseErr1(t *testing.T) {
	conf, buf := genConf("item1: valu\nitem1jfak")

	if err := conf._parse(buf); err == nil {
		t.Errorf("need a EOF error")
	}
}

func TestConfParseErr2(t *testing.T) {
	conf, buf := genConf("item1:  ")

	if err := conf._parse(buf); err == nil {
		t.Errorf("need a EOF error")
	}
}

func TestConfItemsOk(t *testing.T) {
	conf, buf := genConf("a:b\nc:d\ne:f\ng:h")
	expected := map[string]int{"a": 1, "c": 1, "e": 1, "g": 1}

	if err := conf._parse(buf); err != nil {
		t.Errorf("failed to parse, err: %s", err)
	}

	for _, item := range conf.Items() {
		if _, ok := expected[item.Key()]; !ok {
			t.Errorf("extra key '%s'", item.Key())
		}
	}

	for k, _ := range expected {
		if !conf.HasItem(k) {
			t.Errorf("key %s non-exist, items: %s", k, conf.Items())
		}
	}
}

func TestAll(t *testing.T) {
	config := New("conf_sample.conf")

	if err := config.Parse(); err != nil {
		t.Error("failed to Parse, err:", err)
	}

	// iterate items
	t.Log("items:")
	for _, item := range config.Items() {
		t.Log("\t", item.Key())
	}

	strItem, err := config.GetString("StringItem")
	if err == nil {
		t.Log("StringItem =>", strItem)
	}

	intItem, err := config.GetInt("IntItem")
	if err == nil {
		t.Log("IntItem =>", intItem)
	}

	intArray, err := config.GetIntArray("IntArray")
	if err == nil {
		t.Log("IntArray =>", intArray)
	}

	floatArray, err := config.GetFloatArray("FloatArray")
	if err == nil {
		t.Log("FloatArray =>", floatArray)
	}
}

func TestAllByPanicWay(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("failed to load conf, err:", err)
		}
	}()

	config := New("conf_sample.conf")

	config.ParseOrPanic()
	t.Log("items:")
	for _, item := range config.Items() {
		t.Log("\t", item.Key())
	}

	t.Log("StringItem=>", config.ToString("StringItem"))
	t.Log("IntItem=>", config.ToInt("IntItem"))
	t.Log("IntArray=>", config.ToIntArray("IntArray"))
	t.Log("FloatArray=>", config.ToFloatArray("FloatArray"))
}

func TestSection(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error("failed to load conf, err:", err)
		}
	}()

	config := New("conf_sample.conf")
	config.ParseOrPanic()
	config.Section("Section1")

	t.Log(config)
	for _, item := range config.Items() {
		t.Log(item)
	}
}
