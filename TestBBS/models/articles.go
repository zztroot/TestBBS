package models

import "time"

type Article struct {
	Id int64
	Title string `orm:"size(100)"`
	Content string `orm:"size(16000)"`
	User *User `orm:"rel(fk)"`
	Type *Type `orm:"rel(fk)"`
	CreateTime time.Time `orm:"auto_now_add"`
	Comments []*Comments `orm:"reverse(many)"`
}

