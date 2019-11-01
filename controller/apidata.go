package controller

import (
	"encoding/json"
	"fmt"
	"github.com/cjyzwg/forestblog/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HandleDataTest(w http.ResponseWriter, r *http.Request) {
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
				//mardownlist 即可生成json,数据
				// for _, m := range markdownlist {
				// 	fmt.Println("-------------------------------------------------------")
				// 	mark, err := models.GetMarkdownDetails(m.Path)
				// 	fmt.Println(mark, err)
				// 	fmt.Println("-------------------------------------------------------")
				// }

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
			notes := Notes{}
			for _, v := range list {

				markdownlist, err := models.GetMarkdownListByCache("/" + v)
				if err != nil {
					fmt.Println("no markdown files, err:", err)
					return
				}
				//mardownlist 即可生成json,数据
				for _, m := range markdownlist {
					note := Note{}
					mark, err := models.GetMarkdownDetails(m.Path)
					if err != nil {
						fmt.Println("the path is wrong,err:", err)
						return
					}
					note.Date = mark.Date.Format("2006-01-02 15:04:05")
					note.Title = mark.Title
					note.Description = m.Description
					note.Category = mark.Category
					note.Content = mark.Body
					notes = append(notes,note)

				}

			}
			jsonnotes, _ := json.Marshal(notes)
			fmt.Println("-------------------------------------------------------")

			//fmt.Println(string(jsonnotes))
			fmt.Println("-------------------------------------------------------")

			w.Header().Set("Content-Length", strconv.Itoa(len(jsonnotes)))
			//w.Write([]byte(s))
			w.Write(jsonnotes)

			return
		}

	}
}

type Note struct {
	Title       string `json:"title"`
	Date    	string  `json:"date"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Content     string `json:"content"`
}
type Notes []Note
