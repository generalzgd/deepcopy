/*
 * @version: 1.0.0
 * @author: zhangguodong
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: https://code.dobest.com/research-go/deepcopy
 * @software: GoLand
 * @file: getset_test.go.go
 * @time: 2020/5/22 17:03
 * @project: deepcopy
 */

package dcopy

import (
	"reflect"
	"testing"
)

func TestGetFieldValue(t *testing.T) {
	obj := FooFieldTest{
		A0: 11,
		A:  "A",
		B:  "B",
		B2: "B2",
		C:  13,
		C2: 131,
		D:  "D",
		D2: "D2",
		E:  "E",
		F:  "F",
	}
	obj.Code = 1
	objPtr := &obj

	type args struct {
		target interface{}
		field  string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{
			name: "TestGetFieldValue_0",
			args: args{
				target: obj,
				field:  "Code",
			},
			want: 1,
		},
		{
			name: "TestGetFieldValue_1",
			args: args{
				target: obj,
				field:  "A0",
			},
			want: 11,
		},
		{
			name: "TestGetFieldValue_2",
			args: args{
				target: obj,
				field:  "a",
			},
			want: "A",
		},
		{
			name: "TestGetFieldValue_3",
			args: args{
				target: objPtr,
				field:  "c",
			},
			want: 13,
		},
		{
			name: "TestGetFieldValue_4",
			args: args{
				target: objPtr,
				field:  "d2",
			},
			want: "D2",
		},
		{
			name: "TestGetFieldValue_5",
			args: args{
				target: obj,
				field:  "a0",
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFieldValue(tt.args.target, tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	obj := &FooFieldTest{}
	type args struct {
		target         interface{}
		fieldOrTagName string
		value          interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestSetFieldValue_0",
			args: args{
				target:         obj,
				fieldOrTagName: "Code",
				value:          1,
			},
		},
		{
			name: "TestSetFieldValue_1",
			args: args{
				target:         obj,
				fieldOrTagName: "a0",
				value:          "asdf",
			},
		},
		{
			name: "TestSetFieldValue_2",
			args: args{
				target:         obj,
				fieldOrTagName: "A0",
				value:          100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFieldValue(tt.args.target, tt.args.fieldOrTagName, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("SetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFieldsTagName(t *testing.T) {
	type args struct {
		target       interface{}
		fieldType    FieldType
		ignoreFields []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "TestGetFieldsTagName",
			args: args{
				target:       &FooFieldTest{},
				fieldType:    FieldType_Gorm,
				ignoreFields: []string{"d"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFieldsTagName(tt.args.target, tt.args.fieldType, tt.args.ignoreFields); !reflect.DeepEqual(got, tt.want) {
				t.Logf("GetFieldsTagName() = %v", got)
			}
		})
	}
}

func TestGetFieldsValue(t *testing.T) {
	type args struct {
		target       interface{}
		omitempty    bool
		fieldType    FieldType
		ignoreFields []string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			name: "TestGetFieldsValue",
			args: args{
				target: &FooFieldTest{
					A0: 1,
					A:  "asdfa",
					B:  "dd",
					B2: "basdf",
					C:  3,
					C2: 1,
					D:  "ad",
					D2: "zxcv",
					E:  "234",
					F:  "sda",
				},
				omitempty:    false,
				fieldType:    FieldType_Gorm,
				ignoreFields: []string{"d"},
			},
		},
		{
			name: "TestGetFieldsValue_omitempty",
			args: args{
				target: &FooFieldTest{
					A0: 1,
					A:  "asdfa",
					B:  "dd",
					B2: "basdf",
					C:  3,
					C2: 1,
					D:  "ad",
					D2: "zxcv",
					E:  "234",
					F:  "sda",
				},
				omitempty:    true,
				fieldType:    FieldType_Gorm,
				ignoreFields: []string{"d"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFieldsValue(tt.args.target, tt.args.omitempty, tt.args.fieldType, tt.args.ignoreFields); !reflect.DeepEqual(got, tt.want) {
				t.Logf("GetFieldsValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetZeroFields(t *testing.T) {
	type a struct {
		AaXx int32 `json:"aa_xx"`
		BbYy int32 `json:"bb_yy"`
		CcZz int32 `json:"cc_zz"`
	}
	s := a{
		AaXx: 0,
		BbYy: 1,
		CcZz: 0,
	}
	type args struct {
		target    interface{}
		fieldType FieldType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestGetZeroFields",
			args: args{
				target:    s,
				fieldType: FieldType_Json,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetZeroFields(tt.args.target, tt.args.fieldType)
			t.Logf("GetZeroFields() = %v", got)
		})
	}
}
