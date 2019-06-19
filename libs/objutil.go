/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: Gogland
 * @file: objutil.go
 * @time: 2017/9/30 10:38
 */
package libs

import (
	"fmt"
	"reflect"
	"strconv"
)

func Interface2Int(v interface{}) int {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		v, _ := strconv.Atoi(d)
		return v
	case float32, float64:
		return int(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Int64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseInt(d, 10, 64)
		return t
	case float32, float64:
		return int64(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return int64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Uint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseUint(d, 10, 64)
		return t
	case float32, float64:
		return uint64(reflect.ValueOf(d).Float())
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return uint64(reflect.ValueOf(d).Uint())
	}
	return 0
}

// Float64 coerces into a float64
func Interface2Float64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch d := v.(type) {
	case string:
		t, _ := strconv.ParseFloat(d, 64)
		return t
	case float32, float64:
		return reflect.ValueOf(d).Float()
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(d).Uint())
	}
	return 0
}

func Interface2Bool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch d := v.(type) {
	case bool:
		return d
	case string:
		t, _ := strconv.ParseBool(d)
		return t
	case float32, float64:
		return reflect.ValueOf(d).Float() > 0.0
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(d).Int() > 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(d).Uint() > 0
	}
	return false
}

func Interface2String(v interface{}) string {
	if v == nil {
		return ""
	}
	switch d := v.(type) {
	case string:
		return d
	default:
		return fmt.Sprintf("%v", d)
	}
	return ""
}
