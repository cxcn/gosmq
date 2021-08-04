# saimaqi

go 写的赛码器

## 使用

所以文件使用 `utf8` 编码

```shell
saimaqi.exe [-m mb] [-t text] [-o output]

  -h    帮助
  -m string
        码表路径
  -o string
        输出路径
  -t string
        文本路径
``` 

### 普通码表

可以使用 `rime` 格式码表，需要在文件头添加 `smq_conf` 参数

```yaml
# Rime dictionary
# encoding: utf-8
#

---
name: xcxb
version: "1.0"
sort: original
smq_conf: #
  # alter_keys: { #自定义选重键
  #   1: '_',
  #   2: ';',
  #   3: "'",
  # }
  auto_select: 4 #起顶码长
...
```

### 赛码表

也可以使用极速赛码器生成的赛码表，需要转为 `utf8` 编码

### 文本

自动去除空白符号