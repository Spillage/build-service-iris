// build-handler.go
package main

import (
	//"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bytes"

	"github.com/heiing/logs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/context"
)

func buildHandler(ctx context.Context) {
	m := &Income{}

	err := ctx.ReadJSON(m)
	//err := json.NewDecoder(r.Body).Decode(m)
	if nil != err {
		ctx.WriteString("Illegal params: " + m.String(m.AppId, m.Project, m.TaskId, m.Ref))
		logs.Error("Illegal params: ", m)
	}

	logs.Debug("Receive from Deployer: ", m.AppId, m.TaskId, m)

	if "" != m.AppId && "" != m.Project && "" != m.Ref && "" != m.TaskId && 7 < len(m.Project) {
		go startWork(m)
		//w.WriteHeader(201)

	} else {
		ctx.WriteString("Illegal params: " + m.String(m.AppId, m.Project, m.TaskId, m.Ref))
		logs.Error("Illegal param! ", m)
	}

}

func startWork(msg *Income) {
	job := &BuildJob{}

	if err := job.find(msg.AppId); nil != err {
		logs.Error("Find job faild, AppId: ", msg.AppId, "; Error: ", err)
		return
	}

	logs.Debug("Start work for msg: ", msg.AppId, msg.TaskId, msg)

	build(job, msg)
}

func build(job *BuildJob, msg *Income) {

	ext := &BuildExtend{}
	if err := ext.find(job.AppId); nil != err {

		logs.Error("Find Build Extend Faild, AppId: ", job.AppId, "; Error: ", err)
		return
	}

	updateJob(job)

	logs.Debug("The msg are ", msg)

	urlBuild := job.Job + "buildWithParameters?"

	param := url.Values{}

	param.Set("AppId", msg.AppId)
	param.Set("TaskId", msg.TaskId)
	param.Set("Project", msg.Project)
	param.Set("Module", job.Module)
	param.Set("Ref", msg.Ref)

	logs.Info("The param is ", param)

	if strings.Contains(urlBuild, "ci-sz") {
		res, err := basicAuthPost(urlBuild, bytes.NewBufferString(param.Encode()), "Content-Type: application/x-www-form-urlencoded")
		if nil != err {
			logs.Error("Post Request Faild: ", err, ", URL: ", urlBuild)
			return
		}

		logs.Info("The job Build: ", res.Status, "； URL：", urlBuild)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.Info(string(bodyBytes))

	} else {
		res, err := basicAuthPost(urlBuild, bytes.NewBufferString(param.Encode()), "Content-Type: application/x-www-form-urlencoded")
		if nil != err {
			logs.Error("Post Request Faild: ", err, ", URL: ", urlBuild)
			return
		}

		logs.Info("The job Build: ", res.Status, "； URL：", urlBuild)
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		logs.Info(string(bodyBytes))

	}

}

func updateJob(job *BuildJob) {

	if &job == nil {
		logs.Error("The job is nil !!!", job)
		return
	}

	path := config.BuilderTemplateConfigs[job.Type].Path

	logs.Debug("The template path is ", path)
	logs.Debug("The job type is ", job.Type)

	template, err := os.Open(path)
	if nil != err {
		logs.Error("Open template file Faild: ", err)
		return
	}
	defer template.Close()

	url := job.Job + "config.xml"

	if strings.Contains(url, "ci-sz") {
		response, err := basicAuthPost(url, template)
		logs.Debug("The responses of update job is ", response.Status)

		if response.StatusCode != 200 {
			jobCreateInJenkins(job)
		}

		if nil != err {
			logs.Error("Update SZ job, Post Request Faild: ", err, ", URL: ", url)
			return
		}
	} else {
		response, err := basicAuthPost(url, template)
		logs.Debug("The responses of update job is ", response.Status)
		//if response.Status != strconv.Itoa(200) {
		if response.StatusCode != 200 {
			jobCreateInJenkins(job)
		}
		if nil != err {
			logs.Error("Update SH job, Post Request Faild: ", err, ", URL: ", url)
			return
		}
	}

	logs.Info("template ", path, " updated!")
}

func basicAuthPost(url string, body io.Reader, headers ...string) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if nil != err {
		return nil, err
	}
	for _, herder := range headers {
		kvs := strings.SplitN(herder, ":", 2)
		key, val := kvs[0], kvs[1]
		req.Header.Add(strings.TrimSpace(key), strings.TrimSpace(val))
	}
	return basicAuthPostReq(req, body)
}

func basicAuthPostReq(req *http.Request, body io.Reader) (res *http.Response, err error) {
	//取消jenkins用户验证
	//req.SetBasicAuth("admin", "sz8968")
	client := &http.Client{}
	res, err = client.Do(req)
	return
}

/*
func basicPost(url string, body io.Reader) (res *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if nil != err {
		return nil, err
	}
	client := &http.Client{}
	res, err = client.Do(req)
	return
} */

func jobCreateInJenkins(job *BuildJob) {

	path := config.BuilderTemplateConfigs[job.Type].Path
	template, err := os.Open(path)
	if nil != err {
		logs.Error("Open template file Faild: ", err)
		return
	}
	defer template.Close()

	jobName := strings.Replace(job.Project[7:], "/", "-", -1) + "-" + job.AppId + "-Build-CI"

	if strings.Contains(job.Job, "ci-sz") {

		response, err := basicAuthPost(config.CI.CISZ+"/createItem?name="+jobName, template, "Content-Type: application/xml")
		resBytes, _ := ioutil.ReadAll(response.Body)
		logs.Debug("The response of create is ", string(resBytes))
		if nil != err {
			logs.Error("Create SZ job, Post Request Faild: ", err, ", URL: ", "我是马赛克4"+jobName)
			return
		}
	} else {
		response, err := basicAuthPost(config.CI.CISH+"/createItem?name="+jobName, template, "Content-Type: application/xml")
		resBytes, _ := ioutil.ReadAll(response.Body)
		logs.Debug("The response of create is ", string(resBytes))
		if nil != err {
			logs.Error("Create SH job, Post Request Faild: ", err, ", URL: ", "我是马赛克5"+jobName)
			return
		}
	}

}
