package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func createMessage(c *gin.Context) []byte {
	// Retrieve relevant request data from the gin-gonic context
	method := c.Request.Method
	path := c.Request.URL.Path
	nonce := c.Request.Header.Get("Nonce")
	queryParams := c.Request.URL.RawQuery
	requestBody, _ := c.GetRawData() // Assuming gin-gonic context provides a method to get the raw request body

	// Construct the message by concatenating relevant request data
	message := []byte(fmt.Sprintf("%s%s%s%s%s", method, path, nonce, queryParams, requestBody))
	return message
}
