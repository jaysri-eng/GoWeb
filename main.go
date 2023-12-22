package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	// "os"
	// "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	Id       int64
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	// cfg := mysql.Config{
	// 	User:   os.Getenv("root"),
	// 	Passwd: os.Getenv("jayanthsql"),
	// 	Net:    "tcp",
	// 	Addr:   "localhost:3306",
	// 	DBName: "go",
	// }
	var err error
	db, err = sql.Open("mysql", "root:jayanthsql@tcp(localhost:3306)/go")
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/homepage.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, "")
	})
	router.HandleFunc("/allUsers", func(w http.ResponseWriter, r *http.Request) {
		users, err := getUser("jaya")
		if err != nil {
			log.Fatal(err)
		}
		tmpl, err := template.ParseFiles("template/allUsers.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		user, err := allUsers()
		tmpl, err := template.ParseFiles("template/oneUser.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/one", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/oneUser.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")
		fmt.Printf("Received ID: %v\n", id)
		newId, err := strconv.ParseInt(id, 10, 64)
		fmt.Printf("Received ID: %v\n", newId)
		user, err := getOneUser(newId)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/oneUser", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/oneUserTemplate.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/postUser", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/postUser.html")
		username := r.FormValue("username")
		password := r.FormValue("password")
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		newUser, err := addUser(User{
			Username: username,
			Password: string(bytes),
		})
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	http.ListenAndServe(":80", router)
}

func allUsers() ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM usersList")
	if err != nil {
		return nil, fmt.Errorf("users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("users :%v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("users: %v", err)
	}
	return users, nil
}

func getUser(username string) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM usersList WHERE username = ?", username)
	if err != nil {
		return nil, fmt.Errorf("getUser %q: %v", username, err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("getUser %q:%v", username, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getUser %q:%v", username, err)
	}
	return users, nil
}

func getOneUser(id int64) (User, error) {
	var user User
	row := db.QueryRow("SELECT * FROM usersList WHERE id=?", id)
	if err := row.Scan(&user.Id, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("getOneUser %d: no user", id)
		}
		return user, fmt.Errorf("getOneUser %d:%v", id, err)
	}
	return user, nil
}

func addUser(user User) (int64, error) {
	result, err := db.Exec("INSERT INTO usersList(username, passw) VALUES (?,?)", user.Username, user.Password)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}
