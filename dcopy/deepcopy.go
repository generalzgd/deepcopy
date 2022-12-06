// Package dcopy
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
	"time"

	"github.com/sirupsen/logrus"
	"github.com/toolkits/slice"
)

type FieldType int

const (
	FieldType_Idle FieldType = 0 + iota
	FieldType_Origin
	FieldType_Json
	FieldType_Xorm
	FieldType_Gorm
)

const (
	TimeValType_Int64  int8 = 0 + iota // 时间戳
	TimeValType_String                 // 格式化时间字符串
)

type args struct {
	curGetFieldType FieldType           // 字段名获取方式
	omitempty       bool                // 是否忽略0字段
	timeFmtStr      string              // time.Time类型转换格式
	timeValType     int8                // time.Time类型转换成timestamp还是字符串
	ignoreFieldMap  map[string]struct{} // 需要忽略的字段
	log             logrus.StdLogger    // 打印日志
}

var (
	defaultOptArgs = args{
		timeFmtStr:     "2006-01-02 15:04:05",
		timeValType:    TimeValType_String,
		ignoreFieldMap: map[string]struct{}{},
	}
)

type CopyOption func(*args)

// WithFieldType 字段名获取方式
func WithFieldType(tpe FieldType) CopyOption {
	return func(a *args) {
		a.curGetFieldType = tpe
	}
}

// WithOmitempty 原始字段如果空，是否忽略
func WithOmitempty(omitempty bool) CopyOption {
	return func(a *args) {
		a.omitempty = omitempty
	}
}

func WithTimeFormatStr(format string) CopyOption {
	return func(a *args) {
		a.timeFmtStr = format
	}
}

func WithTimeValType(valTpe int8) CopyOption {
	return func(a *args) {
		a.timeValType = valTpe
	}
}

// WithIgnoreFields 自动忽略大小写
func WithIgnoreFields(fieldNames ...string) CopyOption {
	return func(a *args) {
		tmp := make(map[string]struct{}, len(fieldNames))
		for _, name := range fieldNames {
			tmp[strings.ToLower(name)] = struct{}{}
		}
		a.ignoreFieldMap = tmp
	}
}

// WitLog 打印日志
func WitLog() CopyOption {
	return func(a *args) {
		a.log = logrus.StandardLogger()
	}
}

func parseJsonTag(tagStr string) (tagName string, omitempty, ignore bool) {
	if tagStr == "" {
		return
	}
	if tagStr == "-" {
		ignore = true
		return
	}
	tagName = tagStr
	arr := strings.Split(tagStr, ",")
	if len(arr) > 1 {
		tagName = arr[0]
		omitempty = arr[1] == "omitempty"
		return
	}
	return
}

func parseXormTag(tagStr string) (tagName string, omitempty, ignore bool) {
	if tagStr == "" {
		return
	}
	if tagStr == "-" {
		ignore = true
		return
	}
	tagName = tagStr
	// todo 过滤extends等关键字
	return
}

func parseGormTag(tagStr string) (tagName string, omitempty, ignore bool) {
	if tagStr == "" {
		return
	}
	if tagStr == "-" {
		ignore = true
		return
	}
	arr := strings.Split(tagStr, ";")
	for _, it := range arr {
		kv := strings.Split(it, ":")
		if len(kv) == 1 {
			tagName = kv[0]
			continue
		} else if len(kv) == 2 {
			if strings.TrimSpace(kv[0]) == "column" {
				tagName = strings.TrimSpace(kv[1])
				break
			}
		}
	}
	return
}

func parseTagName(field reflect.StructField, tag string) (name string, omitempty, ignore bool) {
	name = field.Tag.Get(tag)
	parseHandle := map[string]func(string) (string, bool, bool){
		"json": parseJsonTag,
		"xorm": parseXormTag,
		"gorm": parseGormTag,
	}
	if handle, ok := parseHandle[tag]; ok {
		return handle(name)
	}
	return
}

