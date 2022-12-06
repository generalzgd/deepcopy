# deepcopy

# 用途：
- 解析json数据到目标结构体上
- 解析map数据到结构体
- 转换结构体到map
- 支持结构体参数复制(相同参数名及类型)
- 支持深度嵌套结构体
- 支持多标签读取。json/xorm/gorm
- 支持指定字段忽略
- 支持0值忽略

# 场景：
- 得到的json数据可能是弱类型语言生成的数据，例如php生成的数字类型的字段，数据可能会带上引号，变成了字符串类型。
而golang是强类型语言，当对这样的json数据进行json解析的时候（字符串类型的数据解析到int类型的字段上），会报错。
故开发此工具，进行深度数据拷贝，同时完成字段类型的转换。
- 结构体与map之间的格式转换
- 结构体之间的参数复制

# 问题点：
过程1：要先对数据进行json解析，并暂存到map[string]interface{}中；
过程2：再根据已定义好的结构体类型，进行字段数据的拷贝和类型转换；
在这两个过程中都使用了反射，性能有所开销。
如果明显知道字段类型一一对应，不推荐使用此方法。
如果来源数据的字段类型不确定，但是字段名一致的情况下，推荐使用。


# usage1:

```
bytes := []byte({...json数据...}) //举例, 具体可参考deepcopy_test.go文件

var testDetail map[string]interface{}
if err := json.Unmarshal(bytes, &testDetail); err != nil {
    fmt.Println("init Unmarshal err:", err)
}

target := &dcopy.CopyStruct{}
dcopy.SetLog(true)
if err := dcopy.InstanceFromMap(target, testDetail); err != nil {
    fmt.Println("deep copy run err.", err)
} else {
    fmt.Println("deep copy run ok.", target)
}
```

# usage2:
```
type Args struct {
    AA int `json:"aa"`
}
jsonStr := `{"aa":"123"}`

data := &Args{}
dcopy.InstanceFromBytes(data, []byte(jsonStr))
```

# useage3:
```
type Args struct {
	AA int `json:"aa"`
}

data := &Args{123}

kvs, err := dcopy.InstanceToMap(data,
    dcopy.WithFieldType(dcopy.FieldType_Json), // 取json标签
    dcopy.WithOmitempty(true), // 忽略0值
)
fmt.Println(kvs) // {aa:123}
```

# usage4
```
type Args struct{
	AA int `json:"aa"`
}
kvs := map[string]interface{}{
    "aa":123,
}
dist := Args{}
err := dcopy.InstanceFromMap(&dist, kvs, dcopy.WithFieldType(dcopy.FieldType_Json))
```

# usage5
```
type Args struct {
    AA int
}

type Src Struct {
    AA int
}

err := dcopy.StructCopy(&Args{}, Src{AA:100})
fmt.Println(err)
```