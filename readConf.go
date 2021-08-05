package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	ak     map[int]string
	as     int  // auto_select
	isConf bool // 是否有配置，没有就当作赛码表
}

func readConf(fp string) config {

	f, err := os.Open(fp)
	errHandler(err)
	defer f.Close()

	// 读文件头
	var conf config
	conf.ak = make(map[int]string)
	buff := bufio.NewReader(f)
	conf.isConf = false
	smq_conf := ""
	for i := 1; true; i++ {
		b, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		if strings.HasPrefix(string(b), "smq_conf:") {
			conf.isConf = true
			continue
		}
		if conf.isConf {
			if strings.HasPrefix(string(b), "  ") {
				smq_conf += string(b) + "\n"
			} else {
				break
			}
		} else if i > 30 { // 30行读不到则退出
			break
		}
	}
	// 解析
	if conf.isConf {
		smq := make(map[string]interface{})
		err = yaml.Unmarshal([]byte(smq_conf), &smq)
		if err != nil {
			conf.isConf = false
			fmt.Println(err)
		}
		tmp, ok := smq["alter_keys"].(map[interface{}]interface{})
		if !ok {
			tmp = map[interface{}]interface{}{2: ";", 3: "'"}
		}
		for k, v := range tmp {
			conf.ak[k.(int)] = v.(string)
		}
		conf.as, ok = smq["auto_select"].(int)
		if !ok {
			conf.as = 4
		}
	}
	return conf
}
