package json

import (
	"encoding/json"
	"reflect"

	"github.com/hdget/utils"
)

var (
	emptyJsonArray  = utils.StringToBytes("[]")
	emptyJsonObject = utils.StringToBytes("{}")
)

// IsEmptyJsonArray 是否是空json array
func IsEmptyJsonArray(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return isEqualSlice(data, emptyJsonArray)
}

// IsEmptyJsonObject 是否是空json object
func IsEmptyJsonObject(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	return isEqualSlice(data, emptyJsonObject)
}

// JsonArray 将slice转换成[]byte数据，如果slice为nil或空则返回空json array bytes
func JsonArray(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return emptyJsonArray
	}

	var jsonData []byte
	switch v := args[0].(type) {
	case string:
		if v != "" {
			jsonData = utils.StringToBytes(v)
		}
	case []byte:
		jsonData = v
	default:
		val := indirect(reflect.ValueOf(args[0]))
		if val.IsValid() && !val.IsZero() && val.Kind() == reflect.Slice {
			jsonData, _ = json.Marshal(args[0])
		}
	}

	if !json.Valid(jsonData) {
		return emptyJsonArray
	}
	return jsonData
}

// JsonObject 将object转换成[]byte数据，如果object为nil或空则返回空json object bytes
func JsonObject(args ...any) []byte {
	if len(args) == 0 || args[0] == nil {
		return emptyJsonObject
	}

	var jsonData []byte
	switch v := args[0].(type) {
	case string:
		if v != "" {
			jsonData = utils.StringToBytes(v)
		}
	case []byte:
		jsonData = v
	default:
		val := indirect(reflect.ValueOf(args[0]))
		if val.IsValid() && !val.IsZero() && (val.Kind() == reflect.Struct || val.Kind() == reflect.Map) {
			jsonData, _ = json.Marshal(args[0])
		}
	}

	if !json.Valid(jsonData) {
		return emptyJsonObject
	}
	return jsonData
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		return reflectValue.Elem()
	}
	return reflectValue
}

// Equals compare elements from the start to the end,
//
// if they are the same is considered the slices are equal if all elements are
// the same is considered the slices are equal
// if each slice == nil is considered that they're equal
//
// if element realizes Equals interface it uses that method, in other way uses
// default compare
func isEqualSlice[T comparable](ss []T, rhs []T) bool {
	if len(ss) != len(rhs) {
		return false
	}

	for i := range ss {
		if ss[i] != rhs[i] {
			return false
		}
	}

	return true
}
