package providers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func FormatDateToString(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("2006-01-02")
}

func FormatStringToDate(date *string) *time.Time {
	if date == nil || *date == "" {
		return nil
	}
	parsedDate, err := time.Parse("2006-01-02", *date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return nil
	}
	return &parsedDate
}

func CheckNullString(str *string) *string {
	if str == nil || *str == "" {
		return nil
	}
	return str
}

func IsValidEmail(email string) bool {
	// Simple regex for email validation
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' && (i == 0 || i == len(email)-1 || email[i-1] == '@') {
			return false
		}
	}
	return at == 1
}

func IsValidPhone(phone string) bool {
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
