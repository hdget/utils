package sql

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/hdget/utils"
	jsonUtils "github.com/hdget/utils/json"
	"github.com/spf13/cast"
	"github.com/sqlc-dev/pqtype"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func GetNullString(filters map[string]string, key string) sql.NullString {
	if len(filters) > 0 {
		if v, ok := filters[key]; ok {
			return sql.NullString{String: v, Valid: true}
		}
	}
	return sql.NullString{}
}

func GetNullInt32(filters map[string]string, key string) sql.NullInt32 {
	if len(filters) > 0 {
		if v, ok := filters[key]; ok {
			return sql.NullInt32{Int32: cast.ToInt32(v), Valid: true}
		}
	}
	return sql.NullInt32{}
}

func GetInt32Slice(filters map[string]string, key string) []int32 {
	if len(filters) > 0 {
		if v, ok := filters[key]; ok {
			return utils.CsvToNumbers[int32](filters["stage"])
		}
	}
	return nil
}

func GetNullInt64(filters map[string]string, key string) sql.NullInt64 {
	if len(filters) > 0 {
		if v, ok := filters[key]; ok {
			return sql.NullInt64{Int64: cast.ToInt64(v), Valid: true}
		}
	}
	return sql.NullInt64{}
}

func GetNullTime(filters map[string]string, key string) sql.NullTime {
	if len(filters) > 0 {
		v, exists := filters[key]
		if !exists || v == "" {
			return sql.NullTime{}
		}

		formats := []string{
			time.DateTime, // 2006-01-02 15:04:05
			time.DateOnly, // 2006-01-02
			time.RFC3339,  // 2006-01-02T15:04:05Z07:00
		}

		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
		}

		// 然后timestamp
		sec, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return sql.NullTime{Time: time.Unix(sec, 0), Valid: true}
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
	if val != nil {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonObject(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

func ToNullJsonArray(val any) pqtype.NullRawMessage {
	if val != nil {
		return pqtype.NullRawMessage{
			RawMessage: jsonUtils.JsonArray(val),
			Valid:      true,
		}
	}
	return pqtype.NullRawMessage{}
}

func ToNullTime[T int64 | string | time.Time](val T) sql.NullTime {
	var v interface{} = val // 先转为空接口

	switch vv := v.(type) { // 对空接口进行类型断言
	case int64:
		if vv > 0 {
			return sql.NullTime{Time: time.Unix(vv, 0), Valid: true}
		}
	case string:
		if vv != "" {
			formats := []string{
				time.DateTime, // 2006-01-02 15:04:05
				time.DateOnly, // 2006-01-02
				time.RFC3339,  // 2006-01-02T15:04:05Z07:00
			}

			for _, format := range formats {
				if t, err := time.Parse(format, vv); err == nil {
					return sql.NullTime{Time: t, Valid: true}
				}
			}
		}
	case time.Time:
		if !vv.IsZero() {
			return sql.NullTime{Time: vv, Valid: true}
		}
	}

	return sql.NullTime{}
}
