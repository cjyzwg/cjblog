package routes

import (
	"net/http"

	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/controller"
)

func initWebRoute() {

	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/categories", controller.Categories)
	http.HandleFunc("/works", controller.Works)
	http.HandleFunc("/about", controller.About)
	http.HandleFunc("/api", controller.HandleData)// old api
	http.HandleFunc("/savedocument", controller.HandleDocumentData)
	http.HandleFunc("/del", controller.HandleDelData)
	http.HandleFunc("/apis/articlelist", controller.HandleArticleListData) //new api
	http.HandleFunc("/apis/articlecontent", controller.HandleArticleContentData) //new api
	http.HandleFunc("/apis/category", controller.HandleCategoryData) //new api
	http.HandleFunc("/apis/categorycontent", controller.HandleCategoryContentData) //new api
	//二级页面
	http.HandleFunc("/article", controller.Article)
	http.HandleFunc("/category", controller.CategoryArticle)
	http.HandleFunc(config.Cfg.DashboardEntrance, controller.Dashboard)
	//静态文件服务器
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("resources/public"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(config.Cfg.DocumentPath+"/assets"))))

}
