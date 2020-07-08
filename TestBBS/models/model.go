package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"log"
)

type Users struct {
	Id int
	Username string
	Pwd string
	Sex string `default:"男"`
}

func init(){
	if err := orm.RegisterDataBase("default", "mysql", "root:123456789@tcp(localhost:3306)/test_bbs?charset=utf8&loc=Local"); err != nil{
		log.Println("数据库连接失败,", err)
	}
	orm.RegisterModel(new(Users))
	orm.RunSyncdb("default", false, true)
}
