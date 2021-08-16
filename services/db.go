package services

import (
	"fmt"
	"log"

	config "github.com/dockboxhq/server/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var schema = `
create table if not exists dockbox (
    id CHAR(36) PRIMARY KEY default uuid(),
    custom_id VARCHAR(64) default null UNIQUE,
    container_id CHAR(64) UNIQUE,
    source TEXT,
    created_at timestamp default now(),
    last_modified timestamp default now() ON UPDATE now()
);
`

var DbConn *sqlx.DB

func Init() {
	host := config.Config.DATABASE_HOST
	port := config.Config.DATABASE_PORT
	database := config.Config.DATABASE_NAME
	user := config.Config.DATABASE_USER
	pass := config.Config.DATABASE_PASSWORD
	DbConn = sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", user, pass, host, port, database))
	log.Println("Connected to DB")
	log.Println("Verifying schema...")
	// DbConn.MustExec(schema)
	log.Println("Database is connected and ready!")
}
