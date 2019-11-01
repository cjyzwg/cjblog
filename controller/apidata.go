package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cjyzwg/forestblog/models"
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("method is:" + r.Method)
	if r.Method == "GET" {
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		fmt.Println(queryForm)
		if err != nil || len(queryForm["category"]) == 0 {
			fmt.Println("query is wrong", err)
			return
		}
		category := queryForm["category"][0]
		fmt.Println(category)
		markdownlist, err := models.GetMarkdownListByCache("/" + category)
		// fmt.Println(markdownlist, err)
		data, err := json.Marshal(markdownlist)
		if err != nil {
			fmt.Println("json.marshal failed, err:", err)
			return
		}
		fmt.Println(string(data))
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)
		fmt.Println(len(result))
		//未知类型的推荐处理方法
		// all_list := make(map[string]string)
		var s string
		if len(result) == 0 {
			list, _ := models.ReadMarkdownDir()
			fmt.Println(list)
			for _, v := range list {
				markdownlist, err := models.GetMarkdownListByCache("/" + v)
				// fmt.Println(markdownlist, err)
				data, err := json.Marshal(markdownlist)
				if err != nil {
					fmt.Println("json.marshal failed, err:", err)
					return
				}

				jsonstr := string(data)
				ss := jsonstr[1 : len(jsonstr)-1]
				s = strings.Join([]string{s, ss}, ",")

				// all_list[v] = string(data)
			}
			fmt.Println(s[1:])
			s = "[" + s[1:] + "]"

			// fmt.Println("-------------------------------------------------------")
			// fmt.Println(s[1:])
			// fmt.Println("-------------------------------------------------------")
			// all_json, _ := json.Marshal(all_list)
			// all_string := string(all_json)
			// fmt.Println(s)
			w.Header().Set("Content-Length", strconv.Itoa(len(s)))
			w.Write([]byte(s))

			return
		}

	}
}
