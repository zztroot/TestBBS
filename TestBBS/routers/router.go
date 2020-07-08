package routers

import (
	"TestBBS/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//首页路由
    beego.Router("/", &controllers.UserController{}, "get:ShowIndex")
    //登录页路由
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandlerLogin")
    //注册页路由
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandlerRegister")

    //个人中心路由
    beego.Router("/centre", &controllers.UserController{}, "get:ShowCentre")

    //退出登录，删除session
    beego.Router("/exit", &controllers.UserController{}, "get:ExitUser")

    //文章详情显示
    beego.Router("/detailArticle", &controllers.ArticleController{}, "get:ShowArticle")
}
