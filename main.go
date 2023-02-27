package main

/*
 * Route messages from HTTP POST to another PST URL or to MQTT.
 */
import (
	"fmt"
	_ "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var exedir string

func handleNothing(w http.ResponseWriter, r *http.Request) {
}

func measureHttpReponseTime(start time.Time, path string) {
	if MeasureTime {
		elapsed := time.Since(start)
		fmt.Printf("Time to respond to %s: %s\n", path, elapsed)
	}
}

func verifyPassword(r *http.Request) bool {
	if Config.Username == "" && Config.Password == "" {
		return true
	}

	u, p, ok := r.BasicAuth()
	if !ok {
		return false
	}

	return u == Config.Username && p == Config.Password
}

func handleRootPostPut(w http.ResponseWriter, r *http.Request) {
	// Verify auth
	defer measureHttpReponseTime(time.Now(), r.URL.Path)

	if !verifyPassword(r) {
		http.Error(w, "Invalid login", 401)
	}

	// Get the body
	body, _ := ioutil.ReadAll(r.Body)

	go routeTraffic(r.URL.Path, string(body))
}

func main() {
	// Get prog dir and thereby html template dir
	exedir = filepath.Dir(os.Args[0])

	readConfig()
	r := mux.NewRouter()
	r.HandleFunc("/{[x|]urlin}", handleRootPostPut).Methods("POST", "PUT")
	r.HandleFunc("/favicon.ico", handleNothing)
	http.Handle("/", r)

	if Config.Address == "" {
		Config.Address = ":8222"
	}

	fmt.Printf("Serve address %v\n", Config.Address)

	e := http.ListenAndServe(Config.Address, r) // Blocking function
	if e != nil {
		fmt.Printf("Serve failed %v\n", e)
	}
}
