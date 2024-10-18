package main

import (
	"database/sql"
	"dem3_demo_v2/pkg/config"
	"dem3_demo_v2/pkg/logging"
	"dem3_demo_v2/pkg/models/mysql"
	"dem3_demo_v2/pkg/models/postgresql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// go build -o ./bin/dem3 dem3_demo_v2/cmd/web/
type application struct {
	//errorLog      *log.Logger
	//infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	profData      *postgresql.ProfDataModel
	templateCache map[string]*template.Template
	logger        logging.Logger
	config        *config.Config
}

func main() {
	logger := logging.GetLogger()
	cfg := *config.GetConfig()

	//mysql
	//db, err := openDB(cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)
	////db, err := openDB(cfg.Storage.Username, cfg.Storage.Password, "", "", cfg.Storage.Database)
	////db, err := openDB(cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database)
	//if err != nil {
	//	logger.Fatal(err)
	//}

	//postgres
	//host, port, user, password, database
	db, err := openDBPostgress(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		snippets:      &mysql.SnippetModel{DB: db},
		profData:      &postgresql.ProfDataModel{DB: db},
		templateCache: templateCache,
		logger:        logger,
		config:        &cfg,
	}

	srv := &http.Server{
		Addr:     cfg.Listen.Port,
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Handler:  app.routes(),
	}

	logger.Info("start app on addr", cfg.Listen.Port)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// mysql
func openDB(user, password, host string, port int, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, database))
	//db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/snip?parseTime=true?")
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, database))
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(3308:3308)/%s?", user, password, database))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// postgresql
func openDBPostgress(host string, port int, user, password, database string) (*sql.DB, error) {
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database))
	//connStr := "user=postgres password=mozene33 dbname=snip sslmode=disable"
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database))
	//db, err := sql.Open("postgres", "user=postgres password=mozene33 dbname=productdb sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
