build service
-------


设计思路
=======

1. 启用端口监听 发布系统 发来的请求:
* 请求包含4个部分
type Income struct {
	AppId   string `json:"app_id"`
	Project string `json:"project"`
	TaskId  string `json:"task_id"`
	Ref     string `json:"ref"`
}
	
2. 进行以下处理：
* 根据 AppId, 搜索该项目是否有特殊配置（构建前/后的命令执行），如果有，不执行刷新 job template 的操作，如果没有, 根据项目类型刷新对应的 job template
* 根据 AppId, 获取该项目的 job url，并执行构建操作
* Job Build 需要的参数
type Param struct {
	Ref     string `json:"Ref"`
	TaskId  string `json:"TaskId"`
	AppId   string `json:"AppId"`
	Project string `json:"project"`
	Module  string `json:"Module"`
}

3. 新增创建job端口，进行以下处理：
* 接收发布系统的请求：
type CreateIncome struct {
	AppId   string `json:"app_id"`
	Project string `json:"project"`
	Type   string `json:"type"`
	Module string `json:"module"`
	SzOrSh string `json:"SzOrSh"`
}
* 目前暂不对 sh 的项目进行处理
* 收到请求后，在 t_app_job 插入记录即可