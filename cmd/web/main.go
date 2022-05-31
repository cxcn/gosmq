package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"sync"

	smq "github.com/cxcn/gosmq"
)

type dictOptions struct {
	Name         string `json:"name"`
	Content      string `json:"content"`
	SingleMode   bool   `json:"singleMode"`
	CommitLeng   int    `json:"commitLeng"`
	CollidedKeys string `json:"collidedKeys"`
	Format       string `json:"format"`
}

type articleOptions struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
type racerOptions struct {
	Algorithm string `json:"algorithm"`
}

type Options struct {
	Dict1   dictOptions    `json:"dict1"`
	Dict2   dictOptions    `json:"dict2"`
	Article articleOptions `json:"article"`
	Racer   racerOptions   `json:"racer"`
}

func parseOptions(src []byte) Options {
	var result Options
	if e := json.Unmarshal(src, &result); e != nil {
		fmt.Println("...Error in Unmarshal: ", e.Error())
	}
	return result
}

func optionsToSmqDict(opt *Options) [2]smq.Dict {
	var result [2]smq.Dict
	result[0].PressSpaceBy = "both"
	result[1].PressSpaceBy = "both"
	result[0].Name = opt.Dict1.Name
	result[1].Name = opt.Dict2.Name
	result[0].Single = opt.Dict1.SingleMode
	result[1].Single = opt.Dict2.SingleMode
	result[0].Format = opt.Dict1.Format
	result[1].Format = opt.Dict2.Format
	result[0].SelectKeys = opt.Dict1.CollidedKeys
	result[1].SelectKeys = opt.Dict2.CollidedKeys
	result[0].PushStart = opt.Dict1.CommitLeng
	result[1].PushStart = opt.Dict2.CommitLeng
	result[0].Algorithm = opt.Racer.Algorithm
	result[1].Algorithm = opt.Racer.Algorithm
	result[0].LoadFromString(opt.Dict1.Content)
	result[1].LoadFromString(opt.Dict2.Content)
	return result
}

func GetResultJson(src []byte) []byte {
	var opts = parseOptions(src)

	var dicts = optionsToSmqDict(&opts)

	s := smq.NewFromString(opts.Article.Name, opts.Article.Content)
	s.Add(&dicts[0])
	s.Add(&dicts[1])
	result, _ := s.ToJson()
	return result
}

//go:embed index.html
var indexPageFile string

//go:embed assets/index.js
var jsBundleFile string

//go:embed assets/index.css
var cssBundleFile string

func HTMLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	fmt.Fprint(w, indexPageFile)
}

func CSSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	fmt.Fprint(w, cssBundleFile)
}

func JSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, jsBundleFile)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println("    post body: ", string(body))
		rjson := GetResultJson(body)
		w.Write(rjson)
		fmt.Println("    returned json: ", string(rjson))
	}
}

func main() {
	http.HandleFunc("/", HTMLHandler)
	http.HandleFunc("/assets/index.js", JSHandler)
	http.HandleFunc("/assets/index.css", CSSHandler)
	http.HandleFunc("/api", PostHandler)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg sync.WaitGroup) {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			panic("...Server failed.")
		}
		wg.Done()
	}(wg)
	doSomething()
	wg.Wait()
}

func doSomething() {
	var name string
	url := "http://localhost:9000"
	switch runtime.GOOS {
	case "windows":
		name = "explorer"
	case "linux":
		name = "xdg-open"
	default:
		name = "open"
	}
	cmd := exec.Command(name, url)
	cmd.Start()
	fmt.Println("Listen and serve: http://localhost:9000")
}
