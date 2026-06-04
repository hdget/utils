package sql

import (
	"database/sql"
	"reflect"
	"strconv"
	"time"

	jsonUtils "github.com/hdget/utils/json"
	"github.com/sqlc-dev/pqtype"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func GetNullString(filters map[string]string, key string) sql.NullString {
	if v, ok := filters[key]; ok {
		return sql.NullString{String: v, Valid: true}
	}

	return sql.NullString{}
}

func GetNullInt32(filters map[string]string, key string) sql.NullInt32 {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseInt(v, 10, 32); err == nil {
			return sql.NullInt32{Int32: int32(vv), Valid: true}
		}
	}

	return sql.NullInt32{}
}

func GetNullBool(filters map[string]string, key string) sql.NullBool {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseBool(v); err == nil {
			return sql.NullBool{Bool: vv, Valid: true}
		}
	}

	return sql.NullBool{}
}

func GetNullFloat64(filters map[string]string, key string) sql.NullFloat64 {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseFloat(v, 64); err == nil {
			return sql.NullFloat64{Float64: vv, Valid: true}
		}
	}

	return sql.NullFloat64{}
}

func GetNullByte(filters map[string]string, key string) sql.NullByte {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseUint(v, 10, 8); err == nil {
			return sql.NullByte{Byte: uint8(vv), Valid: true}
		}

		if len(v) == 1 {
			return sql.NullByte{Byte: v[0], Valid: true}
		}
	}

	return sql.NullByte{}
}

func GetNullInt16(filters map[string]string, key string) sql.NullInt16 {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseInt(v, 10, 16); err == nil {
			return sql.NullInt16{Int16: int16(vv), Valid: true}
		}
	}

	return sql.NullInt16{}
}

func GetNullInt64(filters map[string]string, key string) sql.NullInt64 {
	if v, ok := filters[key]; ok {
		if vv, err := strconv.ParseInt(v, 10, 64); err == nil {
			return sql.NullInt64{Int64: vv, Valid: true}
		}
	}

	return sql.NullInt64{}
}

func GetNullTime(filters map[string]string, key string) sql.NullTime {
	v, exists := filters[key]
	if !exists || v == "" {
		return sql.NullTime{}
	}

	if t, valid := parseTimeFromString(v); valid {
		return sql.NullTime{Time: t, Valid: true}
	}

	// 然后timestamp
	if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
		if ts >= 1e12 && ts < 1e15 { // 毫秒级范围
			return sql.NullTime{Time: time.UnixMilli(ts), Valid: true}
		} else if ts >= 0 && ts < 1e12 { // 秒级范围
			return sql.NullTime{Time: time.Unix(ts, 0), Valid: true}
		}
	}

	return sql.NullTime{}
}

func ToNullString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *val, Valid: true}
}

func ToNullInt32[T Number](val *T) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: int32(*val), Valid: true}
}

func ToNullInt64[T Number](val *T) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*val), Valid: true}
}

func ToNullJsonObject(val any) pqtype.NullRawMessage {
	if !isNilAny(val) {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonObject(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

func ToNullJsonArray(val any) pqtype.NullRawMessage {
	if !isNilAny(val) {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonArray(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

func ToNullTime[T int64 | string | time.Time](val T) sql.NullTime {
	switch vv := any(val).(type) { // 对空接口进行类型断言
	case int64:
		if vv >= 1e12 && vv < 1e15 { // 毫秒级范围
			return sql.NullTime{Time: time.UnixMilli(vv), Valid: true}
		} else if vv >= 0 && vv < 1e12 { // 秒级范围
			return sql.NullTime{Time: time.Unix(vv, 0), Valid: true}
		}
	case string:
		if vv != "" {
			t, valid := parseTimeFromString(vv)
			return sql.NullTime{Time: t, Valid: valid}
		}
	case time.Time:
		if !vv.IsZero() {
			return sql.NullTime{Time: vv, Valid: true}
		}
	}

	return sql.NullTime{}
}

func isNilAny(val any) bool {
	if val == nil {
		return true
	}

	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

func parseTimeFromString(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	formats := []string{time.RFC3339, time.DateTime, time.DateOnly}
	for _, format := range formats {
		// 建议指定本地时区，避免 UTC 时差问题
		if t, err := time.ParseInLocation(format, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
