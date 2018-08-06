package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/G-Node/gin-cli/ginclient"
	glog "github.com/G-Node/gin-cli/ginclient/log"
	"github.com/G-Node/gin-valid/log"
	"github.com/G-Node/gin-valid/resources/templates"
	"github.com/gorilla/mux"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		log.Write("Login page")
		tmpl := template.New("layout")
		tmpl, err := tmpl.Parse(templates.Layout)
		if err != nil {
			log.Write("[Error] failed to parse html layout page")
			http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
			return
		}
		tmpl, err = tmpl.Parse(templates.Login)
		if err != nil {
			log.Write("[Error] failed to render login page")
			http.ServeContent(w, r, "unavailable", time.Now(), bytes.NewReader([]byte("500 Something went wrong...")))
			return
		}
		tmpl.Execute(w, nil)
	} else {
		log.Write("Doing login")
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		client := ginclient.New(serveralias)
		glog.Init("")
		glog.Write("Performing login from gin-valid")
		err := client.Login(username, password, "gin-valid")
		if err != nil {
			log.Write("[error] Login failed: %s", err.Error())
		}
		// TODO: Store user token in session cookie
		// Redirect to repo listing
	}
}

const repostmpl = `
<html>
    <head>
    <title></title>
    </head>
    <body>
        {{ range . }}

			<p><b><a href=/repos/{{.FullName}}/enable>{{.FullName}}</a></b></p>
			<p>{{.Description}} {{.Website}}</p>
			<hr>
		{{ end }}
    </body>
</html>
`

func ListRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
	}
	vars := mux.Vars(r)
	user := vars["user"]
	cl := ginclient.New(serveralias)
	cl.LoadToken()
	fmt.Printf("Requesting repository listing for user %s\n", user)
	fmt.Printf("Server alias: %s\n", serveralias)
	fmt.Println("Server configuration:")
	fmt.Println(cl.Host)

	repos, err := cl.ListRepos(user)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Got %d repos\n", len(repos))

	rl := template.New("repos")
	rl.Parse(repostmpl)
	rl.Execute(w, &repos)
}
