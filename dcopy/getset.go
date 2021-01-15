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

func GetNotZeroFields(target interface{}) []string {
	if target == nil {
		return nil
	}
	instTp := reflect.TypeOf(target)
	if instTp.Kind() == reflect.Ptr {
		instTp = instTp.Elem()
	}
	instVl := reflect.ValueOf(target)
	if instVl.Kind() == reflect.Ptr {
		instVl = instVl.Elem()
	}

	if instTp.Kind() == reflect.Struct {
		num := instTp.NumField()
		out := make([]string,0,num)
		for i:=0;i<num;i++{
			fieldTp := instTp.Field(i)
			fieldVl := instVl.Field(i)
			// fmt.Println(fieldTp.Name,fieldVl.String())
			if fieldVl.IsValid() && !fieldVl.IsZero() {
				out = append(out, fieldTp.Name)
			}
		}
		return out
	}
	return nil
}

// 获取struct对象的字段值
// fieldOrTagName可以是字段名，json/gorm/xorm tag, 或小驼峰字段名
func GetFieldValue(target interface{}, fieldOrTagName string, opts ...CopyOption) interface{} {
	if target == nil {
		return nil
	}
	optArgs := newOpts(opts...)

	inst := reflect.ValueOf(target)
	if inst.Kind() == reflect.Ptr {
		inst = inst.Elem()
	}
	if inst.Kind() == reflect.Struct {
		return getFileValue(inst, fieldOrTagName, &optArgs)
	}
	return nil
}

func getFileValue(from reflect.Value, fieldOrTagName string, optArgs *args) interface{} {
	field := from.FieldByName(fieldOrTagName)
	// 直接字段名获取成功
	if field.IsValid() {
		return field.Interface()
	}
	// 通过 json/gorm/xorm tag 或小驼峰字段名获取
	for i:=0;i<from.NumField();i++ {
		field = from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, _, _ := getFieldTag(fieldType, optArgs)
		if fieldName == fieldOrTagName {
			return field.Interface()
		}
	}
	return nil
}

// 对struct（必须为指针） 对象，设置对应字段的变量
// fieldOrTagName可以是字段名，json/gorm/xorm tag, 或小驼峰字段名
// 如果字段的类型和值的类型对不上，则设置的是0值，不返回错误
func SetFieldValue(target interface{}, fieldOrTagName string, value interface{}, opts ...CopyOption) (err error) {
	if target == nil {
		return
	}
	optArgs := newOpts(opts...)

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
		err = setFieldValue(inst, fieldOrTagName, value, &optArgs)
	}
	return
}

func setFieldValue(from reflect.Value, fieldOrTagName string, value interface{}, optArgs *args) error {

	field := from.FieldByName(fieldOrTagName)

	if field.IsValid() {
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		return valueDeepCopy(field, value, 0, fieldOrTagName, optArgs)
	}

	// 通过 json/gorm/xorm tag 或小驼峰字段名获取
	for i:=0;i<from.NumField();i++ {
		field = from.Field(i)
		fieldType := from.Type().Field(i)
		fieldName, _, _ := getFieldTag(fieldType, optArgs)
		if fieldName == fieldOrTagName {
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			return valueDeepCopy(field, value, 0, fieldOrTagName, optArgs)
		}
	}
	return nil
}