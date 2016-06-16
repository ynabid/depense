package main

import (
	"github.com/ynabid/depense/handlers"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/api/(depense|depenseList)/([a-zA-Z0-9]*)$")

func main() {
	dir := "/var/www/html/depense/"
	http.HandleFunc(
		"/api/depenseList",
		handlers.DepenseListHandler,
	)

	http.HandleFunc(
		"/api/depense/",
		handlers.DepenseHandler,
	)
	http.HandleFunc(
		"/api/depense/month",
		handlers.DepenseMonthHandler,
	)
	http.Handle(
		"/depense/",
		http.StripPrefix(
			"/depense/",
			http.FileServer(http.Dir(dir+"res/")),
		),
	)

	http.ListenAndServe(":8080", nil)
}
