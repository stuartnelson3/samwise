package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/cors"
	"net/http"
	"os"
)

// add consistent hashing
// https://github.com/stathat/consistent

func main() {
	// load config file of servers
	// expose endpoint to add new servers??

	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"http://*", "https://*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin"},
	}))

	m.Get("/stream", func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		if token != os.Getenv("TOKEN") {
			return
		}

		stream := r.FormValue("stream")
		if stream == "" {
			return
		}
		// redirect to server based on consistent hash

		// alternatively, have rails query this endpoint to get the route, then
		// have rails return the found server to the client
		redirect := "http://localhost:9000/stream?stream=" + stream + "&token=" + token
		http.RedirectHandler(redirect, 302).ServeHTTP(w, r)
	})

	m.Post("/update_stream", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// hash r.FormValue('stream') to get server
		http.PostForm("http://localhost:9000/update_stream", r.PostForm)
	})
	m.Run()
}
