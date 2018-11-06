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
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../.."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

var cookie_ string

// TestGet is a sample to run an endpoint test
func TestGetTitle(t *testing.T) {
	var jsonStr = []byte(`{"id":1}`)
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
			So(m["code"].(float64), ShouldEqual, 200)
			//So(m["data"].(string), ShouldEqual, "test")
		})
	})
}

func TestGetPaper(t *testing.T) {
	var jsonStr = []byte(`{"id":2,"user":{"name":"王泽坤","number":"2018211236","school":"w","class":"2018211305","phone":"17801203047","other":{"friend":"2018211234"}}}`)
	r, _ := http.NewRequest("POST", "/api/getPaper", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	result, _ := ioutil.ReadAll(w.Body)
	m := make(map[string]interface{})
	json.Unmarshal(result, &m)
	beego.Informational(fmt.Sprintf("testing TestGetPaper Code[%d]\n%s", w.Code, string(result)))
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(len(result), ShouldBeGreaterThan, 0)
		})
		Convey("The Result Should Be True", func() {
			So(m["code"].(float64), ShouldEqual, 200)
		})
	})
	cookie_ = w.Header().Get("Set-Cookie")
}

func TestAnswer(t *testing.T) {
	var jsonStr = []byte(`[{"oid":879},{"oid":3334},{"oid":3230},{"oid":2264},{"oid":1686},{"oid":957},{"oid":1794},{"oid":234},{"oid":1024},{"oid":2435},{"oid":1337},{"oid":3501},{"oid":248},{"oid":1320},{"oid":3979},{"oid":4022},{"oid":2100},{"oid":2187},{"oid":812},{"oid":2249},{"oid":835},{"oid":4440},{"oid":651},{"oid":343},{"oid":1140},{"oid":2599},{"oid":457},{"oid":1409},{"oid":1657},{"oid":399},{"oid":1152},{"oid":2907},{"oid":1916},{"oid":3169},{"oid":1527},{"oid":3193},{"oid":280},{"oid":1673},{"oid":2030},{"oid":364},{"oid":963},{"oid":3863},{"oid":1280},{"oid":4472},{"oid":3140},{"oid":564},{"oid":501},{"oid":3684},{"oid":1375},{"oid":2669},{"oid":745},{"oid":2321},{"oid":1389},{"oid":237},{"oid":1533},{"oid":2079},{"oid":467},{"oid":851},{"oid":1837},{"oid":1821},{"oid":4734,"context":"q"},{"oid":4819,"context":"1"},{"oid":4766,"context":"1"},{"oid":4572,"context":"1"},{"oid":4648,"context":"1"},{"oid":4716,"context":"1"},{"oid":4635,"context":"1"},{"oid":4655,"context":"1"},{"oid":4793,"context":"1"},{"oid":4569,"context":"1"},{"oid":4723,"context":"1"},{"oid":4700,"context":"1"},{"oid":4561,"context":"1"},{"oid":4545,"context":"1"},{"oid":4707,"context":"我"},{"oid":4784,"context":"1"},{"oid":4687,"context":"1"},{"oid":4680,"context":"1"},{"oid":4543,"context":"1"},{"oid":4726,"context":"1"}]`)
	r, _ := http.NewRequest("POST", "/api/answer", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Cookie", cookie_)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	result, _ := ioutil.ReadAll(w.Body)
	m := make(map[string]interface{})
	json.Unmarshal(result, &m)
	beego.Trace(fmt.Sprintf("testing TestAnswer Code[%d]\n%s", w.Code, string(result)))
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(len(result), ShouldBeGreaterThan, 0)
		})
		Convey("The Result Should Be True", func() {
			So(m["code"].(float64), ShouldEqual, 200)
		})
	})
}
