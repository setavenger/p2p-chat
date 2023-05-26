package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID             string `json:"id"`
	Sender         string `json:"sender"`
	SenderUsername string `json:"sender_username"`
	Recipient      string `json:"recipient"`
	Encrypted      string `json:"encrypted"`
	Signature      string `json:"signature"`
	Timestamp      uint64 `json:"timestamp"`
	Read           bool   `json:"read,omitempty"`
}

type MessageServer struct {
	ID             string         `gorm:"column:id" json:"id,omitempty"`
	Sender         string         `gorm:"column:sender;index" json:"sender,omitempty"`
	SenderUsername string         `gorm:"sender_username;index" json:"sender_username,omitempty"`
	Recipient      string         `gorm:"column:recipient;index" json:"recipient,omitempty"`
	Encrypted      string         `gorm:"column:encrypted" json:"encrypted,omitempty"`
	Signature      string         `gorm:"column:signature" json:"signature,omitempty"`
	Timestamp      uint64         `gorm:"column:timestamp;index" json:"timestamp,omitempty"`
	Read           bool           `gorm:"column:read;index" json:"read,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type MessagePlain struct {
	ID             string `json:"id"`
	Sender         string `json:"sender"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
	Recipient      string `json:"recipient"`
	Timestamp      uint64 `json:"timestamp"`
	Read           bool   `json:"read"`
}

func (m *Message) GetId() error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(bytes)
	m.ID = hex.EncodeToString(hashed[:])
	return err
}
