package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"log"
)

func init(){
	dbHost := beego.AppConfig.String("dbhost")
	dbPort := beego.AppConfig.String("dbport")
	dbUser := beego.AppConfig.String("dbuser")
	dbPwd := beego.AppConfig.String("dbpwd")
	dbName := beego.AppConfig.String("dbname")
	dbUrl := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Local"
	if err := orm.RegisterDataBase("default", "mysql", dbUrl); err != nil{
		log.Println("数据库连接失败,", err)
	}
	orm.RegisterModelWithPrefix(beego.AppConfig.String("dbprefix"), new(User), new(Article), new(Comments), new(Type), new(Sentence), new(Read))
	if err := orm.RunSyncdb("default", false, true); err != nil{
		log.Println("数据库启动失败,", err)
		return
	}
}
