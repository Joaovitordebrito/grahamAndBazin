package controllers

import (
	"graham-bazin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBazin(c *gin.Context) {
	fundo := c.Query("fundo")
	response, err := service.GetBazin(fundo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func GetBazinAndGraham(c *gin.Context) {
	acao := c.Query("acao")
	response, err := service.GetBazinAndGraham(acao)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
