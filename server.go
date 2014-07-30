package main

import (
  "fmt"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "database/sql"
  _ "github.com/lib/pq"
	"net/http"
)

type Person struct {
  Id      int
  Name    string
  Age     int
  Email   string
  Job     string
}

func SetupDB() *sql.DB {
  db, err := sql.Open("postgres", "dbname=martini_example sslmode=disable")
  PanicIf(err)
  return db
}

func PanicIf(err error) {
  if err != nil {
    panic(err)
  }
}

func index (r render.Render, req *http.Request, db *sql.DB) {
  rows, err := db.Query("SELECT * FROM persons")
	PanicIf(err)
  defer rows.Close()
  persons := []Person{}

  for rows.Next() {
    p := Person{}
    err := rows.Scan(&p.Id, &p.Name, &p.Age, &p.Job, &p.Email)
    PanicIf(err)
    persons = append(persons, p)
  }

  r.HTML(200, "index", persons)
}

func newPerson(r render.Render) {
  r.HTML(200, "persons/new", nil)
}

func createPerson(r render.Render, req *http.Request, db *sql.DB){
  fmt.Println(req.FormValue("person"))
  _, err := db.Query("INSERT INTO persons (id, name, age, email, job) VALUES (DEFAULT, $1, $2, $3, $4)", req.FormValue("person[name]"), req.FormValue("person[age]"), req.FormValue("person[email]"), req.FormValue("person[job]"))
	PanicIf(err)
  r.Redirect("/")
}

func main() {
  m := martini.Classic()

  m.Map(SetupDB())

  m.Use(render.Renderer(render.Options {
    Directory: "views",
    Layout: "layouts/layout",
    Charset: "UTF-8",
    IndentJSON: true,
    IndentXML: true,
  }))

  m.Get("/persons/new", newPerson)
  m.Post("/persons", createPerson)
  m.Get("/", index)
  m.Run()
}
