package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	raut "web/Router"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLFiles("./templates/index.html", "./templates/login.html", "./templates/user.html")
	u1 := r.Group("/u1")

	{
		u1.GET("/home", raut.Page)
		u1.POST("/SignUp", raut.SignUp)
		u1.POST("/Login", raut.Login)
		u1.GET("/PageLogin", raut.PageLogin)
		u1.GET("/GotUser")
		u1.POST("/AddPic/:id", raut.AddPic)
	}
	server := &http.Server{
		Addr:    ":8008", // HTTPS port
		Handler: r,       // Gin router as the handler
	}

	err := server.ListenAndServeTLS("./tls/MyCertificate.crt", "./tls/Mykey.key")
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
