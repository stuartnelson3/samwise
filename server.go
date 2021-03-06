package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/cors"
	"github.com/stathat/consistent"
	"net/http"
	"os"
	"time"
)

// write tests you lazy slob
var c = consistent.New()
var servers = map[string]bool{"http://localhost:9000": false, "http://localhost:9090": false}

// load config file of servers
// need to keep a heartbeat on servers and remove dead ones, then add them
// to a queue to try to restart them.

// load servers from config file into servers map
// should this be a slice that values are removed from instead of a map??

// confirm that the server are up??

// TODO: convert to gorilla mux
func main() {
	for server, _ := range servers {
		c.Add(server)
	}

	go MonitorServers(servers)

	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"http://*", "https://*"},
		AllowHeaders: []string{"Origin"},
	}))

	m.Get("/", heartBeat)
	m.Get("/stream", stream)
	m.Post("/add_server", addServer)
	m.Post("/update_stream", updateStream)
	m.Run()
}

func heartBeat(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == os.Getenv("TOKEN") {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

func stream(w http.ResponseWriter, r *http.Request) {
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

	server, err := c.Get(stream)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// what happens if person 1 gets sent to site A, but then site B comes back
	// online. Person 1, if they reconnected, would be directed to site B. All
	// new users are going to site B. Should I send a disconnect message to the
	// server that it needs to dump its users for hash X?
	http.RedirectHandler(server+"/stream?&token="+token+"&stream="+stream, 302).ServeHTTP(w, r)
}

func addServer(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") != os.Getenv("TOKEN") {
		w.WriteHeader(404)
		return
	}
	server := r.FormValue("server")
	if server == "" {
		w.WriteHeader(404)
		return
	}
	_, err := http.Get(server)
	if err != nil {
		servers[server] = false
		w.WriteHeader(404)
	} else {
		servers[server] = true
		if !inHash(server) {
			c.Add(server)
		}
		w.WriteHeader(200)
	}
}

func updateStream(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	stream := r.FormValue("stream")
	if stream == "" {
		w.WriteHeader(404)
		return
	}

	server, err := c.Get(stream)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	http.PostForm(server+"/update_stream", r.PostForm)
}

func MonitorServers(servers map[string]bool) {
	for {
		monitor(servers)
		time.Sleep(5 * time.Second)
	}
}

func monitor(servers map[string]bool) {
	for server, _ := range servers {
		resp, err := checkServer(server)
		if err != nil || resp.StatusCode != 200 {
			servers[server] = false
			if inHash(server) {
				c.Remove(server)
				fmt.Println("Removed from hash: ", server)
			}
		} else {
			servers[server] = true
			if !inHash(server) {
				c.Add(server)
				fmt.Println("Added to hash: ", server)
			}
		}
	}
}

func checkServer(server string) (*http.Response, error) {
	return http.Get(server + "?token=" + os.Getenv("TOKEN"))
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
