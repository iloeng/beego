// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSubstr(t *testing.T) {
	s := `012345`
	if Substr(s, 0, 2) != "01" {
		t.Error("should be equal")
	}
	if Substr(s, 0, 100) != "012345" {
		t.Error("should be equal")
	}
	if Substr(s, 12, 100) != "012345" {
		t.Error("should be equal")
	}
}

func TestHtml2str(t *testing.T) {
	h := `<HTML><style></style><script>x<x</script></HTML><123>  123\n


	\n`
	if HTML2str(h) != "123\\n\n\\n" {
		t.Error("should be equal")
	}
}

func TestDateFormat(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	if ss := DateFormat(tt, "2006-01-02 15:04:05"); ss != "2013-07-01 13:27:42" {
		t.Errorf("2013-07-01 13:27:42 does not equal %v", ss)
	}
}

func TestDate(t *testing.T) {
	ts := "Mon, 01 Jul 2013 13:27:42 CST"
	tt, _ := time.Parse(time.RFC1123, ts)

	if ss := Date(tt, "Y-m-d H:i:s"); ss != "2013-07-01 13:27:42" {
		t.Errorf("2013-07-01 13:27:42 does not equal %v", ss)
	}
	if ss := Date(tt, "y-n-j h:i:s A"); ss != "13-7-1 01:27:42 PM" {
		t.Errorf("13-7-1 01:27:42 PM does not equal %v", ss)
	}
	if ss := Date(tt, "D, d M Y g:i:s a"); ss != "Mon, 01 Jul 2013 1:27:42 pm" {
		t.Errorf("Mon, 01 Jul 2013 1:27:42 pm does not equal %v", ss)
	}
	if ss := Date(tt, "l, d F Y G:i:s"); ss != "Monday, 01 July 2013 13:27:42" {
		t.Errorf("Monday, 01 July 2013 13:27:42 does not equal %v", ss)
	}
}

func TestCompareRelated(t *testing.T) {
	if !Compare("abc", "abc") {
		t.Error("should be equal")
	}
	if Compare("abc", "aBc") {
		t.Error("should be not equal")
	}
	if !Compare("1", 1) {
		t.Error("should be equal")
	}
	if CompareNot("abc", "abc") {
		t.Error("should be equal")
	}
	if !CompareNot("abc", "aBc") {
		t.Error("should be not equal")
	}
	if !NotNil("a string") {
		t.Error("should not be nil")
	}
}

func TestHtmlquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	if Htmlquote(s) != h {
		t.Error("should be equal")
	}
}

func TestHtmlunquote(t *testing.T) {
	h := `&lt;&#39;&nbsp;&rdquo;&ldquo;&amp;&#34;&gt;`
	s := `<' ”“&">`
	if Htmlunquote(h) != s {
		t.Error("should be equal")
	}
}

