package controller

import (
	"encoding/json"
	"fmt"
	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/models"
	"github.com/cjyzwg/forestblog/service"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

//new api article list 
// params: page search
func HandleArticleListData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("query is wrong", err)
		return
	}

	page, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		page = 1
	}
	searchKey := r.Form.Get("search")
	markdownPagination, err := service.GetArticleList(page, "/", searchKey)
	if err != nil {
		fmt.Println("search markdown is wrong", err)
		return
	}
	if markdownPagination.Total > 0 {
		lastpage := markdownPagination.PageNumber[len(markdownPagination.PageNumber)-1]
		if page>lastpage {
			tmpMark := models.MarkdownList{}
			markdownPagination.Markdowns = tmpMark
		}
	}
	rst := ContentResult{}
	rst.Code = 200
	rst.Content = markdownPagination

	jsoncategorylists, _ := json.Marshal(rst)

	w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	w.Write(jsoncategorylists)
}
// params: path
func HandleArticleContentData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("query is wrong", err)
		return
	}

	path := r.Form.Get("path")
	article, _ := models.GetMarkdownDetails(path)

	rst := MarkdowndetailsResult{}
	rst.Code = 200
	rst.Content = article

	jsoncategorylists, _ := json.Marshal(rst)

	w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	w.Write(jsoncategorylists)
}

//category
func HandleCategoryData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("query is wrong", err)
		return
	}


	categories, err := service.GetCategories()
	if err != nil {
		fmt.Println("categoryname markdown is wrong", err)
		return
	}

	rst := CategoryResult{}
	rst.Code = 200
	rst.Content = categories

	jsoncategorylists, _ := json.Marshal(rst)

	w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	w.Write(jsoncategorylists)


}

//categorycontent name
func HandleCategoryContentData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("query is wrong", err)
		return
	}
	categoryName := r.Form.Get("name")
	page, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		page = 1
	}
	content,err := service.GetArticleList(page, "/"+categoryName,"")
	if err != nil {
		fmt.Println("categoryname markdown is wrong", err)
		return
	}
	rst := ContentResult{}
	rst.Code = 200
	rst.Content = content

	jsoncategorylists, _ := json.Marshal(rst)
	w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
	w.Write(jsoncategorylists)


}

//old api 
func HandleData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("method is:" + r.Method)
	if r.Method == "GET" {
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		fmt.Println(queryForm)
		if err != nil {
			fmt.Println("query is wrong", err)
			return
		}
		if len(queryForm["search"]) > 0 {
			// searchKey := queryForm["search"][0]
			page :=1
			markdownPagination, marerr := service.GetArticleList(page, "/", "")
			if marerr != nil {
				fmt.Println("query is wrong", marerr)
				return
			}
			jsoncategorylists, _ := json.Marshal(markdownPagination)
			w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
			w.Write(jsoncategorylists)
			return
		}
		if len(queryForm["category"]) == 0 {
			list, _ := models.ReadMarkdownDir()
			categorylists :=CategoryLists{}
			for _, v := range list {
				categorylist := CategoryList{}
				categorylist.Category = v
				categorylists = append(categorylists,categorylist)
			}
			jsoncategorylists, _ := json.Marshal(categorylists)
			w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
			w.Write(jsoncategorylists)
			return
		}


		category := queryForm["category"][0]
		fmt.Println(category)

		markdownlist, err := models.GetMarkdownListByCache("/" + category)
		// fmt.Println(markdownlist, err)
		if err != nil {
			fmt.Println("no markdown files, err:", err)
			return
		}
		notes := Notes{}
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
		jsonnotes, _ := json.Marshal(notes)
		//fmt.Println(jsonnotes)
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonnotes)))
		w.Write(jsonnotes)


	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)
		fmt.Println(len(result))
		//未知类型的推荐处理方法
		// all_list := make(map[string]string)
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


