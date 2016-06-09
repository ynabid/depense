package main

import (
	"encoding/json"
	//	"html/template"
	"container/list"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Depense struct {
	Id          int64
	Date        int64
	Description string
	Montant     float64
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
		d.Date,
		d.Description,
		d.Montant,
	)
}
func getDepenseList(from, to int64) ([]Depense, error) {
	var depenseList []Depense
	l := list.New()
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(
		"SELECT id,date,montant,description from depense WHERE date >= ? and date <= ? ORDER BY date",
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d Depense
		rows.Scan(&d.Id, &d.Date, &d.Montant, &d.Description)
		l.PushBack(d)
	}
	depenseList = make([]Depense, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		depenseList[i] = e.Value.(Depense)
		i++
	}
	defer rows.Close()
	return depenseList, nil
}

func read(id int64) (*Depense, error) {
	db, err := openMysqlDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(
		"SELECT date,montant,description from depense WHERE id=?",
		id,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rows.Next()
	var d Depense
	rows.Scan(&d.Date, &d.Montant, &d.Description)
	return &d, nil
}

func load(file string) ([]byte, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
func rootHandler(w http.ResponseWriter, r *http.Request) {
  body, err := load("template/index.html")
  if err != nil {
    http.NotFound(w, r)
  }
  fmt.Fprintf(w, "%s", body)
  //	renderTemplate(w, "edit", p)
}
*/
func parseDepense(values url.Values) (*Depense, error) {
	var depense Depense

	t, err := time.Parse("2006-01-01", values["date"][0])
	if err == nil {
		depense.Date = t.Unix()
	} else {
		return nil, err
	}

	depense.Description = values["description"][0]
	depense.Montant, err = strconv.ParseFloat(values["montant"][0], 64)
	if err != nil {
		return nil, err
	}
	if depense.Montant <= 0 {
		return nil, errors.New("Montant doit être >= 0")
	}
	return &depense, nil
}
func depenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		depense, err := parseDepense(r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result, err := depense.insert()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("OK"))
			w.Header().Add("Location", r.URL.Path+strconv.FormatInt(id, 10))
			w.WriteHeader(http.StatusCreated)
		}
	} else if r.Method == "GET" {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil && len(m) > 2 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		d, err := read(id)

		if err != nil {
			http.NotFound(w, r)
		}
		b, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func parseListDepense(values url.Values) (int64, int64, error) {
	from, err := time.Parse("2006-01-01", values["from"][0])
	if err != nil {
		return 0, 0, err
	}
	to, err := time.Parse("2006-01-01", values["to"][0])
	if err != nil {
		return from.Unix(), 0, err
	}
	return from.Unix(), to.Unix(), nil
}

func depenseListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		from, to, err := parseListDepense(r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		depenseList, err := getDepenseList(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			b, err := json.Marshal(depenseList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
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
var validPath = regexp.MustCompile("^/api/(depense|depenseList)/([a-zA-Z0-9]*)$")

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
	dir := "/home/yassine/Dev/godev/src/github.com/ynabid/depense/"
	http.Handle(
		"/depense/",
		http.StripPrefix(
			"/depense/",
			http.FileServer(http.Dir(dir+"res/")),
		),
	)
	http.HandleFunc(
		"/api/depense/list",
		depenseListHandler,
	)

	http.HandleFunc(
		"/api/depense/",
		makeHandler(depenseHandler),
	)

	http.ListenAndServe(":8080", nil)
}
