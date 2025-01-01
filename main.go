package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
    "path/filepath"
    "log"
)

type Counter struct {
    ID    int
    Title string
    Value int
}

type CounterGroup struct {
    ID       int
    Title    string
    counters Counter
}


func main () {
    // initialise initial Counter Group
    counter_groups []CounterGroup := []
    counter_groups = append(counter_groups, CounterGroup{ ID: 1, Title: "Generic", counters: [] })

    // make initial example counters
    counters := make(map[int]Counter)
    counter_tracking_id := 3

    counters[1] = Counter{ID: 1, Title: "Example", Value: 0}
    counters[2] = Counter{ID: 2, Title: "Another Example", Value: 0}

    // add inital example counters to inital counter group
    counter_groups[0].counters = counters

	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

    tmpl := template.Must(template.ParseFiles(files...))

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w,"index.html", nil)
    })

    http.HandleFunc("/counters", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "counters.html", counters)
    })

    http.HandleFunc("/increment/", func(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.URL.Path[len("/increment/"):])
        counter := counters[id]
        counter.Value++
        counters[id] = counter
        fmt.Println("incrementing ", id, counters[id])
        tmpl.ExecuteTemplate(w, "single-counter.html", counter)
    })

    http.HandleFunc("/decrement/", func(w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.URL.Path[len("/decrement/"):])
        counter := counters[id]
        counter.Value--
        counters[id] = counter
        fmt.Println("decrementing ", id, counters[id])
        tmpl.ExecuteTemplate(w, "single-counter.html", counter)
    })

    http.HandleFunc("/new-counter", func (w http.ResponseWriter, r *http.Request) {
        fmt.Println("New counter creating")
        title := r.FormValue("title")

        if title != "" {
            counter := Counter { ID: counter_tracking_id, Title: title, Value: 0 }
            counters[counter_tracking_id] = counter

            counter_tracking_id++

        }

        tmpl.ExecuteTemplate(w, "counters.html", counters)
    })

    http.HandleFunc("/new-group", func(w http.ResponseWriter, r *http.Request) {
        log.Println("New group being made")
        title := r.FormValue("title")

        if title != "" {
            log.Println("New group name: ")


        } else {
            log.Println("No title, no group being made")
        }


    })

    http.HandleFunc("/delete/", func (w http.ResponseWriter, r *http.Request) {
        id, _ := strconv.Atoi(r.URL.Path[len("/delete/"):])

        delete(counters, id)

        tmpl.ExecuteTemplate(w, "counters.html", counters)

    })

    fmt.Println("Server Running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)

}

