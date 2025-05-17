package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseCredentials struct {
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
}
type DatabaseConnection struct {
	databaseInstance    *gorm.DB
	databaseCredentials *DatabaseCredentials
}

func NewDatabaseConnection(databaseCredentials *DatabaseCredentials) *DatabaseConnection {
	return &DatabaseConnection{
		databaseCredentials: databaseCredentials,
	}
}

func (dbConn *DatabaseConnection) GetDatabaseConnection() *gorm.DB {
	if dbConn.databaseInstance == nil {
		sqlDialect := mysql.Open(
			fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				dbConn.databaseCredentials.DatabaseUsername,
				dbConn.databaseCredentials.DatabasePassword,
				dbConn.databaseCredentials.DatabaseHost,
				dbConn.databaseCredentials.DatabaseName))
		gormOpen, err := gorm.Open(sqlDialect, &gorm.Config{})
		dbConn.databaseInstance = gormOpen
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	return dbConn.databaseInstance
}
