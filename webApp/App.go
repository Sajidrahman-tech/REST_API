package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"tawesoft.co.uk/go/dialog"
)

type Employee struct {
	Name     string
	Password string
	Id       int
}

func dbConn() (db *sql.DB) {

	db, err := sql.Open("mysql", "root:Qazxsw#2!@tcp(localhost:3306)/users")
	if err != nil {
		fmt.Println("error validating sql.opem arguments")
		panic(err.Error())
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.ping")
		panic(err.Error())
	} else {
		fmt.Println("sucess")
	}
	return db

}

var t = template.Must(template.ParseGlob("*.html"))

//var t = template.Must(template.New("html-tmpl").ParseGlob("*.tmpl"))

func login(w http.ResponseWriter, r *http.Request) {

	t.ExecuteTemplate(w, "login.html", nil)
}

/////////////-POST user
func postUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside postUser")
	Username := r.FormValue("username")
	Password := r.FormValue("password")

	db := dbConn()
	fmt.Println("username " + Username)
	fmt.Println("password" + Password)

	fmt.Println("inside postUser after dbconn")

	///validation using input
	if Username == "" {
		dialog.Alert("EMpty Username")
		http.Redirect(w, r, "/login", 301)
		return
	}
	if len(Password) < 8 {
		dialog.Alert("Minimum Length should be 8")
		http.Redirect(w, r, "/login", 301)
		return
	}

	fmt.Println(len(Password))

	insert, err := db.Prepare("INSERT INTO `users`.`emp` VALUES (?,?,NULL)")

	if err != nil {
		fmt.Println("inside PostError")
		panic(err.Error())
	}

	insert.Exec(Username, Password)
	defer insert.Close()
	defer db.Close()

	dialog.Alert("Posted USER")
	http.Redirect(w, r, "/login", 301)

}

//////////->GET user
func getUser(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("index.html")
	fmt.Println("inside getUser")
	db := dbConn()
	fmt.Println("inside getuser after dbconn")
	nId := r.FormValue("username")
	selDB, err := db.Query("SELECT * FROM emp WHERE name=?", nId)
	if err != nil {
		panic(err.Error())
	}
	employ := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var name string
		var password string
		var id int
		err = selDB.Scan(&name, &password, &id)
		if err != nil {
			panic(err.Error())
		}
		employ.Name = name
		employ.Password = password
		employ.Id = id
		res = append(res, employ)

	}
	if len(employ.Name) == 0 {
		fmt.Fprintf(w, "no details")
		return
	}

	fmt.Println(res)
	fmt.Println("hi")
	//t.ExecuteTemplate(w, "index.html", res)
	errr := parsedTemplate.Execute(w, res)
	if errr != nil {
		fmt.Println("inside the p error")
		panic(errr.Error())
	}
	defer selDB.Close()
	defer db.Close()
}

/////////->DELETE user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside deleteUser")
	db := dbConn()

	fmt.Println("inside deleteuser after dbconn")
	Username := r.FormValue("username")
	fmt.Println(Username)
	delete, err := db.Prepare("DELETE FROM emp WHERE name= ?")
	if err != nil {
		panic(err.Error())
	}
	delete.Exec(Username)

	dialog.Alert("DELETED USER")
	http.Redirect(w, r, "/login", 301)
	defer db.Close()

}

////->UPDATE user
func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside deleteUser")
	db := dbConn()
	fmt.Println("inside deleteuser after dbconn")
	Username := r.FormValue("username")
	Password := r.FormValue("password")

	fmt.Println(Username)
	fmt.Println(Password)
	update, err := db.Prepare("UPDATE emp SET password=? WHERE name=?")
	if err != nil {
		panic(err.Error())
	}
	update.Exec(Password, Username)
	dialog.Alert("Updated USER")
	defer db.Close()
	http.Redirect(w, r, "/login", 301)

}

