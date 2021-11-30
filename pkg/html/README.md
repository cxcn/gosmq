`NewHTML()` 创建一个对象

`.AddResult(*smq.SmqOut)` 添加一个结果

`.OutputHTML(io.Writer)` 输出 html 到 io 流

`.OutputHTMLFile(string)` 输出 html 到文件

## 例子

```go
si := new(smq.SmqIn)
h := html.NewHTML("文本名")
so,err := si.Smq()
if err != nil {
    panic(err)
}
h.AddResult(so,"码表名")
h.OutputHTMLFile("result.html")
```

## 预览

![](./assets/preview-html.png)
