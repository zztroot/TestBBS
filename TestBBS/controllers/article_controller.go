package controllers

import (
	"TestBBS/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
	"github.com/garyburd/redigo/redis"
)

type ArticleController struct {
	beego.Controller
}

//显示首页
func (a *ArticleController) ShowIndex() {
	type articleUserImg struct {
		Id int64
		Title string
		CreateTime time.Time
		Img string
	}

	//获取用户登录信息
	username := a.GetSession("username")
	a.Data["username"] = username

	//查询所有文章
	o := orm.NewOrm()
	articles := []models.Article{}
	if _, err := o.QueryTable("tb_article").All(&articles); err != nil{
		log.Println("所有文章查询失败:", err)
		return
	}

	//查询文章对应用户图片
	var articlesUsers []articleUserImg

	for _, v := range articles{
		var articleUser articleUserImg
		user := models.User{Id:v.User.Id}
		_ = o.Read(&user)
		articleUser.Id = v.Id
		articleUser.Title = v.Title
		articleUser.CreateTime = v.CreateTime
		articleUser.Img = user.Img
		articlesUsers = append(articlesUsers, articleUser)
	}

	//查询阅读量前5文章
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil{
		log.Println("redis connect failed:", err)
		return
	}

	//文章阅读量结构体
	type articleReads struct {
		articleReadNumberId int64
		articlesId int64
	}
	articleReadList := []articleReads{}
	for _, v := range articles{
		number, err := redis.Int64(conn.Do("get", v.Id))
		if err != nil{
			log.Println("redis 获取文字阅读量错误:", err)
		}
		newArticleReadList := articleReads{}
		newArticleReadList.articleReadNumberId = number
		newArticleReadList.articlesId = v.Id
		articleReadList = append(articleReadList, newArticleReadList)
	}
	conn.Close()
	// 从小到大排序(稳定排序)
	sort.SliceStable(articleReadList, func(i, j int) bool {
		if articleReadList[i].articleReadNumberId < articleReadList[j].articleReadNumberId {
			return true
		}
		return false
	})
	log.Println(articleReadList)
	//取阅读量最多的5篇文章
	if len(articleReadList) <= 5 {
		var articleLists []models.Article
		for _, v := range articleReadList{
			articleList := models.Article{}
			if err := o.QueryTable("tb_article").Filter("id", v.articlesId).One(&articleList); err != nil{
				log.Println("文章查询失败:", err)
			}
			articleLists = append(articleLists, articleList)
		}
		a.Data["articleLists"] = articleLists
	}else {
		var ArticleReadLists []articleReads
		for i:=1; i<6; i++{
			ArticleReadLists = append(ArticleReadLists, articleReadList[len(articleReadList)-i])
		}

		var articleLists []models.Article
		for _, v := range ArticleReadLists{
			articleList := models.Article{}
			if err := o.QueryTable("tb_article").Filter("id", v.articlesId).One(&articleList); err != nil{
				log.Println("文章查询失败:", err)
			}
			articleLists = append(articleLists, articleList)
		}
		a.Data["articleLists"] = articleLists
	}

	//实现文章发布时间顺序显示
	var newArticlesUsers []articleUserImg
	for i:=1; i<=len(articlesUsers); i++{
		index := articlesUsers[len(articlesUsers)-i]
		newArticlesUsers = append(newArticlesUsers, index)
	}

	//取出最新发布的5篇文章
	if len(articles) <= 5{
		a.Data["articlesLists"] = articles
	}else {
		var articlesLists []models.Article
		for i:=1; i<6; i++{
			articlesLists = append(articlesLists, articles[len(articles)-i])
		}
		a.Data["articlesLists"] = articlesLists
	}

	//查询名言展示到首页
	sentence := []models.Sentence{}
	if _, err := o.QueryTable("tb_sentence").All(&sentence); err != nil{
		log.Println("名言查询失败:", err)
		return
	}
	number := rand.Intn(len(sentence))

	//查询发布名言用户
	sentenceUser := models.User{Id:sentence[number].User.Id}
	_ = o.Read(&sentenceUser)

	a.Data["sentence"] = sentence[number].Content
	a.Data["sentenceUser"] = sentenceUser
	a.Data["articles"] = newArticlesUsers
	a.TplName = "index.html"
}

