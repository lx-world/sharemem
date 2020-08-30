# sharemem
> 1. 支持SystemV, Mmap

# 描述
> 1. Go操作共享内存
> 2. 支持key，value 类型的指定添加和删除

# 说明
> 1. 添加时,每块长度不能超过设置的长度
> 2. key 为字符串类型，最长为255字节
> 3. data 为字节切片,最长为4294967295



# 支持
> 1. linux, macos

# 使用
systemV
```go
    m, err := NewSystemV(1, 4, 2)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	m.WriteIdx("lx", []byte("11"))
	m.WriteIdx("lx1", []byte("21"))
	m.DelIdx("lx")
	m.DelIdx("lx1")
	res := m.GetAll()
	for key, value := range res {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

```

mmap
```go
    m, err := NewMMap("test.log",4,2)
	if err != nil{
		t.Errorf("err: %v", err)
		return
	}
	m.WriteIdx("test1",[]byte("31"))
	res,_ := m.GetKey("test1")
    fmt.Printf("%s", res)
    res1 := m.GetAll()
    for key, value := range res1 {
    	fmt.Printf("key: %s, value: %s\n", key, value)
    }
```