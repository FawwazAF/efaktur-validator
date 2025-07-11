package efaktur

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	fiveMBSize = 5 << 20
)

func (h *Handler) HandlerValidateEfaktur(c *gin.Context) {
	var (
		req = c.Request
	)

	// limit to 5 mb.
	if err := req.ParseMultipartForm(fiveMBSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer req.MultipartForm.RemoveAll()

	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tempDir := fmt.Sprintf("/tmp/efaktur/%s", uuid.New().String())
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the file to a temporary location
	tempFilePath := filepath.Join(tempDir, fileHeader.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		tempFile.Close()
		os.RemoveAll(tempDir)
	}()

	// Copy the file contents to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.efaktur.ValidateEfaktur(req.Context(), tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
