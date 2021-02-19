package dvisit

import (
	"reflect"
	"testing"
)

type node struct {
	Key      string
	Value    string
	internal string
}

//SetInternal 设置 internal 值
func (this *node) SetInternal(str string) {
	this.internal = str
}

type Test struct {
	Array       [5]int
	Slice       []string
	Map         map[string]string
	Struct      node
	Ptr         *node
	Key         string
	Interface   interface{}
	Interface02 interface{}
}

func getTestData() Test {
	test := Test{
		Array: [5]int{0, 1, 2, 3, 4},
		Slice: []string{"struct.slice.0", "struct.slice.1", "struct.slice.2"},
		Map:   map[string]string{"key1": "value1", "key2": "value2"},
		Struct: node{
			Key:   "key",
			Value: "value",
		},
		Ptr: &node{
			Key:   "ptrKey",
			Value: "ptrValue",
		},
		Key:         "test.Key",
		Interface:   map[string]string{"hello": "world"},
		Interface02: "interface02",
	}
	test.Struct.SetInternal("internal")
	test.Ptr.SetInternal("ptrInternal")
	return test
}

func TestGet(t *testing.T) {
	test := getTestData()

	type item struct {
		Path    string
		Val     interface{}
		IsError bool
	}
	tests := []item{
		item{Path: "Struct.Key", Val: test.Struct.Key, IsError: false},
		item{Path: "Struct.internal", IsError: true},
		item{Path: "Ptr.Key", Val: test.Ptr.Key, IsError: false},
		item{Path: "Array.0", Val: test.Array[0], IsError: false},
		item{Path: "Array.10", IsError: true},
		item{Path: "Slice.2", Val: test.Slice[2], IsError: false},
		item{Path: "Map.key1", Val: test.Map["key1"], IsError: false},
		item{Path: "Map.key3", IsError: true},
		item{Path: "Key", Val: test.Key, IsError: false},
		item{Path: "Interface.hello", Val: "world", IsError: false},
		item{Path: "Interface.world", IsError: true},
		item{Path: "Interface02", Val: "interface02", IsError: false},
	}

	for k, v := range tests {
		val, err := Get(test, v.Path)
		if err != nil && !v.IsError {
			t.Errorf("%d got unexpected error: %s", k, err)
		}
		if v.IsError {
			continue
		}
		if !reflect.DeepEqual(val, v.Val) {
			t.Errorf("%d expected: %+v, got: %+v", k, v.Val, val)
		}
	}
}

func TestSet(t *testing.T) {
	type item struct {
		Path    string
		Val     interface{}
		IsError bool
	}
	tests := []item{
		item{Path: "Struct.Key", Val: "test.struct.key", IsError: false},
		item{Path: "Struct.internal", Val: "test.struct.internal", IsError: true},
		item{Path: "Ptr.Key", Val: "test.ptr.key", IsError: false},
		item{Path: "Array.0", Val: 100, IsError: false},
		item{Path: "Array.10", IsError: true},
		item{Path: "Slice.2", Val: "test.Slice.2", IsError: false},
		item{Path: "Map.key1", Val: "test.map.key1", IsError: false},
		item{Path: "Map.key3", Val: "test.map,key3", IsError: false},
		item{Path: "Key", Val: "test.key", IsError: false},
		item{Path: "Interface.hello", Val: "test.Interface.hello", IsError: false},
		item{Path: "Interface.world", Val: "test.Interface.world", IsError: false},
		item{Path: "Interface02", Val: "test.Interface02", IsError: false},
	}
	for k, v := range tests {
		test := getTestData()
		err := Set(&test, v.Path, v.Val)
		if err != nil && !v.IsError {
			t.Errorf("%d got unexpected error: %s", k, err)
		}
		if v.IsError {
			continue
		}
		val, err := Get(test, v.Path)
		if err != nil {
			t.Errorf("%d got unexpected error: %s", k, err)
		}
		if !reflect.DeepEqual(val, v.Val) {
			t.Errorf("%d expected: %+v, got: %+v", k, v.Val, val)
		}
	}
}
