package main

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	_ "embed"

	smq "github.com/cxcn/gosmq"

	"github.com/cxcn/gosmq/pkg/html"
)

//go:embed index.html
var index []byte

func main() {

	fmt.Printf("smq-web-server version 0.1 %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("repo address: https://github.com/cxcn/gosmq/\n\n")

	server := http.Server{
		Addr: "localhost:5667",
	}
	http.HandleFunc("/index", func(rw http.ResponseWriter, r *http.Request) {
		// html, _ := ioutil.ReadFile("index.html")
		rw.Write(index)
	})
	http.HandleFunc("/result", func(rw http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		v := r.PostForm
		si := getOptions(v)

		tn := "文本"
		if v.Get("textname") != "" {
			tn = v.Get("textname")
		}
		dn := "码表"
		if v.Get("dictname") != "" {
			dn = v.Get("dictname")
		}

		h := html.NewHTML(tn)
		so, _ := si.Smq()
		h.AddResult(so, dn)
		h.OutputHTML(rw)

		// fmt.Fprintln(rw, r.PostForm)
	})
	fmt.Println("\nhttp://localhost:5667/index\n ")
	server.ListenAndServe()
}

func getOptions(v url.Values) *smq.SmqIn {

	ret := new(smq.SmqIn)
	ret.TextReader = strings.NewReader(v.Get("text"))
	ret.DictReader = strings.NewReader(v.Get("dict"))
	ret.BeginPush, _ = strconv.Atoi(v.Get("ding"))
	ret.SelectKeys = v.Get("csk")
	if v.Get("iss") == "true" {
		ret.IsSingleOnly = true
	}
	if v.Get("as") == "true" {
		ret.IsSpaceDiffHand = true
	}
	return ret
}
