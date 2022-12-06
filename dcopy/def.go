// Package dcopy
/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: def.go
 * @time: 2019/7/23 11:11
 */
package dcopy

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
	TestData   = map[string]interface{}{}
	testDetail map[string]interface{}
)