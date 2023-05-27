package main

import (
	"fmt"
	"github.com/setavenger/p2p-chat/common"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type NonceCounter struct {
	PublicKey string `gorm:"column:public_key;index"`
	Nonce     uint64 `gorm:"column:nonce"`
}

func ConnectToPostgres() (db *gorm.DB, err error) {
	host := "db" // The service name of the PostgreSQL container defined in the docker-compose.yml file
	port := 5432 // Default PostgreSQL port
	user := "main"
	dbname := "p2p"
	password := ".znb6PF_yHWzCtbF6sYfNC_3CB!yoq"
	sslmode := "disable" // or "require" if SSL is enabled

	// Create connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)

	// Open a connection to the database
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func CreateDB(path string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

func Migrate(dbPath string) error {
	db, err := ConnectToPostgres()
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&common.MessageServer{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&NonceCounter{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&common.UserWellKnown{})
	if err != nil {
		return err
	}
	return nil
}

func UpdateNonceCounter(db *gorm.DB, publicKey string, newNonce uint64) error {
	conditions := NonceCounter{PublicKey: publicKey}
	updates := map[string]interface{}{
		"nonce": newNonce,
	}

	result := db.Where(conditions).Assign(updates).FirstOrCreate(&NonceCounter{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RetrieveNonce(db *gorm.DB, publicKey string) (*NonceCounter, error) {
	var nonceCounter NonceCounter
	result := db.Where("public_key = ?", publicKey).Find(&nonceCounter)
	if result.Error != nil {
		return nil, result.Error
	}
	return &nonceCounter, nil
}

func SaveEntry(db *gorm.DB, message *common.MessageServer) error {
	result := db.Create(message)

	if result.Error != nil {
		return result.Error
	}

	fmt.Println("persisted message to db")
	return nil
}

func retrieveAllByPublicKey(db *gorm.DB, publicKey string) ([]common.MessageServer, error) {
	var messages []common.MessageServer

	result := db.Where("recipient = ? OR sender = ?", publicKey, publicKey).Order("timestamp DESC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func retrieveAllByRecipient(db *gorm.DB, recipient string) ([]common.MessageServer, error) {
	var messages []common.MessageServer

	result := db.Where("recipient = ?", recipient).Order("timestamp DESC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func retrieveByRecipient(db *gorm.DB, recipient string, read bool) ([]common.MessageServer, error) {
	var messages []common.MessageServer

	result := db.Where("recipient = ? AND read = ?", recipient, read).Order("timestamp DESC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func retrieveBySender(db *gorm.DB, sender string, read bool) ([]common.MessageServer, error) {
	var messages []common.MessageServer

	result := db.Where("sender = ? AND read = ?", sender, read).Order("timestamp DESC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func retrieveByRecipientAndSender(db *gorm.DB, sender, recipient string, read bool) ([]common.MessageServer, error) {
	var messages []common.MessageServer

	result := db.Where("sender = ? AND recipient = ? AND read = ?", sender, recipient, read).Order("timestamp DESC").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

func retrieveUsersWellKnown(db *gorm.DB, username string) (*common.UserWellKnown, error) {
	var user common.UserWellKnown

	result := db.Where("username = ?", username).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
