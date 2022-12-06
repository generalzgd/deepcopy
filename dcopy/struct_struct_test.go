/**
 * @version: v0.1.0
 * @author: zhangguodong
 * @license: LGPL v3
 * @contact: zhangguodong@dobest.com
 * @site: https://code.dobest.com/research-go
 * @software: GoLand
 * @file: struct_struct_test.go
 * @time: 2022/12/6 10:01
 * @project: deepcopy
 */

package dcopy

import "testing"

func TestStructFromStruct(t *testing.T) {

	type Foo struct {
		NotExist string
		Str      string
		Num      int
		Tmp      map[int]int
		Arr      []uint
		//
		Tmp2 map[int]Foo
		Arr2 []Foo
	}

	type Eoo struct {
		Str string
		Num int64
		Tmp map[int]int
		Arr []uint
		//
		Tmp2 map[int]Eoo
		Arr2 []Eoo
	}

	type args struct {
		dest interface{}
		from interface{}
		opts []CopyOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestStructFromStruct",
			args: args{
				dest: &Foo{},
				from: Eoo{
					Str:  "asdf",
					Num:  123,
					Tmp:  map[int]int{100: 200},
					Arr:  []uint{1, 2, 3},
					Tmp2: map[int]Eoo{1000: {Str: "inner"}},
					Arr2: []Eoo{{Str: "slice"}, {Str: "Arr"}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StructFromStruct(tt.args.dest, tt.args.from, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("StructFromStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
