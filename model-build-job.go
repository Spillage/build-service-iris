package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BuildJob struct {
	Id      int `gorm:"primary_key"`
	AppId   string
	Job     string
	Project string
	Module  string
	Type    string
}

/*
type MultiJob struct {
	Jobs []*BuildJob
}
*/

func (this *BuildJob) TableName() string {
	return "t_app_job"
}

//返回一条记录
func (this *BuildJob) find(appId string) (err error) {
	db, err := gorm.Open("mysql", config.Mysql.Conn)
	if nil != err {
		return
	}
	defer db.Close()
	db.Where("app_id = ?", appId).First(this)
	return
}

func (this *BuildJob) create() {
	db, err := gorm.Open("mysql", config.Mysql.Conn)
	if nil != err {
		return
	}
	defer db.Close()

	db.Create(this)

}

//更新joburl
func (this *BuildJob) update() {
	db, err := gorm.Open("mysql", config.Mysql.Conn)
	if nil != err {
		return
	}
	defer db.Close()

	db.Save(this)
}

func (this *BuildJob) exists() bool {
	return this.AppId != ""
}
