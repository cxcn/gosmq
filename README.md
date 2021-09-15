# saimaqi

go 写的赛码器

## 基本用法

所有文件使用 `utf8` 编码

```shell
smq-cli.exe [-i mb] [-n int] [-d] [-w] [-t text] [-s] [-k string] [-o output]

  -h    显示帮助
  -i string
        码表路径，可以是`rime`格式码表或速跟打器赛码表
  -n int
        普通码表起顶码长，码长大于等于此数，首选不会追加空格
  -d    是否只跑单字
  -w    是否输出赛码表(保存在.\\smb\\文件夹下)
  -t string
        文本路径，utf8编码格式文本，会自动去除空白符
  -f    是否关闭手感统计
  -s    空格是否互击
  -k string
        自定义选重键(2重开始) (default ";'")
  -o string
        输出路径
``` 

## 码表

码表顺序和匹配顺序无关（最长顺序匹配）  
字词有重复时取较短码长

通用: 可选 `-d`, `-k`  
码表: 必须 `-n`, 可选 `-w`  
赛码表: 不要 `-n`, `-w`

### rime 格式的码表

只支持编码在后的格式  
必须指定 `-n` 参数（起顶码长，码长大于等于此数，首选不会追加空格 `_`）  
可选 `-w` 参数（是否输出赛码表，保存在`.\\smb\\`文件夹下）

例：
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

`-n=4` 转换结果
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

例：
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
`-f` 是否关闭手感统计  
`-s` 空格是否互击  
`-k` 自定义选重键(2重开始) (default ";'")  
将编码中末尾数字替换，只支持10重以内

## 例子

赛码表：`-i mbpath -t textpath`

普通码表：
- 四码定长：`-i mbpath -n=4 [-w] -t textpath`
- 二码顶功：`-i mbpath -n=2 [-w] -t textpath`
- 不定长：`-i mbpath -n=99 [-w] -t textpath`
