package routes

import (
	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/controller"
	"net/http"
)

func initApiRoute()  {

	http.HandleFunc(config.Cfg.GitHookUrl, controller.GithubHook)

}