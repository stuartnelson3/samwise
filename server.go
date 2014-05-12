package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/cors"
	"github.com/stathat/consistent"
	"net/http"
	"os"
	"time"
)

// write tests you lazy slob
var c = consistent.New()

func main() {
	// load config file of servers
	// need to keep a heartbeat on servers and remove dead ones, then add them
	// to a queue to try to restart them.

	// load servers from config file into servers map
	// should this be a slice that values are removed from instead of a map??

	// servers := map[string]bool{"server1": true, "server2": true, "server3": true}
	servers := map[string]bool{"http://localhost:9000": true}
	// confirm that the server are up??
	for server, _ := range servers {
		c.Add(server)
	}

	go MonitorServers(servers)

	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"http://*", "https://*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin"},
	}))

	m.Get("/stream", Stream)

	m.Post("/add_server", func(w http.ResponseWriter, r *http.Request) {
		server := r.FormValue("server")
		_, err := http.Get(server)
		if err != nil {
			servers[server] = false
			w.WriteHeader(404)
			// write unsuccessful response
		} else {
			servers[server] = true
			w.WriteHeader(200)
			// make sure this isn't a duplicate!!
			c.Add(server)
			// write successful response
		}
	})

	m.Post("/update_stream", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		stream := r.FormValue("stream")
		if stream == "" {
			return
		}

		server, err := c.Get(stream)
		if err != nil {
			return
		}

		http.PostForm(server+"/update_stream", r.PostForm)
	})
	m.Run()
}

func Stream(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token != os.Getenv("TOKEN") {
		w.WriteHeader(404)
		return
	}

	stream := r.FormValue("stream")
	if stream == "" {
		w.WriteHeader(404)
		return
	}
	// redirect to server based on consistent hash

	// alternatively, have rails query this endpoint to get the route, then
	// have rails return the found server to the client
	server, err := c.Get(stream)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	http.RedirectHandler(server+"/stream?&token="+token+"stream="+stream, 302).ServeHTTP(w, r)
}

func MonitorServers(servers map[string]bool) {
	for {
		monitor(servers)
		time.Sleep(5 * time.Minute)
	}
}

func monitor(servers map[string]bool) {
	for server, _ := range servers {
		resp, err := http.Get(server)
		if err != nil || resp.StatusCode != 200 {
			servers[server] = false
			c.Remove(server)
		} else {
			servers[server] = true
			if !inHash(server) {
				c.Add(server)
			}
		}
	}
}

func inHash(server string) bool {
	for _, s := range c.Members() {
		if s == server {
			return true
		}
	}
	return false
}

func Restart(server string) {
	// execute code to restart the server
	// need to re-add the server to the servers slice
}
