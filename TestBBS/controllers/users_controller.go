package controllers

import (
	"TestBBS/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
)

type UserController struct {
	beego.Controller
}

//显示首页
func (c *UserController) ShowIndex() {
	username := c.GetSession("username")
	c.Data["username"] = username
	c.TplName = "index.html"
}

//显示登录页
func (c *UserController) ShowLogin() {
	c.TplName = "login.html"
}

//处理登录
func (c *UserController) HandlerLogin(){
	username := c.GetString("username")
	pwd := c.GetString("pwd")
	if username == "" || pwd == ""{
		log.Println("用户名或密码为空！")
		c.Data["json"] = "用户名或密码为空！"
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	user := models.Users{}
	user.Username = username
	if err := o.Read(&user, "username"); err != nil{
		log.Println("用户名错误或数据库查询失败!:", err)
		c.Data["json"] = "用户名或密码错误!"
		c.ServeJSON()
		return
	}
	if user.Pwd != pwd{
		log.Println("密码错误！")
		c.Data["json"] = "用户名或密码错误!"
		c.ServeJSON()
		return
	}
	//设置cookie值
	//c.Ctx.SetCookie("username", username, time.Hour)
	//u := c.Ctx.GetCookie("usernames")
	//log.Println(u)

	log.Println("登录成功！")
	c.Data["json"] = "0"
	c.ServeJSON()

	//设置session值
	c.SetSession("username", username)
	c.TplName = "/"
}

//显示注册页
func (c *UserController) ShowRegister() {
	c.TplName = "register.html"
}

//处理注册业务
func (c *UserController) HandlerRegister() {
	sex0 := c.GetString("sex0")
	sex1 := c.GetString("sex1")
	username := c.GetString("username")
	pwd := c.GetString("pwd")
	cpwd := c.GetString("cpwd")
	if username == "" || pwd == "" || cpwd == ""{
		log.Println("数据不能为空，请重新输入!")
		c.Data["json"] = "数据不能为空，请重新输入!"
		c.ServeJSON()
		return
	}
	if pwd != cpwd{
		log.Println("两次密码不同，请重新输入!")
		c.Data["json"] = "两次密码不同，请重新输入!"
		c.ServeJSON()
		return
	}
	log.Println(sex0)
	log.Println(sex1)
	if sex1 == "false"{
		log.Println("请选择性别！")
		////返回ajax数据
		//resp := make(map[string]interface{})
		//c.Data["json"] = resp
		c.Data["json"] = "请选择性别！"
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	user := models.Users{}
	user.Username = username
	user.Pwd = pwd
	user.Sex = sex0
	_, err := o.Insert(&user)
	if err != nil{
		log.Println("数据插入失败,", err)
		return
	}

	log.Println("注册成功")
	c.Data["json"] = "恭喜你，注册成功！赶快去登录吧！"
	c.ServeJSON()
}

//显示个人中心
func (c *UserController) ShowCentre(){
	username := c.GetSession("username")
	if username == nil{
		c.Redirect("/", 302)
		return
	}
	c.Data["username"] = username
	c.TplName = "centre.html"
}

//退出登录，删除session
func (c *UserController) ExitUser(){
	c.DelSession("username")
	c.Redirect("/", 302)
}