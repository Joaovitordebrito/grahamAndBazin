package main

import (
	"graham-bazin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/fiis", controllers.GetBazin)
	r.GET("/acoes", controllers.GetBazinAndGraham)
	r.Run()
}