// 获取字段名优先级, json tag -> gorm tag -> xorm tag -> FileName, 如果没有则使用字段名的小驼峰格式
// return fieldname, omitempty, ignore
func getFieldTag(fieldType reflect.StructField, optArgs *args) (fieldName string, omitempty bool, ignore bool) {
	switch optArgs.curGetFieldType {
	case FieldType_Origin:
		fieldName = fieldType.Name
		omitempty = optArgs.omitempty
		return
	case FieldType_Json:
		fieldName, omitempty, ignore = parseTagName(fieldType, "json")
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}
		fieldName = littleCamelCase(fieldType.Name)
		return
	case FieldType_Gorm:
		fieldName, omitempty, ignore = parseTagName(fieldType, "gorm")
		omitempty = optArgs.omitempty
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}
		fieldName = littleCamelCase(fieldType.Name)
		return
	case FieldType_Xorm:
		fieldName, omitempty, ignore = parseTagName(fieldType, "xorm")
		omitempty = optArgs.omitempty
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}
		fieldName = littleCamelCase(fieldType.Name)
		return
	default:
		fieldName, omitempty, ignore = parseTagName(fieldType, "json")
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}

		fieldName, omitempty, ignore = parseTagName(fieldType, "gorm")
		omitempty = optArgs.omitempty
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}

		fieldName, omitempty, ignore = parseTagName(fieldType, "xorm")
		omitempty = optArgs.omitempty
		if len(fieldName) > 0 || ignore {
			if ignore {
				fieldName = littleCamelCase(fieldType.Name)
			}
			return
		}

		fieldName = littleCamelCase(fieldType.Name)
		return // fieldName, omitempty, false
	}

}

func getDeepIndent(deep int) string {
	if deep > 0 {
		return strings.Repeat("    ", deep)
	}
	return ""
}

func printLog(optArgs *args, deep int, args ...interface{}) {
	if optArgs.log != nil {
		tmp := make([]interface{}, 0, len(args)+1)
		tmp = append(tmp, getDeepIndent(deep))
		tmp = append(tmp, args...)
		fmt.Println(tmp...)
		//
		optArgs.log.Println(tmp...)
	}
}

func newOpts(opts ...CopyOption) args {
	opt := args{
		timeFmtStr:     "2006-01-02 15:04:05",
		timeValType:    TimeValType_String,
		ignoreFieldMap: map[string]struct{}{},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func InstanceFromBytes(dest interface{}, from []byte, opts ...CopyOption) (err error) {
	// optArgs := newOpts(opts...)
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(from, &tmp); err != nil {
		return err
	}
	return InstanceFromMap(dest, tmp, opts...)
}

func InstanceFromMap(dest interface{}, from interface{}, opts ...CopyOption) (err error) {
	optArgs := newOpts(opts...)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(interface2String(r))
			printLog(&optArgs, 0, r)
		}
	}()

	inst := reflect.ValueOf(dest)
	if inst.Kind() == reflect.Ptr {
		err = valueDeepCopy(inst.Elem(), from, 0, "", &optArgs)
	} else {
		err = errors.New("not ptr type")
	}
	return
}

func InstanceValueFromMap(dest reflect.Value, from interface{}, opts ...CopyOption) (err error) {
	optArgs := newOpts(opts...)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(interface2String(r))
			printLog(&optArgs, 0, r)
		}
	}()
	if dest.Kind() == reflect.Ptr {
		err = valueDeepCopy(dest.Elem(), from, 0, "", &optArgs)
	} else {
		err = valueDeepCopy(dest, from, 0, "", &optArgs)
	}
	return
}

