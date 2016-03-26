package main

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"text/template"
)

type Inventory struct {
	Material string
	Count    uint
}

// landing page
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "<h1><b>Welcome to your Raspberry Pi!</b></h1><br />")
}

// says hello to name param
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

/*
	cpu stats, thanks to
	http://www.josephspurrier.com/how-to-use-template-blocks-in-go-1-6/
*/
func Stats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	stats, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	for i, s := range stats.CPUStats {
		if i < len(stats.CPUStats)-1 {
			output := fmt.Sprintf("%v\n\n", s)
			w.Write([]byte(output))
		} else {
			output := fmt.Sprintf("%v", s)
			w.Write([]byte(output))
		}
	}
}

/*
	advanced testing of text/template, thanks to
	http://www.josephspurrier.com/how-to-use-template-blocks-in-go-1-6/
*/
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("tmpl/base.tmpl", "tmpl/about.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func Template(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sweaters := Inventory{"wool", 17}

	template.Must(template.ParseGlob("tmpl/*"))

	t, err := template.ParseFiles("tmpl/welcome.html") //create a new template
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, sweaters) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		panic(err)
	}
}

func main() {
	router := httprouter.New()

	router.GET("/", Index)

	router.GET("/hello/:name", Hello)

	router.GET("/stats", Stats)

	router.GET("/template", Template)

	router.GET("/about", About)

	log.Fatal(http.ListenAndServe(":3000", router))
}
