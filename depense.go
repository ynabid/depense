package main

import (
	"github.com/ynabid/depense/db"
	"github.com/ynabid/depense/handlers"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/api/(depense|depenseList)/([a-zA-Z0-9]*)$")

func main() {
	db.OpenMysqlDB()
	dir := "/var/www/html/depense/"
	http.HandleFunc(
		"/api/depenseList",
		handlers.MakeHandler(handlers.DepenseListHandler),
	)

	http.HandleFunc(
		"/api/depense",
		handlers.MakeHandler(handlers.DepenseHandler),
	)
	http.HandleFunc(
		"/api/depense/month",
		handlers.MakeHandler(handlers.DepenseMonthHandler),
	)
	http.HandleFunc(
		"/api/auth",
		handlers.AuthHandler,
	)

	http.HandleFunc(
		"/api/depense/account/list",
		handlers.MakeHandler(handlers.AccountListHandler),
	)
	http.HandleFunc(
		"/api/depense/category/list",
		handlers.MakeHandler(handlers.CategoryListHandler),
	)
	http.HandleFunc(
		"/api/depense/bycategory",
		handlers.MakeHandler(handlers.DepenseCatHandler),
	)
	http.HandleFunc(
		"/api/depense/accountTR",
		handlers.MakeHandler(handlers.AccountTRHandler),
	)
	http.HandleFunc(
		"/api/depense/all",
		handlers.MakeHandler(handlers.DepenseDataHandler),
	)

	http.Handle(
		"/depense/",
		http.StripPrefix(
			"/depense/",
			handlers.CompresserHandlerFunc(http.FileServer(http.Dir(dir+"res/"))),
		),
	)

	http.ListenAndServe(":8080", nil)
}
