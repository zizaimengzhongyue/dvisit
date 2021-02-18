package dvisit

import (
	"errors"
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
	if err != nil || data.Len() <= index {
		return nil, fmt.Errorf(errorData, nextPath)
	}
	field := data.Index(index)
	return get(field, paths[1:], nextPath)
}

//Set data 为原始数据，path 是要访问的数据路径，path 应该是“p1.p2.0.key”这种格式的字符串
func Set(data interface{}, path string, value interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.CanAddr() {
		return errors.New("value is not addressable")
	}
	paths := strings.Split(path, ".")
	v := reflect.ValueOf(value)
	return set(val, paths, v, "")
}

func set(data reflect.Value, paths []string, v reflect.Value, path string) error {
	if len(paths) == 0 {
		if data.CanSet() && data.Kind() == v.Kind() {
			data.Set(v)
			return nil
		}
		return fmt.Errorf(errorData, path)
	}
	kind := data.Kind()
	switch kind {
	case reflect.Struct:
		return setStruct(data, paths, v, path)
	case reflect.Map:
		return setMap(data, paths, v, path)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return setArray(data, paths, v, path)
	case reflect.Ptr:
		elem := data.Elem()
		return set(elem, paths, v, path)
	default:
		return fmt.Errorf(errorPath, path)
	}
}

func setStruct(data reflect.Value, paths []string, v reflect.Value, path string) error {
	nextPath := path + "." + paths[0]
	field := data.FieldByName(paths[0])
	if !field.IsValid() {
		return fmt.Errorf(errorPath, path)
	}
	return set(field, paths[1:], v, nextPath)
}

func setMap(data reflect.Value, paths []string, v reflect.Value, path string) error {
	nextPath := path + "." + paths[0]
	// 长度为 1 的时候 map 数据可以直接写了
	if len(paths) == 1 {
		k := reflect.ValueOf(paths[0])
		typK, typV := data.Type().Key(), data.Type().Elem()
		if typK.Kind() == k.Kind() && typV.Kind() == v.Kind() {
			data.SetMapIndex(k, v)
			return nil
		}
		return fmt.Errorf(errorData, path)
	}

	key := reflect.ValueOf(paths[0])
	field := data.MapIndex(key)
	if !field.IsValid() {
		return fmt.Errorf(errorPath, path)
	}
	return set(field, paths[1:], v, nextPath)
}

func setArray(data reflect.Value, paths []string, v reflect.Value, path string) error {
	nextPath := path + "." + paths[0]
	index, err := strconv.Atoi(paths[0])
	if err != nil || data.Len() <= index {
		return fmt.Errorf(errorData, path)
	}
	field := data.Index(index)
	return set(field, paths[1:], v, nextPath)
}
