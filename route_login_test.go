package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// <test code>
	t.Run(fmt.Sprintf("Test empty body"), func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		w := httptest.NewRecorder()
		Login(w, req)
		if w.Result().StatusCode != 400 {
			t.Fatal("Failed TestLogin empty body")
		}
	})

	t.Run(fmt.Sprintf("Test body with wrong parameters"), func(t *testing.T) {
		params := []byte(`{"test1":"test"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "required")
	})

	t.Run(fmt.Sprintf("Test missing email"), func(t *testing.T) {
		params := []byte(`{"email":"", "password":"test"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "required")
	})

	t.Run(fmt.Sprintf("Test missing password"), func(t *testing.T) {
		params := []byte(`{"email":"test@test.com", "password":""}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "required")
	})

	t.Run(fmt.Sprintf("Test invalid email"), func(t *testing.T) {
		params := []byte(`{"email":"testtest.com", "password":""}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "email")
	})

	t.Run(fmt.Sprintf("Test too short password"), func(t *testing.T) {
		params := []byte(`{"email":"test@test.com", "password":"test"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "min")
	})

	t.Run(fmt.Sprintf("Test missing special char password"), func(t *testing.T) {
		params := []byte(`{"email":"test@test.com", "password":"testtesttest"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "containsany")
	})

	t.Run(fmt.Sprintf("Test missing numbers password"), func(t *testing.T) {
		params := []byte(`{"email":"test@test.com", "password":"!testtesttest"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "containsany")
	})

	t.Run(fmt.Sprintf("Test email not found"), func(t *testing.T) {
		params := []byte(`{"email":"thisisjustatesttest@test.com", "password":"!1testtesttest"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, string(resp), "User not found")
	})

	t.Run(fmt.Sprintf("Test unsuccessfull login"), func(t *testing.T) {
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
		params := []byte(`{"email":"t@t.t", "password":"!1wrongpassword"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		if w.Result().StatusCode != 403 {
			t.Fatal("Failed test unsuccessfull login")
		}
	})

	t.Run(fmt.Sprintf("Test successfull login"), func(t *testing.T) {
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
		params := []byte(`{"email":"x@x.x", "password":"!1abcabcabc"}`)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(params))
		w := httptest.NewRecorder()
		Login(w, req)
		temp_resp := w.Result()
		body, _ := io.ReadAll(temp_resp.Body)
		test_ok_login_user := struct {
			Code string
		}{}
		json.Unmarshal([]byte(body), &test_ok_login_user)
		assert.Equal(t, test_ok_login_user.Code, test_ok_login_user_code)
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
