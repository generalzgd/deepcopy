# deepcopy

# 场景：
得到的json数据可能是弱类型语言生成的数据，例如php生成的数据，数据可能会带上引号，变成了字符串。
而golang是强类型语言，当对这样的数据进行json解析的时候（字符串类型的数据解析到int类型的字段上），会报错。
估开发此工具，进行深度数据拷贝，同时完成字段类型的转发。

# 问题点：
过程1：要先对数据进行json解析，并暂存到map[string]interface{}中；
过程2：再根据已定义好的结构体类型，进行拷贝数据和转换类型；
在两个过程中都使用了一次反射，性能有所开销

#usage:

```
bytes := []byte({...json数据...}) //举例, 具体可参考deepcopy_test.go文件

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
```



