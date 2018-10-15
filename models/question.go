package models

type Question struct {
	QId     int       `orm:"pk;column(id)"`
	Title   string    `orm:"size(255)"`
	Type    int8      `orm:"type(tinyint)"`
	Paper   *Paper    `orm:"rel(fk)"`
	Options []*Option `orm:"-"`
}
