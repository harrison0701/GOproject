package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	registerRoutes(r)
	return r
}

func TestPingHandler(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestCreateUserHandler(t *testing.T) {
	r := setupRouter()

	user := map[string]string{"name": "testuser"}
	jsonValue, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expectedResponse := `{"message":"User created","name":"testuser"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetUserFromDBHandler(t *testing.T) {
	// 初始化資料庫連接
	var err error
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// 插入測試數據
	_, err = db.Exec("INSERT INTO users (name) VALUES ($1)", "testuser")
	if err != nil {
		t.Fatal(err)
	}

	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dbuser/testuser", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	expectedResponse := `{"message":"User found","name":"testuser"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	// 清理測試數據
	_, err = db.Exec("DELETE FROM users WHERE name = $1", "testuser")
	if err != nil {
		t.Fatal(err)
	}
}
