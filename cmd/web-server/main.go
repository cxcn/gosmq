package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	_ "embed"

	smq "github.com/cxcn/gosmq"

	"github.com/cxcn/gosmq/pkg/html"
)

//go:embed index.html
var index []byte

func main() {

	server := http.Server{
		Addr: "localhost:5667",
	}
	http.HandleFunc("/index", func(rw http.ResponseWriter, r *http.Request) {
		// html, _ := ioutil.ReadFile("index.html")
		rw.Write(index)
	})
	http.HandleFunc("/result", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		si := getOptions(r.PostForm)
		h := html.NewHTML()
		so := si.Smq()
		h.AddResult(so)
		h.OutputHTML(rw)

		// fmt.Fprintln(rw, r.PostForm)
	})
	fmt.Println("\nhttp://localhost:5667/index\n ")
	server.ListenAndServe()
}

func getOptions(v url.Values) *smq.SmqIn {
	ret := new(smq.SmqIn)

	textname := "文本"
	if v.Get("textname") != "" {
		textname = v.Get("textname")
	}
	ret.Fpt = "text/" + textname + ".txt"
	_ = os.Mkdir("text", 0666)
	err := ioutil.WriteFile(ret.Fpt, []byte(v.Get("text")), 0666)
	if err != nil {
		panic(err)
	}

	dictname := "码表"
	if v.Get("dictname") != "" {
		dictname = v.Get("dictname")
	}

	ret.Fpd = "dict/" + dictname + ".txt"
	_ = os.Mkdir("dict", 0666)
	err = ioutil.WriteFile(ret.Fpd, []byte(v.Get("dict")), 0666)
	if err != nil {
		panic(err)
	}

	ret.Ding, _ = strconv.Atoi(v.Get("ding"))
	ret.Csk = v.Get("csk")
	if v.Get("iss") == "true" {
		ret.IsS = true
	}
	if v.Get("as") == "true" {
		ret.As = true
	}
	if v.Get("iso") == "true" {
		ret.IsO = true
	}

	return ret
}
