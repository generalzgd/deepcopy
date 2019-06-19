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
	"fmt"
	"testing"
)

type InnerStruct struct {
	A int    `json:"aa"`
	B string `json:"bb"`
}

type CopyStruct struct {
	Int       int         `json:"int"`
	Float     float64     `json:"float"`
	Bool      bool        `json:"bool"`
	String    string      `json:"string"`
	Interface interface{} `json:"interface"`
	Struct    struct {
		AA int    `json:"aa"`
		BB string `json:"bb"`
	} `json:"struct"`
	InnerStu InnerStruct `json:"inner_stu"`
	//
	IntPtr    *int         `json:"int_ptr"`
	FloatPtr  *float64     `json:"float_ptr"`
	BoolPtr   *bool        `json:"bool_ptr"`
	StringPtr *string      `json:"string_ptr"`
	StructPtr *InnerStruct `json:"struct_ptr"`
	//
	MapInt       map[string]int         `json:"map_int"`
	MapFloat     map[string]float64     `json:"map_float"`
	MapBool      map[string]bool        `json:"map_bool"`
	MapString    map[string]string      `json:"map_string"`
	MapInterface map[string]interface{} `json:"map_interface"`
	MapStruct    map[string]struct {
		AA int    `json:"aa"`
		BB string `json:"bb"`
	} `json:"map_struct"`
	MapInnerStu map[string]InnerStruct  `json:"map_inner_stu"`
	MapIntPtr   map[string]*int         `json:"map_int_ptr"`
	MapFloatPtr map[string]*float64     `json:"map_float_ptr"`
	MapBoolPtr  map[string]*bool        `json:"map_bool_ptr"`
	MapStrPtr   map[string]*string      `json:"map_str_ptr"`
	MapStuPtr   map[string]*InnerStruct `json:"map_stu_ptr"`
	//
	ArrInt       []int         `json:"arr_int"`
	ArrFloat     []float64     `json:"arr_float"`
	ArrBool      []bool        `json:"arr_bool"`
	ArrString    []string      `json:"arr_string"`
	ArrInterface []interface{} `json:"arr_interface"`
	ArrStruct    []struct {
		AA int    `json:"aa"`
		BB string `json:"bb"`
	} `json:"arr_struct"`
	ArrInnerStu []InnerStruct  `json:"arr_inner_stu"`
	ArrIntPtr   []*int         `json:"arr_int_ptr"`
	ArrFloatPtr []*float64     `json:"arr_float_ptr"`
	ArrBoolPtr  []*bool        `json:"arr_bool_ptr"`
	ArrStrPtr   []*string      `json:"arr_str_ptr"`
	ArrStuPtr   []*InnerStruct `json:"arr_stu_ptr"`
	//
	MapMap       map[string]map[string]interface{}  `json:"map_map"`
	MapMapInt    map[string]map[string]int          `json:"map_map_int"`
	MapMapStu    map[string]map[string]InnerStruct  `json:"map_map_stu"`
	MapMapStuPtr map[string]map[string]*InnerStruct `json:"map_map_stu_ptr"`
	MapArr       map[string][]interface{}           `json:"map_arr"`
	MapArrInt    map[string][]int                   `json:"map_arr_int"`
	MapArrStu    map[string][]InnerStruct           `json:"map_arr_stu"`
	MapArrStuPtr map[string][]*InnerStruct          `json:"map_arr_stu_ptr"`
	//
	ArrMap       []map[string]interface{}  `json:"arr_map"`
	ArrMapInt    []map[string]int          `json:"arr_map_int"`
	ArrMapStu    []map[string]InnerStruct  `json:"arr_map_stu"`
	ArrMapStuPtr []map[string]*InnerStruct `json:"arr_map_stu_ptr"`
	ArrArr       [][]int                   `json:"arr_arr"`
	ArrArrStu    [][]InnerStruct           `json:"arr_arr_stu"`
	ArrArrStuPtr [][]*InnerStruct          `json:"arr_arr_stu_ptr"`
}

