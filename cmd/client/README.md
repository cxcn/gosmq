# gosmq 的命令行程序

- 命令行输入
- 命令行输出
- html 文件输出

> `SmqIn.Isw`字段为`true`，即为普通码表时，会自动输出赛码表

## 预览

![](./assets/preview-cli.png)

## 基本用法

```shell
smq-cli.exe [OPTIONS]

Application Options:
  /i, /input:    []string       码表路径，可设置多个
  /d, /ding:     int    普通码表起顶码长，码长大于等于此数，首选不会追加空格
  /s, /single    bool   是否只跑单字
  /t, /text:     string 文本
  /c:            string 自定义选重键(2重开始) (default: ";'")
  /k             bool   空格是否互击
  /o, /output    bool   是否输出结果
  /v, /version   bool   查看版本信息

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```

## 例子

赛码表：`-i mbpath -t textpath`

普通码表：

- 四码定长：`-i mbpath -d=4 -t textpath`
- 二码顶功：`-i mbpath -d=2 -t textpath`
- 不定长：`-i mbpath -d=99 -t textpath`

多个码表同时测试：`-i mb1 -i mb2 -i mb3 -t textpath`
