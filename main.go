package main

//go:generate go-bindata-assetfs static/... template
// # // go:generate go-bindata template

import (
	"fmt"
	"log"
	"net/http"

	"minoris.se/rabbitmq/camq"

	"github.com/BurntSushi/toml"
	"github.com/arschles/go-bindata-html-template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var cookiestore *sessions.CookieStore
var store = sessions.NewCookieStore([]byte("something very secret this time"))

type Config struct {
	AMQ     camq.AMQConfig
	AuthURL string

	Lights []Light
}

var config Config

func readTemplateFile(filename string) *template.Template {
	return template.Must(template.New("base", Asset).ParseFiles("template/base.html", filename))
}

func failOnErr(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		http.Error(w, "Sorry", 500)
		log.Panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	// session, _ := cookiestore.Get(r, "session-name")

	type Info struct {
		Title string
	}

	info := Info{"Index"}

	tmpl := readTemplateFile("template/index.html")
	err := tmpl.Execute(w, info)
	failOnErr(err, w, r)
}

func isAuthoriative(groups *[]string) bool {
	for _, group := range *groups {
		if group == "admin" || group == "light_controller" {
			return true
		}
	}
	return false
}
func auth(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/getgroupsfromtoken?token=%s", config.AuthURL, r.FormValue("token"))

	resp, err := http.Post(url, "nil", nil)

	if err != nil {
		fmt.Fprintf(w, "Sanity strikes and I get an error\n")
		fmt.Fprintf(w, "This is it: %q\n", err)
		return
	}

	defer resp.Body.Close()

	type Info struct {
		Groups []string
	}

	var info Info
	_, err = toml.DecodeReader(resp.Body, &info)
	failOnErr(err, w, r)
	// body, err := ioutil.ReadAll(resp.Body)

	if isAuthoriative(&info.Groups) {
		fmt.Fprintf(w, "Weeeee!!")
	} else {
		fmt.Fprintf(w, "Dream on")
	}
}

func main() {
	toml.DecodeFile("/opt/lightweb/config.toml", &config)

	cookiestore = sessions.NewCookieStore([]byte("bla bla bla"))

	r := mux.NewRouter()

	StaticFS(r)
	http.Handle("/", r)
	r.HandleFunc("/auth", auth)
	//r.HandleFunc("/", index)
	lightInit(r)

	log.Print("Server started")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}
