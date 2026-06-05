package utils

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// StringToBytes converts text to byte slice without memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to text without memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToString 尝试将值转换成字符串
func ToString(value any) (string, error) {
	switch reply := value.(type) {
	case string:
		return reply, nil
	case []byte:
		return BytesToString(reply), nil
	}

	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return BytesToString(bs), nil
}

func ToBytes(value any) ([]byte, error) {
	var data []byte
	switch t := value.(type) {
	case string:
		data = StringToBytes(t)
	case []byte:
		data = t
	default:
		v, err := json.Marshal(value)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal value, value: %v", value)
		}
		data = v
	}
	return data, nil
}

// ToSlice 将传过来的数据转换成[]any
func ToSlice(data any) []any {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	sliceLength := v.Len()
	sliceData := make([]any, sliceLength)
	for i := 0; i < sliceLength; i++ {
		sliceData[i] = v.Index(i).Interface()
	}

	return sliceData
}

func CsvToNullNumbers[T Numeric](s string) []T {
	numbers := CsvToNumbers[T](s)
	if len(numbers) > 0 {
		return numbers
	}
	return nil
}

// CsvToNumbers converts a comma-separated string to a slice of the specified numeric type.
// Returns the slice and any conversion error encountered.
func CsvToNumbers[T Numeric](s string) []T {
	// Handle empty string case
	if s == "" {
		return []T{}
	}

	// Split the string by commas
	strValues := strings.Split(s, ",")
	result := make([]T, len(strValues))

	// Convert each string element to the numeric type T
	switch any(*new(T)).(type) {
	case int, int8, int16, int32, int64:
		for i, str := range strValues {
			val, _ := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
			result[i] = T(val)
		}
	case uint, uint8, uint16, uint32, uint64:
		for i, str := range strValues {
			val, _ := strconv.ParseUint(strings.TrimSpace(str), 10, 64)
			result[i] = T(val)
		}
	case float32, float64:
		for i, str := range strValues {
			val, _ := strconv.ParseFloat(strings.TrimSpace(str), 64)
			result[i] = T(val)
		}
	}
	return result
}

// NumbersToCsv 将int64 slice转换成用逗号分隔的字符串: 1,2,3
func NumbersToCsv[T Numeric](numbers []T) string {
	return strings.Join(NumbersToStrings(numbers), ",")
}

// NumbersToStrings converts a slice of numeric types to a slice of strings.
func NumbersToStrings[T Numeric](numbers []T) []string {
	result := make([]string, len(numbers))
	for i, v := range numbers {
		// Use a type switch for efficient, type-specific formatting.
		switch any(v).(type) {
		case int, int8, int16, int32, int64:
			result[i] = strconv.FormatInt(int64(v), 10)
		case uint, uint8, uint16, uint32, uint64:
			result[i] = strconv.FormatUint(uint64(v), 10)
		case float32, float64:
			result[i] = strconv.FormatFloat(float64(v), 'f', -1, 64)
		}
	}
	return result
}

// StringsToNumbers converts a string slice to numeric slice with error handling
func StringsToNumbers[T Numeric](strSlice []string) []T {
	result := make([]T, len(strSlice))

	for i, str := range strSlice {
		// Trim whitespace from each string element
		str = strings.TrimSpace(str)

		// Use type switching based on the target type T
		switch any(*new(T)).(type) {
		case int, int8, int16, int32, int64:
			val, _ := strconv.ParseInt(str, 10, 64)
			result[i] = T(val)
		case uint, uint8, uint16, uint32, uint64:
			val, _ := strconv.ParseUint(str, 10, 64)
			result[i] = T(val)
		case float32, float64:
			val, _ := strconv.ParseFloat(str, 64)
			result[i] = T(val)
		}
	}

	return result
}

func SafeGet[FieldType comparable](obj any, field FieldType) FieldType {
	if obj == nil {
		var ret FieldType
		return ret
	}

	return field
}
