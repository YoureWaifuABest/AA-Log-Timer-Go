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

func logsHandler(w http.ResponseWriter, r *http.Request) {
	var p Logs

	p.Nui, _         = client.Get("nui").Result()
	p.Forest, _      = client.Get("forest").Result()
	p.AboveCastle, _ = client.Get("aboveCastle").Result()
 	p.BelowCastle, _ = client.Get("belowCastle").Result()
	t, _ := template.ParseFiles("logs.html")
	t.Execute(w, p)
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
	go countDown("nui")
	go countDown("forest")
	go countDown("aboveCastle")
	go countDown("belowCastle")
	log.Fatal(http.ListenAndServe(":80", nil))
}
