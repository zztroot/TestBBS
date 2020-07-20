package models

type Comments struct {
	Id int64
	Content string `orm:"size(16000)"`
	User *User `orm:"rel(fk)"`
	Article *Article `orm:"rel(fk)"`
	ToUser int64
}
