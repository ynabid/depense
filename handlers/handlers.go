package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ynabid/depense/db"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	"regexp"
	"strconv"
	"time"
)

type DepenseString struct {
	Date        string
	Description string
	Montant     string
}

var validPath = regexp.MustCompile("^/api/(depense|depenseList)/([a-zA-Z0-9]*)$")

func DepenseMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		f, err := db.DepenseMonth(r.Form.Get("month"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
func DepenseHandler(w http.ResponseWriter, r *http.Request) {
	/*	data, _ := httputil.DumpRequest(r, true)
		log.Println(string(data))
	*/
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {
			depense, err := parseDepense(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if depense != nil {
				log.Println(depense)

			}
			result, err := depense.Insert()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				id, err := result.LastInsertId()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Add("Location", r.URL.Path+strconv.FormatInt(id, 10))
				w.WriteHeader(http.StatusCreated)
			}

		} else {
			http.Error(w, errors.New(r.Header.Get("Content-Type")+" not maintained").Error(), http.StatusBadRequest)
			return
		}
	} else if r.Method == "GET" {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil || len(m) < 3 {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		d, err := db.Read(id)

		if err != nil {
			http.NotFound(w, r)
			return
		}
		b, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
func DepenseListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		from, to, err := db.ParseMonth(r.Form.Get("month"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		depenseList, err := db.DepenseList(from, to)
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

func MakeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}
func parseDepense(body io.ReadCloser) (*db.Depense, error) {
	var depense db.Depense
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var dStr DepenseString
	json.Unmarshal(b, &dStr)
	t, err := time.Parse("2006-01-02", dStr.Date)

	if err == nil {
		depense.Date = t.Unix()
	} else {
		return nil, err
	}

	depense.Description = dStr.Description
	depense.Montant, err = strconv.ParseFloat(dStr.Montant, 64)
	if err != nil {
		return nil, err
	}
	if depense.Montant <= 0 {
		return nil, errors.New("Montant doit être >= 0")
	}
	return &depense, nil
}
