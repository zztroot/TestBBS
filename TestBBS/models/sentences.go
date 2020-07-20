package models

type Sentence struct {
	Id int64
	Content string `orm:"size(999)"`
	User *User `orm:"rel(fk)"`
}
