package controller

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/helper"
)

func GithubHook(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helper.SedResponse(w, err.Error())
		return
	}

	if config.Cfg.WebHookSecret == "" || r.Header.Get("x-github-event") != "push" {
		helper.SedResponse(w, "No Configuration WebHookSecret Or Not Pushing Events")
		log.Println("No Configuration WebHookSecret Or Not Pushing Events")
		return
	}

	sign := r.Header.Get("X-Hub-Signature")

	bodyContent, err := ioutil.ReadAll(r.Body)

	if err != nil {
		helper.SedResponse(w, err.Error())
		log.Println("WebHook err:" + err.Error())
		return
	}

	if err = r.Body.Close(); err != nil {
		helper.SedResponse(w, err.Error())
		log.Println("WebHook err:" + err.Error())
		return
	}

	mac := hmac.New(sha1.New, []byte(config.Cfg.WebHookSecret))
	mac.Write(bodyContent)
	expectedHash := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	if sign != expectedHash {
		helper.SedResponse(w, "WebHook err:Signature does not match")
		log.Println("WebHook err:Signature does not match")
		return
	}

	helper.SedResponse(w, "ok")

	helper.UpdateArticle()
}

func GiteeHook(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helper.SedResponse(w, err.Error())
		return
	}

	if config.Cfg.WebHookSecret == "" {
		helper.SedResponse(w, "No Configuration WebHookSecret ")
		log.Println("No Configuration WebHookSecret")
		return
	}
	if config.Cfg.WebHookSecret != r.Header.Get("X-Gitee-Token") {
		helper.SedResponse(w, "Configuration WebHookSecret Is Wrong")
		log.Println("Configuration WebHookSecret Is Wrong")
		return
	}

	if r.Header.Get("X-Gitee-Event") != "Push Hook" {
		helper.SedResponse(w, "No Push Events ")
		log.Println("No Push Events ")
		return
	}

	helper.SedResponse(w, "ok")

	helper.UpdateArticle()
}
