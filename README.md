# go 赛码器

一个高性能的赛码器库，多码表支持，多种码表格式支持，多匹配算法支持，可自定义码表格式和匹配算法。

## 开始之前

下载：`go get -u github.com/cxcn/gosmq`  
导入：

```go
import (
    smq "github.com/cxcn/gosmq"
)
```

## 使用

### 创建赛码器

你可以从字节流、字符串、文件路径创建赛码器

```go
func main(){
    // 从字节流
    f := os.Open(`filepath`) // 实现了 io.Reader 的方法
    s1 := smq.New("文本名",f)
    // 从字符串
    s2:= smq.NewFromString("文本名","这里是文本内容")
    // 从文件路径
    // 第一个参数可以填空，这时的文本名从路径推断
    s3 := smq.NewFromPath("","这里填写路径")
}
```

### 创建一个码表

```go
func main(){
    // 先定义一些基本参数
    dict := &smq.Dict{
        Name:         "", // 码表名
        Single:       false, /* 单字模式
        注意单字模式最好用多多格式的码表，
        因为其他码表带有 order 信息，转换后可以用本格式码表
        */
        Format:       "", /* 码表格式
        default: 默认 本程序赛码表 词\t编码选重\t选重
        jisu:js 极速赛码表 词\t编码选重
        duoduo:dd 多多格式码表 词\t编码
        jidian:jd 极点格式 编码\t词1 词2 词3
        bingling:bl 冰凌格式码表 编码\t词
        */
        Transfer   Transfer // 自定义码表格式转换
        SavePath   string // 读取非默认码表格式时自动转换并保存的路径，默认保存在 dict 目录下
        SelectKeys:   "", // 普通码表自定义选重键(默认为_;')
        PushStart:    4, // 普通码表起顶码长(码长大于等于此数，首选不会追加空格)
        Algorithm:    "longest", // 匹配算法 trie:前缀树 order:顺序匹配（极速跟打器） longest:最长匹配
        Matcher:      nil, // 自定义匹配算法
        PressSpaceBy: "both", // 空格按键方式 left|right|both
        Details      bool,   // 输出详细数据
    }

    // 载入码表，同样提供 3 种方式
    // 从字节流
    f := os.Open(`filepath`) // 实现了 io.Reader 的方法
    dict.Load(f)
    // 从字符串
    dict.LoadFromString("这里是码表内容")
    // 从文件路径
    dict.LoadFromPath("这里填写路径")
}
```

### 仅转换码表

`dict.Convert()` 这会将其它格式转换为本编码格式，并输出保存

> 如果你还要进行后续操作，则不需要手动转换

### 添加码表到赛码器

`s.Add(dict)`

### 开始比赛

```go
func main() {
    s := smq.New(...)
    dict := &smq.Dict{...}
    dict.Load(...)
    s.Add(dict)
    res := s.Run() // 他返回一个 []*smq.Result 结构体指针数组，具体定义可查看 struct.go 文件
    // 你可以输出为 json
    s.ResToJson(res)
    // 如果你不需要结构体
    s.ToJson()
}

```

## 自定义码表转换

```go
// 需要从 dict 生成本赛码器格式码表的字节数组
type Transfer interface {
    Read(*Dict) []byte
}
```

一个例子：

```go
type newFormat struct {}

func (n *newFormat) Read(dict *Dict) []byte {
    // 推荐使用 bytes.Buffer
    var buf bytes.Buffer
    buf.Grow(1e6)
    // 逐行读取
    scan := bufio.NewScanner(dict.reader)
    for scan.Scan() {
        // 把每一行的格式转为
        // 本赛码器格式 word\tcode\torder
    }
    return buf.Bytes()
}
```

## 自定义匹配算法

程序定义了一个接口，在创建 `smq.Dict` 时传入 `Matcher`，这时 `Algorithm` 失效

```go
type Matcher interface {
    // 插入一个词条 word code order
    Insert(string, string, int)
    // 读取完码表后的操作
    Handle()
    // 匹配下一个词 text point -> 匹配到的词长，code，order
    Match([]rune, int) (int, string, int)
}
```

## Benchmark

> 这是老版本的 Benchmark，新版要慢一点，不会差太多  
> 配置：windows 10, r5 3600 4.2g, 8g\*2 2933Mhz c18

| 文本         | 文本字数 |             码表 | 码表词条数 |  耗时 |
| ------------ | -------: | ---------------: | ---------- | ----: |
| 心情决定事情 |     9.6w | 091 点儿 2021 夏 | 20w        | 120ms |
| 红楼梦原著   |      87w | 091 点儿 2021 夏 | 20w        | 250ms |
| 极品全能高手 |    2164w | 091 点儿 2021 夏 | 20w        |  3.8s |
| 极品全能高手 |    2164w |   红辣椒五笔码表 | 880w       | 12.6s |
| 心情决定事情 |     9.6w |   红辣椒五笔码表 | 880w       |  6.7s |
