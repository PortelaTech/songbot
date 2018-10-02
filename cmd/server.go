package main

import (
	"fmt"
	"github.com/PortelaTech/songbot/internal/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"time"
)

func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}


func main() {
	var conf = app.GetConfig()
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Db.Host,  conf.Db.Port, conf.Db.User,  conf.Db.Pass, conf.Db.Name)
	db, err := sqlx.Connect("postgres",dataSourceName)
	if (err != nil) {
		panic(err)
	}
	defer db.Close()
	var a = app.GetApp(db,conf)

	myHandler := http.HandlerFunc(a.SearchSongs)
	chain := alice.New(timeoutHandler, nosurf.NewPure).Then(myHandler)

	fmt.Println("Server starting at port 8080.")
	log.Fatal(http.ListenAndServe(":8080", chain))
}
