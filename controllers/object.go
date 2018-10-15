package controllers

import (
	"encoding/json"
	"questionaire/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

// @Title GetTitle
// @Description 拿到题目的Title
// @Param id body json true "题目id"
// @Success 200 {json} {code:200,data:{title:题目}}
// @Failure 200 {code: ,message: }
// @router /getTitle [post]
func (c *AnswerController) GetTitle() {
	var id int
	var ob map[string]json.RawMessage

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
		} else if err == orm.ErrNoRows {
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
	var id int
	var user models.User
	var ob map[string]json.RawMessage

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
	a, err := json.Marshal(ob["user"])
	err = json.Unmarshal(a, &user)
	user.PaperId = id
	if err != nil {
		c.Ctx.WriteString(BadRequest)
		return
	}
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
		} else if err == orm.ErrNoRows {
			c.Ctx.WriteString(NotFound)
		} else {
			c.Ctx.WriteString(InternalServerError)
		}
		return
	}
	err = p.Read()
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(InternalServerError)
		return
	}
	c.SetSession("number", user.Number)
	c.SetSession("pid", user.PaperId)
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
	var pid int
	var oid []int
	var ob map[string]json.RawMessage

	number, ok := c.GetSession("number").(string)
	pid_1, ok := c.GetSession("pid").(int)
	if !ok {
		c.Ctx.WriteString(ForbiddionTimeOut)
		return
	}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	err = json.Unmarshal(ob["pid"], &pid)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(BadRequest)
		return
	}
	if pid != pid_1 {
		beego.Informational("试卷错误！")
		c.Ctx.WriteString(BadRequest)
		return
	}
	err = json.Unmarshal(ob["oid"], &oid)
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
	_, num, err := models.Check(oid)
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
	c.ServeJSON()
}

func addCode(in interface{}) map[string]interface{} {
	a := make(map[string]interface{})
	a["data"] = in
	a["code"] = 200
	return a
}
