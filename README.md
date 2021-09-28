# saimaqi

go 写的赛码器

## 预览

![](/assets/preview-cli.png)
![](/assets/preview-web.png)

## 基本用法

所有文件使用 `utf8` 编码

```shell
smq-cli.exe [OPTIONS]

Application Options:
  /i, /input:    []string       码表路径，可设置多个
  /d, /ding:     int    普通码表起顶码长，码长大于等于此数，首选不会追加空格
  /s, /single    bool   是否只跑单字
  /w             bool   是否输出赛码表(保存在.\smb\文件夹下)
  /t, /text:     string utf8编码格式文本
  /k             bool   空格是否互击
  /c:            string 自定义选重键(2重开始) (default: ;\')
  /o, /output:   string 输出编码路径
  /v, /version   bool   查看版本信息

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```

## 码表

码表顺序和匹配顺序无关（最长顺序匹配）  
词条有重复时取前者

### rime 格式的码表

只支持编码在后的格式  
必须指定 `-d` 参数（起顶码长，码长大于等于此数，首选不会追加空格 `_`）  
可选 `-w` 参数（是否输出赛码表，保存在`.\\smb\\`文件夹下）

```
# Rime dictionary
# encoding: utf-8
#

---
name: xcxb
version: "1.0"
sort: original
...

人	a
如果	a
瑞	aa
睿	aaa
仍然是	aad
瑞士	aadi
...
```

`-d=4` 转换结果

```
人	a_
如果	a2
瑞	aa_
睿	aaa_
仍然是	aad_
瑞士	aadi
```

### 极速赛码器格式

> 注意：极速赛码器转换的格式编码为 utf-16，需要手动转为 utf-8 才能使用

```
工	a_
花	a2
華	a3
式	aa_
区区	aa2
工艺	aa3
恭恭敬敬	aaaa
花花草草	aaaa2
```

## 手感

可选:  
`-k` 空格是否互击  
`-c` 自定义选重键(2 重开始) (default ";'")  
将编码中末尾数字替换，只支持 10 重以内

## 例子

赛码表：`-i mbpath -t textpath`

普通码表：

- 四码定长：`-i mbpath -d=4 [-w] -t textpath`
- 二码顶功：`-i mbpath -d=2 [-w] -t textpath`
- 不定长：`-i mbpath -d=99 [-w] -t textpath`

多个码表同时测试：`-i mb1 -i mb2 -i mb3 -t textpath`

## 性能

> 配置：windows 10, r5 3600 4.2g, 8g\*2 2933Mhz c18

| 文本         | 文本字数 |          码表 | 码表词条数 |  耗时 |
| ------------ | -------: | ------------: | ---------- | ----: |
| 心情决定事情 |     9.6w | 091点儿2021夏 | 20w        | 120ms |
| 红楼梦原著   |      87w | 091点儿2021夏 | 20w        | 250ms |
| 极品全能高手 |    2164w | 091点儿2021夏 | 20w        |  3.8s |
| 极品全能高手 |    2164w | 红辣椒五笔码表 | 880w       | 12.6s |
| 心情决定事情 |     9.6w | 红辣椒五笔码表 | 880w       |  6.7s |
