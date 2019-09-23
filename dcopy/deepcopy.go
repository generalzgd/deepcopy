/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: deepcopy.go
 * @time: 2019/6/11 10:13
 */
package dcopy

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/toolkits/slice"

	// libs "github.com/generalzgd/deepcopy/libtools"

	libs `github.com/generalzgd/comm-libs`
)

var (
	log bool
)

func SetLog(status bool) {
	log = status
}

// 获取字段名，优先使用json tag, 然后使用xorm tag, 如果没有则使用字段名的小驼峰格式
func getFieldTag(fieldType reflect.StructField) (string, bool) {
	fieldName := fieldType.Tag.Get("json")
	omitempty := false
	if len(fieldName) < 1 || fieldName == "-" {
		//
		fieldName = fieldType.Tag.Get("xorm")
		if len(fieldName) > 0 && fieldName != "extends" {
			return fieldName, false
		}
		//
		fieldName = fieldType.Name
		return libs.LowCaseString(fieldName), false
	}

	if trr := strings.Split(fieldName, ","); len(trr) > 1 {
		fieldName = strings.TrimSpace(trr[0]) // 过滤掉omitempty
		omitempty = trr[1] == "omitempty"
	}
	return fieldName, omitempty
}

func getDeepIndent(deep int) string {
	if deep > 0 {
		return strings.Repeat("    ", deep)
	}
	return ""
}

func printLog(deep int, args ...interface{}) {
	if log {
		tmp := make([]interface{}, 0, len(args)+1)
		tmp = append(tmp, getDeepIndent(deep))
		tmp = append(tmp, args...)
		fmt.Println(tmp...)
	}
}

func InstanceFromBytes(dest interface{}, from []byte) (err error) {
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(from, &tmp); err != nil {
		return err
	}
	return InstanceFromMap(dest, tmp)
}

func InstanceFromMap(dest interface{}, from interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(libs.Interface2String(r))
			printLog(0, r)
		}
	}()

	inst := reflect.ValueOf(dest)
	if inst.Kind() == reflect.Ptr {
		err = valueDeepCopy(inst.Elem(), from, 0, "")
	} else {
		err = errors.New("not ptr type")
	}
	return
}

func InstanceValueFromMap(dest reflect.Value, from interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(libs.Interface2String(r))
			printLog(0, r)
		}
	}()
	if dest.Kind() == reflect.Ptr {
		err = valueDeepCopy(dest.Elem(), from, 0, "")
	} else {
		err = valueDeepCopy(dest, from, 0, "")
	}
	return
}

// 将泛型数据map[string]interface{}, 通过reflect深度拷贝到对应的结构体中
// 如果直接调用此方法，需要外部捕获panic
func valueDeepCopy(inst reflect.Value, from interface{}, deep int, fieldName string) (err error) {
	if !inst.CanSet() {
		return errors.New("target cannt be set")
	}

	// printlog("target name>>:", inst.Type().String(), inst.Kind())
	switch inst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := libs.Interface2Int64(from)
		inst.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := libs.Interface2Uint64(from)
		inst.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v := libs.Interface2Float64(from)
		inst.SetFloat(v)
	case reflect.String:
		v := libs.Interface2String(from)
		inst.SetString(v)
	case reflect.Bool:
		v := libs.Interface2Bool(from)
		inst.SetBool(v)
	case reflect.Interface:
		val := reflect.ValueOf(from)
		inst.Set(val)
	case reflect.Ptr:
		it := reflect.New(inst.Type().Elem())
		printLog(deep, "Ptr>>:", it.String())

		err = valueDeepCopy(it.Elem(), from, deep+1, fieldName)
		if err != nil {
			return
		}
		inst.Set(it)
	case reflect.Struct:
		if mp, ok := from.(map[string]interface{}); ok {
			printLog(deep, "Struct>>:", inst.String())

			tpe := inst.Type()
			for i := 0; i < tpe.NumField(); i += 1 {
				fieldType := inst.Type().Field(i)
				field := inst.Field(i)

				fieldName, _ := getFieldTag(fieldType)
				fieldValue, ok := mp[fieldName]
				if !ok || fieldValue == nil {
					continue
				}
				err = valueDeepCopy(field, fieldValue, deep+1, fieldName)
				if err != nil {
					return
				}
				// printlog(getDeepIndident(deep+1),"field name:", fieldName, "value:", field.Interface(), "kind:",field.Kind())
			}
		}
		return
	case reflect.Map:
		if vv, ok := from.(map[string]interface{}); ok {
			mp := reflect.MakeMap(inst.Type())
			printLog(deep, "Map>>:", mp.String())

			err = mapValueDeepCopy(mp, vv, deep+1, fieldName)
			if err != nil {
				return
			}
			inst.Set(mp)
		}
	case reflect.Slice:
		if vv, ok := from.([]interface{}); ok {
			sl := reflect.MakeSlice(inst.Type(), len(vv), cap(vv))
			printLog(deep, "Slice>>:", sl.String())

			err = sliceValueDeepCopy(sl, vv, deep+1, fieldName)
			if err != nil {
				return
			}
			inst.Set(sl)
		}
	case reflect.Array, reflect.Chan, reflect.Func, reflect.UnsafePointer: // 不处理
	}
	printLog(deep, "field:", fieldName, "value:", inst.Interface(), "kind:", inst.Kind())
	return
}

