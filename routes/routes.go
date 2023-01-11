package routes

import (
	"redistest/handlers"

	"github.com/gin-gonic/gin"
)

func Routes() {
	r := gin.Default()
	r.GET("/blocks", handlers.GetAllBlock)
	r.GET("/blocks/:id", handlers.GetBlockID)
	r.GET("/blocks/tree/:id", handlers.GetTreeID)
	r.POST("/blocks/", handlers.CreateBlock)
	r.PUT("/block/:id", handlers.PutBlock)
	r.DELETE("/blocks/:id", handlers.DeleteBlockID)
	r.Run()
}
