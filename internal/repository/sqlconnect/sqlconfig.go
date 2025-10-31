package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/joho/godotenv"
)

func ConnectDb() (*sql.DB, error) {
	fmt.Println("Trying to Connect to Maria DB")
	// err := godotenv.Load("../../cmd/.env")			// Already loaded in the main function
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbhost := os.Getenv("HOST")
	dbport := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, dbhost, dbport, dbname)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		// panic(err)
		return nil, err
	}
	fmt.Println("Connected to Maria DB")
	return db, nil
}
