package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shurcooL/go/gzip_file_server"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

func loadTemplates() (*template.Template, error) {
	t := template.New("").Funcs(template.FuncMap{})
	t, err := vfstemplate.ParseGlob(assets, t, "/assets/*.tmpl")
	return t, err
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	t, err := loadTemplates()
	if err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ct := fmt.Sprintf("%9x", time.Now().Unix())
	tt := make([]string, 9)
	for i, _ := range ct {
		tt[i] = ct[i : i+1]
	}

	ny := (time.Now().Unix() & 0xffffff000000) + 0x1000000
	dny := time.Unix(ny, 0).UTC()

	var data = struct {
		CurrentTime  []string
		NextYearUnix int64
		NextYear     string
	}{
		CurrentTime:  tt,
		NextYearUnix: ny,
		NextYear:     dny.Format("2 January 2006, 15:04:05 UTC"),
	}

	err = t.ExecuteTemplate(w, "index.html.tmpl", data)
	if err != nil {
		log.Println("t.Execute:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", mainHandler)
	http.Handle("/assets/", gzip_file_server.New(assets))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	printServingAt(*httpFlag)
	err := http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}

func printServingAt(addr string) {
	hostPort := addr
	if strings.HasPrefix(hostPort, ":") {
		hostPort = "localhost" + hostPort
	}
	fmt.Printf("serving at http://%s/\n", hostPort)
}
