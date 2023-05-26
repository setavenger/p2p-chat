package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"log"
	"net/http"
	"p2p-chat/common"
	"strconv"
)

type Daemon struct {
	DBPath string
	DB     *gorm.DB
	Users  map[string]*common.UserWellKnown
}

func (d *Daemon) GetUserPubKey(username string) *common.UserWellKnown {
	data := d.Users[username]
	fmt.Println(data)
	return data
}
func (d *Daemon) LoadUsers(path string) error {
	usersMap := make(map[string]*common.UserWellKnown)
	users, err := common.LoadConfig(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, user := range users {
		usersMap[user.Username] = &user
	}
	d.Users = usersMap
	return err
}

func (d *Daemon) Authentication(c *gin.Context) {
	publicKey := c.GetHeader("Public-Key")

	// check nonce
	nonceForeignStr := c.Request.Header.Get("Nonce")
	nonceForeign, err := strconv.ParseUint(nonceForeignStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	nonce, err := RetrieveNonce(d.DB, publicKey)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if nonceForeign <= nonce.Nonce {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nonce"})
		c.Abort()
		return
	}
	// now that nonce is accepted authenticity was cleared, overwrite to new nonce state

	message := createMessage(c)
	err = common.VerifySignature(message, c.GetHeader("Signature"), publicKey)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err = UpdateNonceCounter(d.DB, publicKey, nonceForeign)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
}

// ForwardMessage - Endpoint to forward the message
func (d *Daemon) ForwardMessage(c *gin.Context) {
	var message common.MessageServer
	if err := c.ShouldBindJSON(&message); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward the message to the designated server here
	// todo
	//

	err := SaveEntry(d.DB, &message)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Forwarding message: %+v", message)
	c.Status(http.StatusOK)
}

// MarkMessageAsRead - Endpoint to mark a message as read
func (d *Daemon) MarkMessageAsRead(c *gin.Context) {
	// todo
	c.Status(http.StatusOK)
}

// GetUnreadMessages - Endpoint to fetch all unread messages
func (d *Daemon) GetUnreadMessages(c *gin.Context) {
	publicKey := c.Param("public-key")

	messages, err := retrieveByRecipient(d.DB, publicKey, false)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// GetAllMessages - Endpoint to fetch all unread messages
func (d *Daemon) GetAllMessages(c *gin.Context) {
	publicKey := c.Param("public-key")

	messages, err := retrieveAllByRecipient(d.DB, publicKey)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// GetReadMessages - Endpoint to fetch all read messages
func (d *Daemon) GetReadMessages(c *gin.Context) {
	publicKey := c.Param("public-key")

	messages, err := retrieveByRecipient(d.DB, publicKey, true)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// GetUsernameWellKnown - Endpoint to show the users public key information
func (d *Daemon) GetUsernameWellKnown(c *gin.Context) {
	username := c.Param("username")

	user, err := retrieveUsersWellKnown(d.DB, username)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//user := d.GetUserPubKey(username)
	c.JSON(http.StatusOK, user)
}

func runServer() {

	dbPath := "./data/data.db"

	err := Migrate(dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//db, err := CreateDB(dbPath)
	db, err := ConnectToPostgres()
	if err != nil {
		fmt.Println(err)
		return
	}
	daemon := Daemon{DBPath: dbPath, DB: db}
	err = daemon.LoadUsers("./data/users.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	router := gin.Default()

	// well-known
	router.GET("/.well-known/p2pchat/:username", daemon.GetUsernameWellKnown)

	// public
	router.POST("/api/forward", daemon.ForwardMessage)

	// "private" route
	privateGroup := router.Group("/api", daemon.Authentication)
	privateGroup.PUT("/messages/:messageid/read", daemon.MarkMessageAsRead)
	privateGroup.PUT("/messages/:messageid/unread")
	privateGroup.GET("/users/:public-key/messages", daemon.GetAllMessages)
	privateGroup.GET("/users/:public-key/messages/read", daemon.GetReadMessages)
	privateGroup.GET("/users/:public-key/messages/unread", daemon.GetUnreadMessages)

	if err = router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
