package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	_ "questionaire/routers"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

var cookie_ string

// TestGet is a sample to run an endpoint test
func TestGetTitle(t *testing.T) {
	var jsonStr = []byte(`{"id":5}`) //unvaild id
	r, _ := http.NewRequest("POST", "/api/getTitle", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	result, _ := ioutil.ReadAll(w.Body)
	m := make(map[string]interface{})
	json.Unmarshal(result, &m)
	beego.Trace(fmt.Sprintf("testing TestGetTitle Code[%d]\n%s", w.Code, string(result)))

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(len(result), ShouldBeGreaterThan, 0)
		})
		Convey("The Result Should Be True", func() {
			So(m["code"].(float64), ShouldEqual, 404)
			//So(m["data"].(string), ShouldEqual, "test")
		})
	})
}

func TestGetPaper(t *testing.T) {
	var jsonStr = []byte(`{"id":5,"user":{"Name":"test","Number":"2012211516","School":"123456","Class":"2018211203","Phone":"123132113"}}`) //unvaild id
	r, _ := http.NewRequest("POST", "/api/getPaper", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	result, _ := ioutil.ReadAll(w.Body)
	m := make(map[string]interface{})
	json.Unmarshal(result, &m)
	beego.Trace(fmt.Sprintf("testing TestGetPaper Code[%d]\n%s", w.Code, string(result)))
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(len(result), ShouldBeGreaterThan, 0)
		})
		Convey("The Result Should Be True", func() {
			So(m["code"].(float64), ShouldEqual, 404)
		})
	})
	cookie_ = w.Header().Get("Set-Cookie")
}
