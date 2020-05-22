/*
 * @version: 1.0.0
 * @author: zhangguodong
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: https://code.dobest.com/research-go/deepcopy
 * @software: GoLand
 * @file: getset.go
 * @time: 2020/5/22 15:31
 * @project: deepcopy
 */

package dcopy

import (
	"fmt"
	"reflect"
)

// 获取struct对象的字段值
// fieldOrTagName可以是字段名，json/gorm/xorm tag, 或小驼峰字段名
func GetFieldValue(target interface{}, fieldOrTagName string) interface{} {
	if target == nil {
		return nil
	}
	inst := reflect.ValueOf(target)
	if inst.Kind() == reflect.Ptr {
		inst = inst.Elem()
	}
	if inst.Kind() == reflect.Struct {
		return getFileValue(inst, fieldOrTagName)
	}
	return nil
}

func getFileValue(from reflect.Value, fieldOrTagName string) interface{} {
	field := from.FieldByName(fieldOrTagName)
	// 直接字段名获取成功
	if field.IsValid() {
		return field.Interface()
	}
	// 通过 json/gorm/xorm tag 或小驼峰字段名获取
	for i:=0;i<from.NumField();i++ {
		field = from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, _, _ := getFieldTag(fieldType)
		if fieldName == fieldOrTagName {
			return field.Interface()
		}
	}
	return nil
}

// 对struct（必须为指针） 对象，设置对应字段的变量
// fieldOrTagName可以是字段名，json/gorm/xorm tag, 或小驼峰字段名
// 如果字段的类型和值的类型对不上，则设置的是0值，不返回错误
func SetFieldValue(target interface{}, fieldOrTagName string, value interface{}) (err error) {
	if target == nil {
		return
	}
	inst := reflect.ValueOf(target)
	if inst.Kind() != reflect.Ptr {
		err = fmt.Errorf("not pointer target")
		return
	}
	inst = inst.Elem()
	if inst.Kind() == reflect.Struct {
		defer func() {
			if r := recover(); r!=nil {
				err = fmt.Errorf("set field value err=[%v]", r)
			}
		}()
		err = setFieldValue(inst, fieldOrTagName, value)
	}
	return
}

func setFieldValue(from reflect.Value, fieldOrTagName string, value interface{}) error {

	field := from.FieldByName(fieldOrTagName)

	if field.IsValid() {
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		return valueDeepCopy(field, value, 0, fieldOrTagName)
	}

	// 通过 json/gorm/xorm tag 或小驼峰字段名获取
	for i:=0;i<from.NumField();i++ {
		field = from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, _, _ := getFieldTag(fieldType)
		if fieldName == fieldOrTagName {
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			return valueDeepCopy(field, value, 0, fieldOrTagName)
		}
	}
	return nil
}