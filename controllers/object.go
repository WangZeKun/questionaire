package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"questionaire/models"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// Operations about object
type AnswerController struct {
	beego.Controller
}

const (
	BadRequest          string = `{"code":400,"message":"错误的输入!"}`
	InternalServerError string = `{"code":500,"message":"服务器错误!"}`
	ForbiddionTimeEarly string = `{"code":403,"message":"问卷还未发布!"}`
	ForbiddionTimeLate  string = `{"code":403,"message":"问卷已经结束!"}`
	ForbiddionTimeOut   string = `{"code":403,"message":"答题超时!"}`
	NotFound            string = `{"code":404,"message":"未找到试卷!"}`
	UserHasExist        string = `{"code":400,"message":"用户已答卷!"}`
)

func (c *AnswerController) Prepare() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
}

func (c *AnswerController) Options() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
	c.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	c.ServeJSON()
}

// @Title GetTitle
// @Description 拿到题目的Title
// @Param id body json true "题目id"
// @Success 200 {json} {code:200,data:{title:题目}}
// @Failure 200 {code: ,message: }
// @router /getTitle [post]
func (c *AnswerController) GetTitle() {
	var id int
	var ob map[string]json.RawMessage

	fmt.Println(string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob)

	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	err = json.Unmarshal(ob["id"], &id)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	p, err := models.GetPaper(id)
	if err != nil {
		beego.Error(err)
		if err == models.ErrTimeEarly {
			c.Ctx.WriteString(ForbiddionTimeEarly)
		} else if err == models.ErrTimeLate {
			c.Ctx.WriteString(ForbiddionTimeLate)
		} else if err == sql.ErrNoRows {
			c.Ctx.WriteString(NotFound)
		} else {
			c.Ctx.WriteString(InternalServerError)
		}
		return
	}
	c.Data["json"] = addCode(p)
	c.ServeJSON()
}

// @Title GetPaper
// @Description 拿到试卷
// @Param user body object true models.user.User
// @Success 200 {json} {code:200,data:{}}
// @Failure 200 {code: ,message: }
// @router /getPaper [post]
func (c *AnswerController) GetPaper() {

	end, ok := c.GetSession("end").(bool)
	if ok && end {
		beego.Informational("用户已答卷")
		c.Ctx.WriteString(UserHasExist)
		return
	}

	var id int
	var user models.User
	var ob map[string]json.RawMessage

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	if err != nil {
		beego.Informational("用户输入错误：", err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	err = json.Unmarshal(ob["id"], &id)
	if err != nil {
		beego.Informational("用户输入错误：", err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	a, err := json.Marshal(ob["user"])
	err = json.Unmarshal(a, &user)
	if err != nil {
		beego.Informational("用户输入错误：", err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	if !checkNumber(user.Number) {
		beego.Informational(user.Number, " 学号匹配错误！")
		c.Ctx.WriteString(BadRequest)
		return
	}
	user.PaperId = id
	err = user.Insert()
	if err != nil {
		beego.Error(err)
		if strings.HasPrefix(err.Error(), "Error 1062") {
			c.Ctx.WriteString(UserHasExist)
		} else {
			c.Ctx.WriteString(InternalServerError)
		}
		return
	}
	p, err := models.GetPaper(id)
	if err != nil {
		beego.Error(err)
		if err == models.ErrTimeEarly {
			c.Ctx.WriteString(ForbiddionTimeEarly)
		} else if err == models.ErrTimeLate {
			c.Ctx.WriteString(ForbiddionTimeLate)
		} else if err == sql.ErrNoRows {
			c.Ctx.WriteString(NotFound)
		} else {
			c.Ctx.WriteString(InternalServerError)
		}
		return
	}
	err = p.RandomQuestion(10, 5, 5)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(InternalServerError)
		return
	}
	c.SetSession("number", user.Number)
	c.SetSession("pid", user.PaperId)
	c.SetSession("time", time.Now())
	c.Data["json"] = addCode(p)
	c.ServeJSON()
}

// @Title Answer
// @Description 提交答案
// @Param pid body json true "试卷id  答案信息"
// @Success 200 {json} {code:200,data:{}}
// @Failure 200 {code: ,message: }
// @router /answer [post]
func (c *AnswerController) Answer() {
	var ob []models.Option

	number, ok := c.GetSession("number").(string)
	if !ok {
		c.Ctx.WriteString(ForbiddionTimeOut)
		return
	}
	pid, ok := c.GetSession("pid").(int)
	t, ok := c.GetSession("time").(time.Time)
	beego.Debug(time.Now().Sub(t).Minutes())
	if time.Now().Sub(t).Minutes() > 15.0 {
		c.Ctx.WriteString(ForbiddionTimeOut)
		return
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	//pid_f, ok := ob["pid"].(float64) //cnmb
	//if !ok {
	//	c.Ctx.WriteString(BadRequest)
	//	return
	//}
	//pid := int(pid_f)
	u := models.User{Number: number, PaperId: pid}
	//oid_f, ok := ob["oid"].([]float64)
	//if !ok {
	//	c.Ctx.WriteString(BadRequest)
	//	return
	//}
	//data, num, err := models.Check(oid)
	num, err := models.Check(ob)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(InternalServerError)
		return
	}
	u.Score = int(num)
	err = u.Update()
	if err == models.ErrHasScore {
		beego.Error(err)
		c.Ctx.WriteString(UserHasExist)
		return
	} else if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(InternalServerError)
		return
	}
	a := make(map[string]interface{})
	// a["question"] = data
	a["score"] = num //TODO 分数计算
	c.Data["json"] = addCode(a)
	c.DelSession("pid")
	c.DelSession("number")
	c.DelSession("time")
	c.SetSession("end", true)
	c.ServeJSON()
}

func addCode(in interface{}) map[string]interface{} {
	a := make(map[string]interface{})
	a["data"] = in
	a["code"] = 200
	return a
}

func checkNumber(in string) bool {
	exp := regexp.MustCompile(`201[5678]\d{6}`)
	return exp.MatchString(in)
}
