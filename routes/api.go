package routes

import (
	"net/http"

	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/controller"
)

func initApiRoute() {

	http.HandleFunc(config.Cfg.GitHookUrl, controller.GiteeHook)

}