func TestParseForm(t *testing.T) {
	type ExtendInfo struct {
		Hobby []string `form:"hobby"`
		Memo  string
	}

	type OtherInfo struct {
		Organization string `form:"organization"`
		Title        string `form:"title"`
		ExtendInfo
	}

	type user struct {
		ID      int         `form:"-"`
		tag     string      `form:"tag"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text"`
		Email   string
		Intro   string    `form:",textarea"`
		StrBool bool      `form:"strbool"`
		Date    time.Time `form:"date,2006-01-02"`
		OtherInfo
	}

	u := user{}
	form := url.Values{
		"ID":           []string{"1"},
		"-":            []string{"1"},
		"tag":          []string{"no"},
		"username":     []string{"test"},
		"age":          []string{"40"},
		"Email":        []string{"test@gmail.com"},
		"Intro":        []string{"I am an engineer!"},
		"strbool":      []string{"yes"},
		"date":         []string{"2014-11-12"},
		"organization": []string{"beego"},
		"title":        []string{"CXO"},
		"hobby":        []string{"", "Basketball", "Football"},
		"memo":         []string{"nothing"},
	}
	if err := ParseForm(form, u); err == nil {
		t.Fatal("nothing will be changed")
	}
	if err := ParseForm(form, &u); err != nil {
		t.Fatal(err)
	}
	if u.ID != 0 {
		t.Errorf("ID should equal 0 but got %v", u.ID)
	}
	if len(u.tag) != 0 {
		t.Errorf("tag's length should equal 0 but got %v", len(u.tag))
	}
	if u.Name.(string) != "test" {
		t.Errorf("Name should equal `test` but got `%v`", u.Name.(string))
	}
	if u.Age != 40 {
		t.Errorf("Age should equal 40 but got %v", u.Age)
	}
	if u.Email != "test@gmail.com" {
		t.Errorf("Email should equal `test@gmail.com` but got `%v`", u.Email)
	}
	if u.Intro != "I am an engineer!" {
		t.Errorf("Intro should equal `I am an engineer!` but got `%v`", u.Intro)
	}
	if !u.StrBool {
		t.Errorf("strboll should equal `true`, but got `%v`", u.StrBool)
	}
	y, m, d := u.Date.Date()
	if y != 2014 || m.String() != "November" || d != 12 {
		t.Errorf("Date should equal `2014-11-12`, but got `%v`", u.Date.String())
	}
	if u.Organization != "beego" {
		t.Errorf("Organization should equal `beego`, but got `%v`", u.Organization)
	}
	if u.Title != "CXO" {
		t.Errorf("Title should equal `CXO`, but got `%v`", u.Title)
	}
	if u.Hobby[0] != "" {
		t.Errorf("Hobby should equal ``, but got `%v`", u.Hobby[0])
	}
	if u.Hobby[1] != "Basketball" {
		t.Errorf("Hobby should equal `Basketball`, but got `%v`", u.Hobby[1])
	}
	if u.Hobby[2] != "Football" {
		t.Errorf("Hobby should equal `Football`, but got `%v`", u.Hobby[2])
	}
	if len(u.Memo) != 0 {
		t.Errorf("Memo's length should equal 0 but got %v", len(u.Memo))
	}
}

func TestRenderForm(t *testing.T) {
	type user struct {
		ID      int         `form:"-"`
		Name    interface{} `form:"username"`
		Age     int         `form:"age,text,年龄："`
		Sex     string      `form:"sex,select"`
		Email   []string
		Intro   string `form:",textarea"`
		Ignored string `form:"-"`
	}

	u := user{Name: "test", Intro: "Some Text", Sex: "Male,Female"}
	output := RenderForm(u)
	if output != template.HTML("") {
		t.Errorf("output should be empty but got %v", output)
	}
	output = RenderForm(&u)
	result := template.HTML(
		`Name: <input name="username" type="text" value="test"></br>` +
			`年龄：<input name="age" type="text" value="0"></br>` +
			`Sex: <select name="sex"></br>` +
			`  <option value="Male"> Male </option></br>` +
			`  <option value="Female"> Female </option></br>` +
			`</select></br>` +
			`Intro: <textarea name="Intro">Some Text</textarea>`)
	if output != result {
		t.Errorf("output should equal `%v` but got `%v`", result, output)
	}
}

func TestRenderFormField(t *testing.T) {
	html := renderFormField("Label: ", "Name", "text", "Value", "", "", false)
	if html != `Label: <input name="Name" type="text" value="Value">` {
		t.Errorf("Wrong html output for input[type=text]: %v ", html)
	}

	html = renderFormField("Label: ", "Name", "textarea", "Value", "", "", false)
	if html != `Label: <textarea name="Name">Value</textarea>` {
		t.Errorf("Wrong html output for textarea: %v ", html)
	}

	html = renderFormField("Label: ", "Name", "textarea", "Value", "", "", true)
	if html != `Label: <textarea name="Name" required>Value</textarea>` {
		t.Errorf("Wrong html output for textarea: %v ", html)
	}
}

