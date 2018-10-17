package models

import (
	"fmt"
)

type Option struct {
	OId    int    `json:"oid"`
	Answer string `json:"context"`
}

func Check(os []Option) (num int64, err error) {
	//	sql := fmt.Sprintf("select question.title, `option`.context "+
	//		"from question "+
	//		"inner join `option` on `option`.question_id = question.id "+
	//		"where `option`.is_true = false "+
	//		"and question.id in (select question.id "+
	//		"	from `option` "+
	//		"	inner join question on `option`.question_id = question.id "+
	//		"	where `option`.is_true = false "+
	//		"	and (`option`.id,`option`.answer) in %s) ", trans(os))
	sql := fmt.Sprintf("select count(id) from option where (id,is_true,answer) in %s", trans(os))
	fmt.Println(sql)
	row := DB.QueryRow(sql)
	err = row.Scan(&num)
	if err != nil {
		return
	}
	return num * 5, err
}

func trans(os []Option) string {
	str := "("
	for i, o := range os {
		if i == len(os)-1 {
			str += fmt.Sprintf("(%d,true,'%s')", o.OId, o.Answer)
		} else {
			str += fmt.Sprintf("(%d,true,'%s'),", o.OId, o.Answer)
		}
	}
	str += ")"
	return str
}
