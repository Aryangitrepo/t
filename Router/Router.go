package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"web/db"
)

func AddPic(c *gin.Context) {
	userid := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(), // Corrected typo in message
		})
	}

	filepat := filepath.Join("./user profile pic", userid+"-profile.jpg")
	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(), // Corrected typo in message
		})
	}
	defer src.Close()

	outFile, err := os.Create(filepat)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(), // Corrected typo in message
		})
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, src)
	c.JSON(200, gin.H{
		"message": "picture accepted",
	})
}

func Page(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func uid() string {
	uuidValue := uuid.New()
	return uuidValue.String()[:5]
}

var validInputRegex = regexp.MustCompile("^[a-zA-Z]+$")

func validateInput(username, password, email string) error {
	// Check if any field is empty or contains only spaces
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(email) == "" {
		return fmt.Errorf("fields cannot be empty or contain only spaces")
	}

	// Check if fields contain only alphabetic characters
	if !validInputRegex.MatchString(username) {
		return fmt.Errorf("username must only contain alphabetic characters (a-z, A-Z)")
	}
	if !validInputRegex.MatchString(password) {
		return fmt.Errorf("password must only contain alphabetic characters (a-z, A-Z)")
	}

	// You can add additional email validation here if needed

	return nil
}

func SignUp(c *gin.Context) {
	pass := c.PostForm("password")
	email := c.PostForm("email")
	username := c.PostForm("username")
	userid := uid()
	err := validateInput(username, pass, email)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}

	errChan1 := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	d, err := db.Con()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal server error", // Corrected typo in message
		})
		panic(err)
	}

	go db.AddUser(errChan1, d, &wg, userid, username, pass, email)

	wg.Wait()

	// Read errors from channels
	res := <-errChan1

	if res != nil {
		c.JSON(200, gin.H{
			"message": res.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/u1/home")
}

func PageLogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func Login(c *gin.Context) {
	// email := c.PostForm("email")
	// pass := c.PostForm("password")

	c.Redirect(http.StatusSeeOther, "/u1/GotUser")
}
func GotUser(c *gin.Context) {
	c.HTML(200, "user.html", nil)
}
