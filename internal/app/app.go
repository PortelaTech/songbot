package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DbConfig struct {
	Name string
	Host string
	Port int
	User string
	Pass string
}

type TelegramConfig struct {
	BotId  string
	ApiKey string
}

type Config struct {
	Db DbConfig
	Telegram TelegramConfig
}

type App struct {
	db   *sqlx.DB
	conf Config
	loc *time.Location
}

func setCors(w http.ResponseWriter) {
	frontendUrl := "http://localhost:3000"
	w.Header().Set("Access-Control-Allow-Origin", frontendUrl)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (ac *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	setCors(w)
	fmt.Fprintf(w, "This is the RESTful api")
}

func (ac *App) SearchSongs(w http.ResponseWriter, r *http.Request) {
	setCors(w)
	fmt.Fprintf(w, "Write some songs")
}

func (ac *App) CorsHandler(w http.ResponseWriter, r *http.Request) {
	setCors(w)
}

func GetConfig() Config {
	port,_ := strconv.Atoi(os.Getenv("DB_PORT"))
	conf := Config {
		DbConfig {
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			port,
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
		},
		TelegramConfig {
			os.Getenv("TELEGRAM_BOTID"),
			os.Getenv("TELEGRAM_APIKEY"),
		},
	}
	fmt.Printf("%+v\n", conf)
	return conf;
}


func GetApp(db *sqlx.DB,conf Config) App {
	loc, _ := time.LoadLocation("America/Bahia")
	return App{db, conf, loc}
}

