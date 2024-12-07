package main

import (
	_ "GoProject/docs" // 自动生成的文档
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for Swagger with Gin.
// @host localhost:8080
// @BasePath /

var db *sql.DB

func main() {
	var err error
	// 連接 PostgreSQL 資料庫
	connStr := "user=postgres dbname=my_database sslmode=disable password=mysecretpassword"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	// 注册 Swagger 路由
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 註冊 API 路由
	registerRoutes(r)

	r.Run(":8080")
}

// registerRoutes 註冊所有 API 路由
func registerRoutes(r *gin.Engine) {
	r.GET("/ping", pingHandler)
	r.POST("/user", createUserHandler)
	r.GET("/user/:name", getUserHandler)
	r.GET("/dbuser/:name", getUserFromDBHandler)
}

// pingHandler 處理 /ping 路由
// @Summary Ping API
// @Description Test the API
// @Success 200 {string} string "pong"
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// createUserHandler 處理 /user 路由
// @Summary Create User
// @Description Create a new user
// @Accept  json
// @Produce  json
// @Param   name     body    string     true  "User Name"
// @Success 200 {object} map[string]interface{}
// @Router /user [post]
func createUserHandler(c *gin.Context) {
	var json struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "User created",
		"name":    json.Name,
	})
}

// getUserHandler 處理 /user/:name 路由
// @Summary Get User
// @Description Get user by name
// @Produce  json
// @Param   name     path    string     true  "User Name"
// @Success 200 {object} map[string]interface{}
// @Router /user/{name} [get]
func getUserHandler(c *gin.Context) {
	name := c.Param("name")
	c.JSON(200, gin.H{
		"message": "User found",
		"name":    name,
	})
}

// getUserFromDBHandler 從資料庫中取用戶資料
// @Summary Get User from DB
// @Description Get user by name from PostgreSQL database
// @Produce  json
// @Param   name     path    string     true  "User Name"
// @Success 200 {object} map[string]interface{}
// @Router /dbuser/{name} [get]
func getUserFromDBHandler(c *gin.Context) {
	name := c.Param("name")
	var userName string
	err := db.QueryRow("SELECT name FROM users WHERE name = $1", name).Scan(&userName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(200, gin.H{
		"message": "User found",
		"name":    userName,
	})
}
