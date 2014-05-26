package main

import (
	"github.com/HouzuoGuo/tiedot/db"
)

const COLLECTION = "names"

func main() {
	database := GetDb()
	defer database.Close()
//	database.Drop(COLLECTION)
	database.Create(COLLECTION, 1)
	collection := database.Use(COLLECTION)
	collection.Index([]string{"firstname"})
	collection.Index([]string{"lastname"})
}

func GetDb() *db.DB {
	myDb, err := db.OpenDB("C:\\tmp\\tiedotgui2")
	if err != nil {
		panic(err)
	}
	return myDb
}