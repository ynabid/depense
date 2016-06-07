package main

import (
	"fmt"
	//	"html/template"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Depense struct {
	date        int64
	description string
	montant     float64
}

func openMysqlDB() (*sql.DB, error) {
	return sql.Open("mysql", "root:3811469@tcp(localhost:3306)/mm")
}
func (d *Depense) insert() (sql.Result, error) {
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	return db.Exec(
		"INSERT INTO depense (date,description,montant) VALUES (?,?,?)",
		d.date,
		d.description,
		d.montant,
	)
}
func read(id int64) (Depense, error) {
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.query(
		"SELECT date,montant,description from depense WHERE id=?",
		id,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	var d Depense
	rows.scan(&d.date, &d.montant, &d.description)
	return d
}

func load(file string) ([]byte, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := load("template/index.html")
	if err != nil {
		http.NotFound(w, r)
	}
	fmt.Fprintf(w, "%s", body)
	//	renderTemplate(w, "edit", p)
}
func depenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var depense Depense

		t, err := time.Parse("2006-01-01", r.FormValue("date"))
		if err == nil {
			depense.date = t.Unix()
		}
		depense.description = r.FormValue("description")
		depense.montant, err = strconv.ParseFloat(r.FormValue("montant"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if depense.montant <= 0 {
			http.Error(w, "Montant doit être >= 0", http.StatusBadRequest)
		}
		result, err := depense.insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Location", "/depense/"+string(id))
			w.WriteHeader(http.StatusCreated)
		}
	} else if r.Method == "GET" {

	}
}

/*func renderTemplate(w http.ResponseWriter, file string) {
	t, err := template.ParseFiles("template/" + file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
*/
var validPath = regexp.MustCompile("^/(depense)/[a-zA-Z0-9]*")

func makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func main() {
	dir := "/home/yassine/Dev/godev/src/github.com/ynabid/mesrouf/"
	http.Handle(
		"/mesrouf/res/",
		http.StripPrefix(
			"/mesrouf/res/",
			http.FileServer(http.Dir(dir+"res/")),
		),
	)
	http.HandleFunc(
		"/depense/",
		makeHandler(depenseHandler),
	)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Test")
}
