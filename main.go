package main

/*
Post:
Inout format: post json
Output format: Post same body as input

Get:
Web interface for monitoring the action.

*/
import (
	"container/list"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var config map[string]interface{}
var urls []map[string]interface{}

var exedir string

var webmess = list.New()

func logCall(in string, out_log string, res string) {
	fmt.Printf("In:  %v Out: %v Res: %v\n", in, out_log, res)
	// Attach to weebmess.
	webmess.PushFront(in + out_log)
}

func handleNothing(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("DO NOTHING for %v\n", r.URL.Path)
}

func pubMqtt(postString string, url_config map[string]interface{}) {
	var root_ca_file string
	var cert_file string
	var private_key_file string
	var out string
	var mqtt_topic string
	var err error

	// Mandatory parameters
	temp := url_config["out"]
	if temp == nil {
		log.Fatal("No 'out' parameter (broker url)")
	}
	out = temp.(string)
	opts := MQTT.NewClientOptions().AddBroker(out)

	temp = url_config["mqtt_topic"]
	if temp == nil {
		log.Printf("No 'mqtt_topic' parameter. Using input path")
		temp = url_config["in"]
	}
	mqtt_topic = temp.(string)

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
		fmt.Printf("Using root CA")
		_, err := os.Stat(root_ca_file)
		if err != nil {
			fmt.Printf("root_ca err%v\n", err)
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
		fmt.Printf("postMqtt %v\n%v\n%v\n%v\n", postString, root_ca_file, cert_file, private_key_file)
		tlsConf, err := makeTlsConfig(root_ca_file, cert_file, private_key_file)
		fmt.Printf("TLS CFG ERR: %v\n", err)
		opts.SetTLSConfig(tlsConf)
		opts.SetCleanSession(true)
	}

	c := MQTT.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	} else {
		token = c.Publish(mqtt_topic, 0, false, postString)
		token.Wait()
		if token.Error() != nil {
			fmt.Println("posted token err:", token.Error())
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

func routeTraffic(path string, jbody string) {
	//var val float32
	var routed int
	fmt.Printf(path)
	// fmt.Printf("Trying to route %v\n", path)
	for _, b := range urls {
		tmp := b["in"]
		inurl, ok := tmp.(string)
		if ok && strings.Index(path, inurl) == 0 {
			//fmt.Println("Found match")

			newBody := TransformBody(jbody, b)

			tmp, exist := b["protocol"]
			var prot string
			if exist {
				prot = tmp.(string)
			}
			if exist && strings.Index("mqtt", prot) == 0 {
				pubMqtt(string(newBody), b)
				routed += 1
				return
			}

			// HTTP Post is default protocol
			tmp = b["out"]
			outurl, _ := tmp.(string)
			outlog, _ := tmp.(string)
			fmt.Printf("Route %v to %v using http POST %v\n", inurl, outurl[0:20], outlog)
			req, err := http.NewRequest("POST", outurl, strings.NewReader(newBody))
			hk, hk_exist := b["header_key"]
			hv, hv_exist := b["header_value"]
			if hk_exist && hv_exist {
				key := hk.(string)
				val := hv.(string)
				req.Header.Set(key, val)
			}
			req.Header.Set("Content-Type", "application/json")
			temp := b["root_ca"]

			var client *http.Client
			if temp != nil {
				fmt.Println("Found root_ca")
				ca := temp.(string)
				temp = b["cert"]
				cert := temp.(string)
				temp = b["private_key"]
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

			logCall(inurl, outurl, resp.Status)
			routed += 1
		}
	}
	if routed == 0 {
		fmt.Printf("URL %v not routed.\n", path)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Verify auth
	bu, bu_exist := config["in_username"]
	bp, bp_exist := config["in_password"]
	if bu_exist && bp_exist {
		u, p, ok := r.BasicAuth()
		if !ok || u != bu || p != bp {
			http.Error(w, "Invalid login", 401)
		}
	}

	// Get the body
	body, _ := ioutil.ReadAll(r.Body)

	vars := mux.Vars(r)
	urlarg := vars["urlin"]
	fmt.Fprintf(w, "OK\n")
	if r.Method == http.MethodPost {
		go routeTraffic(r.URL.Path, string(body))
	} else {
		fmt.Printf("Only post supported. %v\n", urlarg /*r.URL.Path*/)
	}
}


func handleWeb(w http.ResponseWriter, r *http.Request) {
	htmltemp := filepath.Join(exedir, "templates", "index.html")
	t, err := template.ParseFiles(htmltemp)
	if (err != nil) {
		fmt.Printf("Cannot find html template\n")
		fmt.Fprintf(w, "Oh, no! Cannot show web page!\n")
		return
	}
	t.Execute(w, nil)
}

func readConfig() {
	// Get prog dir and thereby html template dir
	exedir = filepath.Dir(os.Args[0])

	// TODO: Support folder with many json files in.
	confdir := flag.String("conf", "", "Configuration file")
	flag.Parse()
	if len(*confdir) < 1 {
		u, _ := user.Current()
		*confdir = u.HomeDir + "/.route2cloud.json"
	}
	raw, err := ioutil.ReadFile(*confdir)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Route data to cloud, config %v\n", *confdir)
	u2 := config["urls"].([]interface{})
	for _, v := range u2 {
		urls = append(urls, v.(map[string]interface{}))
	}
}

func main() {
	readConfig()
	r := mux.NewRouter()
	r.HandleFunc("/", handleWeb).Methods("GET")
	r.HandleFunc("/{urlin}", handleRoot).Methods("POST")
	//http.HandleFunc("/favicon.ico", handleNothing)
	r.HandleFunc("/favicon.ico", handleNothing)
	//http.HandleFunc("/", handleRoot)
	http.Handle("/", r)

	addr := ":8080"
	tmp, addr_cfg_exist := config["address"]
	if addr_cfg_exist {
		addr = tmp.(string)
	}

	fmt.Printf("Serve address %v\n", addr)

	log.Fatal(http.ListenAndServe(addr, r))
	fmt.Printf("Bye!")
}
