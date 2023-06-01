package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Daemon struct {
	DB     *gorm.DB
	Client Client
}

type MetaData struct {
	PublicKey string `json:"public_key,omitempty"`
}

type SendMessageData struct {
	Recipient string `json:"recipient,omitempty"`
	Body      string `json:"body,omitempty"`
	ParentId  string `json:"parent_id,omitempty"`
}

type NewPrivateKeyBody struct {
	Key string `json:"key,omitempty"`
}

type NewHostBody struct {
	Host string `json:"host,omitempty"`
}

func (d *Daemon) GetMetaData(c *gin.Context) {
	c.JSON(http.StatusOK, MetaData{PublicKey: d.Client.PublicKey})
}

func (d *Daemon) SetHost(c *gin.Context) {
	var data NewHostBody
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Printf("Error parsing: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d.Client.BaseURL = data.Host
	c.Status(http.StatusOK)
}

func (d *Daemon) SetPrivateKey(c *gin.Context) {
	var data NewPrivateKeyBody
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Printf("Error parsing: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d.Client.SetPrivateKey(data.Key)
	c.Status(http.StatusOK)
}

func (d *Daemon) GetMessages(c *gin.Context) {
	messagesRaw, err := d.Client.GetEveryMessage()
	if err != nil {
		fmt.Printf("Error receiving messages: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messages := DecryptAllMessages(d.Client.PrivateKey, messagesRaw)
	//fmt.Printf("%+v", messages)
	c.JSON(http.StatusOK, messages)

	return
}

func (d *Daemon) GetUnreadMessages(c *gin.Context) {
	messagesRaw, err := d.Client.GetUnreadMessages()
	if err != nil {
		fmt.Printf("Error receiving messages: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	messages := DecryptAllMessages(d.Client.PrivateKey, messagesRaw)

	c.JSON(http.StatusOK, messages)

	return
}

func (d *Daemon) SendMessage(c *gin.Context) {
	var data SendMessageData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Printf("Error parsing: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = SendMessageByUsername(d.Client, d.Client.PrivateKey, d.Client.PublicKey, data.Recipient, data.Body, data.ParentId)
	if err != nil {
		fmt.Printf("Error sending message: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
	return
}

func runServer() {

	daemon := Daemon{Client: Client{
		BaseURL: "http://localhost:8889",
		//PrivateKey: "b5ecbb76d605b0d9025bf7cdd830bf9c01a0a1967d89462aa4016d7fe897f63e",
	}}
	//daemon.Client.SetPublicKey()

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}   // Add the allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"} // Add the allowed HTTP methods
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}          // Add the allowed headers
	router.Use(cors.New(config))

	// "public" to local machine
	router.POST("/send", daemon.SendMessage)
	router.POST("/set-key", daemon.SetPrivateKey)
	router.POST("/set-host", daemon.SetHost)

	router.GET("/meta", daemon.GetMetaData)

	router.GET("/messages", daemon.GetMessages)
	router.GET("/messages/unread", daemon.GetUnreadMessages)

	if err := router.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
