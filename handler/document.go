package handler

import (
	"ai-pdb/config"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	// "github.com/pdfcpu/pdfcpu/pkg/api"
	// "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Document represents metadata about an uploaded file
// type Document struct {
// 	ID       int       `json:"id"`
// 	Name     string    `json:"name"`
// 	Path     string    `json:"path"`
// 	Content  string    `json:"content"`
// 	PageSize int       `json:"page_size"`
// 	Uploaded time.Time `json:"uploaded"`
// }

// type Document struct {
// 	ID       int       `json:"id"`
// 	Name     string    `json:"name"`
// 	Path     string    `json:"path"`
// 	Content  string    `json:"content"`
// 	PageSize int       `json:"page_size"`
// 	Uploaded time.Time `json:"uploaded"`

// 	// PDF Metadata
// 	Title            string    `json:"title"`
// 	Author           string    `json:"author"`
// 	Subject          string    `json:"subject"`
// 	Keywords         string    `json:"keywords"`
// 	Creator          string    `json:"creator"`
// 	Producer         string    `json:"producer"`
// 	CreationDate     time.Time `json:"creation_date"`
// 	ModificationDate time.Time `json:"modification_date"`

// 	// Technical Properties
// 	PDFVersion  string `json:"pdf_version"`
// 	IsEncrypted bool   `json:"is_encrypted"`
// 	IsTagged    bool   `json:"is_tagged"` // PDF/A compliance
// 	PageLayout  string `json:"page_layout"`
// 	PageMode    string `json:"page_mode"`

// 	// Content Statistics
// 	WordCount int             `json:"word_count"`
// 	CharCount int             `json:"char_count"`
// 	PageSizes []PageDimension `json:"page_sizes"` // Individual page dimensions
// 	HasImages bool            `json:"has_images"`
// 	HasForms  bool            `json:"has_forms"`
// 	HasLinks  bool            `json:"has_links"`

// 	// Extraction Details
// 	ExtractionDuration time.Duration `json:"extraction_duration"`
// 	ExtractionSuccess  bool          `json:"extraction_success"`
// 	ExtractionErrors   []string      `json:"extraction_errors"`

// 	// File Information
// 	FileSize int64  `json:"file_size"`
// 	FileHash string `json:"file_hash"` // For duplicate detection
// 	MimeType string `json:"mime_type"`
// }

// type PageDimension struct {
// 	PageNumber int     `json:"page_number"`
// 	Width      float64 `json:"width"`
// 	Height     float64 `json:"height"`
// 	Rotation   int     `json:"rotation"`
// 	Units      string  `json:"units"` // points, inches, mm, etc.
// }

// type ExtractionResult struct {
// 	Content        string
// 	PageCount      int
// 	Metadata       map[string]string
// 	PageDimensions []PageDimension
// 	HasImages      bool
// 	HasForms       bool
// 	HasLinks       bool
// 	IsEncrypted    bool
// 	PDFVersion     string
// }

type Document struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Content  string    `json:"content"`
	PageSize int       `json:"page_size"`
	Uploaded time.Time `json:"uploaded"`

	// PDF Metadata (limited with ledongthuc/pdf)
	Title    string `json:"title"`
	Author   string `json:"author"`
	Subject  string `json:"subject"`
	Keywords string `json:"keywords"`
	Creator  string `json:"creator"`
	Producer string `json:"producer"`

	// Technical Properties
	IsEncrypted bool   `json:"is_encrypted"`
	PDFVersion  string `json:"pdf_version"`

	// Content Statistics
	WordCount int  `json:"word_count"`
	CharCount int  `json:"char_count"`
	HasForms  bool `json:"has_forms"` // Basic detection

	// Extraction Details
	ExtractionDuration time.Duration `json:"extraction_duration"`
	ExtractionSuccess  bool          `json:"extraction_success"`
	ExtractionErrors   []string      `json:"extraction_errors"`

	// File Information
	FileSize int64  `json:"file_size"`
	FileHash string `json:"file_hash"`
	MimeType string `json:"mime_type"`
}

type ExtractionResult struct {
	Content     string
	PageCount   int
	Metadata    map[string]interface{}
	IsEncrypted bool
	PDFVersion  string
	HasForms    bool
	Duration    time.Duration
	Errors      []string
}

func extractTextFromPDF(filePath string) (ExtractionResult, error) {
	result := ExtractionResult{
		Metadata: make(map[string]interface{}),
		Errors:   []string{},
	}

	startTime := time.Now()

	// Open PDF file
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("could not open PDF: %v", err)
	}
	defer file.Close()

	// Basic PDF info
	result.PageCount = reader.NumPage()

	// Get PDF version from the reader's trailer
	if trailer := reader.Trailer(); !trailer.IsNull() {
		// Try to extract basic metadata
		if info := trailer.Key("Info"); !info.IsNull() {
			extractMetadata(info, result.Metadata)
		}

		// Check encryption
		if encrypt := trailer.Key("Encrypt"); !encrypt.IsNull() {
			result.IsEncrypted = true
		}
	}

	// Extract content
	text := ""
	for pageIndex := 1; pageIndex <= result.PageCount; pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			result.Errors = append(result.Errors, fmt.Sprintf("Page %d is null", pageIndex))
			continue
		}

		// Extract text content
		content, err := page.GetPlainText(nil)
		if err != nil {
			errMsg := fmt.Sprintf("Page %d: %v", pageIndex, err)
			result.Errors = append(result.Errors, errMsg)
			log.Printf("Warning: %s", errMsg)
			continue
		}
		text += content + "\n"

		if !result.HasForms {
			result.HasForms = detectForms(page)
		}
	}

	result.Content = text
	result.Duration = time.Since(startTime)

	return result, nil
}

