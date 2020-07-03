package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	log "github.com/sirupsen/logrus"
)

var Db *sql.DB

func InitDB(host, port, dbName, user, pass string) *sql.DB {
	if Db != nil {
		return Db
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s", user, pass, host, port, dbName))
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	Db = db
	return Db
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}

	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://migration/mysql",
		"db",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}