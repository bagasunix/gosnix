package utils

import (
	"regexp"
	"strings"

	"github.com/bagasunix/gosnix/pkg/errors"
)

func ValidatePhone(phone string) (*string, error) {
	phoneRegex := regexp.MustCompile(`^(?:62|08)[0-9]{8,13}$`)
	if !phoneRegex.MatchString(phone) {
		return nil, errors.CustomError("Nomor telepon tidak valid")
	} else {
		if strings.HasPrefix(phone, "0") {
			return &phone, nil
		} else if strings.HasPrefix(phone, "8") {
			// Menambahkan 0 di depan jika digit pertama adalah 8
			phone = "0" + phone
		} else if strings.HasPrefix(phone, "62") {
			// Menambahkan 0 di depan jika digit pertama adalah 6 dan kedua bukan 2
			if phone[2] != '8' {
				return nil, errors.CustomError("Nomor telepon tidak valid")
			}
			phone = "0" + phone[2:]
		} else {
			return nil, errors.CustomError("Nomor telpone tidak valid")
		}
	}
	return &phone, nil
}
