/**
 * @version: v0.1.0
 * @author: zhangguodong
 * @license: LGPL v3
 * @contact: zhangguodong@dobest.com
 * @site: https://code.dobest.com/research-go
 * @software: GoLand
 * @file: struct_struct.go
 * @time: 2022/12/5 17:09
 * @project: deepcopy
 */

package dcopy

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/toolkits/slice"
)

//
// StructCopy 结构体同字段名（同字段类型）拷贝
//  @Description:
//  @param dest
//  @param from
func StructCopy(dest interface{}, from interface{}, opts ...CopyOption) error {
	optArgs := newOpts(opts...)
	//
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.New("dest not ptr type")
	} else {
		destValue = destValue.Elem()
		if destValue.Kind() != reflect.Struct {
			return errors.New("dest not struct type")
		}
	}

	fromValue := reflect.ValueOf(from)
	if fromValue.Kind() == reflect.Ptr {
		fromValue = fromValue.Elem()
	}
	if fromValue.Kind() != reflect.Struct {
		return errors.New("from not struct type")
	}

	hit, miss := structCopy(destValue, fromValue, optArgs)
	printLog(&optArgs, 0, fmt.Sprintf("struct copy complete: hit(%d) miss(%d)", hit, miss))
	return nil
}

func structCopy(dest, from reflect.Value, optArgs args) (hit, mis int) {

	for i := 0; i < dest.NumField(); i++ {
		destFieldType := dest.Type().Field(i)
		destField := dest.Field(i)
		fieldName := destFieldType.Name
		if !destField.CanSet() {
			mis += 1
			continue
		}

		if destField.Kind() == reflect.Ptr {
			destField = destField.Elem()
		}

		fromField := from.FieldByName(fieldName)
		if !fromField.IsValid() { // 找不到字段
			mis += 1
			continue
		}

		if fromField.Kind() == reflect.Ptr {
			fromField = fromField.Elem()
		}

		// 如果数据类型不匹配则返回错误
		if !isFieldTypeMatch(destField, fromField) {
			//return errors.New("field type not match")
			mis += 1
			continue
		}

		switch destField.Kind() {
		case reflect.Slice:
			if e := sliceCopy(destField, fromField, optArgs); e == nil {
				hit += 1
			} else {
				mis += 1
			}
		case reflect.Map:
			if e := mapCopy(destField, fromField, optArgs); e == nil {
				hit += 1
			} else {
				mis += 1
			}
		case reflect.Struct:
			h, m := structCopy(destField, fromField, optArgs)
			hit += h
			mis += m
		default:
			basicCopy(destField, fromField, optArgs)
			hit += 1
		}
	}
	return
}

func mapCopy(dest, from reflect.Value, optArgs args) error {
	makeMap := reflect.MakeMap(dest.Type())
	iter := from.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		makeMap.SetMapIndex(k, v)
	}
	dest.Set(makeMap)
	return nil
}

func sliceCopy(dest, from reflect.Value, optArgs args) error {
	makeSlice := reflect.MakeSlice(dest.Type(), from.Len(), from.Cap())
	reflect.Copy(makeSlice, from)
	dest.Set(makeSlice)
	return nil
}

func basicCopy(dest, from reflect.Value, optArgs args) {
	switch dest.Kind() {
	case reflect.String:
		v := from.String()
		dest.SetString(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := from.Int()
		dest.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := from.Uint()
		dest.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v := from.Float()
		dest.SetFloat(v)
	case reflect.Bool:
		v := from.Bool()
		dest.SetBool(v)
	}
}

// isFieldTypeMatch 相识类型判断  int8,int16,int32,int64; uint8,uint16,uint32,uint64,int 归一类
//  @Description:
//  @param destType
//  @param fromType
//  @return bool
func isFieldTypeMatch(dest, from reflect.Value) bool {
	intKind := []interface{}{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}
	uintKind := []interface{}{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64}
	floatKind := []interface{}{reflect.Float64, reflect.Float32}

	switch dest.Kind() {
	case reflect.String:
		return from.Kind() == reflect.String
	case reflect.Float32, reflect.Float64:
		return slice.Contains(floatKind, from.Kind())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return slice.Contains(intKind, from.Kind())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return slice.Contains(uintKind, from.Kind())
	case reflect.Slice:
		if from.Kind() == reflect.Slice {
			return dest.Type().Elem() == from.Type().Elem()
		}
		return false
	case reflect.Map:
		if from.Kind() == reflect.Map {
			return dest.Type().Elem() == from.Type().Elem()
		}
		return false
	default:
		return dest.Kind() == from.Kind()
	}
}
