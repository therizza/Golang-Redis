package handlers

import (
	"fmt"
	"net/http"
	"redistest/models"
	"strings"

	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
)

type GetBlockIDRequest struct {
	ID string `uri:"id" binding:"required"`
}
type GetTreeIDRequest struct {
	ID string `uri:"id" binding:"required"`
}
type DeleteBlockIDRequest struct {
	ID string `uri:"id" binding:"required"`
}
type CreateBlockRequest struct {
	ID       string           `uri:"block_id"`
	Name     string           `uri:"name"`
	ParentID string           `uri:"parent_id"`
	CentroID geojson.Geometry `uri:"centro_id"`
	Value    string           `uri:"value"`
}

func GetAllBlock(c *gin.Context) {
	value, err := models.BringAllBlocks()
	if err == models.ErrNotFound {
		c.JSON(http.StatusNotFound, value)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, false)

	}

	c.JSON(http.StatusOK, value)

}
func GetBlockID(c *gin.Context) {
	var req GetBlockIDRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	value, err := models.GetBlockID(strings.ToUpper(req.ID))
	if err == models.PassingWrong {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passing Wrong Parameter"})
		return
	}
	if err == models.ParameterNotFound {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if err == models.ParameterNotID {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, value)

}

func DeleteBlockID(c *gin.Context) {
	var req DeleteBlockIDRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	value := models.DeleteBlockID(strings.ToUpper(req.ID))

	fmt.Println(value)
	if value == models.ParameterNotFound {
		c.JSON(http.StatusBadRequest, gin.H{})
		return

	}
	if value == models.BlockTree {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Block that counted child cannot be deleted"})

	}
	if value != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func PutBlock(c *gin.Context) {
	var req CreateBlockRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	argBlock := models.Block{
		ID:       req.ID,
		Name:     req.Name,
		ParentID: req.ParentID,
		CentroID: req.CentroID,
		Value:    req.Value,
	}

	value, err := models.PutBlock(argBlock)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, value)
}

func CreateBlock(c *gin.Context) {
	var req CreateBlockRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	argBlock := models.Block{
		ID:       req.ID,
		Name:     req.Name,
		ParentID: req.ParentID,
		CentroID: req.CentroID,
		Value:    req.Value,
	}

	value, err := models.CreateBlock(argBlock)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, value)
}

func GetTreeID(c *gin.Context) {
	var req GetTreeIDRequest

	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	value, err := models.GetTreeID(strings.ToUpper(req.ID))

	if err == models.PassingWrong {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if err == models.ParameterNotID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passing Wrong Parameter"})
		return
	}
	if err == models.ParameterNotID {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, value)

}
