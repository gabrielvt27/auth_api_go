package initializers

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbCharSet := os.Getenv("DB_CHARSET")
	dbLoc := os.Getenv("DB_LOC")
	containerName := os.Getenv("CONTAINER_NAME")

	DB, err = gorm.Open(mysql.Open(dbUser+":"+dbPass+"@tcp("+containerName+")/"+dbName+"?charset="+dbCharSet+"&parseTime=True&loc="+dbLoc), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to DB")
	}
}