// 将泛型数据map[string]interface{}, 通过reflect深度拷贝到对应的结构体中
// 如果直接调用此方法，需要外部捕获panic
func valueDeepCopy(inst reflect.Value, from interface{}, deep int, fieldName string, optArgs *args) (err error) {
	if !inst.CanSet() {
		return errors.New("target cannt be set")
	}

	// printlog("target name>>:", inst.Type().String(), inst.Kind())
	switch inst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := interface2Int64(from)
		inst.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := interface2Uint64(from)
		inst.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v := interface2Float64(from)
		inst.SetFloat(v)
	case reflect.String:
		v := interface2String(from)
		inst.SetString(v)
	case reflect.Bool:
		v := interface2Bool(from)
		inst.SetBool(v)
	case reflect.Interface:
		val := reflect.ValueOf(from)
		inst.Set(val)
	case reflect.Ptr:
		it := reflect.New(inst.Type().Elem())
		printLog(optArgs, deep, "Ptr>>:", it.String())

		err = valueDeepCopy(it.Elem(), from, deep+1, fieldName, optArgs)
		if err != nil {
			return
		}
		inst.Set(it)
	case reflect.Struct:
		if inst.Type().String() == "time.Time" {
			if optArgs.timeValType == TimeValType_String {
				timeStr := interface2String(from)
				if t, err := time.ParseInLocation(optArgs.timeFmtStr, timeStr, time.Local); err == nil {
					inst.Set(reflect.ValueOf(t))
				}
			} else if optArgs.timeValType == TimeValType_Int64 {
				timestamp := interface2Int64(from)
				t := time.Unix(timestamp, 0)
				inst.Set(reflect.ValueOf(t))
			}
			return nil
		}
		if mp, ok := from.(map[string]interface{}); ok {
			printLog(optArgs, deep, "Struct>>:", inst.String())

			tpe := inst.Type()
			for i := 0; i < tpe.NumField(); i += 1 {
				fieldType := inst.Type().Field(i)
				field := inst.Field(i)

				fieldName, _, ignore := getFieldTag(fieldType, optArgs)
				if ignore {
					continue
				}

				if fieldType.Anonymous {
					valueDeepCopy(field, mp, deep+1, fieldName, optArgs)
				} else {
					fieldValue, ok := mp[fieldName]
					if !ok || fieldValue == nil {
						continue
					}
					err = valueDeepCopy(field, fieldValue, deep+1, fieldName, optArgs)
					if err != nil {
						return
					}
				}
				// printlog(getDeepIndident(deep+1),"field name:", fieldName, "value:", field.Interface(), "kind:",field.Kind())
			}
		}
		return
	case reflect.Map:
		if vv, ok := from.(map[string]interface{}); ok {
			mp := reflect.MakeMap(inst.Type())
			printLog(optArgs, deep, "Map>>:", mp.String())

			err = mapValueDeepCopy(mp, vv, deep+1, fieldName, optArgs)
			if err != nil {
				return
			}
			inst.Set(mp)
		}
	case reflect.Slice:
		if vv, ok := from.([]interface{}); ok {
			sl := reflect.MakeSlice(inst.Type(), len(vv), cap(vv))
			printLog(optArgs, deep, "Slice>>:", sl.String())

			err = sliceValueDeepCopy(sl, vv, deep+1, fieldName, optArgs)
			if err != nil {
				return
			}
			inst.Set(sl)
		}
	case reflect.Array, reflect.Chan, reflect.Func, reflect.UnsafePointer: // 不处理
	}
	printLog(optArgs, deep, "field:", fieldName, "value:", inst.Interface(), "kind:", inst.Kind())
	return
}

func mapValueDeepCopy(inst reflect.Value, data map[string]interface{}, deep int, fieldName string, optArgs *args) (err error) {
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
			val = reflect.ValueOf(interface2Int(v))
		case reflect.Float32, reflect.Float64:
			val = reflect.ValueOf(interface2Float64(v))
		case reflect.String:
			val = reflect.ValueOf(interface2String(v))
		case reflect.Bool:
			val = reflect.ValueOf(interface2Bool(v))
		case reflect.Interface:
			val = reflect.ValueOf(v)
		case reflect.Struct: // map[string]TestStruct
			val = reflect.New(inst.Type().Elem()).Elem()
			printLog(optArgs, deep, "Struct>>:", val.String())

			err = valueDeepCopy(val, v, deep+1, "fieldName", optArgs)
			if err != nil {
				return
			}
		case reflect.Ptr: // map[string]*TestStruct
			val = reflect.New(inst.Type().Elem().Elem())
			printLog(optArgs, deep, "Ptr>>:", val.String())

			err = valueDeepCopy(val.Elem(), v, deep+1, fieldName, optArgs)
			if err != nil {
				return
			}
		case reflect.Map: // map[string]map[string]interface{}
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printLog(optArgs, deep, "Map>>:", val.String())

				err = mapValueDeepCopy(val, vv, deep+1, fieldName, optArgs)
				if err != nil {
					return
				}
			} else {
				continue
			}
		case reflect.Slice: // map[string][]interface{}
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printLog(optArgs, deep, "Slice>>:", val.String())

				err = sliceValueDeepCopy(val, vv, deep+1, fieldName, optArgs)
				if err != nil {
					return
				}
			} else {
				continue
			}
		}
		inst.SetMapIndex(reflect.ValueOf(k), val)
		printLog(optArgs, deep, "Map key:", k, "value:", val.Interface())
	}
	return
}