// Helper function to extract metadata from PDF Info dictionary
func extractMetadata(info pdf.Value, metadata map[string]interface{}) {
	// Common metadata keys in PDF
	keys := []string{"Title", "Author", "Subject", "Keywords", "Creator", "Producer", "CreationDate", "ModDate"}

	for _, key := range keys {
		if value := info.Key(key); !value.IsNull() {
			if str := value.String(); str != "" {
				metadata[key] = str
			}
		}
	}
}

// Basic form detection by looking for Annots and Form elements
func detectForms(page pdf.Page) bool {
	// Check for annotations
	annots := page.V.Key("Annots")
	if !annots.IsNull() {
		return true
	}

	// Check for form fields in resources
	resources := page.V.Key("Resources")
	if resources.IsNull() {
		return false
	}

	// Look for AcroForm or XFA forms
	acroForm := resources.Key("AcroForm")
	if !acroForm.IsNull() {
		return true
	}

	xfa := resources.Key("XFA")
	return !xfa.IsNull()
}

func processPDFFile(filePath string) (Document, error) {
	doc := Document{
		Path:     filePath,
		Name:     filepath.Base(filePath),
		Uploaded: time.Now(),
	}

	// Get file info for size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return doc, fmt.Errorf("could not get file info: %v", err)
	}
	doc.FileSize = fileInfo.Size()

	// Extract content and metadata
	result, err := extractTextFromPDF(filePath)
	if err != nil {
		doc.ExtractionSuccess = false
		doc.ExtractionErrors = []string{err.Error()}
		return doc, err
	}

	// Populate document fields
	doc.Content = result.Content
	doc.PageSize = result.PageCount
	doc.ExtractionSuccess = true
	doc.ExtractionDuration = result.Duration
	doc.ExtractionErrors = result.Errors
	doc.IsEncrypted = result.IsEncrypted
	doc.HasForms = result.HasForms

	// Extract metadata from result
	if title, ok := result.Metadata["Title"].(string); ok {
		doc.Title = title
	}
	if author, ok := result.Metadata["Author"].(string); ok {
		doc.Author = author
	}
	if subject, ok := result.Metadata["Subject"].(string); ok {
		doc.Subject = subject
	}
	if keywords, ok := result.Metadata["Keywords"].(string); ok {
		doc.Keywords = keywords
	}
	if creator, ok := result.Metadata["Creator"].(string); ok {
		doc.Creator = creator
	}
	if producer, ok := result.Metadata["Producer"].(string); ok {
		doc.Producer = producer
	}

	// Content statistics
	doc.WordCount = countWords(result.Content)
	doc.CharCount = len(result.Content)

	// Generate file hash
	if hash, err := calculateFileHash(filePath); err == nil {
		doc.FileHash = hash
	}

	doc.MimeType = "application/pdf"

	return doc, nil
}

// Helper functions
func countWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// func extractTextFromPDF(filePath string) (string, int, error) {
// 	file, reader, err := pdf.Open(filePath)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("could not open PDF: %v", err)
// 	}
// 	// Ensure the file is closed after processing
// 	defer file.Close()

// 	text := ""
// 	totalPage := reader.NumPage()

// 	// Iterate through each page
// 	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
// 		page := reader.Page(pageIndex)
// 		if page.V.IsNull() {
// 			continue
// 		}
// 		content, err := page.GetPlainText(nil)
// 		if err != nil {
// 			log.Printf("Warning: could not extract text from page %d: %v", pageIndex, err)
// 			continue
// 		}
// 		text += content + "\n"
// 	}
// 	return text, totalPage, nil
// }

// UploadDocument handles file uploads and stores metadata + extracted content
func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	os.MkdirAll("./uploads", os.ModePerm)
	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	var doc Document
	if strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		doc, err = processPDFFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract PDF content"})
			return
		}
	}

	// doc := Document{
	// 	Name:     file.Filename,
	// 	Path:     filePath,
	// 	Content:  content,
	// 	PageSize: size,
	// 	Uploaded: time.Now(),
	// }

	// Insert into PostgreSQL using pgx
	// err = config.Db.QueryRow(
	// 	context.Background(),
	// 	`INSERT INTO documents (name, path, content, uploaded, page_size) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
	// 	doc.Name, doc.Path, doc.Content, doc.Uploaded, doc.PageSize,
	// ).Scan(&doc.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document to DB"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": " File uploaded and stored successfully",
	// 	"doc":     doc,
	// })

	err = config.Db.QueryRow(
		context.Background(),
		`INSERT INTO documents (
        name, path, content, uploaded, page_size,
        title, author, subject, keywords, creator, producer,
        is_encrypted, pdf_version,
        word_count, char_count, has_forms,
        extraction_duration, extraction_success, extraction_errors,
        file_size, file_hash, mime_type
    ) VALUES (
        $1, $2, $3, $4, $5,
        $6, $7, $8, $9, $10, $11,
        $12, $13,
        $14, $15, $16, $17,
        $18, $19, $20,
        $21, $22
    ) RETURNING id`,
		doc.Name, doc.Path, doc.Content, doc.Uploaded, doc.PageSize,
		doc.Title, doc.Author, doc.Subject, doc.Keywords, doc.Creator, doc.Producer,
		doc.IsEncrypted, doc.PDFVersion,
		doc.WordCount, doc.CharCount, doc.HasForms,
		doc.ExtractionDuration, doc.ExtractionSuccess, doc.ExtractionErrors,
		doc.FileSize, doc.FileHash, doc.MimeType,
	).Scan(&doc.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document to DB: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and stored successfully",
		"doc":     doc,
	})
}

// GetDocument serves a file by name
func GetDocument(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file name"})
		return
	}

	filePath := "./uploads/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
