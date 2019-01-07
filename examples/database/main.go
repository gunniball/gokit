package main

import (
	"fmt"

	"github.com/finalist736/gokit/database"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dbCfg := &database.DBConfig{}

	dbCfg.ConnectionName = "test"
	dbCfg.Driver = "mysql"
	dbCfg.Dsn = "test:test@tcp(127.0.0.1:3306)/test"

	err := database.Add(dbCfg)
	if err != nil {
		panic(err)
	}

	session := database.GetSession("test")
	if session == nil {
		panic("no session returned")
	}
	var ids []int
	n, err := session.Select("id").From("test").Load(&ids)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ids loaded: %d - %v\n", n, ids)

	database.Close()
}
