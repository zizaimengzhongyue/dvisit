package dvisit

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	errorData = "data is error: %s"
	errorPath = "path is error: %s"
)

//Get data 为原始数据，path 是要访问的数据路径，path 应该是“p1.p2.0.key”这种格式的字符串
func Get(data interface{}, path string) (interface{}, error) {
	val := reflect.ValueOf(data)
	paths := strings.Split(path, ".")
	return get(val, paths, "")
}

func get(data reflect.Value, paths []string, path string) (interface{}, error) {
	if len(paths) == 0 {
		if data.CanInterface() {
			return data.Interface(), nil
		}
		return nil, fmt.Errorf(errorData, path)
	}
	kind := data.Kind()
	switch kind {
	case reflect.Struct:
		return getStruct(data, paths, path)
	case reflect.Map:
		return getMap(data, paths, path)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return getSlice(data, paths, path)
	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		elem := data.Elem()
		return get(elem, paths, path)
	default:
		return data.Interface(), nil
	}
}

func getStruct(data reflect.Value, paths []string, path string) (interface{}, error) {
	nextPath := path + "." + paths[0]
	field := data.FieldByName(paths[0])
	if !field.IsValid() {
		return nil, fmt.Errorf(errorData, nextPath)
	}
	return get(field, paths[1:], nextPath)
}

func getMap(data reflect.Value, paths []string, path string) (interface{}, error) {
	nextPath := path + "." + paths[0]
	key := reflect.ValueOf(paths[0])
	field := data.MapIndex(key)
	if !field.IsValid() {
		return nil, fmt.Errorf(errorData, nextPath)
	}
	return get(field, paths[1:], nextPath)
}

func getSlice(data reflect.Value, paths []string, path string) (interface{}, error) {
	nextPath := path + "." + paths[0]
	index, err := strconv.Atoi(paths[0])
	if err != nil {
		return nil, fmt.Errorf(errorData, nextPath)
	}
	if data.Len() <= index {
		return nil, fmt.Errorf(errorData, nextPath)
	}
	field := data.Index(index)
	return get(field, paths[1:], nextPath)
}
