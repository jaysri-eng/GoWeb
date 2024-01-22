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
	"github.com/gorilla/securecookie"
	// "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

type User struct {
	Id       int64
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Cart struct {
	Id     int64
	Item   string `form:"item" binding:"required"`
	Price  int64  `form:"price" binding:"required"`
	UserId int64  `form:"userId" binding:"required"`
}

func main() {
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	// cfg := mysql.Config{
	// 	User:   os.Getenv("root"),
	// 	Passwd: os.Getenv(""),
	// 	Net:    "tcp",
	// 	Addr:   "localhost:3306",
	// 	DBName: "go",
	// }
	var err error
	db, err = sql.Open("mysql", "root:jayanthsql@tcp(localhost:3306)/Cart")
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")
	router := mux.NewRouter()
	router.HandleFunc("/loginHandler", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, "")
	})
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// tmpl, err := template.ParseFiles("template/login.html")
		username := r.FormValue("username")
		password := r.FormValue("password")
		redirectTarget := "/"
		dbPwd := "jayanth"
		dbUser := "jayanth"
		if username == dbUser && password == dbPwd {
			SetCookie(username, w)
			redirectTarget = "/home"
			fmt.Fprintln(w, "Login successfull!")
		} else {
			fmt.Fprintln(w, "Login failed!")
		}
		http.Redirect(w, r, redirectTarget, 80)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// err = tmpl.Execute(w, "")
	})
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		ClearCookie(w)
		fmt.Fprintln(w, "Logout successful")
		http.Redirect(w, r, "/loginHandler", 80)
	})
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/signup.html")
		username := r.FormValue("username")
		password := r.FormValue("password")
		// bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		newUser, err := addUser(User{
			Username: username,
			Password: string(password),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, newUser)
	})
	router.HandleFunc("/forgotPassword", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/forgotPassword.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, "")
	})
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/main.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, "")
	})
	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		cart, err := allItems()
		tmpl, err := template.ParseFiles("template/allItems.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/one", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/oneItem.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := r.FormValue("id")
		fmt.Printf("Received ID: %v\n", id)
		newId, err := strconv.ParseInt(id, 10, 64)
		fmt.Printf("Received ID: %v\n", newId)
		item, err := getOneItem(newId)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandleFunc("/oneItem", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/oneItemTemplate.html")
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
	router.HandleFunc("/addItems", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("template/addItems.html")
		item := r.FormValue("item")
		price := r.FormValue("price")
		newPrice, err := strconv.ParseInt(price, 10, 64)
		userId := r.FormValue("userid")
		newId, err := strconv.ParseInt(userId, 10, 64)
		// fmt.Printf("Received ID: %v\n", newId)
		// bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		newItem, err := addItems(Cart{
			Item:   item,
			Price:  newPrice,
			UserId: newId,
		})
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, newItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	http.ListenAndServe(":80", router)
}

func allUsers() ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM users")
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
func allItems() ([]Cart, error) {
	var items []Cart
	rows, err := db.Query("SELECT * FROM cart")
	if err != nil {
		return nil, fmt.Errorf("cart: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cart Cart
		if err := rows.Scan(&cart.Id, &cart.Item, &cart.Price, &cart.UserId); err != nil {
			return nil, fmt.Errorf("cart :%v", err)
		}
		items = append(items, cart)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cart: %v", err)
	}
	return items, nil
}

func getUser(username string) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM users WHERE username = ?", username)
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
	row := db.QueryRow("SELECT * FROM users WHERE id=?", id)
	if err := row.Scan(&user.Id, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("getOneUser %d: no user", id)
		}
		return user, fmt.Errorf("getOneUser %d:%v", id, err)
	}
	return user, nil
}
func getOneItem(id int64) (Cart, error) {
	var cart Cart
	row := db.QueryRow("SELECT * FROM cart WHERE id=?", id)
	if err := row.Scan(&cart.Id, &cart.Item, &cart.Price, &cart.UserId); err != nil {
		if err == sql.ErrNoRows {
			return cart, fmt.Errorf("getOneItem %d: no item", id)
		}
		return cart, fmt.Errorf("getOneItem %d:%v", id, err)
	}
	return cart, nil
}

func addUser(user User) (int64, error) {
	result, err := db.Exec("INSERT INTO users(username, password) VALUES (?,?)", user.Username, user.Password)
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addUser: %v", err)
	}
	return id, nil
}
func addItems(cart Cart) (int64, error) {
	result, err := db.Exec("INSERT INTO cart(item, price, user_id) VALUES (?,?,?)", cart.Item, cart.Price, cart.UserId)
	if err != nil {
		return 0, fmt.Errorf("addItems: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addItems: %v", err)
	}
	return id, nil
}

func SetCookie(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func ClearCookie(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// for POST
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	ClearCookie(response)
	http.Redirect(response, request, "/", 80)
}