func sliceValueDeepCopy(inst reflect.Value, slice []interface{}, deep int, fieldName string, optArgs *args) (err error) {
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
			val = reflect.ValueOf(interface2Int(v))
		case reflect.Float32, reflect.Float64:
			val = reflect.ValueOf(interface2Float64(v))
		case reflect.String:
			val = reflect.ValueOf(interface2String(v))
		case reflect.Bool:
			val = reflect.ValueOf(interface2Bool(v))
		case reflect.Interface:
			val = reflect.ValueOf(v)
		case reflect.Struct: // []struct{}
			printLog(optArgs, deep, "Struct>>:", item.String())

			err = valueDeepCopy(item, v, deep+1, fieldName, optArgs)
			if err != nil {
				return
			}
			continue
		case reflect.Ptr: // []Ptr(*int|*struct)
			val = reflect.New(inst.Type().Elem().Elem())
			printLog(optArgs, deep, "Ptr>>:", val.String())

			err = valueDeepCopy(val.Elem(), v, deep+1, fieldName, optArgs)
			if err != nil {
				return
			}
		case reflect.Map: // []map[string]interface
			if vv, ok := v.(map[string]interface{}); ok {
				val = reflect.MakeMap(inst.Type().Elem())
				printLog(optArgs, deep, "Map>>:", val.String())

				err = mapValueDeepCopy(val, vv, deep+1, fieldName, optArgs)
				if err != nil {
					return
				}
			}
		case reflect.Slice: // [][]interface
			if vv, ok := v.([]interface{}); ok {
				val = reflect.MakeSlice(inst.Type().Elem(), len(vv), cap(vv))
				printLog(optArgs, deep, "Slice>>:", val.String())

				err = sliceValueDeepCopy(val, vv, deep+1, fieldName, optArgs)
				if err != nil {
					return
				}
			}
		}
		item.Set(val)
		printLog(optArgs, deep, "Slice index:", i, "value:", val.Interface())
	}
	return
}

// InstanceToMap 结构体转map
func InstanceToMap(from interface{}, opts ...CopyOption) (out map[string]interface{}, err error) {
	optArgs := newOpts(opts...)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(interface2String(r))
			printLog(&optArgs, 0, r)
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
	err = instanceToMap(out, inst, 0, &optArgs)
	return out, err
}

var (
	structTypes = []interface{}{reflect.Struct, reflect.Map, reflect.Slice}
	// NothingTypes = []interface{}{reflect.Array, reflect.Chan, reflect.Func, reflect.UnsafePointer}
)

func instanceToMap(dest map[string]interface{}, from reflect.Value, deep int, optArgs *args) (err error) {
	if from.Kind() == reflect.Ptr {
		return instanceToMap(dest, from.Elem(), deep, optArgs)
	}

	for i := 0; i < from.NumField(); i++ {
		field := from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, omitempty, ignore := getFieldTag(fieldType, optArgs)
		if ignore {
			continue
		}
		// 指定需要忽略的字段
		if _, ok := optArgs.ignoreFieldMap[fieldName]; ok {
			continue
		}
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		printLog(optArgs, deep, "kind:", field.Kind(), "fieldName:", fieldName, "value:", field.Interface(), "omitempty:", omitempty, "anonymous", fieldType.Anonymous)

		// 提前过来time解析

		switch field.Kind() {
		case reflect.Struct:
			if t, ok := field.Interface().(time.Time); ok {
				if t.IsZero() && omitempty {
					continue
				}
				if optArgs.timeValType == TimeValType_String {
					dest[fieldName] = t.Format(optArgs.timeFmtStr)
				} else if optArgs.timeValType == TimeValType_Int64 {
					dest[fieldName] = t.Unix()
				}
				continue
			}
			subMap := dest
			if !fieldType.Anonymous {
				subMap = make(map[string]interface{}, field.NumField())
				dest[fieldName] = subMap
			}
			if err = instanceToMap(subMap, field, deep+1, optArgs); err != nil {
				return
			}
		case reflect.Map:
			keys := field.MapKeys()
			if len(keys) == 0 && omitempty {
				continue
			}
			subMap := dest
			if !fieldType.Anonymous {
				subMap = make(map[string]interface{}, len(keys))
				dest[fieldName] = subMap
			}
			if err = instanceMapToMap(subMap, field, deep+1, optArgs); err != nil {
				return
			}
		case reflect.Slice:
			if field.Len() == 0 && omitempty {
				continue
			}
			subSlice := make([]interface{}, field.Len())
			dest[fieldName] = subSlice
			if err = instanceSliceToArr(subSlice, field, deep+1, optArgs); err != nil {
				return
			}
		default:
			if valueEmpty(field.Interface()) && omitempty {
				continue
			}
			dest[fieldName] = getBasicValue(field.Interface())
		}
	}
	return
}

