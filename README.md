# saimaqi

go 写的赛码器

## 使用

所有文件使用 `utf8` 编码

```shell
saimaqi.exe [-i mb] [-d int] [-w] [-t text] [-s] [-k string] [-o output]

  -h    显示帮助
  -i string
        码表路径，可以是rime格式码表 或 极速跟打器赛码表
  -d int
        普通码表起顶码长，码长大于等于此数，首选不会追加空格
  -w    输出赛码表
  -t string
        文本路径，utf8编码格式文本，会自动去除空白符
  -s    空格是否互击
  -k string
        自定义选重键(2重开始) (default ";'")
  -o string
        输出路径
``` 