//显示文章详细界面
func (a *ArticleController) ShowArticle(){
	//获取用户登录信息
	username := a.GetSession("username")
	a.Data["username"] = username

	//获取文章ID
	articleId, _ := a.GetInt64("id")
	a.Data["articleId"] = articleId
	o := orm.NewOrm()

	//进入文章详细界面，redis加1
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer conn.Close()
	_, err = conn.Do("incr", articleId)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	number, err := redis.Int64(conn.Do("get", articleId))
	if err != nil{
		log.Println("redis阅读量获取失败：", err)
	}


	//进入文章详细界面，mysql--read加1
	read := models.Read{}
	read.ArticleId = articleId
	err = o.Read(&read, "article_id")
	if err != nil{
		log.Printf("数据库read表读取失败:", err)
		//根据文章ID，设置mysql--read表的文章ID和阅读量
		read := models.Read{}
		read.ArticleId = articleId
		read.Read = 0
		_, err = o.Insert(&read)
		if err != nil{
			log.Printf("文章阅读量插入失败")
		}
	}
	read.Read += 1
	_, err = o.Update(&read, "read")
	if err != nil{
		log.Printf("数据库插入失败！")
	}

	//判断mysql 和 redis数据是否相同
	if number == read.Read{
		a.Data["number"] = number
	}else {
		a.Data["number"] = read.Read

		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("Connect to redis error", err)
			return
		}
		_, err = conn.Do("set", articleId, read.Read)
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
		conn.Close()
	}

	//查询文章
	article := models.Article{Id:articleId}
	err = o.Read(&article)
	if err != nil{
		log.Println("文章查询失败:", err)
		return
	}

	//查询用户
	user := models.User{Id:article.User.Id}
	if err = o.Read(&user); err != nil{
		log.Println("用户查询失败:", err)
		return
	}

	//查询文章评论
	type userAndComment struct {
		UserId int64
		UserName string
		Content string
		ToUser string
		ToUserId string
	}
	comments := []models.Comments{}
	_, err = o.QueryTable("tb_comments").Filter("article_id", articleId).All(&comments)
	if err != nil{
		log.Println("文章评论查询失败:", err)
		return
	}
	var userAndCommentList []userAndComment
	for _, v := range comments{
		user := models.User{Id: v.User.Id}
		_ = o.Read(&user)
		toUser := models.User{Id: v.ToUser}
		_ = o.Read(&toUser)
		userandcom := userAndComment{}
		userandcom.UserId = user.Id
		userandcom.UserName = user.Username
		userandcom.Content = v.Content
		userandcom.ToUser = toUser.Username
		strId := strconv.FormatInt(toUser.Id, 10)
		userandcom.ToUserId = strId
		userAndCommentList = append(userAndCommentList, userandcom)
	}

	a.Data["comments"] = userAndCommentList
	a.Data["user"] = user
	a.Data["article"] = article
	a.TplName = "detailArticle.html"
}

//文章留言
func (a *ArticleController) HandlerComments(){
	key := a.GetString("key")
	if key == "0"{
		username := a.GetSession("username")
		comContent := a.GetString("comContent")
		articleId := a.GetString("articleId")
		o := orm.NewOrm()

		if comContent == ""{
			log.Println("评论内容为空")
			a.Data["json"] = "评论内容为空"
			a.ServeJSON()
			return
		}

		//查询用户
		user := models.User{}
		user.Username = username.(string)
		if err := o.Read(&user, "username"); err != nil{
			log.Println("用户查询失败:", err)
			return
		}

		//查询文章
		Id, _ := strconv.ParseInt(articleId, 10, 64)
		article := models.Article{Id:Id}
		if err := o.Read(&article); err != nil{
			log.Println("文章查询失败:", err)
			return
		}

		//插入评论
		comment := models.Comments{}
		comment.Content = comContent
		comment.User = &user
		comment.Article = &article
		_, err := o.Insert(&comment)
		if err != nil{
			log.Println("评论保存失败:", err)
			a.Data["json"] = "评论失败"
			a.ServeJSON()
			return
		}
		log.Println("评论成功")
		a.Data["json"] = "评论成功"
		a.ServeJSON()
	}else {
		username := a.GetSession("username")
		comContent := a.GetString("comContent")
		articleId := a.GetString("articleId")
		toUserId := a.GetString("toUserId")
		o := orm.NewOrm()

		if comContent == ""{
			log.Println("评论内容为空")
			a.Data["json"] = "评论内容为空"
			a.ServeJSON()
			return
		}
		//查询用户
		user := models.User{}
		user.Username = username.(string)
		if err := o.Read(&user, "username"); err != nil{
			log.Println("用户查询失败:", err)
			return
		}

		//判断用户是不是自己回复自己
		toUserIds, _ := strconv.ParseInt(toUserId, 10, 64)
		if user.Id == toUserIds{
			log.Println("用户自己回复自己，错误")
			a.Data["json"] = "不能自己回复自己"
			a.ServeJSON()
			return
		}

		//查询出to用户的名称
		toUser := models.User{Id: toUserIds}
		err := o.Read(&toUser)
		if err != nil{
			log.Println("to用户不存在:", err)
			return
		}

		//查询文章
		Id, _ := strconv.ParseInt(articleId, 10, 64)
		article := models.Article{Id:Id}
		if err := o.Read(&article); err != nil{
			log.Println("文章查询失败:", err)
			return
		}

		//插入评论
		comment := models.Comments{}
		comment.Content = comContent
		comment.Article = &article
		comment.User = &user
		comment.ToUser = toUser.Id
		_, err = o.Insert(&comment)
		if err != nil{
			log.Println("评论回复失败:", err)
			a.Data["json"] = "评论回复失败"
			a.ServeJSON()
			return
		}
		a.Data["json"] = "评论回复成功"
		a.ServeJSON()
	}
}

