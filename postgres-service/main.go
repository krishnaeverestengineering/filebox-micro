package main

import "Filebox-Micro/postgres-service/migrations"

func main() {
	db := migrations.InitDatabase()
	db.CreateTable()
}
