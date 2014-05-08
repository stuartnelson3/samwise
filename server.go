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
func main() {
	// load config file of servers
	// need to keep a heartbeat on servers and remove dead ones, then add them
	// to a queue to try to restart them.
	c := consistent.New()

	// load servers from config file into servers map
	// should this be a slice that values are removed from instead of a map??
	servers := map[string]bool{"server1": true, "server2": true, "server3": true}
	// confirm that the server are up??
	for server, _ := range servers {
		c.Add(server)
	}

	go MonitorServers(servers, c)

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
		server, err := c.Get(stream)
		if err != nil {
			return
		}
		http.RedirectHandler(server+"/stream?&token="+token+"stream="+stream, 302).ServeHTTP(w, r)
	})

	m.Post("/add_server", func(w http.ResponseWriter, r *http.Request) {
		server := r.FormValue("server")
		_, err := http.Get(server + "/pulse")
		if err != nil {
			servers[server] = false
			// write unsuccessful response
		} else {
			servers[server] = true
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

func MonitorServers(servers map[string]bool, c *consistent.Consistent) {
	for {
		for server, _ := range servers {
			_, err := http.Get(server + "/pulse")
			if err != nil {
				servers[server] = false
				// remove from consistent hashing
				c.Remove(server)
			} else {
				servers[server] = true
				c.Add(server)
				// add to consistent hashing if not already there
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func Restart(server string) {
	// execute code to restart the server
	// need to re-add the server to the servers slice
}
