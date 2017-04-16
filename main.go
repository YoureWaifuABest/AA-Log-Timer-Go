package main

import (
	"github.com/go-redis/redis"
	"net/http"
	"html/template"
	"strconv"
	"time"
	"fmt"
	"log"
)

var client = redis.NewClient(&redis.Options{
          Addr:     "localhost:6379",
          Password: "",
          DB:       0,
})

type Logs struct {
	Nui string
	Forest string
	AboveCastle string
	BelowCastle string
}

var logTemplate = template.Must(template.ParseFiles("logs.html"))

func logsHandler(w http.ResponseWriter, r *http.Request) {
	var p Logs
	var err error

	p.Nui, err       = client.Get("nui").Result()
	if err != nil {
		fmt.Println(err)
	}
	p.Forest, _      = client.Get("forest").Result()
	p.AboveCastle, _ = client.Get("aboveCastle").Result()
 	p.BelowCastle, _ = client.Get("belowCastle").Result()
	err = logTemplate.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nuiHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := client.Get("nui").Result() 
	fmt.Fprint(w, str) 
}

func forestHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := client.Get("forest").Result() 
	fmt.Fprint(w, str) 
}

func aboveHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := client.Get("aboveCastle").Result() 
	fmt.Fprint(w, str) 
}

func belowHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := client.Get("belowCastle").Result() 
	fmt.Fprint(w, str) 
}

func countDown(key string) {
	str, _ := client.Get(key).Result()
	val, _ := strconv.Atoi(str)
	for ; val >= 0; val-- {
		if val % 5 == 0 || val == 0 {
			client.Set(key, val, 0).Err()
		}
		time.Sleep(time.Second)
	}
	return
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var key string
	// There has to be a better way
	for key, _ = range r.Form {
		;
	}
	val := r.FormValue(key)
	newVal, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	if newVal == 0 {
		http.Redirect(w, r, "/logs/", http.StatusFound)	
		return
	}
	client.Set(key, (newVal*60), 0).Err()
	go countDown(key)
	http.Redirect(w, r, "/logs/", http.StatusFound)
}

func main() {
	http.HandleFunc("/logs/", logsHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/nuitimer/", nuiHandler)
	http.HandleFunc("/foresttimer/", forestHandler)
	http.HandleFunc("/abovetimer/", aboveHandler)
	http.HandleFunc("/belowtimer/", belowHandler)
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("/scripts/"))))
	go countDown("nui")
	go countDown("forest")
	go countDown("aboveCastle")
	go countDown("belowCastle")
	log.Fatal(http.ListenAndServe(":80", nil))
}