///////////->SHOW USER
func showUser(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("index.html")
	fmt.Println("inside showUser")
	db := dbConn()
	fmt.Println("inside showuser after dbconn")

	selDB, err := db.Query("SELECT * FROM emp ")
	if err != nil {
		panic(err.Error())
	}
	employ := Employee{}

	res := []Employee{}
	for selDB.Next() {
		var name string
		var password string
		var id int
		//selDB.Scan(&\)
		err = selDB.Scan(&name, &password, &id)
		if err != nil {
			panic(err.Error())
		}
		employ.Name = name
		employ.Password = password
		employ.Id = id
		res = append(res, employ)

	}
	if len(employ.Name) == 0 {
		fmt.Fprintf(w, "no details")
		return
	}

	fmt.Println(res)
	fmt.Println("hi")

	errr := parsedTemplate.Execute(w, res)
	if errr != nil {
		fmt.Println("inside the p error")
		panic(errr.Error())
	}
	defer selDB.Close()
	defer db.Close()
}

////////->VALIDATE USER

func validate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside validateUser")
	db := dbConn()
	fmt.Println("inside validateuser after dbconn")
	Username := r.FormValue("username")
	Password := r.FormValue("password")
	valDb, err := db.Prepare("SELECT * FROM emp WHERE name=? AND password=?")

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("username is:" + Username)
	fmt.Println("password is:" + Password)
	//valDb.Exec(Username, Password)
	fmt.Println(valDb)
	employ := Employee{}
	var name string
	var password string
	var id int
	errr := valDb.QueryRow(Username, Password).Scan(&name, &password, &id)
	if errr != nil {
		if errr == sql.ErrNoRows {

		}

	}

	employ.Name = name
	employ.Password = password
	employ.Id = id
	fmt.Println("Inside errr")
	fmt.Println("the data name is" + employ.Name)
	fmt.Println("the data password is" + employ.Password)

	if name == "" {
		fmt.Println("INVALID USER")
		dialog.Alert("INVALID USER")
		defer valDb.Close()
		defer db.Close()
		http.Redirect(w, r, "/login", 301)
		return
	}
	dialog.Alert("VALID USER")

	http.Redirect(w, r, "/login", 301)

	defer valDb.Close()
	defer db.Close()
}

///////////
func novalidate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside novalidateUser")
	db := dbConn()
	fmt.Println("inside novalidateuser after dbconn")
	Username := r.FormValue("username")
	Password := r.FormValue("password")
	valDb, err := db.Prepare("SELECT * FROM emp WHERE name='" + Username + "' AND password='" + Password + "'")

	if err != nil {
		// 	panic(err.Error())
	}
	fmt.Println("username is:" + Username)
	fmt.Println("password is:" + Password)

	fmt.Println(valDb)
	employ := Employee{}
	var name string
	var password string
	var id int
	valDb.QueryRow().Scan(&name, &password, &id)

	employ.Name = name
	employ.Password = password
	employ.Id = id
	fmt.Println("Inside errr")
	fmt.Println("the data name is" + employ.Name)
	fmt.Println("the data password is" + employ.Password)

	if name == "" {
		fmt.Println("INVALID USER")
		dialog.Alert("INVALID USER")
		defer valDb.Close()
		defer db.Close()
		http.Redirect(w, r, "/login", 301)
		return
	}
	dialog.Alert("VALID USER")

	http.Redirect(w, r, "/login", 301)

	defer valDb.Close()
	defer db.Close()
}
func main() {

	r := mux.NewRouter()
	//dbConn()
	r.HandleFunc("/postUser", postUser).Methods("POST")
	r.HandleFunc("/getUser", getUser).Methods("POST")
	r.HandleFunc("/deleteUser", deleteUser).Methods("POST")
	r.HandleFunc("/updateUser", updateUser).Methods("POST")
	r.HandleFunc("/showUser", showUser).Methods("POST")
	r.HandleFunc("/validate", validate).Methods(("POST"))
	r.HandleFunc("/novalidate", novalidate).Methods(("POST"))
	r.HandleFunc("/login", login)

	errr := http.ListenAndServe(":9090", r) // setting listening port
	if errr != nil {
		log.Fatal("ListenAndServe: ", errr)
	}
}
