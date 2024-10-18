package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/dem/show", app.show)
	//mux.HandleFunc("/dem/create", app.create)

	mux.HandleFunc("/dem/read", app.ExportFromProfstroi)

	//static
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
