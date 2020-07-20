package controllers

import (
	"TestBBS/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type UserController struct {
	beego.Controller
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
	user := models.User{}
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
	if sex1 == "false"{
		log.Println("请选择性别！")
		////返回ajax数据
		//resp := make(map[string]interface{})
		//c.Data["json"] = resp
		c.Data["json"] = "请选择性别！"
		c.ServeJSON()
		return
	}
	log.Println(sex0)
	log.Println(sex1)
	 var imgPath string
	if sex0 == "男生"{
		imgPath = "../static/img/man.jpg"
	}else {
		imgPath = "../static/img/woman.jpg"
	}

	o := orm.NewOrm()
	user := models.User{}
	user.Username = username
	user.Pwd = pwd
	user.Sex = sex0
	user.Img = imgPath
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
	//查询用户
	c.Data["username"] = username
	var userName string
	userName = username.(string)
	o := orm.NewOrm()
	user := models.User{Username:userName}
	if err := o.Read(&user, "username"); err != nil{
		log.Println("数据查询失败:", err)
		return
	}

	//查询用户发布的名言
	sentence := []models.Sentence{}
	if _, err :=o.QueryTable("tb_sentence").Filter("user_id", user.Id).All(&sentence); err != nil{
		log.Println("没有发布名言", err)
	}
	if len(sentence) != 0{
		number := rand.Intn(len(sentence))
		c.Data["sentence"] = sentence[number]
	}

	//查询用户发布的文章
	articles := []models.Article{}
	if _, err := o.QueryTable("tb_article").Filter("user_id", user.Id).All(&articles); err != nil{
		log.Println("该用户没有发布文章:", err)
	}
	//查询文章类型
	type articleAndTypes struct {
		Id int64
		Title string
		Type string
		Time time.Time
	}
	articleAndTypeList := []articleAndTypes{}
	for _, v := range articles{
		types := models.Type{Id:v.Type.Id}
		if err := o.Read(&types); err != nil{
			log.Println("文章类型查询失败:", err)
		}
		articleAndTypes := articleAndTypes{}
		articleAndTypes.Id = v.Id
		articleAndTypes.Title = v.Title
		articleAndTypes.Type = types.Name
		articleAndTypes.Time = v.CreateTime
		articleAndTypeList = append(articleAndTypeList, articleAndTypes)
	}
	c.Data["articleAndTypes"] = articleAndTypeList
	c.Data["sentenceList"] = sentence
	c.Data["user"] = user
	c.TplName = "centre.html"
}

//处理个人中心Post请求
//个人中心--个人资料修改Post请求
func (c *UserController) data(){
	usernameID := c.GetSession("username")
	o := orm.NewOrm()
	usernames := c.GetString("username")
	sexs := c.GetString("sex")
	phones := c.GetString("phone")

	//找到用户
	user := models.User{}
	var usernameId string
	usernameId = usernameID.(string)
	user.Username = usernameId
	if err := o.Read(&user, "username"); err != nil{
		log.Println("这个用户不存在!")
		return
	}
	if usernames == ""{
		log.Println("用户名不能为空")
		c.Data["json"] = "用户名不能为空！"
		c.ServeJSON()
		return
	}
	//修改此用户数据
	user.Username = usernames
	user.Sex = sexs
	phonesInt, _ := strconv.ParseInt(phones, 10, 64)
	user.Phone = phonesInt
	if _, err := o.Update(&user, "username", "sex", "phone"); err != nil{
		log.Println("数据修改失败！", err)
		c.Data["json"] = "修改资料失败！"
		c.ServeJSON()
		return
	}

	log.Println("修改数据成功！")
	c.Data["json"] = "修改资料成功！"
	c.ServeJSON()
}

//个人中心--添加分类Post请求
func (c *UserController) add_type(){
	username := c.GetSession("username")
	o := orm.NewOrm()

	//查询当前用户ID
	user := models.User{}
	userName := username.(string)
	user.Username = userName
	_ = o.Read(&user, "username")


	title := c.GetString("type_title")

	ifTitle := models.Type{}
	ifTitle.Name = title
	if err := o.Read(&ifTitle, "name"); err != nil{
		//保存文章类型
		articleType := models.Type{}
		articleType.Name = title
		articleType.User = &user
		if _, err := o.Insert(&articleType); err != nil{
			log.Println("文章类型添加失败!", err)
			c.Data["json"] = "章类型添加失败!"
			c.ServeJSON()
			return
		}
		log.Println("文章类型添加成功！")
		c.Data["json"] = "文章类型添加成功！"
		c.ServeJSON()
		return
	}else {
		log.Println("文章类型添加失败，重复！")
		c.Data["json"] = "类型已经存在！"
		c.ServeJSON()
		return
	}
}

//个人中心--添加名言Post请求
func (c *UserController) add_sentence(){
	username := c.GetSession("username")
	sentenceContent := c.GetString("sentence")
	o := orm.NewOrm()

	//获取该用户
	user := models.User{}
	userName := username.(string)
	user.Username = userName
	if err := o.Read(&user, "username"); err != nil{
		log.Println("没有此用户")
		return
	}

	ifSentence := models.Sentence{}
	ifSentence.Content = sentenceContent
	if err := o.Read(&ifSentence, "content"); err != nil{
		//插入名言
		sentence := models.Sentence{}
		sentence.Content = sentenceContent
		sentence.User = &user
		if _, err := o.Insert(&sentence); err != nil{
			log.Println("名言警句添加失败")
			c.Data["json"] = "名言警句添加失败"
			c.ServeJSON()
			return
		}
		log.Println("名言警句添加成功")
		c.Data["json"] = "名言警句添加成功，2秒后刷新页面"
		c.ServeJSON()
		return
	}else {
		log.Println("名言警句添加失败,重复")
		c.Data["json"] = "您添加的名言已经存在，请重新添加"
		c.ServeJSON()
		return
	}
}

//个人中心--名言管理Post请求
func(c *UserController) handler_sentence(){
	val := c.GetString("val")
	o := orm.NewOrm()
	if val == "获取"{
		id := c.GetString("id")
		Id, _ := strconv.ParseInt(id, 10, 64)
		log.Println(Id)
		sentence := models.Sentence{}
		sentence.Id = Id
		if err := o.Read(&sentence); err != nil{
			log.Println("名言编辑查询失败！")
			return
		}
		c.Data["json"] = sentence
		c.ServeJSON()
	}else {
		id := c.GetString("id")
		Id, _ := strconv.ParseInt(id, 10, 64)
		input := c.GetString("input")
		sentence := models.Sentence{}
		sentence.Id = Id
		if err := o.Read(&sentence); err != nil{
			log.Println("名言查询失败")
			c.Data["json"] = "名言修改失败"
			c.ServeJSON()
			return
		}
		sentence.Content = input
		if _, err := o.Update(&sentence, "content"); err != nil{
			log.Println("名言修改失败")
			c.Data["json"] = "名言修改失败"
			c.ServeJSON()
			return
		}
		c.Data["json"] = "名言修改成功"
		c.ServeJSON()
	}
}

//处理个人中心数据提交
func (c *UserController) HandlerCentre(){
	key := c.GetString("key")
	switch key {
	case "0":
		c.data()
	case "2":
		c.add_sentence()
	case "3":
		log.Println("文章管理")
	case "4":
		c.handler_sentence()
	case "5":
		c.add_type()
	}
	c.TplName = "centre.html"
}

//退出登录，删除session
func (c *UserController) ExitUser(){
	c.DelSession("username")
	c.Redirect("/", 302)
}