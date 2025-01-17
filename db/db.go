package db

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type users struct {
	UserID   string `gorm:"primaryKey;not null"` // Primary key and NOT NULL
	Username string `gorm:"not null"`            // NOT NULL and UNIQUE constraint
	Password string `gorm:"not null"`            // NOT NULL constraint
	Email    string `gorm:"not null;unique"`     // NOT NULL and UNIQUE constraint
}

func Con() (*gorm.DB, error) {
	dsn := "admin:aryan0011@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&users{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connection initialized and migrations ran successfully.")
	return db, nil
}

func AddUser(errChan1 chan error, db *gorm.DB, wg *sync.WaitGroup, userid string, username string, pass string, email string) {
	defer wg.Done() // Ensure WaitGroup is decremented

	user := users{UserID: userid, Username: username, Password: pass, Email: email}

	if result := db.Create(&user); result.Error != nil {
		fmt.Println("Error creating user:", result.Error)

		errChan1 <- result.Error
		return
	}

	errChan1 <- nil
	close(errChan1)
}

// func FindUser(email string, password string) bool {
// 	db, err := con()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var count int
// 	fmt.Print(email, password)
// 	err = db.con.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND password = ?", email, password).Scan(&count)
// 	if count == 1 {
// 		return true
// 	}

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			fmt.Println("No user found with that email.")
// 			return false

// 		}
// 		return false
// 	}
// 	fmt.Printf("Total Rows: %d\n", count)
// 	return false
// }
