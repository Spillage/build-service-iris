package main

import (
	"strings"

	"github.com/heiing/logs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/context"
)

type CreateIncome struct {
	AppId   string `json:"app_id"`
	Project string `json:"project"`
	//PreBuildCommand   string `json:"PreBuildCommand"`
	//AfterBuildCommand string `json:"AfterBuildCommand"`
	Type   string `json:"type"`
	Module string `json:"module"`
	SzOrSh string `json:"SzOrSh"`
}

func (*CreateIncome) String(appId string, project string, Type string, module string, SzOrSh string) string {
	return "AppId is " + appId + " ,Project is " + project + " ,Type is " + Type + " ,Module is " + module + " ,SzOrSh is " + SzOrSh
}

func jobUrlCreate(msg *CreateIncome) (jobUrl string) {

	jobUrl = ""
	var jobName string

	switch msg.SzOrSh {
	case "sz":
		jobName = strings.Replace(msg.Project[7:], "/", "-", -1) + "-" + msg.AppId
		jobUrl = config.CI.CISZ + "/job/" + jobName + "-Build-CI/"

		return

	case "sh":
		jobName = strings.Replace(msg.Project[7:], "/", "-", -1) + "-" + msg.AppId
		jobUrl = config.CI.CISH + "/job/" + jobName + "-Build-CI/"

		return

	default:
		return
	}
}

func createHandler(ctx context.Context) {
	m := &CreateIncome{}
	//err := json.NewDecoder(r.Body).Decode(m)
	err := ctx.ReadJSON(m)

	if nil != err {
		ctx.WriteString("Illegal params: " + m.String(m.AppId, m.Project, m.Type, m.Module, m.SzOrSh))
		logs.Error("Illegal params: ", m)
	}

	logs.Debug("Receive from Deployer: ", m)

	if "" != m.AppId && "" != m.Module && "" != m.Project && "" != m.SzOrSh && "" != m.Type && 7 < len(m.Project) {
		//go createJob(m)

		jobUrl := jobUrlCreate(m)
		newBuildJob := BuildJob{AppId: m.AppId, Project: m.Project, Job: jobUrl, Type: m.Type, Module: m.Module}
		oldBuildJob := &BuildJob{}
		oldBuildJob.find(m.AppId)
		if !oldBuildJob.exists() {

			newBuildJob.create()
			logs.Debug("The record created! ", newBuildJob)
			ctx.WriteString("Created")

		} else {
			oldBuildJob.Job = newBuildJob.Job
			oldBuildJob.update()
			logs.Debug("The record updated!", oldBuildJob)
			ctx.WriteString("Updated")
		}

		//w.WriteHeader(201)

	} else {
		ctx.WriteString("Illegal params: " + m.String(m.AppId, m.Project, m.Type, m.Module, m.SzOrSh))
		logs.Error("Illegal params! ", m)
	}
}
