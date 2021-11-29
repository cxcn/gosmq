package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "embed"

	smq "github.com/cxcn/gosmq"

	"github.com/cxcn/gosmq/pkg/html"
)

//go:embed index.html
var index string

func main() {

	type names struct {
		DictNames []string
		TextNames []string
	}

	theNames := new(names)

	// 读取dict/目录中的所有文件和子目录
	files, err := ioutil.ReadDir(`dict/`)
	if err != nil {
		panic(err)
	}
	fmt.Println("检测到以下赛码表：")
	for _, file := range files {
		theNames.DictNames = append(theNames.DictNames, file.Name())
		fmt.Println("  ", file.Name())
	}
	fmt.Println()
	// 读取text/目录中的所有文件和子目录
	files, err = ioutil.ReadDir(`text/`)
	if err != nil {
		panic(err)
	}
	fmt.Println("检测到以下文本：")
	for _, file := range files {
		theNames.TextNames = append(theNames.TextNames, file.Name())
		fmt.Println("  ", file.Name())
	}

	funcMap := template.FuncMap{"getName": getName}
	t := template.New("index.html").Funcs(funcMap)
	_, err = t.Parse(index)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr: "localhost:8080",
	}
	http.HandleFunc("/index", func(rw http.ResponseWriter, r *http.Request) {
		// html, _ := ioutil.ReadFile("index.html")
		// rw.Write(html)
		t.Execute(rw, theNames)
	})
	http.HandleFunc("/result", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		options := getOptions(r.PostForm)
		h := html.NewHTML()
		for _, v := range options {
			so := v.Smq()
			h.AddResult(so)
		}
		h.OutputHTML(rw)

		// fmt.Fprintln(rw, r.PostForm)
	})
	fmt.Println("\nhttp://localhost:8080/index\n ")
	server.ListenAndServe()
}

func getOptions(v url.Values) []*smq.SmqIn {
	var ret []*smq.SmqIn
	if v.Get("fpd") != "" {
		tmp := new(smq.SmqIn)
		tmp.Fpt = "text/" + v.Get("fpt")
		tmp.Fpd = "dict/" + v.Get("fpd")
		tmp.Ding, _ = strconv.Atoi(v.Get("ding"))
		tmp.Csk = v.Get("csk")
		if v.Get("iss") == "true" {
			tmp.IsS = true
		}
		if v.Get("as") == "true" {
			tmp.As = true
		}
		if v.Get("iso") == "true" {
			tmp.IsO = true
		}
		ret = append(ret, tmp)
	}
	if v.Get("fpd1") != "" {
		tmp := new(smq.SmqIn)
		tmp.Fpt = "text/" + v.Get("fpt")
		tmp.Fpd = "dict/" + v.Get("fpd1")
		tmp.Ding, _ = strconv.Atoi(v.Get("ding1"))
		tmp.Csk = v.Get("csk1")
		if v.Get("iss1") == "true" {
			tmp.IsS = true
		}
		if v.Get("as1") == "true" {
			tmp.As = true
		}
		if v.Get("iso1") == "true" {
			tmp.IsO = true
		}
		ret = append(ret, tmp)
	}
	return ret
}

func getName(s string) string {
	s = strings.TrimSuffix(s, ".txt")
	return strings.TrimSuffix(s, "_赛码表")
}
