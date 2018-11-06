package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type User struct {
	Name      string
	Number    string //学号
	School    string
	Class     string
	Phone     string
	Other     map[string]interface{}
	Score     int
	PaperId   int
	Answer    string
	StartTime time.Time
	EndTime   time.Time
}

var ErrHasScore error = errors.New("用户已答卷")

func (u *User) Insert() (err error) {
	other, err := json.Marshal(u.Other)
	_, err = DB.Exec("insert into user (name, number, school, class, phone, other,score, paper_id,start_time) "+
		"select ?, ?, ?, ?, ?, ?, ?,?,? ",
		u.Name, u.Number, u.School, u.Class, u.Phone, other, -1, u.PaperId, getTime())
	if err != nil && strings.HasPrefix(err.Error(), "Error 1062") {
		stmt, err := DB.Prepare("select count(number) from user where paper_id=? and number = ? and score=-1")
		if err != nil {
			return err
		}
		row := stmt.QueryRow(u.PaperId, u.Number)
		var count int
		err = row.Scan(&count)
		if err != nil {
			return err
		}
		if count == 1 {
			_, err = DB.Exec("update user set name=?,school=?,class=?,phone=?,other=?,start_time=? where number=? and paper_id = ?",
				u.Name, u.School, u.Class, u.Phone, other, getTime(), u.Number, u.PaperId)
			if err != nil {
				return err
			}
		} else {
			return ErrHasScore
		}
	}
	return nil
}

func (u *User) Update() error {
	result, err := DB.Exec("update user set score = ?,end_time=?,answer=? where number=? and paper_id=? and score = -1", u.Score, getTime(), u.Answer, u.Number, u.PaperId)
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

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
