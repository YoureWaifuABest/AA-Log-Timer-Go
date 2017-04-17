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

// Change these settings as needed
var client = redis.NewClient(&redis.Options{
          Addr:     "localhost:6379",
          Password: "",
          DB:       0,
})

// Load template into RAM; only one request to logs.html is necessary per server
var logTemplate = template.Must(template.ParseFiles("logs.html"))

// Main handler; displays web page.
func logsHandler(w http.ResponseWriter, r *http.Request) {
	// Execute values into template
	err := logTemplate.Execute(w, nil)
	// If the server can't access the template, bail
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// I don't like that there are four handlers doing essentially the same thing;
// Could at least use just 1 function.
// Clean-up later.
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
	// Get the value of key
	str, err := client.Get(key).Result()
	if err != nil {
		fmt.Println(err)
	}
	// Convert key's value to an integer, so it can be 
	// decremented and compared to other ints.
	val, err := strconv.Atoi(str)
	// If an error occurs (if the string is not a number),
	// return. Should never happen, since the value has
	// already been parsed by the save function
	if err != nil {
		return
	}
	for ; val >= 0; val-- {
		// Save the value to the database every interval of 5
		// Can cause some weird behaviour, but saves insignificant 
		// amounts of stress on server :^)
		if val % 5 == 0 || val == 0 {
			// Error value is ignored for now.
			client.Set(key, val, 0).Err()
		}
		// Slightly inaccurate if I remember correctly
		// For our purposes, should work fine
		time.Sleep(time.Second)
	}
	return
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	// Initializes r to be parsed
	r.ParseForm()
	var key string
	// There has to be a better way
	for key, _ = range r.Form {
		;
	}
	val := r.FormValue(key)
	newVal, err := strconv.Atoi(val)
	// If the value is not able to be converted to an int
	// (if it's not a number)
	// return. 
	// Error is already handed client side; anyone entering
	// invalid input is bypassing said error
	if err != nil {
		return
	}
	// If no value is entered
	// Should also be illegal through normal methods
	if newVal == 0 {
		// Redirect back to main page
		http.Redirect(w, r, "/logs/", http.StatusFound)
		return
	}
	// Save the value in the database, converting minutes to seconds
	client.Set(key, (newVal*60), 0).Err()
	// Start countdown for newly saved value
	go countDown(key)
	http.Redirect(w, r, "/logs/", http.StatusFound)
}

// Serves the JS. Will be replaced if later necessary
func sendLogsJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "scripts/logs.js")
}

func main() {
	http.HandleFunc("/logs/", logsHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/nuitimer/", nuiHandler)
	http.HandleFunc("/foresttimer/", forestHandler)
	http.HandleFunc("/abovetimer/", aboveHandler)
	http.HandleFunc("/belowtimer/", belowHandler)
	http.HandleFunc("/scripts/logs.js", sendLogsJs)
	// Starts the countdown for all four functions on server start
	// these will self-terminate if unnecessary
	go countDown("nui")
	go countDown("forest")
	go countDown("aboveCastle")
	go countDown("belowCastle")
	// Start listening on port 80 (default port for http), logging
	// Fatal errors (and closing upon error)
	log.Fatal(http.ListenAndServe(":80", nil))
}
