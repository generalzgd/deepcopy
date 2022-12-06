/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: useage.go
 * @time: 2019/6/19 18:10
 */
package main

import (
	"encoding/json"
	"fmt"

	"github.com/generalzgd/deepcopy/dcopy"
)

type InnerObj struct {
	Sun int
}

type Inner struct {
	Foo string
	*InnerObj
}

type Outer struct {
	Inner
	Abc string
}

func testAnonymousEncode()  {
	obj := Outer{}
	obj.InnerObj = &InnerObj{}
	obj.Foo = "foo"
	obj.Abc = "abc"
	obj.Sun = 123

	out, err := dcopy.InstanceToMap(obj)
	fmt.Println(err)
	fmt.Println(out)
}

func testAnonymousDecode() {
	tmp := map[string]interface{}{
		"abc": "abc",
		"foo": "foo",
		"sun": 123,
	}
	tar := Outer{}
	err := dcopy.InstanceFromMap(&tar, tmp)
	fmt.Println(err)
	fmt.Println("out:", tar)
}

func main() {
	// testAnonymousEncode()
	// testAnonymousDecode()


	bytes, err := json.Marshal(dcopy.TestData)
	if err != nil {
		fmt.Println("init Marshal err:", err)
		return
	}
	var testDetail map[string]interface{}
	if err := json.Unmarshal(bytes, &testDetail); err != nil {
		fmt.Println("init Unmarshal err:", err)
	}

	target := &dcopy.CopyStruct{}
	if err := dcopy.InstanceFromMap(target, testDetail); err != nil {
		fmt.Println("deep copy run err.", err)
	} else {
		fmt.Println("deep copy run ok.", target)
	}
}