// 结构体转成map，暂时不支持map/slice类型的字段
func instanceMapToMap(dest map[string]interface{}, field reflect.Value, deep int, optArgs *args) error {
	// inst := reflect.ValueOf(from)
	if field.Kind() != reflect.Map {
		return errors.New("field type is not map")
	}

	keys := field.MapKeys()
	for _, key := range keys {
		keyStr := interface2String(key.Interface())
		subField := field.MapIndex(key)
		if subField.Kind() == reflect.Ptr {
			subField = subField.Elem()
		}
		printLog(optArgs, deep, "kind:", subField.Kind(), "key:", keyStr, "value:", subField.Interface())
		switch subField.Kind() {
		case reflect.Struct:
			subMap := make(map[string]interface{}, subField.NumField())
			dest[keyStr] = subMap
			if err := instanceToMap(subMap, subField, deep+1, optArgs); err != nil {
				return err
			}
		case reflect.Map:
			keys := subField.MapKeys()
			subMap := make(map[string]interface{}, len(keys))
			dest[keyStr] = subMap
			if err := instanceMapToMap(subMap, subField, deep+1, optArgs); err != nil {
				return err
			}
		case reflect.Slice:
			subSlice := make([]interface{}, subField.Len())
			dest[keyStr] = subSlice
			if err := instanceSliceToArr(subSlice, subField, deep+1, optArgs); err != nil {
				return err
			}
		default:
			dest[keyStr] = subField.Interface()
		}
	}
	return nil
}

func instanceSliceToArr(dest []interface{}, field reflect.Value, deep int, optArgs *args) error {
	if field.Kind() != reflect.Slice {
		return errors.New("field type is not slice")
	}

	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		printLog(optArgs, deep, "kind:", item.Kind(), "index:", i, "value:", item.Interface())
		switch item.Kind() {
		case reflect.Struct:
			subMap := make(map[string]interface{}, item.NumField())
			dest[i] = subMap
			if err := instanceToMap(subMap, item, deep+1, optArgs); err != nil {
				return err
			}
		case reflect.Map:
			keys := item.MapKeys()
			subMap := make(map[string]interface{}, len(keys))
			dest[i] = subMap
			if err := instanceMapToMap(subMap, item, deep+1, optArgs); err != nil {
				return err
			}
		case reflect.Slice:
			subSlice := make([]interface{}, item.Len())
			dest[i] = subSlice
			if err := instanceSliceToArr(subSlice, item, deep+1, optArgs); err != nil {
				return err
			}
		default:
			dest[i] = item.Interface()
		}
	}
	return nil
}

func valueEmpty(v interface{}) bool {
	// 自定义类型，需要查看基础类型
	t := reflect.TypeOf(v)
	// fmt.Println(t.Kind()) // map
	switch t.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(v).Bool() == false
	case reflect.String:
		return reflect.ValueOf(v).String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v).Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(v).Uint() == 0
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(v).Float() == 0.0
	default:
		return v == nil
	}
}

// 将自定义类型的数据转换成基础类型数据
func getBasicValue(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	// fmt.Println(t.Kind()) // map
	switch t.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(v).Bool()
	case reflect.String:
		return reflect.ValueOf(v).String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(v).Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(v).Uint()
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(v).Float()
	default:
		return v
	}
}
