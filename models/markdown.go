package models

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cjyzwg/forestblog/config"
	"github.com/go-yaml/yaml"
)

func readMarkdown(path string) (Markdown, MarkdownDetails, error) {
	//path=>/categoryName/xxx.md
	fullPath := config.Cfg.DocumentPath + "/content" + path

	categoryName := strings.Replace(path, "/", "", 1)

	if strings.Index(categoryName, "/") == -1 { //文件在根目录下(content/)没有分类名称
		categoryName = ""
	} else {
		categoryName = strings.Split(categoryName, "/")[0]
	}

	var (
		content     Markdown
		fullContent MarkdownDetails
	)

	markdownFile, fileErr := os.Stat(fullPath)

	if fileErr != nil {
		return content, fullContent, fileErr
	}
	if markdownFile.IsDir() {
		return content, fullContent, errors.New("this path is Dir")
	}
	markdown, mdErr := ioutil.ReadFile(fullPath)

	if mdErr != nil {
		return content, fullContent, mdErr
	}
	markdown = bytes.TrimSpace(markdown)
	dateString, body, _ := ParseMetaData(markdown)
	date, _ := time.ParseInLocation("2006-01-02 15:04:05", dateString, time.Local)
	content.Path = path
	content.Category = categoryName
	content.Title = markdownFile.Name()
	content.Date = Time(date)

	fullContent.Markdown = content
	fullContent.Body = string(body)

	if !bytes.HasPrefix(markdown, []byte("```json")) {
		content.Description = cropDesc(markdown)
		return content, fullContent, nil
	}

	markdown = bytes.Replace(markdown, []byte("```json"), []byte(""), 1)
	markdownArrInfo := bytes.SplitN(markdown, []byte("```"), 2)

	content.Description = cropDesc(markdownArrInfo[1])

	if err := json.Unmarshal(bytes.TrimSpace(markdownArrInfo[0]), &content); err != nil {
		return content, fullContent, err
	}

	content.Path = path //保证Path不被用户json赋值，json不能添加`json:"-"`忽略，否则编码到缓存的时候会被忽悠。
	fullContent.Markdown = content
	fullContent.Body = string(markdownArrInfo[1])

	return content, fullContent, nil
}

func cropDesc(c []byte) string {
	content := []rune(string(c))
	contentLen := len(content)

	if contentLen <= config.Cfg.DescriptionLen {
		return string(content[0:contentLen])
	}

	return string(content[0:config.Cfg.DescriptionLen])
}

// 读取路径下的md文件的部分信息json
func GetMarkdown(path string) (Markdown, error) {

	content, _, err := readMarkdown(path)

	if err != nil {
		return content, err
	}
	return content, nil
}

// 读取路径下的md文件完整信息
func GetMarkdownDetails(path string) (MarkdownDetails, error) {

	_, content, err := readMarkdown(path)

	if err != nil {
		return content, err
	}

	return content, nil
}

// 递归获取md文件信息
func getMarkdownList(dir string) (MarkdownList, error) {
	//path=>categoryName
	var fullDir string
	fullDir = config.Cfg.DocumentPath + "/content" + dir

	fileOrDir, err := ioutil.ReadDir(fullDir)

	var mdList MarkdownList

	if err != nil {
		return mdList, err
	}

	for _, fileInfo := range fileOrDir {
		var subDir string
		if "/" == dir {
			subDir = "/" + fileInfo.Name()
		} else {
			subDir = dir + "/" + fileInfo.Name()
		}
		if fileInfo.IsDir() {

			subMdList, err := getMarkdownList(subDir)
			if err != nil {
				return mdList, err
			}
			mdList = append(mdList, subMdList...)
		} else if strings.HasSuffix(strings.ToLower(fileInfo.Name()), "md") {
			markdown, err := GetMarkdown(subDir)
			if err != nil {
				return mdList, err
			}
			mdList = append(mdList, markdown)
		}
	}
	return mdList, nil

}

func GetMarkdownListByCache(dir string) (MarkdownList, error) {

	cacheFileName := fmt.Sprintf("%x", md5.Sum([]byte(dir)))

	cacheFilePath := config.CurrentDir + "/cache/" + cacheFileName + ".json"

	var content MarkdownList

	cacheFile, cacheErr := ioutil.ReadFile(cacheFilePath)

	if cacheErr == nil && json.Unmarshal(cacheFile, &content) == nil {
		return content, nil
	}

	content, err := getMarkdownList(dir)

	if err != nil {
		return content, err
	}

	sort.Sort(content)
	markdownListJson, err := json.Marshal(content)

	if err != nil {
		return content, err
	}

	cacheDir := config.CurrentDir + "/cache"
	cacheInfo, err := os.Stat(cacheDir)

	if err != nil || !cacheInfo.IsDir() {
		if os.Mkdir(cacheDir, os.ModePerm) != nil {
			return content, err
		}
	}

	err = ioutil.WriteFile(cacheFilePath, markdownListJson, os.ModePerm)

	if err != nil {
		return content, err
	}

	return content, nil
}

func ReadMarkdownBody(path string) (string, error) {

	fullPath := config.Cfg.DocumentPath + path

	markdown, err := ioutil.ReadFile(fullPath)

	if err != nil {
		return "", err
	}

	return string(markdown), nil
}

func ReadMarkdownDir() ([]string, error) {
	dir_list, err := ioutil.ReadDir(config.Cfg.DocumentPath + "/content")
	if err != nil {
		fmt.Println("read dir error")
		return nil, err
	}
	var list []string
	for _, v := range dir_list {
		list = append(list, v.Name())
	}
	// fmt.Println(list)
	return list, err
}

type Meta struct {
	Title      string
	Date       string
	Tags       []string
	Categories []string
}

func ParseMetaData(fileData []byte) (string, string, error) {
	buf := bufio.NewReader(bytes.NewReader(fileData))
	metaData := []byte{}
	content := []byte{}
	metaStart := false
	metaFinished := false
	// get the data between the line only have "---"
	for {
		line, isPrefix, lineErr := buf.ReadLine()
		if lineErr != nil {
			break
		}
		for _, s := range []string{"---", "--- ", " --- ", " ---", "----"} {
			if string(line) == s {
				if metaStart {
					metaFinished = true

				}
				metaStart = true
			}
		}

		if !isPrefix {
			line = append(line, []byte("\n")...)
		}

		if metaStart && !metaFinished {
			metaData = append(metaData, line...)
		}

		if metaFinished {
			content = append(content, line...)
		}
	}
	meta := &Meta{}
	err := yaml.Unmarshal(metaData, meta)
	if err != nil {
		return "", "", err
	}

	return meta.Date, string(content), nil
}
