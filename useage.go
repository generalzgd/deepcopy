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
	"dcopy"
)

func main() {
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
	dcopy.SetLog(true)
	if err := dcopy.DeepCopy(target, testDetail); err != nil {
		fmt.Println("deep copy run err.", err)
	} else {
		fmt.Println("deep copy run ok.", target)
	}
}
