/**
 * @version: v0.1.0
 * @author: zhangguodong
 * @license: LGPL v3
 * @contact: zhangguodong@dobest.com
 * @site: https://code.dobest.com/research-go
 * @software: GoLand
 * @file: conv.go
 * @time: 2022/8/17 17:16
 * @project: deepcopy
 */

package dcopy

import (
	"fmt"
	"reflect"
	"strconv"
)

func interface2String(v interface{}) string {
	if v == nil {
		return ""
	}
	switch d := v.(type) {
	case string:
		return d
	default:
		return fmt.Sprintf("%v", d)
	}
	// return ""
}

func interface2Int(v interface{}) int {
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

func interface2Int64(v interface{}) int64 {
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

func interface2Uint64(v interface{}) uint64 {
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
func interface2Float64(v interface{}) float64 {
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

func interface2Bool(v interface{}) bool {
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
