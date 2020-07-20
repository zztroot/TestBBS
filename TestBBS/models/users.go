package models

import "time"

type User struct {
	Id int64
	Username string
	Pwd string
	Phone int64
	Sex string
	Img string
	CreateTime time.Time `orm:"auto_now_add"`
	Sentence []*Sentence `orm:"reverse(many)"`
	Article []*Article `orm:"reverse(many)"`
	Comments []*Comments `orm:"reverse(many)"`
	Type []*Type `orm:"reverse(many)"`
}
