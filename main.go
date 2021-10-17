package main

/*
 * Route messages from HTTP POST to another PST URL or to MQTT.
 */
import (
	"encoding/json"
	"fmt"
	_ "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var exedir string

func handleNothing(w http.ResponseWriter, r *http.Request) {
}

func handleRootPost(w http.ResponseWriter, r *http.Request) {
	// Verify auth
	if Config.Username != "" && Config.Password != "" {
		u, p, ok := r.BasicAuth()
		if !ok || u != Config.Username || p != Config.Password {
			http.Error(w, "Invalid login", 401)
		}
	}

	// Get the body
	body, _ := ioutil.ReadAll(r.Body)

	go routeTraffic(r.URL.Path, string(body))
}

/*
 * api/routes Get a list of all configured inputs
 * Optional query args: (api/routes?in=inId&out=outId)
 * in/out: Show only config for one input/output.
 */
func handleApiRoutes(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(Config.Routes)

	/* For debug logCall("inUrl", "somewhere/else", "ftp", "OK"); */
	w.Write(js)
}

func main() {
	// Get prog dir and thereby html template dir
	exedir = filepath.Dir(os.Args[0])

	fmt.Printf("Exec dir %v\n", exedir)

	readConfig()
	r := mux.NewRouter()
	r.HandleFunc("/api/routes", handleApiRoutes).Methods("GET")
	r.HandleFunc("/{[x|]urlin}", handleRootPost).Methods("POST", "PUT")
	r.HandleFunc("/favicon.ico", handleNothing)
	http.Handle("/", r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(exedir + "/web/")))

	if Config.Address == "" {
		Config.Address = ":8222"
	}

	fmt.Printf("Serve address %v\n", Config.Address)

	e := http.ListenAndServe(Config.Address, r) // Blocking function
	if e != nil {
		fmt.Printf("Serve failed %v\n", e)
	}
}
