package models

import "errors"
import "github.com/astaxie/beego/orm"

type User struct {
	Name    string
	Number  string //学号
	School  string
	Class   string
	Phone   string
	Other   string
	Score   int
	PaperId int
}

var ErrHasScore error = errors.New("用户已答卷")

func (u *User) Insert() (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("insert into user (name, number, school, class, phone, other, paper_id) "+
		"select ?, ?, ?, ?, ?, ?, ? "+
		"from dual where not exists(select * from user where user.paper_id = ? and user.number = ? and user.score = 0)",
		u.Name, u.Number, u.School, u.Class, u.Phone, u.Other, u.PaperId, u.PaperId, u.Number).Exec()
	return
}

func (u *User) Update() error {
	o := orm.NewOrm()
	result, err := o.Raw("update user set score = ? where number=? and paper_id=? and score !=0", u.Score, u.Number, u.PaperId).Exec()
	num, _ := result.RowsAffected()
	if num == 0 {
		return ErrHasScore
	}
	//_, err := o.Update(u, "score")
	return err
}
