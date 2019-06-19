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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"svr-frame/libs"
)

var (
	log bool
)

func SetLog(status bool) {
	log = status
}

func getFieldTag(fieldType reflect.StructField) string {
	fieldName := fieldType.Tag.Get("json")
	if len(fieldName) < 1 || fieldName == "-" {
		fieldName = fieldType.Name
		return libs.LowCaseString(fieldName)
	}

	if trr := strings.Split(fieldName, ","); len(trr) > 1 {
		fieldName = strings.TrimSpace(trr[0]) // 过滤掉omitempty
	}
	return fieldName
}

func getDeepIndident(deep int) string {
	if deep > 0 {
		return strings.Repeat("    ", deep)
	}
	return ""
}

func printlog(args ...interface{}) {
	if log {
		fmt.Println(args...)
	}
}

func DeepCopy(i interface{}, from interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(libs.Interface2String(r))
			printlog(r)
		}
	}()

	inst := reflect.ValueOf(i)
	if inst.Kind() == reflect.Ptr {
		err = ValueDeepCopy(inst.Elem(), from, 0, "")
	} else {
		err = errors.New("not ptr type")
	}
	return
}

// 将泛型数据map[string]interface{}, 通过reflect深度拷贝到对应的结构体中
// 如果直接调用此方法，需要外部捕获panic
func ValueDeepCopy(inst reflect.Value, from interface{}, deep int, fieldName string) (err error) {
	if !inst.CanSet() {
		return errors.New("target cannt be set")
	}

	// printlog("target name>>:", inst.Type().String(), inst.Kind())
	switch inst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := int64(libs.Interface2Int(from))
		inst.SetInt(v)
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
		printlog(getDeepIndident(deep), "Ptr>>:", it.String())

		ValueDeepCopy(it.Elem(), from, deep+1, fieldName)
		inst.Set(it)
	case reflect.Struct:
		if mp, ok := from.(map[string]interface{}); ok {
			printlog(getDeepIndident(deep), "Struct>>:", inst.String())

			tpe := inst.Type()
			for i := 0; i < tpe.NumField(); i += 1 {
				fieldType := inst.Type().Field(i)
				field := inst.Field(i)

				fieldName := getFieldTag(fieldType)
				fieldValue, ok := mp[fieldName]
				if !ok || fieldValue == nil {
					continue
				}
				ValueDeepCopy(field, fieldValue, deep+1, fieldName)
				// printlog(getDeepIndident(deep+1),"field name:", fieldName, "value:", field.Interface(), "kind:",field.Kind())
			}
		}
		return
	case reflect.Map:
		if vv, ok := from.(map[string]interface{}); ok {
			mp := reflect.MakeMap(inst.Type())
			printlog(getDeepIndident(deep), "Map>>:", mp.String())

			mapValueDeepCopy(mp, vv, deep+1, fieldName)
			inst.Set(mp)
		}
	case reflect.Slice:
		if vv, ok := from.([]interface{}); ok {
			sl := reflect.MakeSlice(inst.Type(), len(vv), cap(vv))
			printlog(getDeepIndident(deep), "Slice>>:", sl.String())

			sliceValueDeepCopy(sl, vv, deep+1, fieldName)
			inst.Set(sl)
		}
	case reflect.Array, reflect.Chan, reflect.Func, reflect.UnsafePointer: // 不处理
	}
	printlog(getDeepIndident(deep), "field:", fieldName, "value:", inst.Interface(), "kind:", inst.Kind())
	return
}

func mapValueDeepCopy(inst reflect.Value, data map[string]interface{}, deep int, fieldName string) {
	if !inst.IsValid() || inst.Kind() != reflect.Map {
		return
	}

	kind := inst.Type().Elem().Kind()
	// printlog(inst.String(), kind)

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
			printlog(getDeepIndident(deep), "Struct>>:", val.String())

			ValueDeepCopy(val, v, deep+1, "fieldName")
		case reflect.Ptr: // map[string]*TestStruct
			val = reflect.New(inst.Type().Elem().Elem())
			printlog(getDeepIndident(deep), "Ptr>>:", val.String())

			ValueDeepCopy(val.Elem(), v, deep+1, fieldName)
		case reflect.Map: // map[string]map[string]interface{}
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printlog(getDeepIndident(deep), "Map>>:", val.String())

				mapValueDeepCopy(val, vv, deep+1, fieldName)
			} else {
				continue
			}
		case reflect.Slice: // map[string][]interface{}
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printlog(getDeepIndident(deep), "Slice>>:", val.String())

				sliceValueDeepCopy(val, vv, deep+1, fieldName)
			} else {
				continue
			}
		}
		inst.SetMapIndex(reflect.ValueOf(k), val)
		printlog(getDeepIndident(deep), "Map key:", k, "value:", val.Interface())
	}
}

func sliceValueDeepCopy(inst reflect.Value, slice []interface{}, deep int, fieldName string) {
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
		case reflect.Struct:
			printlog(getDeepIndident(deep), "Struct>>:", item.String())

			ValueDeepCopy(item, v, deep+1, fieldName)
			continue
		case reflect.Ptr:
			val = reflect.New(inst.Type().Elem().Elem())
			printlog(getDeepIndident(deep), "Ptr>>:", val.String())

			ValueDeepCopy(val.Elem(), v, deep+1, fieldName)
		case reflect.Map:
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printlog(getDeepIndident(deep), "Map>>:", val.String())

				mapValueDeepCopy(val, vv, deep+1, fieldName)
			}
		case reflect.Slice:
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printlog(getDeepIndident(deep), "Slice>>:", val.String())

				sliceValueDeepCopy(val, vv, deep+1, fieldName)
			}
		}
		item.Set(val)
		printlog(getDeepIndident(deep), "Slice index:", i, "value:", val.Interface())
	}
}
