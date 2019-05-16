// main.go
package main

import (
	"github.com/kataras/iris"
	//"github.com/kataras/iris/context"
)

type Income struct {
	AppId   string `json:"app_id"`
	Project string `json:"project"`
	TaskId  string `json:"task_id"`
	Ref     string `json:"ref"`
}

func (*Income) String(appId string, project string, taskId string, ref string) string {
	return "AppId is " + appId + " ,Project is " + project + " ,TaskId is " + taskId + " ,Ref is " + ref
}

func main() {
	app := iris.New()
	app.Post("/build", buildHandler)
	app.Post("/create", createHandler)
	app.Run(iris.Addr(":8007"))
}