func TestParseFormTag(t *testing.T) {
	// create struct to contain field with different types of struct-tag `form`
	type user struct {
		All            int `form:"name,text,年龄："`
		NoName         int `form:",hidden,年龄："`
		OnlyLabel      int `form:",,年龄："`
		OnlyName       int `form:"name" id:"name" class:"form-name"`
		Ignored        int `form:"-"`
		Required       int `form:"name" required:"true"`
		IgnoreRequired int `form:"name"`
		NotRequired    int `form:"name" required:"false"`
	}

	objT := reflect.TypeOf(&user{}).Elem()

	label, name, fType, _, _, ignored, _ := parseFormTag(objT.Field(0))
	if !(name == "name" && label == "年龄：" && fType == "text" && !ignored) {
		t.Errorf("Form Tag with name, label and type was not correctly parsed.")
	}

	label, name, fType, _, _, ignored, _ = parseFormTag(objT.Field(1))
	if !(name == "NoName" && label == "年龄：" && fType == "hidden" && !ignored) {
		t.Errorf("Form Tag with label and type but without name was not correctly parsed.")
	}

	label, name, fType, _, _, ignored, _ = parseFormTag(objT.Field(2))
	if !(name == "OnlyLabel" && label == "年龄：" && fType == "text" && !ignored) {
		t.Errorf("Form Tag containing only label was not correctly parsed.")
	}

	label, name, fType, id, class, ignored, _ := parseFormTag(objT.Field(3))
	if !(name == "name" && label == "OnlyName: " && fType == "text" && !ignored &&
		id == "name" && class == "form-name") {
		t.Errorf("Form Tag containing only name was not correctly parsed.")
	}

	_, _, _, _, _, ignored, _ = parseFormTag(objT.Field(4))
	if !ignored {
		t.Errorf("Form Tag that should be ignored was not correctly parsed.")
	}

	_, name, _, _, _, _, required := parseFormTag(objT.Field(5))
	if !(name == "name" && required) {
		t.Errorf("Form Tag containing only name and required was not correctly parsed.")
	}

	_, name, _, _, _, _, required = parseFormTag(objT.Field(6))
	if !(name == "name" && !required) {
		t.Errorf("Form Tag containing only name and ignore required was not correctly parsed.")
	}

	_, name, _, _, _, _, required = parseFormTag(objT.Field(7))
	if !(name == "name" && !required) {
		t.Errorf("Form Tag containing only name and not required was not correctly parsed.")
	}
}

