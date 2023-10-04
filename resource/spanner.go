package resource

import (
	"fmt"
	"log"
	"websiteMonitor/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados() *gorm.DB {
	stringDeConexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.SPANNER_HOST, config.SPANNER_USER, config.SPANNER_PASSWORD, config.SPANNER_DATABASE, config.SPANNER_PORT)

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