func mapValueDeepCopy(inst reflect.Value, data map[string]interface{}, deep int, fieldName string) (err error) {
	if !inst.IsValid() || inst.Kind() != reflect.Map {
		return
	}

	kind := inst.Type().Elem().Kind()
	// printLog(inst.String(), kind)

	for k, v := range data {
		var val reflect.Value
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = reflect.ValueOf(libs.Interface2Int(v))
		case reflect.Float32, reflect.Float64:
			val = reflect.ValueOf(libs.Interface2Float64(v))
		case reflect.String:
			val = reflect.ValueOf(libs.Interface2String(v))
		case reflect.Bool:
			val = reflect.ValueOf(libs.Interface2Bool(v))
		case reflect.Interface:
			val = reflect.ValueOf(v)
		case reflect.Struct: // map[string]TestStruct
			val = reflect.New(inst.Type().Elem()).Elem()
			printLog(deep, "Struct>>:", val.String())

			err = valueDeepCopy(val, v, deep+1, "fieldName")
			if err != nil {
				return
			}
		case reflect.Ptr: // map[string]*TestStruct
			val = reflect.New(inst.Type().Elem().Elem())
			printLog(deep, "Ptr>>:", val.String())

			err = valueDeepCopy(val.Elem(), v, deep+1, fieldName)
			if err != nil {
				return
			}
		case reflect.Map: // map[string]map[string]interface{}
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printLog(deep, "Map>>:", val.String())

				err = mapValueDeepCopy(val, vv, deep+1, fieldName)
				if err != nil {
					return
				}
			} else {
				continue
			}
		case reflect.Slice: // map[string][]interface{}
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printLog(deep, "Slice>>:", val.String())

				err = sliceValueDeepCopy(val, vv, deep+1, fieldName)
				if err != nil {
					return
				}
			} else {
				continue
			}
		}
		inst.SetMapIndex(reflect.ValueOf(k), val)
		printLog(deep, "Map key:", k, "value:", val.Interface())
	}
	return
}

func sliceValueDeepCopy(inst reflect.Value, slice []interface{}, deep int, fieldName string) (err error) {
	if !inst.IsValid() || inst.Kind() != reflect.Slice {
		return
	}
	kind := inst.Type().Elem().Kind()
	// printlog(inst.String(), kind)

	for i, v := range slice {
		item := inst.Index(i)
		var val reflect.Value
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = reflect.ValueOf(libs.Interface2Int(v))
		case reflect.Float32, reflect.Float64:
			val = reflect.ValueOf(libs.Interface2Float64(v))
		case reflect.String:
			val = reflect.ValueOf(libs.Interface2String(v))
		case reflect.Bool:
			val = reflect.ValueOf(libs.Interface2Bool(v))
		case reflect.Interface:
			val = reflect.ValueOf(v)
		case reflect.Struct: // []struct{}
			printLog(deep, "Struct>>:", item.String())

			err = valueDeepCopy(item, v, deep+1, fieldName)
			if err != nil {
				return
			}
			continue
		case reflect.Ptr: // []Ptr(*int|*struct)
			val = reflect.New(inst.Type().Elem().Elem())
			printLog(deep, "Ptr>>:", val.String())

			err = valueDeepCopy(val.Elem(), v, deep+1, fieldName)
			if err != nil {
				return
			}
		case reflect.Map: // []map[string]interface
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printLog(deep, "Map>>:", val.String())

				err = mapValueDeepCopy(val, vv, deep+1, fieldName)
				if err != nil {
					return
				}
			}
		case reflect.Slice: // [][]interface
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printLog(deep, "Slice>>:", val.String())

				err = sliceValueDeepCopy(val, vv, deep+1, fieldName)
				if err != nil {
					return
				}
			}
		}
		item.Set(val)
		printLog(deep, "Slice index:", i, "value:", val.Interface())
	}
	return
}

