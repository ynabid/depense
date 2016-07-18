package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/ynabid/depense/db"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"time"
)

type DepenseString struct {
	Date        string
	Description string
	//Type        int
	Montant    float64
	CategoryId int64
	AccountId  int64
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
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
func DepenseDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		from, to, err := db.ParseMonth(r.Form.Get("month"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		all, err := db.DepenseData(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		b, err := json.Marshal(all)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)

	}
}

func AccountTRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		from, to, err := db.ParseMonth(r.Form.Get("month"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		accountsList, err := db.ReadAccountsTR(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		b, err := json.Marshal(accountsList)
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
	data, _ := httputil.DumpRequest(r, true)
	log.Println(string(data))
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {
			depense, err := parseDepense(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			depense.UserId, err = getUserId(r)
			if err != nil {
				http.Error(
					w,
					err.Error(),
					http.StatusUnauthorized,
				)
				return
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
func AccountListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		accountList, err := db.ReadAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			b, err := json.Marshal(accountList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}
}

func CategoryListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		categoryList, err := db.ReadAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			b, err := json.Marshal(categoryList)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}
}
func DepenseCatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		from, to, err := db.ParseMonth(r.Form.Get("month"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		depenseList, err := db.DepenseByCategory(from, to)
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

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		session, err := db.Login(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		uid := strconv.FormatInt(session.UserId, 10)
		http.SetCookie(
			w,
			&http.Cookie{
				Name:  "uid",
				Value: uid,
				Path:  "/",
				//	Expires: time.Unix(session.ExpiryDate, 0),
				MaxAge: int(session.ExpiryDate - time.Now().Unix()),
			},
		)
		http.SetCookie(
			w,
			&http.Cookie{
				Name:  "sessionid",
				Value: session.Id,
				Path:  "/",
				//	Expires: time.Unix(session.ExpiryDate, 0),
				MaxAge: int(session.ExpiryDate - time.Now().Unix()),
			},
		)
		w.WriteHeader(http.StatusOK)
	}
}
func getUserId(r *http.Request) (int64, error) {
	uidC, err := r.Cookie("uid")
	if err != nil {
		return -1, err
	}
	userId, err := strconv.ParseInt(uidC.Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return userId, nil
}

func MakeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := getUserId(r)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusUnauthorized,
			)
			return
		}
		sidC, err := r.Cookie("sessionid")
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusUnauthorized,
			)
			return

		}

		sid := sidC.Value
		_, err = db.IsAuthentified(userId, sid)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusUnauthorized,
			)
			return
		}
		CompresserFunc(fn, w, r)
	}
}

func CompresserFunc(fn func(w http.ResponseWriter, r *http.Request), w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	fn(gzr, r)
}

type CompresserHandler struct {
	Handler http.Handler
}

func (ch *CompresserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	CompresserFunc(ch.Handler.ServeHTTP, w, r)
}
func CompresserHandlerFunc(handler http.Handler) http.Handler {
	var ch CompresserHandler
	ch.Handler = handler
	return &ch
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
	depense.Montant = dStr.Montant
	//	depense.Type = dStr.Type
	depense.AccountId = dStr.AccountId
	depense.CategoryId = dStr.CategoryId
	return &depense, nil
}
