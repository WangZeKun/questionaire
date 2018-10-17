package models

import (
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/tealeg/xlsx"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:745521@/zekin?charset=utf8")
	if err != nil {
		panic(err)
	}
}

func CreatePaper(filename string) (qs []InsertQ, err error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) == 0 {
				break
			}
			type_ := row.Cells[0].String()
			if type_ == "判断" {
				var j Judge
				row.ReadStruct(&j)
				qs = append(qs, j)
			} else if type_ == "填空" {
				var q QA
				row.ReadStruct(&q)
				qs = append(qs, q)
			} else if type_ == "选择" {
				var s Select
				row.ReadStruct(&s)
				qs = append(qs, s)
			}
		}
	}
	return
}

type InsertQ interface {
	Insert() error
}

type Select struct {
	Title  string `xlsx:"2"`
	IsTrue int    `xlsx:"3"`
	A      string `xlsx:"4"`
	B      string `xlsx:"5"`
	C      string `xlsx:"6"`
	D      string `xlsx:"7"`
}

func (s Select) Insert() (err error) {
	stmt, err := DB.Prepare("insert question set title = ?,type=1,paper_id=1")
	if err != nil {
		return
	}
	res, err := stmt.Exec(s.Title)
	if err != nil {
		return
	}
	qid, err := res.LastInsertId()
	if err != nil {
		return
	}
	var body []map[string]interface{}
	stmt, err = DB.Prepare("insert option set question_id=?, is_true=?")
	if err != nil {
		return
	}
	res, err = stmt.Exec(qid, s.IsTrue == 0)
	if err != nil {
		return
	}
	oid, _ := res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": s.A})

	res, err = stmt.Exec(qid, s.IsTrue == 1)
	if err != nil {
		return
	}
	oid, _ = res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": s.B})

	res, err = stmt.Exec(qid, s.IsTrue == 2)
	if err != nil {
		return
	}
	oid, _ = res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": s.C})

	res, err = stmt.Exec(qid, s.IsTrue == 3)
	if err != nil {
		return
	}
	oid, _ = res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": s.D})

	stmt, err = DB.Prepare("update question set body=? where id = ?")
	body_, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err = stmt.Exec(body_, qid)
	return nil
}

type Judge struct {
	Title  string `xlsx:"2"`
	IsTrue bool   `xlsx:"3"`
}

func (j Judge) Insert() (err error) {
	stmt, err := DB.Prepare("insert question set title = ?,type=0,paper_id=1")
	if err != nil {
		return
	}
	res, err := stmt.Exec(j.Title)
	if err != nil {
		return
	}
	qid, err := res.LastInsertId()
	if err != nil {
		return
	}
	var body []map[string]interface{}
	stmt, err = DB.Prepare("insert option set question_id=?, is_true=?")
	if err != nil {
		return
	}
	res, err = stmt.Exec(qid, j.IsTrue)
	if err != nil {
		return
	}
	oid, _ := res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": "正确"})

	res, err = stmt.Exec(qid, !j.IsTrue)
	if err != nil {
		return
	}
	oid, _ = res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": "错误"})

	stmt, err = DB.Prepare("update question set body=? where id = ?")
	body_, err := json.Marshal(body)
	if err != nil {
		return
	}
	res, err = stmt.Exec(body_, qid)
	return
}

type QA struct {
	Title  string `xlsx:"2"`
	Answer string `xlsx:"3"`
}

func (q QA) Insert() (err error) {
	stmt, err := DB.Prepare("insert question set title = ?,type=2,paper_id=1")
	if err != nil {
		return
	}
	res, err := stmt.Exec(q.Title)
	if err != nil {
		return
	}
	qid, err := res.LastInsertId()
	if err != nil {
		return
	}
	var body []map[string]interface{}
	stmt, err = DB.Prepare("insert option set question_id=?, answer=?")
	if err != nil {
		return
	}
	res, err = stmt.Exec(qid, q.Answer)
	if err != nil {
		return
	}
	oid, _ := res.LastInsertId()
	body = append(body, map[string]interface{}{"oid": oid, "context": ""})

	stmt, err = DB.Prepare("update question set body=? where id = ?")
	body_, err := json.Marshal(body)
	if err != nil {
		return err
	}
	res, err = stmt.Exec(body_, qid)
	return
}

//func main() {
//	qs, _ := CreatePaper("/home/zekin/女生节题目.xlsx")
//	for _, q := range qs {
//		fmt.Println(q.Insert())
//	}
//}
