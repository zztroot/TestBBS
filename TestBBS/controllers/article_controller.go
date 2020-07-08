package controllers

import "github.com/astaxie/beego"

type ArticleController struct {
	beego.Controller
}

func (a *ArticleController) ShowArticle(){
	a.TplName = "detailArticle.html"
}
