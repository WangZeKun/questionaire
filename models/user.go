package models

import "errors"

type User struct {
	Name      string
	Number    string //学号
	School    string
	Class     string
	Phone     string
	Other     string
	Score     int
	PaperId   int
	questions string
	answer    string
}

var ErrHasScore error = errors.New("用户已答卷")

func (u *User) Insert() (err error) {
	_, err = DB.Exec("insert into user (name, number, school, class, phone, other,score, paper_id) "+
		"select ?, ?, ?, ?, ?, ?, ?,? "+
		"from dual where not exists(select * from user where user.paper_id = ? and user.number = ? and user.score = -1)",
		u.Name, u.Number, u.School, u.Class, u.Phone, u.Other, -1, u.PaperId, u.PaperId, u.Number)
	return
}

func (u *User) Update() error {
	result, err := DB.Exec("update user set score = ? where number=? and paper_id=? and score = -1", u.Score, u.Number, u.PaperId)
	if err != nil {
		return err
	}
	num, _ := result.RowsAffected()
	if num == 0 && u.Score != 0 {
		return ErrHasScore
	}
	//_, err := o.Update(u, "score")
	return err
}
