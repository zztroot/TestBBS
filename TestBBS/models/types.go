package models

type Type struct {
	Id int64
	Name string `orm:"size(500)"`
	User *User `orm:"rel(fk)"`
	Article []*Article `orm:"reverse(many)"`
}
