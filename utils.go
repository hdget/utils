package utils

import (
	"encoding/json"
	"fmt"
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

// CsvToNumbers converts a comma-separated string to a slice of the specified numeric type.
// Returns the slice and any conversion error encountered.
func CsvToNumbers[T Numeric](s string) ([]T, error) {
	// Handle empty string case
	if s == "" {
		return []T{}, nil
	}

	// Split the string by commas
	strValues := strings.Split(s, ",")
	result := make([]T, len(strValues))

	// Convert each string element to the numeric type T
	for i, str := range strValues {
		// Trim any surrounding whitespace from each element
		str = strings.TrimSpace(str)

		var err error
		switch any(*new(T)).(type) {
		case int, int8, int16, int32, int64:
			var val int64
			val, err = strconv.ParseInt(str, 10, 64)
			result[i] = T(val)
		case uint, uint8, uint16, uint32, uint64:
			var val uint64
			val, err = strconv.ParseUint(str, 10, 64)
			result[i] = T(val)
		case float32, float64:
			var val float64
			val, err = strconv.ParseFloat(str, 64)
			result[i] = T(val)
		}
		if err != nil {
			return nil, fmt.Errorf("error converting '%s' at index %d: %v", str, i, err)
		}
	}
	return result, nil
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
func StringsToNumbers[T Numeric](strSlice []string) ([]T, error) {
	result := make([]T, len(strSlice))

	var zeroValue T
	for i, str := range strSlice {
		// Trim whitespace from each string element
		str = strings.TrimSpace(str)

		var err error
		var value interface{}

		// Use type switching based on the target type T
		switch any(*new(T)).(type) {
		case int, int8, int16, int32, int64:
			var val int64
			val, err = strconv.ParseInt(str, 10, 64)
			value = val
		case uint, uint8, uint16, uint32, uint64:
			var val uint64
			val, err = strconv.ParseUint(str, 10, 64)
			value = val
		case float32, float64:
			var val float64
			val, err = strconv.ParseFloat(str, 64)
			value = val
		}

		if err != nil {
			return nil, fmt.Errorf("error converting '%s' at index %d: %v", str, i, err)
		}

		// Type assertion to convert interface{} to T
		if v, ok := value.(T); ok {
			result[i] = v
		} else {
			result[i] = zeroValue
		}
	}

	return result, nil
}