func HandleDocumentData(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("method is:" + r.Method)
	if r.Method == "POST" {
		mr,err := r.MultipartReader()
		if err != nil{
			fmt.Println("r.MultipartReader() err,",err)
			return
		}
		//r.Body.Close()
		form ,_ := mr.ReadForm(128)
		m := make(map[string]string)
		for k,v := range form.Value{
			fmt.Println("value,k,v = ",k,",",v)
			if k=="title" {
				m["title"] = v[0]
			}
			if k=="body"{
				m["body"] = v[0]
			}
			if k=="category"{
				m["category"] = v[0]
			}
		}
		blogPath := config.CurrentDir + "/" + config.Cfg.DocumentPath
		filename :=blogPath+"/content/"+m["category"]+"/"+m["title"]+".md"
		newfile, error := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766)
		if error != nil {
			fmt.Println(error)
		}
		io.WriteString(newfile, m["body"])

		newfile.Close()

		_, err = exec.LookPath("git")

		if err != nil {
			fmt.Println("请先安装git并克隆博客文档到" + blogPath)

		}
		fmt.Println("blogpath is:"+blogPath)
		fmt.Println("filename is:"+filename)

		//cmd := exec.Command("git", "pull")
		//cmd.Dir = blogPath
		//_, err = cmd.CombinedOutput()
		//if err != nil {
		//	fmt.Println("cmd.Run() failed with", err)
		//	return
		//}

		cmd1 := exec.Command("git", "add",filename)
		cmd1.Dir = blogPath
		_, err = cmd1.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}
		cmd2 := exec.Command("git", "commit","-m","add")
		cmd2.Dir = blogPath
		_, err = cmd2.CombinedOutput()
		if err != nil {
			fmt.Println("cmd2.Run() failed with", err)
			return
		}
		cmd3 := exec.Command("git", "push","origin","master")
		cmd3.Dir = blogPath
		_, err = cmd3.CombinedOutput()
		if err != nil {
			fmt.Println("cmd3.Run() failed with", err)
			return
		}

		var categorylists = "1"
		jsoncategorylists, _ := json.Marshal(categorylists)
		fmt.Println(string(jsoncategorylists))
		//返回的这个是给json用的，需要去掉
		//w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
		w.Write(jsoncategorylists)
		return


	}
}

func HandleDelData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("method is:" + r.Method)
	if r.Method == "GET" {
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		fmt.Println(queryForm)
		if err != nil && len(queryForm["category"]) == 0 && len(queryForm["name"]) == 0 {
			fmt.Println("query is wrong", err)
			return
		}
		category := queryForm["category"][0]
		fmt.Println(category)
		name := queryForm["name"][0]
		fmt.Println(name)
		blogPath := config.CurrentDir + "/" + config.Cfg.DocumentPath
		filename :=blogPath+"/content/"+category+"/"+name
		fmt.Println(filename)
		cmd := exec.Command("rm", "-rf",filename)
		_, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}

		cmd1 := exec.Command("git", "add",filename)
		cmd1.Dir = blogPath
		_, err = cmd1.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}
		cmd2 := exec.Command("git", "commit","-m","add")
		cmd2.Dir = blogPath
		_, err = cmd2.CombinedOutput()
		if err != nil {
			fmt.Println("cmd2.Run() failed with", err)
			return
		}
		cmd3 := exec.Command("git", "push","origin","master")
		cmd3.Dir = blogPath
		_, err = cmd3.CombinedOutput()
		if err != nil {
			fmt.Println("cmd3.Run() failed with", err)
			return
		}

		var categorylists = "1"
		jsoncategorylists, _ := json.Marshal(categorylists)
		fmt.Println(string(jsoncategorylists))
		//返回的这个是给json用的，需要去掉
		w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
		w.Write(jsoncategorylists)
		return


	}
}

type MarkdowndetailsResult struct {
	Code     int  `json:"code"`
	Content  models.MarkdownDetails `json:"content"`
}

type CategoryResult struct {
	Code     int `json:"code"`
	Content  models.Categories `json:"content"`
}

type ContentResult struct {
	Code     int `json:"code"`
	Content  models.MarkdownPagination `json:"content"`
}


type Note struct {
	Title       string `json:"title"`
	Date    	string  `json:"date"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Content     string `json:"content"`
}
type Notes []Note

type CategoryList struct {
	Category       string `json:"category"`
}
type CategoryLists []CategoryList



