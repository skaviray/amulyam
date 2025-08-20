package main

import (
	"amulyam/api"
	db "amulyam/db/sqlc"
	"amulyam/utils"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("unable to load the config fie %e", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db %e", err)
	}
	log.Println("successfully connected to database..")

	store := db.NewStore(conn)
	server, err := api.New(store, config)
	if err != nil {
		log.Println(err)
	}
	err = server.CreateDefaultAdminUser()
	if err != nil {
		log.Fatal("cannot create default admin user: ", err)
	}
	server.Start(config.ListenAddress)
}
