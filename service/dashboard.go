package service

import "github.com/cjyzwg/forestblog/config"

func SetThemeColor(index int)  {

	config.Cfg.ThemeColor = config.Cfg.ThemeOption[index]
	//需要将配置写入app.json吗?
}
