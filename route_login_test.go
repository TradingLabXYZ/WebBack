package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {

	// Test empty body
	req := httptest.NewRequest("POST", "/login", nil)
	w := httptest.NewRecorder()
	Login(w, req)
	if w.Result().StatusCode != 400 {
		t.Fatal("Failed TestLogin empty body")
	}

	// Test body with wrong parameters
	params := []byte(`{"test1":"test"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res := w.Result()
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "required")

	// Test missing email
	params = []byte(`{"email":"", "password":"test"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "required")

	// Test missing password
	params = []byte(`{"email":"test@test.com", "password":""}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "required")

	// Test unvalid email
	params = []byte(`{"email":"testtest.com", "password":""}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "email")

	// Test too short password
	params = []byte(`{"email":"test@test.com", "password":"test"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "min")

	// Test missing special char password
	params = []byte(`{"email":"test@test.com", "password":"testtesttest"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "containsany")

	// Test missing numbers password
	params = []byte(`{"email":"test@test.com", "password":"!testtesttest"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "containsany")

	// Test email not found
	params = []byte(`{"email":"thisisjustatesttest@test.com", "password":"!1testtesttest"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	res = w.Result()
	defer res.Body.Close()
	resp, _ = ioutil.ReadAll(res.Body)
	assert.Equal(t, string(resp), "User not found")

	// Test unsuccessfull login
	test_ko_login_pass, _ := Encrypt("!1cdcdcdcd")
	test_ko_login_user_code := "DEFTGT"
	Db.Exec(`
		INSERT INTO users (
			code,email,username,password,
			privacy,plan,createdat,updatedat)
		VALUES (
			$1,'t@t.t','t',$2,'all','basic',
			current_timestamp,current_timestamp);`,
		test_ko_login_user_code, test_ko_login_pass)
	params = []byte(`{"email":"t@t.t", "password":"!1wrongpassword"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	if w.Result().StatusCode != 403 {
		t.Fatal("Failed test unsuccessfull login")
	}

	// Test successfull login
	test_ok_login_pass, _ := Encrypt("!1abcabcabc")
	test_ok_login_user_code := "ABCDEFT"
	Db.Exec(`
		INSERT INTO users (
			code,email,username,password,
			privacy,plan,createdat,updatedat)
		VALUES (
			$1,'x@x.x','x',$2,'all','basic',
			current_timestamp,current_timestamp);`,
		test_ok_login_user_code, test_ok_login_pass)
	params = []byte(`{"email":"x@x.x", "password":"!1abcabcabc"}`)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
	w = httptest.NewRecorder()
	Login(w, req)
	temp_resp := w.Result()
	body, _ := io.ReadAll(temp_resp.Body)
	test_ok_login_user := struct {
		Code string
	}{}
	json.Unmarshal([]byte(body), &test_ok_login_user)
	assert.Equal(t, test_ok_login_user.Code, test_ok_login_user_code)
}