func TestMapGet(t *testing.T) {
	// test one level map
	m1 := map[string]int64{
		"a": 1,
		"1": 2,
	}

	if res, err := MapGet(m1, "a"); err == nil {
		if res.(int64) != 1 {
			t.Errorf("Should return 1, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, "1"); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, 1); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 2 level map
	m2 := M{
		"1": map[string]float64{
			"2": 3.5,
		},
	}

	if res, err := MapGet(m2, 1, 2); err == nil {
		if res.(float64) != 3.5 {
			t.Errorf("Should return 3.5, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 5 level map
	m5 := M{
		"1": M{
			"2": M{
				"3": M{
					"4": M{
						"5": 1.2,
					},
				},
			},
		},
	}

	if res, err := MapGet(m5, 1, 2, 3, 4, 5); err == nil {
		if res.(float64) != 1.2 {
			t.Errorf("Should return 1.2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// check whether element not exists in map
	if res, err := MapGet(m5, 5, 4, 3, 2, 1); err == nil {
		if res != nil {
			t.Errorf("Should return nil, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}
}

func Test_eq(t *testing.T) {
	tests := []struct {
		a      interface{}
		b      interface{}
		result bool
	}{
		{uint8(1), int(1), true},
		{uint8(3), int(1), false},
		{uint16(1), int(1), true},
		{uint16(3), int(1), false},
		{uint32(1), int(1), true},
		{uint32(3), int(1), false},
		{uint(1), int(1), true},
		{uint(3), int(1), false},
		{uint64(1), int(1), true},
		{uint64(3), int(1), false},
		{int8(-1), uint(1), false},
		{int16(-2), uint(1), false},
		{int32(-3), uint(1), false},
		{int64(-4), uint(1), false},
		{int8(1), uint(1), true},
		{int16(1), uint(1), true},
		{int32(1), uint(1), true},
		{int64(1), uint(1), true},
		{int8(-1), uint8(1), false},
		{int16(-2), uint8(1), false},
		{int32(-3), uint8(1), false},
		{int64(-4), uint8(1), false},
		{int8(1), uint8(1), true},
		{int16(1), uint8(1), true},
		{int32(1), uint8(1), true},
		{int64(1), uint8(1), true},
		{int8(-1), uint16(1), false},
		{int16(-2), uint16(1), false},
		{int32(-3), uint16(1), false},
		{int64(-4), uint16(1), false},
		{int8(1), uint16(1), true},
		{int16(1), uint16(1), true},
		{int32(1), uint16(1), true},
		{int64(1), uint16(1), true},
		{int8(-1), uint32(1), false},
		{int16(-2), uint32(1), false},
		{int32(-3), uint32(1), false},
		{int64(-4), uint32(1), false},
		{int8(1), uint32(1), true},
		{int16(1), uint32(1), true},
		{int32(1), uint32(1), true},
		{int64(1), uint32(1), true},
		{int8(-1), uint64(1), false},
		{int16(-2), uint64(1), false},
		{int32(-3), uint64(1), false},
		{int64(-4), uint64(1), false},
		{int8(1), uint64(1), true},
		{int16(1), uint64(1), true},
		{int32(1), uint64(1), true},
		{int64(1), uint64(1), true},
	}

	for _, test := range tests {
		if res, err := eq(test.a, test.b); err != nil {
			if res != test.result {
				t.Errorf("a:%v(%T) equals b:%v(%T) should be %v", test.a, test.a, test.b, test.b, test.result)
			}
		}
	}
}

func Test_lt(t *testing.T) {
	tests := []struct {
		a      interface{}
		b      interface{}
		result bool
	}{
		{uint8(1), int(3), true},
		{uint8(1), int(1), false},
		{uint8(3), int(1), false},
		{uint16(1), int(3), true},
		{uint16(1), int(1), false},
		{uint16(3), int(1), false},
		{uint32(1), int(3), true},
		{uint32(1), int(1), false},
		{uint32(3), int(1), false},
		{uint(1), int(3), true},
		{uint(1), int(1), false},
		{uint(3), int(1), false},
		{uint64(1), int(3), true},
		{uint64(1), int(1), false},
		{uint64(3), int(1), false},
		{int(-1), int(1), true},
		{int(1), int(1), false},
		{int(1), int(3), true},
		{int(3), int(1), false},
		{int8(-1), uint(1), true},
		{int8(1), uint(1), false},
		{int8(1), uint(3), true},
		{int8(3), uint(1), false},
		{int16(-1), uint(1), true},
		{int16(1), uint(1), false},
		{int16(1), uint(3), true},
		{int16(3), uint(1), false},
		{int32(-1), uint(1), true},
		{int32(1), uint(1), false},
		{int32(1), uint(3), true},
		{int32(3), uint(1), false},
		{int64(-1), uint(1), true},
		{int64(1), uint(1), false},
		{int64(1), uint(3), true},
		{int64(3), uint(1), false},
		{int(-1), uint(1), true},
		{int(1), uint(1), false},
		{int(1), uint(3), true},
		{int(3), uint(1), false},
		{int8(-1), uint8(1), true},
		{int8(1), uint8(1), false},
		{int8(1), uint8(3), true},
		{int8(3), uint8(1), false},
		{int16(-1), uint8(1), true},
		{int16(1), uint8(1), false},
		{int16(1), uint8(3), true},
		{int16(3), uint8(1), false},
		{int32(-1), uint8(1), true},
		{int32(1), uint8(1), false},
		{int32(1), uint8(3), true},
		{int32(3), uint8(1), false},
		{int64(-1), uint8(1), true},
		{int64(1), uint8(1), false},
		{int64(1), uint8(3), true},
		{int64(3), uint8(1), false},
		{int(-1), uint8(1), true},
		{int(1), uint8(1), false},
		{int(1), uint8(3), true},
		{int(3), uint8(1), false},
		{int8(-1), uint16(1), true},
		{int8(1), uint16(1), false},
		{int8(1), uint16(3), true},
		{int8(3), uint16(1), false},
		{int16(-1), uint16(1), true},
		{int16(1), uint16(1), false},
		{int16(1), uint16(3), true},
		{int16(3), uint16(1), false},
		{int32(-1), uint16(1), true},
		{int32(1), uint16(1), false},
		{int32(1), uint16(3), true},
		{int32(3), uint16(1), false},
		{int64(-1), uint16(1), true},
		{int64(1), uint16(1), false},
		{int64(1), uint16(3), true},
		{int64(3), uint16(1), false},
		{int(-1), uint16(1), true},
		{int(1), uint16(1), false},
		{int(1), uint16(3), true},
		{int(3), uint16(1), false},
		{int8(-1), uint32(1), true},
		{int8(1), uint32(1), false},
		{int8(1), uint32(3), true},
		{int8(3), uint32(1), false},
		{int16(-1), uint32(1), true},
		{int16(1), uint32(1), false},
		{int16(1), uint32(3), true},
		{int16(3), uint32(1), false},
		{int32(-1), uint32(1), true},
		{int32(1), uint32(1), false},
		{int32(1), uint32(3), true},
		{int32(3), uint32(1), false},
		{int64(-1), uint32(1), true},
		{int64(1), uint32(1), false},
		{int64(1), uint32(3), true},
		{int64(3), uint32(1), false},
		{int(-1), uint32(1), true},
		{int(1), uint32(1), false},
		{int(1), uint32(3), true},
		{int(3), uint32(1), false},
		{int8(-1), uint64(1), true},
		{int8(1), uint64(1), false},
		{int8(1), uint64(3), true},
		{int8(3), uint64(1), false},
		{int16(-1), uint64(1), true},
		{int16(1), uint64(1), false},
		{int16(1), uint64(3), true},
		{int16(3), uint64(1), false},
		{int32(-1), uint64(1), true},
		{int32(1), uint64(1), false},
		{int32(1), uint64(3), true},
		{int32(3), uint64(1), false},
		{int64(-1), uint64(1), true},
		{int64(1), uint64(1), false},
		{int64(1), uint64(3), true},
		{int64(3), uint64(1), false},
		{int(-1), uint64(1), true},
		{int(1), uint64(1), false},
		{int(1), uint64(3), true},
		{int(3), uint64(1), false},
	}

	for _, test := range tests {
		if res, err := lt(test.a, test.b); err != nil {
			if res != test.result {
				t.Errorf("a:%v(%T) lt b:%v(%T) should be %v", test.a, test.a, test.b, test.b, test.result)
			}
		}
	}
}

func TestRenderFormSecurity(t *testing.T) {
	type UserProfile struct {
		DisplayName string `form:"displayName,text,Name:"`
		Bio         string `form:",textarea"`
	}

	// Test case 1: Test proper escaping of special characters in attributes
	specialCharsProfile := UserProfile{
		DisplayName: `Special " ' < > & Characters`,
		Bio:         "Normal text content",
	}

	output := string(RenderForm(&specialCharsProfile))

	// Verify the output has all special characters properly escaped
	if strings.Contains(output, `"Special "`) {
		t.Errorf("Quotation mark not properly escaped in attribute")
	}

	if !strings.Contains(output, `value="Special`) {
		t.Errorf("Expected escaped attribute value not found")
	}

	// Test case 2: Test proper escaping of HTML-like content
	htmlContentProfile := UserProfile{
		DisplayName: "Normal Name",
		Bio:         `<div>Sample HTML content</div>`,
	}

	output = string(RenderForm(&htmlContentProfile))

	// Verify the output has HTML content properly escaped
	if strings.Contains(output, `<div>`) {
		t.Errorf("HTML tags not properly escaped in content")
	}

	if !strings.Contains(output, `&lt;div&gt;`) {
		t.Errorf("Expected escaped HTML not found")
	}
}
