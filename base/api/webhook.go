package api

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	"github.com/ngaut/log"
	"github.com/pingcap/github-base/pkg/types"
)

// HookBody for parsing webhook
type HookBody struct {
	Repository struct {
		FullName string `json:"full_name"`
	}
}

// Webhook process webhook
func (hdl *ManagerHandler) Webhook(ctx iris.Context) {
	r := ctx.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("body read error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	// restore body for iris ReadJSON use
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	hookBody := HookBody{}
	if err := ctx.ReadJSON(&hookBody); err != nil {
		// body parse error
		log.Errorf("read json error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	repoInfo := strings.Split(hookBody.Repository.FullName, "/")
	if len(repoInfo) != 2 {
		// invalid repo name
		log.Errorf("invalid repo name")
		ctx.StatusCode(iris.StatusInternalServerError)
		if err != nil {
			_, _ = ctx.WriteString(err.Error())
		} else {
			_, _ = ctx.WriteString("invalid repo name")
		}
		return
	}
	repo := &types.Repo{
		Owner: repoInfo[0],
		Repo:  repoInfo[1],
	}

	// restore body for github ValidatePayload use
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	payload, err := github.ValidatePayload(r, []byte(hdl.mgr.GetConfig().GithubSecret))
	if err != nil {
		// invalid payload
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		log.Errorf("invalid payload %v", errors.ErrorStack(err))
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		// event parse err
		log.Errorf("webhook parse error %v", errors.ErrorStack(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	_, _ = ctx.WriteString("ok")
	hdl.mgr.Webhook(repo, event)
}