var (
	TestData = map[string]interface{}{
		"int":           1,
		"float":         1.2,
		"bool":          true,
		"string":        "abc",
		"interface":     5465.656,
		"struct":        map[string]interface{}{"aa": 11, "bb": "bb",},
		"inner_stu":     map[string]interface{}{"aa": 11, "bb": "bb",},
		"int_ptr":       "11",
		"float_ptr":     "23.65",
		"bool_ptr":      true,
		"string_ptr":    "poi",
		"struct_ptr":    map[string]interface{}{"aa": 11, "bb": "bb",},
		"map_int":       map[string]interface{}{"aa": 1,},
		"map_float":     map[string]interface{}{"aa": 23.54,},
		"map_bool":      map[string]interface{}{"aa": true,},
		"map_string":    map[string]interface{}{"aa": "map_string",},
		"map_interface": map[string]interface{}{"aa": "map_interface",},
		"map_struct":    map[string]map[string]interface{}{"key": {"aa": 11, "bb": "bb",}},
		"map_inner_stu": map[string]map[string]interface{}{"key": {"aa": 11, "bb": "bb"}},
		"map_int_ptr":   map[string]interface{}{"aa": 1,},
		"map_float_ptr": map[string]interface{}{"aa": 23.54,},
		"map_bool_ptr":  map[string]interface{}{"aa": true,},
		"map_str_ptr":   map[string]interface{}{"aa": "map_str_ptr",},
		"map_stu_ptr":   map[string]map[string]interface{}{"key": {"aa": 11, "bb": "bb"}},
		"arr_int":       []interface{}{1, 2, 3},
		"arr_float":     []interface{}{1, 2, 3},
		"arr_bool":      []interface{}{1, 2, true},
		"arr_string":    []interface{}{1, 2, true},
		"arr_interface": []interface{}{"arr_interface", 2, true},
		"arr_struct":    []map[string]interface{}{{"aa": 11, "bb": "bb",}},
		"arr_inner_stu": []map[string]interface{}{{"aa": 11, "bb": "bb",}},
		"arr_int_ptr":   []interface{}{1, 2, 3},
		"arr_float_ptr": []interface{}{1, 2, 3},
		"arr_bool_ptr":  []interface{}{1, 2, 3},
		"arr_str_ptr":   []interface{}{1, 2, true},
		"arr_stu_ptr":   []map[string]interface{}{{"aa": 11, "bb": "bb",}},
		"map_map":       map[string]map[string]interface{}{"key": {"aa": 11, "bb": "bb",}},
		"map_map_int":   map[string]map[string]interface{}{"key": {"aa": 111}},
		"map_map_stu": map[string]map[string]map[string]interface{}{"key": {"aa": {"aa": 11, "bb": "bb"}}},
		"map_map_stu_ptr": map[string]map[string]map[string]interface{}{"key": {"aa": {"aa": 11, "bb": "bb"}}},
		"map_arr":     map[string][]interface{}{"key": {1, "arr"}},
		"map_arr_int": map[string][]interface{}{"key": {1, 2}},
		"map_arr_stu": map[string][]map[string]interface{}{"key": {{"aa": 11, "bb": "bb1"}, {"aa": 11, "bb": "bb2"}}},
		"map_arr_stu_ptr": map[string][]map[string]interface{}{"key": {{"aa": 11, "bb": "bb1"}, {"aa": 11, "bb": "bb2"}}},
		"arr_map":     []map[string]interface{}{{"aa": 11, "bb": "bb1"}, {"aa": 11, "bb": "bb1"}},
		"arr_map_int": []map[string]interface{}{{"aa": 11, "bb": 22}, {"aa": 11, "bb": 33}},
		"arr_map_stu": []map[string]map[string]interface{}{{"key": {"aa": 11, "bb": 22}, "key2": {"aa": 11, "bb": 22}}},
		"arr_map_stu_ptr": []map[string]map[string]interface{}{{"key": {"aa": 11, "bb": 22}, "key2": {"aa": 11, "bb": 22}}},
		"arr_arr": [][]interface{}{{1, 2}, {11, 22}},
		"arr_arr_stu": [][]map[string]interface{}{{{"aa": 11, "bb": 22}}, {{"aa": 11, "bb": 22}}},
		"arr_arr_stu_ptr": [][]map[string]interface{}{{{"aa": 11, "bb": 22}}, {{"aa": 11, "bb": 22}}},
	}

	testDetail map[string]interface{}
)

func init() {
	bytes, err := json.Marshal(TestData)
	if err != nil {
		fmt.Println("init Marshal err:", err)
		return
	}
	if err := json.Unmarshal(bytes, &testDetail); err != nil {
		fmt.Println("init Unmarshal err:", err)
	}
}

func TestDeepCopy(t *testing.T) {
	SetLog(true)
	// var a int
	type args struct {
		i    interface{}
		from interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t2",
			args: args{
				i:    &CopyStruct{},
				from: testDetail,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.i, tt.args.from); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}






