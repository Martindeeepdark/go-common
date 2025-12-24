package conv

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ToString converts any value to string
func ToString(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case *time.Time:
		if v != nil {
			return v.Format(time.RFC3339)
		}
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts any value to int64
func ToInt(value any) (int64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ToInt64 converts any value to int64, returns 0 on error
func ToInt64(value any) int64 {
	result, err := ToInt(value)
	if err != nil {
		return 0
	}
	return result
}

// ToFloat64 converts any value to float64
func ToFloat64(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToFloat64Default converts any value to float64, returns 0 on error
func ToFloat64Default(value any) float64 {
	result, err := ToFloat64(value)
	if err != nil {
		return 0
	}
	return result
}

// ToBool converts any value to bool
func ToBool(value any) (bool, error) {
	if value == nil {
		return false, nil
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(value).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(value).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(value).Float() != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ToBoolDefault converts any value to bool, returns false on error
func ToBoolDefault(value any) bool {
	result, err := ToBool(value)
	if err != nil {
		return false
	}
	return result
}

// ToEphemeralTime converts any value to time.Time
func ToTime(value any) (time.Time, error) {
	if value == nil {
		return time.Time{}, nil
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case *time.Time:
		if v != nil {
			return *v, nil
		}
		return time.Time{}, nil
	case string:
		// Try common time formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"2006-01-02",
		}
		for _, format := range formats {
			t, err := time.Parse(format, v)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time string: %s", v)
	case int64:
		return time.Unix(v, 0), nil
	case float64:
		sec := int64(v)
		nsec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", value)
	}
}

// ToTimeDefault converts any value to time.Time, returns zero time on error
func ToTimeDefault(value any) time.Time {
	result, err := ToTime(value)
	if err != nil {
		return time.Time{}
	}
	return result
}
