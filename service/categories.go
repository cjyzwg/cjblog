package service

import (
	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/models"
	"io/ioutil"
)

func GetCategories() (models.Categories, error) {

	var content models.Categories

	categoriesDir, err := ioutil.ReadDir(config.Cfg.DocumentPath + "/content")

	if err != nil {
		return content, err
	}

	for _, category := range categoriesDir {

		if !category.IsDir() {
			continue
		}
		var categoryContent models.Category
		markdownList, err := models.GetMarkdownListByCache("/" + category.Name())

		if err != nil {
			return content, err
		}

		listLen := len(markdownList)
		categoryListFileNumber := listLen

		if listLen >= config.Cfg.CategoryListFileNumber {
			categoryListFileNumber = config.Cfg.CategoryListFileNumber
		}

		categoryContent.Name = category.Name()
		categoryContent.Path = "/" + category.Name()
		categoryContent.Number = listLen
		categoryContent.MarkdownFileList = markdownList[0:categoryListFileNumber]

		content = append(content, categoryContent)
	}

	return content, nil
}
