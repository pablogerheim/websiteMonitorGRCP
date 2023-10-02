package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados() *gorm.DB {
	stringDeConexao := "host=localhost user=root password=root dbname=root port=5432 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(stringDeConexao), &gorm.Config{})
	if err != nil {
		log.Panic("Erro ao conectar com banco de dados", err)
	}
	if DB == nil {
		fmt.Println("erro DB server")
		return nil
	}
	return DB
}
