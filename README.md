goconf
======

A config parser for Golang. Config item with Int, Float, String and Array Type is supported.

####An example of config file:
     > StringItem: value
     > IntItem: 1000
     > FloatItem: 90.5
     >
     > [@IntArray]: 10 12 13
     > [@IntArray1@,]: 1, 2, 3, 4, 5

####Note:
Int, Float and String item is easy to specify, just using format of 'key: value'.
And the rule of defining an Array is a little complex:
    1) [@ARRAY_KEY]: ELEMENTS_OF_ARRAY
    2) [@ARRAY_KEY@ELEMENT_SEPARATOR]: ELEMENTS_OF_ARRAY
The first way uses the default element separator ' '. And it's possible to specify a customed separator using the latter way.

####Sample code:
Sample code can be found in 'conf_test.go'. There are two mode to use the conf.Conf:
    1) Error mode which is idiomatic way in Go, but also tedious.
    2) Panic mode which just like exception in Java.

