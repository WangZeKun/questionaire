package models

type Question struct {
	Id int `orm:"pk"`
	Title string `orm:"size(255)"`
	Type int8 `orm:"type(tinyint)"`
	Paper *Paper `orm:"rel(fk)"`
	Options []*Option `orm:"-"`
}

