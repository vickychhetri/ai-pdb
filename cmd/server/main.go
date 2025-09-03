/*
This is a Go web backend for your AI-PDB (Powered Document Brain) project.
It lets you store documents (title + content)
# It lets you list all stored documents
# Storage is in-memory (RAM only, disappears when server restarts) — fast for MVP testing
Later, we’ll plug in LLM features (summarization, search, Q&A)
*/

package main

import (
	"ai-pdb/config"
	"ai-pdb/internal/api"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitDB()

	// Make sure uploads folder exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}

	r := gin.Default()

	// Register routes
	api.RegisterRoutes(r)

	log.Println("🚀 AI-PDB server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
