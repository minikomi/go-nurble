package main

import (
	"github.com/gorilla/pat"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	index   *template.Template
	nurb    *template.Template
	err     error
	nouns   []string
	cleaner *regexp.Regexp
)

func init() {
	index, err = template.ParseFiles("base.html", "index.html")
	if err != nil {
		panic(err)
	}
	nurb, err = template.ParseFiles("base.html", "nurb.html")
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile("nouns.txt")
	if err != nil {
		panic(err)
	}
	nouns = strings.Split(string(content), "\n")

	cleaner, err = regexp.Compile("[^a-z ]")
	if err != nil {
		panic(err)
	}
}

func nounCheck(word string) bool {
	for _, noun := range nouns {
		if word == noun {
			return true
		}
	}
	return false
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

func nurbleHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Input Error", 400)
		return
	}
	rawtext := r.FormValue("text")
	text := strings.ToUpper(rawtext)
	words := strings.Fields(cleaner.ReplaceAllString(strings.ToLower(rawtext), ""))

	for _, word := range words {
		if !nounCheck(word) {
			r, err := regexp.Compile("(\\b)(?i)" + word + "(\\b)")
			if err != nil {
				http.Error(w, "Parse Error", 400)
				return
			}
			text = r.ReplaceAllString(text, "$1<span class='nurble'>nurble</span>$2")
		}
	}
	text = strings.Replace(text, "\n", "<br>", -1)

	nurb.Execute(w, &map[string]template.HTML{"Text": template.HTML(text)})
}

func main() {
	r := pat.New()
	r.Get("/", indexHandler)
	r.Post("/nurble", nurbleHandler)
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.Handle("/", r)
	http.ListenAndServe(":9000", nil)
}
