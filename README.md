## forestblog
拷贝xusenlin/ForestBlog
改成在config/main.go  
Cfg.GitHookUrl = "/api/git_push_hook"  
Cfg.AppRepository = "https://github.com/cjyzwg/forestblog"
在github的settings,Webhooks添加，playloadurl后面追加/api/git_push_hook  
config.json 中secret就是webhook的secret需要更改
- [<font color=#0099ff>startmyblog</font>](https://github.com/cjyzwg/startmyblog) 根据这个加了几个接口。
