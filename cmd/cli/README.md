# go 赛码器

可以使用多个 `-i` 指定多个码表，此时其它参数都是共享的，所以其它码表的格式最好先转换过来。

不指定 `-t` 只指定 `-i` 时，可以初始化码表并转换输出。

```
Application Options:
  /t, /text:     string 文本路径
  /i, /input:    []string 码表路径
  /s, /single    bool   单字模式
  /f, /format:   string 码表格式 default|jisu|duoduo|jidian
      /select:   string 自定义选重键
  /p, /push:     int    普通码表起顶码长，码长大于等于此数，首选不会追加空格
  /a, /alg:      string 匹配算法 trie|order|longest
  /k, /space:    string 空格按键方式 left|right|both
  /v, /version   bool   查看版本信息

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```