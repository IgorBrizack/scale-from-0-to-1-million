package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	dbUser       string
	dbPass       string
	dbName       string
	dbHostMaster string
	dbHostSlave  string
}

func NewDatabase() *Database {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env.")
	}

	db := &Database{
		dbUser:       os.Getenv("MYSQL_USER"),
		dbPass:       os.Getenv("MYSQL_PASSWORD"),
		dbName:       os.Getenv("MYSQL_DATABASE"),
		dbHostMaster: os.Getenv("DB_HOST_MASTER"),
		dbHostSlave:  os.Getenv("DB_HOST_SLAVE"),
	}

	return db
}

func (d *Database) MasterDB() *gorm.DB {
	masterDB, err := gorm.Open(mysql.Open(d.GetMasterDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the master database:", err)
	}

	return masterDB
}

func (d *Database) SlaveDB() *gorm.DB {
	slaveDB, err := gorm.Open(mysql.Open(d.GetSlaveDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the slave database:", err)
	}

	return slaveDB
}

func (d *Database) GetMasterDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		d.dbUser, d.dbPass, d.dbHostMaster, d.dbName)
}

func (d *Database) GetSlaveDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		d.dbUser, d.dbPass, d.dbHostSlave, d.dbName)
}
