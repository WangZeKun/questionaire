package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/astaxie/beego"
)

type Paper struct {
	Id        int                      `json:"id"`
	Title     string                   `json:"title"`
	StartTime time.Time                `json:"-"`
	EndTime   time.Time                `json:"-"`
	Num       []int                    `json:"-"`
	Other     []map[string]string      `json:"other"`
	Questions []map[string]interface{} `json:"questions"`
	TimeLimit float64                  `json:"timeLimit"` //秒为单位
}

func GetPaper(id int) (p Paper, err error) {
	p = Paper{}
	var sTime, eTime string
	var num, other []byte
	row := DB.QueryRow("select id,title,start_time,end_time,num,other,time_limit from paper where id = ?", id)
	row.Scan(&p.Id, &p.Title, &sTime, &eTime, &num, &other, &p.TimeLimit)
	p.StartTime, _ = time.Parse("2006-01-02 15:04:05", sTime)
	p.EndTime, _ = time.Parse("2006-01-02 15:04:05", eTime)
	json.Unmarshal(num, &p.Num)
	json.Unmarshal(other, &p.Other)
	if err != nil {
		return Paper{}, err
	}
	if p.StartTime.After(time.Now()) {
		return Paper{}, ErrTimeEarly
	} else if p.EndTime.Before(time.Now()) {
		return Paper{}, ErrTimeLate
	} else {
		return
	}
}

func (p *Paper) RandomQuestion(select_, judge, text int) (err error) {
	beego.Informational(judge)
	stmt, err := DB.Prepare("select id,title,type,body from question where type = ?  and paper_id = ? order by rand() limit ?")
	if err != nil {
		return
	}
	if judge != 0 {
		rows, err := stmt.Query(0, p.Id, judge)
		if err != nil {
			return err
		}
		p.Questions = append(p.Questions, map[string]interface{}{"data": getData(rows), "name": "判断"})
	}
	if select_ != 0 {
		rows, err := stmt.Query(1, p.Id, select_)
		if err != nil {
			return err
		}
		p.Questions = append(p.Questions, map[string]interface{}{"data": getData(rows), "name": "选择"})
	}
	if text != 0 {
		rows, err := stmt.Query(2, p.Id, text)
		if err != nil {
			return err
		}
		p.Questions = append(p.Questions, map[string]interface{}{"data": getData(rows), "name": "问答"})
	}
	return
}

func getData(rows *sql.Rows) (qs []*Question) {
	for rows.Next() {
		var q Question
		var body []byte
		rows.Scan(&q.QId, &q.Title, &q.Type, &body)
		json.Unmarshal(body, &q.Options)
		qs = append(qs, &q)
	}
	return
}

var (
	ErrTimeEarly error = errors.New("问卷还未发布")
	ErrTimeLate  error = errors.New("问卷已经结束")
)
