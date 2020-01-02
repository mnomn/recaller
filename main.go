package main

/*
 * Route messages from HTTP POST to another PST URL or to MQTT.
 */
import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	_ "html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//var config map[string]interface{}
var address string
var main_username string
var main_password string
var globalDebug bool

var routes []map[string]interface{}

var exedir string

type JsonTime struct {
	time.Time
}

func (t JsonTime)MarshalJSON() ([]byte, error) {
	//	stamp := fmt.Sprintf("\"%s\"", t.Format("Mon Jan _2"))
	stamp := "\"\""
	now := time.Now();
	if (t.Year() < 2000) {
		stamp = "\"Old times\""
	} else if t.Year() != now.Year() {
		stamp = "\"last year\""
	} else if t.Day() == now.Day() -1 { // Todo: New month.
		stamp = fmt.Sprintf("\"Yesterday %s\"", t.Format("15:04")) 
	} else if now.Day() == t.Day() {
		stamp = fmt.Sprintf("\"Today %s\"", t.Format("15:04"))
	} else { // Generic time stamp
		stamp = fmt.Sprintf("%s", t.Format("01-02 15:04"))
	}
	return []byte(stamp), nil
}

// // For every "in", save time and body
// var latestPosts = make(map[string]LatestPost)

type OldPost struct {
	Time JsonTime
	Input string
	Output string
	OutProtocol string
}

type OldPosts struct {
	posts []OldPost
}
// type oldPosts []OldPost

var oldPostsLists = make(map[string]*OldPosts)

func logCall(in string, out_log string, outProt string, res string) {
	fmt.Printf("In:  %v Out: %v Res: %v\n", in, out_log, res)
	// Attach to weebmess.

	ll, ok := oldPostsLists[in]
	if !ok || ll == nil {
		ll = &OldPosts{}
		ll.posts=make([]OldPost,0)
		// newL := make([]OldPost, 0)
		oldPostsLists[in] = ll
	}
	op := OldPost{JsonTime{time.Now()}, in, out_log, outProt}
	ll.posts = append(ll.posts, op)

	if len(ll.posts) > 20 {
		ll.posts = ll.posts[1:]
	}

}

// Remove leading "/"
func normalizeIn(in *string) {
	if (len(*in)>0 && (*in)[0] == '/') {
		s2 := (*in)[1:]
		*in=s2
	}
}

func handleNothing(w http.ResponseWriter, r *http.Request) {
}

func pubMqtt(postString string, url_config map[string]interface{}, debug bool) {
	var root_ca_file string
	var cert_file string
	var private_key_file string
	var out string
	var mqtt_topic string
	var err error

	// Mandatory parameters
	temp := url_config["out"]
	if temp == nil {
		fmt.Printf("No 'out' parameter (broker url) \n")
		os.Exit(1)
	}
	out = temp.(string)

	inurl := url_config["in"]

	opts := MQTT.NewClientOptions().AddBroker(out)

	temp = url_config["mqtt_topic"]
	if temp == nil {
		if globalDebug || debug {
			fmt.Println("No 'mqtt_topic' parameter. Using input path as topic")
		}
		temp = url_config["in"]
	}
	mqtt_topic = temp.(string)

	logCall(inurl.(string), mqtt_topic, "mqtt", "")
	fmt.Printf("Route to mqtt broker %v, topic %v\n", out, mqtt_topic)
	// Optional parameters
	temp = url_config["username"]
	temp2 := url_config["password"]
	if temp != nil && temp2 != nil {
		opts.SetUsername(temp.(string))
		opts.SetPassword(temp2.(string))
	}

	temp = url_config["root_ca"]
	if temp != nil {
		root_ca_file = temp.(string)
		if globalDebug || debug {
			fmt.Println("Using root CA")
		}
		_, err := os.Stat(root_ca_file)
		if err != nil {
			fmt.Println("root_ca err%v\n", err)
		}
		return
	}
	temp = url_config["cert"]
	if temp != nil {
		cert_file = temp.(string)
		_, err = os.Stat(cert_file)
		if err != nil {
			fmt.Printf("cert err%v\n", err)
		}
	}
	temp = url_config["private_key"]
	if temp != nil {
		private_key_file = temp.(string)
		_, err = os.Stat(private_key_file)
		if err != nil {
			fmt.Printf("private_key err%v\n", err)
		}
	}

	cid := "ClientID"
	opts.SetClientID(cid)
	if root_ca_file != "" || cert_file != "" || private_key_file != "" {
		if globalDebug || debug {
			fmt.Printf("Mqtt cert files%v\n%v\n%v\n", root_ca_file, cert_file, private_key_file)
		}
		tlsConf, err := makeTlsConfig(root_ca_file, cert_file, private_key_file)
		fmt.Printf("TLS CFG ERR: %v\n", err)
		opts.SetTLSConfig(tlsConf)
		opts.SetCleanSession(true)
	}

	c := MQTT.NewClient(opts)

	if globalDebug || debug {
		fmt.Printf("Mqtt data: %v\n", postString)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		val, retain := url_config["retain"]
		if retain {
			retain = val.(bool)
		}
		token = c.Publish(mqtt_topic, 0, retain, postString)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("posted token err: %v\n", token.Error())
		}
	}
	c.Disconnect(250)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// From https://github.com/manamanmana/aws-mqtt-chat-example