//显示发布文章界面
func (a *ArticleController) ShowAddArticle(){
	username := a.GetSession("username")
	if username == nil{
		a.Redirect("/", 302)
		return
	}
	a.Data["username"] = username
	var userName string
	userName = username.(string)
	o := orm.NewOrm()
	user := models.User{Username:userName}
	if err := o.Read(&user, "username"); err != nil{
		log.Println("数据查询失败:", err)
		return
	}
	//查询出该用户所有的文章类型
	types := []models.Type{}
	_, err := o.QueryTable("tb_type").Filter("user_id", user.Id).All(&types)
	if err != nil{
		log.Println("你还没有添加文章类型或数据查询失败！")
		return
	}
	a.Data["types"] = types
	a.TplName = "addArticle.html"
}

//个人中心--发布文章处理
func (c *ArticleController) AddArticle(){
	articleTitle := c.GetString("title")
	articleType := c.GetString("types")
	articleContent := c.GetString("content")

	if articleTitle == "" || articleContent == ""{
		log.Println("标题和内容为空")
		c.Data["json"] = "请输入标题和内容"
		c.ServeJSON()
		c.TplName = "addArticle.html"
		return
	}

	//查询出type名称
	o := orm.NewOrm()
	types := models.Type{}
	typeInt, _ := strconv.ParseInt(articleType, 10, 64) //字符串转int64
	types.Id = typeInt
	if err := o.Read(&types); err != nil{
		log.Println("文章类型查询失败:", err)
		return
	}

	//查询出当前用户
	username := c.GetSession("username")
	userName := username.(string)
	user := models.User{}
	user.Username = userName
	if err := o.Read(&user, "username"); err != nil{
		log.Println("查询用户失败:", err)
		return
	}

	//添加文章
	article := models.Article{}
	article.Title = articleTitle
	article.Content = articleContent
	article.Type = &types
	article.User = &user
	if _, err := o.Insert(&article); err != nil{
		log.Println("文章发布失败:", err)
		c.Data["json"] = "文章发布失败，请重试!"
		c.ServeJSON()
		c.TplName = "addArticle.html"
		return
	}

	//根据文章ID，设置redis的key和value
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	_, err = conn.Do("set", article.Id, 0)
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	conn.Close()

	//根据文章ID，设置mysql--read表的文章ID和阅读量
	read := models.Read{}
	read.ArticleId = article.Id
	read.Read = 0
	_, err = o.Insert(&read)
	if err != nil{
		log.Printf("文章阅读量插入失败")
	}

	log.Println("文章发布成功！")
	c.Data["json"] = "文章发布成功"
	c.ServeJSON()
	c.TplName = "addArticle.html"
}

//个人中心--显示名言编辑 get
func (c *ArticleController) ShowModifySentence(){
	username := c.GetSession("username")
	if username == nil{
		c.Redirect("/", 302)
		return
	}
	c.Data["username"] = username

	//通过ID获取需要修改的名言
	sentenceId, _ := c.GetInt64("id")
	o := orm.NewOrm()
	sentence := models.Sentence{Id:sentenceId}
	if err := o.Read(&sentence); err != nil{
		log.Println("查询名言失败:", err)
		return
	}
	c.Data["sentence"] = sentence
	c.TplName = "modifySentence.html"
}

