package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type User struct {
	Model
	Fullname     string   `json:"fullname"`
	Username     string   `json:"username"`
	Email        string   `json:"email" gorm:"unique"`
	Password     []byte   `json:"-"`
	Isambassador bool     `json:"-"`
	Revenue      *float64 `json:"revenue,omitempty" gorm:"-"`
}

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func (user *User) SetPassword(password string) {
	salt := generateSalt()

	hashedPassword := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hashedPassword)

	user.Password = []byte(fmt.Sprintf("$argon2id$v=19$m=65536,t=3,p=4$%s$%s", encodedSalt, encodedHash))
}

func (user *User) ComparePassword(inputPassword string) bool {
	storedHash := string(user.Password)
	parts := strings.Split(storedHash, "$")
	if len(parts) != 6 {
		return false
	}

	encodedSalt := parts[4]
	encodedStoredHash := parts[5]

	salt, err := base64.StdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false
	}

	newHash := argon2.IDKey([]byte(inputPassword), salt, 3, 64*1024, 4, 32)

	return base64.StdEncoding.EncodeToString(newHash) == encodedStoredHash
}

type Admin User

func (admin *Admin) CalculateRevenue(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   admin.Id,
		Complete: true,
	})

	var revenue float64 = 0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AdminRevenue
		}
	}

	admin.Revenue = &revenue
}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenue(db *gorm.DB) {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, &Order{
		UserId:   ambassador.Id,
		Complete: true,
	})

	var revenue float64 = 0.0

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}

	ambassador.Revenue = &revenue
}
