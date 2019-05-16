// pre-handler.go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BuildExtend struct {
	Id                int `gorm:"primary_key"`
	AppId             string
	PreBuildCommand   string
	AfterBuildCommand string
}

func (this *BuildExtend) TableName() string {
	return "t_app_build_extend"
}

func (this *BuildExtend) find(appId string) (err error) {
	db, err := gorm.Open("mysql", config.Mysql.Conn)
	if nil != err {
		return err
	}
	defer db.Close()
	db.Where("app_id = ?", appId).First(this)
	return
}

func (this *BuildExtend) exists() bool {
	return this.AppId != ""
}
