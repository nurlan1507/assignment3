package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nurlan1507/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
)

import _ "github.com/jackc/pgx/v4"
import _ "github.com/go-sql-driver/mysql"

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "postgres://postgres:admin@localhost:5432/postgres", "dsn for postgresql")
	flag.Parse()

	//file where to write logs
	file, err := os.OpenFile("serverLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//logs
	infoLog := log.New(file, "INFO	\t", log.Ldate|log.Ltime)
	errorLog := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	//database connection
	//db, err := OpenDb(cfg.dsn)
	db, err := ConnectToDb(cfg.dsn)
	if err != nil {
		errorLog.Println(err)
	}
	templateCache, err := newTemplateCache()
	app := &application{errorLog: errorLog, infoLog: infoLog, snippets: &models.SnippetModel{Db: db}, templateCache: templateCache}
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("LOGLOGLOGLOGLOG")

	srcError := srv.ListenAndServe()
	infoLog.Println("staring server in port %v", cfg.addr)
	if srcError != nil {
		srv.ErrorLog.Fatal(err)
	}

}

func ConnectToDb(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Println("Unable to connect to database")
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Unable to connect to database")
		return nil, err
	}
	return pool, nil
}

//func OpenDb(dsn string) (*sql.DB, error) {
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		return nil, err
//	}
//	if err = db.Ping(); err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