func makeTlsConfig(cafile, cert, key string) (*tls.Config, error) {
	var TLSConfig *tls.Config = &tls.Config{InsecureSkipVerify: false}

	var certPool *x509.CertPool
	var err error
	var tlsCert tls.Certificate
	if cafile != "" {
		certPool, err = getCertPool(cafile)
		if err != nil {
			return nil, err
		}
		TLSConfig.RootCAs = certPool
	}
	if cert != "" {
		certPool, err = getCertPool(cert)
		if err != nil {
			return nil, err
		}
		TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		TLSConfig.ClientCAs = certPool
	}
	if key != "" {
		if cert == "" {
			return nil, fmt.Errorf("key specified but cert is not specified")
		}
		tlsCert, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		TLSConfig.Certificates = []tls.Certificate{tlsCert}
	}
	return TLSConfig, nil
}

func getCertPool(pemPath string) (*x509.CertPool, error) {
	var certs *x509.CertPool = x509.NewCertPool()
	var pemData []byte
	var err error

	pemData, err = ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}
	certs.AppendCertsFromPEM(pemData)
	return certs, nil
}

func routeTraffic(path string, body string) {
	//var val float32
	var routed int
	if globalDebug {
		fmt.Printf("Incomming path: %v\n", path)
	}
	for _, route := range routes {
		tmp := route["in"]
		inurl, ok := tmp.(string)
		_, debug := route["debug"] // Optional. Debug print

		normalizeIn(&path)
		normalizeIn(&inurl)
		if ok && strings.Index(path, inurl) == 0 {

			if globalDebug || debug {
				fmt.Printf("Path %v configured\n", path)
			}

			newBody := TransformBody(body, route)

			opList, ok := oldPostsLists[inurl]
			fmt.Printf("OLD: %v, %v \n", opList, ok)
			tmp, ok = route["out"]
			if !ok {
				// If no out url: We are done.
				return
			}
			outurl, _ := tmp.(string)
			outlog, _ := tmp.(string)

			tmp, exist := route["protocol"]
			var prot string
			if exist {
				prot = tmp.(string)
			}
			if exist && strings.Index("mqtt", prot) == 0 {
				pubMqtt(string(newBody), route, debug)
				routed += 1
				return
			}

			// HTTP Post is default protocol
			fmt.Printf("Route %v to %v using http POST %v\n", inurl, outurl[0:20], outlog)
			req, err := http.NewRequest("POST", outurl, strings.NewReader(newBody))
			hk, hk_exist := route["header_key"]
			hv, hv_exist := route["header_value"]
			if hk_exist && hv_exist {
				key := hk.(string)
				val := hv.(string)
				req.Header.Set(key, val)
			}
			req.Header.Set("Content-Type", "application/json")
			temp := route["root_ca"]

			var client *http.Client
			if temp != nil {
				if globalDebug || debug {
					fmt.Println("Found root_ca")
				}
				ca := temp.(string)
				temp = route["cert"]
				cert := temp.(string)
				temp = route["private_key"]
				private_key := temp.(string)
				tc, _ := makeTlsConfig(ca, cert, private_key)

				client = &http.Client{
					Transport: &http.Transport{TLSClientConfig: tc},
				}
			} else {
				client = &http.Client{}
			}

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			logCall(inurl, outurl, "POST", resp.Status)
			routed += 1
		}
	}
	if routed == 0 {
		if globalDebug {
			fmt.Printf("URL %v not routed.\n", path)
		}
	}
}

func handleRootPost(w http.ResponseWriter, r *http.Request) {
	// Verify auth
	if len(main_username) > 0 && len(main_password) > 0 {
		u, p, ok := r.BasicAuth()
		if !ok || u != main_username || p != main_password {
			http.Error(w, "Invalid login", 401)
		}
	}

	// Get the body
	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	urlarg := vars["urlin"]
	if r.Method == http.MethodPost {
		go routeTraffic(r.URL.Path, string(body))
	} else {
		fmt.Printf("Only post supported. %v\n", urlarg /*r.URL.Path*/)
	}
}

/*
 * api/log: Get a list of the latest routes made.
 * Optional uery arguments (api/log?count=5&in=inId&out=outId)
 * in: Show only log from a certain input
 * out: Show only log from a certain output
 * count: How many messages to get (max 20)
 */
func handleApiLog(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["in"]
	if ok && len(keys) > 0 {
		fmt.Println("Got parameter " , keys[0])
		inp := keys[0]
		logs, ok := oldPostsLists[inp]
		if !ok {
			w.Write([]byte("[]"))
			return
		}
		js, _ := json.Marshal(logs.posts)
		w.Write(js)
		return
	} else {
		fmt.Println("No 'in' parameter")
	}

	w.Write([]byte("[]"))
}

/*
 * api/routes Get a list of all configured inputs
 * Optional query args: (api/routes?in=inId&out=outId)
 * in/out: Show only config for one input/output.
 */
func handleApiRoutes(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(routes)

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
	r.HandleFunc("/api/log", handleApiLog).Methods("GET")
	r.HandleFunc("/{[x|]urlin}", handleRootPost).Methods("POST")
	r.HandleFunc("/favicon.ico", handleNothing)
	http.Handle("/", r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir( exedir + "/web/")))

	if address == "" {
		address = ":8222"
	}

	fmt.Printf("Serve address %v\n", address)

	// Add some weblogging for test
	logCall("Start system", "route2cloud", "SYS", "OK")
	logCall("/testA", "SOME.URL", "POST", "OK")
	logCall("/testA", "SOME.URL", "POST", "OK")
	logCall("/testA", "SOME.URL", "POST", "OK")

	e := http.ListenAndServe(address, r) // Blocking function
	if (e != nil) {
		fmt.Printf( "Serve failed %v\n", e)
	}
}
