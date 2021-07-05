gin_admin_demo
这个项目是以Gin框架为基础搭建的后台管理平台，虽然很多人都认为go是用来开发高性能服务端项目的，但是也难免有要做web管理端的需求，总不能再使用别的语言来开发吧。所以整合出了GinAdmin项目，请大家多提意见指正！


依赖
golang > 1.8

项目目录
|--api  // Api接口控制器

|--comment // 封装的公共方法

|--conf // 配置文件

|--controllers // Admin控制器存在目录

|--logs // 日志存放目录

|--middleware //中间件

|--models //Gorm中的model类

|--router //自定义路由目录

|--statics //css js等静态文件目录

|--uploadfile //上传文件目录

|--views //视图模板目录

分页
使用 comment/util.go 里面的 PageOperation 进行分页
adminDb := models.Db.Table("admin_users").Select("nickname","username").Where("uid != ?", 1)
adminUserData := comment.PageOperation(c, adminDb, 1, &adminUserList)
在html中使用
{{ .adminUserData.PageHtml }}
日志
自定义日志 在 comment/loggers 目录下新建logger
参考 userlog.go 文件

调用自定义的的logger写日志
loggers.UserLogger.Info("无法获取网址",
zap.String("url", "http://www.baidu.com"),
zap.Int("attempt", 3),
zap.Duration("backoff", time.Second),)
数据库

models下定义的文件均需要实现 TableName() string 方法，并将实现该结构体的指针写入到 GetModels 方法中

func GetModels() []interface{} {
	return []interface{}{
		&AdminUsers{},
		&AdminGroup{},
	}
}

数据库迁移,在 cli\cmd 执行命令行工具

go run ginadmin-cli.go db migrate
数据填充，需在相应目录下实现 FillData() 方法执行如下命令

go run ginadmin-cli.go db seed
定时任务
在 comment/cron/cron.go 添加定时执行任务
配置文件
现在 conf/conf.go 添加配置项的 struct 类型，例如

type AppConf struct {
	BaseConf `ini:"base"`
}
type BaseConf struct {
	Port string `ini:"port"`
}
在 conf/conf.ini 添加配置信息

[base]
port=:8091
在代码中调用配置文件的信息

conf.App.BaseConf.Port
模板页面
所有的后台模板都写到 views/template 目录下面，并且分目录存储，调用时按照 目录/模板名称 的方式调用
用户权限
菜单权限定义到 comment/menu/menu.go 文件下，定义完之后在用户组管理里面编辑权限

在控制器中可用从 gin.context 获取权限

privs,_ := c.Get("userPrivs")
template 中判断权限的函数 judgeContainPriv 定义在 comment/template/default.go 文件下

"judgeContainPriv": func(privMap map[string]interface{},priv string)bool {
	//判断权限是all的全通过
	_,o :=privMap["all"]
	if o {
		return true
	}
	_,ok := privMap[priv]
	return ok
},
