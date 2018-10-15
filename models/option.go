package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Option struct {
	Id       int `orm:"pk"`
	Context  string
	IsTrue   bool      `orm:"type(bool)"`
	QuestionId int
}

func Check(os []int) (maps []orm.Params, num int64, err error) {
	o := orm.NewOrm()
	a, _ := json.Marshal(os)
	str := "(" + string(a)[1:len(string(a))-1] + ")"
	sql := fmt.Sprintf("select question.title, `option`.context "+
		"from question "+
		"inner join `option` on `option`.question_id = question.id "+
		"where `option`.is_true = false "+
		"and question.id in (select question.id "+
		"	from `option` "+
		"	inner join question on `option`.question_id = question.id "+
		"	where `option`.is_true = false "+
		"	and `option`.id in %s) ", str)
	num, err = o.Raw(sql).Values(&maps)
	return
}
