package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/anil/golang-imagestorage/handlers"
	"github.com/anil/golang-imagestorage/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize initialise db and router
func (a *App) Initialize(user, password, dbname, host string, port int) error {

	//initialise log
	initLog()

	var err error
	a.DB, err = a.openDb(user, password, dbname, host, port)
	if err != nil {
		log.Fatal("error while creating connection to the database: %#v", err)
		return err
	}
	a.Router = mux.NewRouter()
	//initalise tables
	store.InitialiseTables(a.DB)

	return nil
}

func (a *App) initializeRoutes() {
	mh := handlers.NewMetafestHandler(a.DB)
	//Albums
	a.Router.HandleFunc("/albums", mh.GetAllAlbumHandler).Methods("GET")
	a.Router.HandleFunc("/album/{id}", mh.GetAlbumHandler).Methods("GET")
	a.Router.HandleFunc("/album", mh.InsertAlbumHandler).Methods("PUT")
	a.Router.HandleFunc("/album/{id}", mh.DeleteAlbumHandler).Methods("DELETE")

	//Images
	a.Router.HandleFunc("/album/{id}/image", mh.InsertImageHandler).Methods("PUT")
	a.Router.HandleFunc("/album/{id}/image/{}", mh.InsertImageHandler).Methods("PUT")
	a.Router.HandleFunc("/album/{id}/image", mh.InsertImageHandler).Methods("PUT")
	a.Router.HandleFunc("/album/{id}/image", mh.InsertImageHandler).Methods("PUT")

	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func (a *App) openDb(user, password, database, host string, port int) (*sql.DB, error) {
	log.Println("opening db")
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, database)
	d, err := sql.Open("mysql", connectionStr)
	if err != nil {
		log.Fatal("Unable to connect database", err)
		return nil, err
	}
	d.SetConnMaxLifetime(time.Minute * 20)
	d.SetMaxIdleConns(4)
	d.SetMaxOpenConns(4)
	return d, nil
}

func initLog() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func (a *App) Run(addr string) {}
