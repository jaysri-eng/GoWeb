package main

import (
	"database/sql"
	"fmt"

	// "html/template"
	"log"
	"net/http"

	// "os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/gorilla/mux"
	// "golang.org/x/crypto/bcrypt"
)

var dbb *sql.DB

type UserDec struct {
	Id       int64
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

var userss = []UserDec{
	{Id: 1, Username: "John Coltrane", Password: "summa"},
	{Id: 2, Username: "Gerry Mulligan", Password: "animal"},
}

func all(c *gin.Context) {
	allusers, err := getUser("jaya")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("data: %v\n", allusers)
	c.IndentedJSON(http.StatusOK, allusers)
}
func one(c *gin.Context) {
	oneuser, err := getOneUser(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found: %v\n", oneuser)
	c.IndentedJSON(http.StatusOK, oneuser)
}
func add(c *gin.Context) {
	var newuser User
	if err := c.BindJSON(&newuser); err != nil {
		return
	}
	userss = append(userss, UserDec(newuser))
	c.IndentedJSON(http.StatusCreated, newuser)
}

// func main() {
// 	router := gin.Default()
// 	router.GET("/albums", all)
// 	router.POST("/albums", add)
// 	router.GET("/albums/:id", one)
// 	router.Run("localhost:80")
// }
