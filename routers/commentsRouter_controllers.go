package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["questionaire/controllers:AnswerController"] = append(beego.GlobalControllerRouter["questionaire/controllers:AnswerController"],
        beego.ControllerComments{
            Method: "Answer",
            Router: `/answer`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["questionaire/controllers:AnswerController"] = append(beego.GlobalControllerRouter["questionaire/controllers:AnswerController"],
        beego.ControllerComments{
            Method: "GetPaper",
            Router: `/getPaper`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["questionaire/controllers:AnswerController"] = append(beego.GlobalControllerRouter["questionaire/controllers:AnswerController"],
        beego.ControllerComments{
            Method: "GetTitle",
            Router: `/getTitle`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