//
func InstanceToMap(from interface{}) (out map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(libs.Interface2String(r))
			printLog(0, r)
		}
	}()

	inst := reflect.ValueOf(from)
	kind := inst.Kind()
	numField := 0
	if kind == reflect.Ptr {
		kind = inst.Elem().Kind()
		numField = inst.Elem().NumField()
	} else {
		numField = inst.NumField()
	}
	// handleType := []interface{}{reflect.Struct, reflect.Map, reflect.Slice}
	if !slice.Contains(structTypes, kind) {
		return nil, errors.New("only process struct/map/slice type")
	}
	out = make(map[string]interface{}, numField)
	err = instanceToMap(out, inst, 0)
	return out, err
}

var (
	structTypes = []interface{}{reflect.Struct, reflect.Map, reflect.Slice}
	// NothingTypes = []interface{}{reflect.Array, reflect.Chan, reflect.Func, reflect.UnsafePointer}
)

func instanceToMap(dest map[string]interface{}, from reflect.Value, deep int) (err error) {
	if from.Kind() == reflect.Ptr {
		return instanceToMap(dest, from.Elem(), deep)
	}

	for i := 0; i < from.NumField(); i++ {
		field := from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, omitempty := getFieldTag(fieldType)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		printLog(deep, "kind:", field.Kind(), "fieldName:", fieldName, "value:", field.Interface(), "omitempty:", omitempty)
		switch field.Kind() {
		case reflect.Struct:
			subMap := make(map[string]interface{}, field.NumField())
			dest[fieldName] = subMap
			if err = instanceToMap(subMap, field, deep+1); err != nil {
				return
			}
		case reflect.Map:
			keys := field.MapKeys()
			if len(keys) == 0 && omitempty {
				continue
			}
			subMap := make(map[string]interface{}, len(keys))
			dest[fieldName] = subMap
			if err = instanceMapToMap(subMap, field, deep+1); err != nil {
				return
			}
		case reflect.Slice:
			if field.Len() == 0 && omitempty {
				continue
			}
			subSlice := make([]interface{}, field.Len())
			dest[fieldName] = subSlice
			if err = instanceSliceToArr(subSlice, field, deep+1); err != nil {
				return
			}
		default:
			if valueEmpty(field.Interface()) && omitempty {
				continue
			}
			dest[fieldName] = field.Interface()
		}
	}
	return
}

// 结构体转成map，暂时不支持map/slice类型的字段
func instanceMapToMap(dest map[string]interface{}, field reflect.Value, deep int) error {
	// inst := reflect.ValueOf(from)
	if field.Kind() != reflect.Map {
		return errors.New("field type is not map")
	}

	keys := field.MapKeys()
	for _, key := range keys {
		keyStr := libs.Interface2String(key.Interface())
		subField := field.MapIndex(key)
		if subField.Kind() == reflect.Ptr {
			subField = subField.Elem()
		}
		printLog(deep, "kind:", subField.Kind(), "key:", keyStr, "value:", subField.Interface())
		switch subField.Kind() {
		case reflect.Struct:
			subMap := make(map[string]interface{}, subField.NumField())
			dest[keyStr] = subMap
			if err := instanceToMap(subMap, subField, deep+1); err != nil {
				return err
			}
		case reflect.Map:
			keys := subField.MapKeys()
			subMap := make(map[string]interface{}, len(keys))
			dest[keyStr] = subMap
			if err := instanceMapToMap(subMap, subField, deep+1); err != nil {
				return err
			}
		case reflect.Slice:
			subSlice := make([]interface{}, subField.Len())
			dest[keyStr] = subSlice
			if err := instanceSliceToArr(subSlice, subField, deep+1); err != nil {
				return err
			}
		default:
			dest[keyStr] = subField.Interface()
		}
	}
	return nil
}

func instanceSliceToArr(dest []interface{}, field reflect.Value, deep int) error {
	if field.Kind() != reflect.Slice {
		return errors.New("field type is not slice")
	}

	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		printLog(deep, "kind:", item.Kind(), "index:", i, "value:", item.Interface())
		switch item.Kind() {
		case reflect.Struct:
			subMap := make(map[string]interface{}, item.NumField())
			dest[i] = subMap
			if err := instanceToMap(subMap, item, deep+1); err != nil {
				return err
			}
		case reflect.Map:
			keys := item.MapKeys()
			subMap := make(map[string]interface{}, len(keys))
			dest[i] = subMap
			if err := instanceMapToMap(subMap, item, deep+1); err != nil {
				return err
			}
		case reflect.Slice:
			subSlice := make([]interface{}, item.Len())
			dest[i] = subSlice
			if err := instanceSliceToArr(subSlice, item, deep+1); err != nil {
				return err
			}
		default:
			dest[i] = item.Interface()
		}
	}
	return nil
}

func valueEmpty(v interface{}) bool {
	switch d := v.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() == 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 0
	case string:
		return reflect.ValueOf(v).String() == ""
	case float64, float32:
		return reflect.ValueOf(v).Float() == 0.0
	case bool:
		return reflect.ValueOf(v).Bool() == false
	default:
		return d == nil
	}
}