//个人中心--名言编辑 post
func (c *ArticleController) HandlerModifySentence(){
	//通过ID查询需要修改的名言
	sentenceId := c.GetString("sentenceId")
	Id, _ := strconv.ParseInt(sentenceId, 10, 64)
	sentenceContent := c.GetString("sentenceContent")

	o := orm.NewOrm()
	sentence := models.Sentence{Id:Id}
	if err := o.Read(&sentence); err != nil{
		log.Println("查询名言失败:", err)
		return
	}
	sentence.Content = sentenceContent
	if _, err := o.Update(&sentence, "content"); err != nil{
		log.Println("修改名言失败:", err)
		c.Data["json"] = "改名言失败"
		c.ServeJSON()
		return
	}
	c.Data["json"] = "修改名言成功，3秒后自动跳转个人中心"
	c.ServeJSON()
}

//个人中心--名言删除 get
func (c *ArticleController) DelSentence(){
	username := c.GetSession("username")
	if username == nil{
		c.Redirect("/", 302)
		return
	}

	Id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	sentence := models.Sentence{Id:Id}
	_, err := o.Delete(&sentence)
	if err != nil{
		log.Println("删除名言失败，请稍后重试！")
		return
	}
	log.Println("删除名言成功")
	c.Redirect("/centre", 302)
}

//个人中心--文章删除 get
func (c *ArticleController) DelArticle(){
	username := c.GetSession("username")
	if username == nil{
		c.Redirect("/", 302)
		return
	}

	Id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	article := models.Article{Id:Id}
	if _, err := o.Delete(&article); err != nil{
		log.Println("文章删除失败")
	}
	log.Println("文章删除成功")

	//文章删除后，并删除redis数据
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil{
		log.Println("redis connect failed:", err)
	}

	_, err = conn.Do("del", Id)
	if err != nil{
		log.Println("redis data delete failed:", err)
	}
	conn.Close()
	log.Println("redis删除成功")
	c.Redirect("/centre", 302)

}

//个人中心--显示文章编辑 get
func (c *ArticleController) ShowModifyArticle(){
	username := c.GetSession("username")
	if username == nil{
		c.Redirect("/", 302)
		return
	}
	Id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	article := models.Article{Id:Id}
	if err := o.Read(&article); err != nil{
		log.Println("文章查询失败!:", err)
	}

	type articleAndType struct {
		Id int64
		Title string
		Type string
		Content string
	}
	//查询此文章类型
	types := models.Type{Id: article.Type.Id}
	if err := o.Read(&types); err != nil{
		log.Println("文章类型查询失败:", err)
	}

	//给自定义结构体赋值
	articleAndTypeList := articleAndType{}
	articleAndTypeList.Id = types.Id
	articleAndTypeList.Title = article.Title
	articleAndTypeList.Type = types.Name
	articleAndTypeList.Content = article.Content

	c.Data["articleId"] = Id
	c.Data["articleAndTypeList"] = articleAndTypeList
	c.TplName = "modifyArticle.html"
}

//个人中心--文章编辑 post
func (c *ArticleController) HandlerModifyArticle(){
	articleTitle := c.GetString("title")
	articleType := c.GetString("types")
	articleContent := c.GetString("content")
	Id := c.GetString("id")
	articleId, _ := strconv.ParseInt(Id, 10, 64)

	if articleTitle == "" || articleContent == ""{
		log.Println("标题和内容为空")
		c.Data["json"] = "请输入标题和内容"
		c.ServeJSON()
		c.TplName = "modifyArticle.html"
		return
	}

	//查询出type名称
	o := orm.NewOrm()
	types := models.Type{}
	typeInt, _ := strconv.ParseInt(articleType, 10, 64) //字符串转int64
	types.Id = typeInt
	if err := o.Read(&types); err != nil{
		log.Println("文章类型查询失败:", err)
		return
	}

	//查询出当前用户
	username := c.GetSession("username")
	userName := username.(string)
	user := models.User{}
	user.Username = userName
	if err := o.Read(&user, "username"); err != nil{
		log.Println("查询用户失败:", err)
		return
	}

	//修改文章
	article := models.Article{Id:articleId}
	if err := o.Read(&article); err != nil{
		log.Println("没有找到该文章:", err)
	}
	article.Title = articleTitle
	article.Content = articleContent
	article.Type = &types
	if _, err := o.Update(&article, "title", "content", "type_id"); err != nil{
		log.Println("文章修改失败:", err)
		c.Data["json"] = "文章修改失败，请重试!"
		c.ServeJSON()
		c.TplName = "modifyArticle.html"
		return
	}
	log.Println("文章修改成功！")
	c.Data["json"] = "文章修改成功"
	c.ServeJSON()
	c.TplName = "addArticle.html"
}