package migrations

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repository interface {
	CreateTable()
}

type Connection struct {
	*gorm.DB
}

type User struct {
	UId      string
	Name     string
	Root_dir string
}

func InitDatabase() Repository {
	db, err := gorm.Open("postgres", "host=localhost user=postgres password=krishna009 sslmode=disable port=5432")
	if err != nil {
		return nil
	}

	db.Exec(`CREATE DATABASE filebox`)
	filebox, err := gorm.Open("postgres", "host=localhost user=postgres password=krishna009 dbname=filebox sslmode=disable port=5432")
	if err != nil {
		return nil
	}
	return &Connection{
		DB: filebox,
	}
}

func (c *Connection) CreateTable() {
	var user User
	c.DB.CreateTable(&user)
}
