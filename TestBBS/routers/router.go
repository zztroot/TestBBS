package routers

import (
	"TestBBS/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//首页路由
    beego.Router("/", &controllers.ArticleController{}, "get:ShowIndex")

    //用户相关
    //登录页路由
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandlerLogin")
    //注册页路由
    beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandlerRegister")
    //个人中心路由
    beego.Router("/centre", &controllers.UserController{}, "get:ShowCentre;post:HandlerCentre")

    //退出登录，删除session
    beego.Router("/exit", &controllers.UserController{}, "get:ExitUser")

    //文章、名言、类型
    //文章详情显示
    beego.Router("/detailArticle", &controllers.ArticleController{}, "get:ShowArticle;post:HandlerComments")
    //发布文章
    beego.Router("/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:AddArticle")
    //名言编辑
    beego.Router("/modifySentence", &controllers.ArticleController{}, "get:ShowModifySentence;post:HandlerModifySentence")
    //名言删除
    beego.Router("/delSentence", &controllers.ArticleController{}, "get:DelSentence")
    //文章删除
    beego.Router("/delArticle", &controllers.ArticleController{}, "get:DelArticle")
    //文章编辑
    beego.Router("/modifyArticle", &controllers.ArticleController{}, "get:ShowModifyArticle;post:HandlerModifyArticle")
}
