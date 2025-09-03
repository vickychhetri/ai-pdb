package api

import (
	"ai-pdb/handler"
	"ai-pdb/internal/core"
	"ai-pdb/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

var store = storage.NewMemoryStore()

func RegisterRoutes(r *gin.Engine) {
	r.POST("/documents", uploadDocument)
	r.GET("/documents", listDocuments)
	r.POST("/documents/upload", handler.UploadDocument)
	r.GET("/documents/get", handler.GetDocument)

}

func uploadDocument(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := core.NewDocument(req.Title, req.Content)
	store.Save(doc)

	c.JSON(http.StatusCreated, gin.H{"id": doc.ID})
}

func listDocuments(c *gin.Context) {
	c.JSON(http.StatusOK, store.List())
}
