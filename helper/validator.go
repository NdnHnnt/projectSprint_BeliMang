package helpers

import (
	// "fmt"

	"regexp"
	"strconv"
	"time"

	"github.com/NdnHnnt/projectSprint_BeliMang/db"
)

func ValidateNIP(nip int64) bool {
	nipStr := strconv.FormatInt(nip, 10)

	// Check length
	if len(nipStr) < 13 || len(nipStr) > 15 {
		return false
	}

	// Check first three digits
	if len(nipStr) < 3 || (nipStr[:3] != "615" && nipStr[:3] != "303") {
		return false
	}

	// Check fourth digit
	if len(nipStr) < 4 || (nipStr[3] != '1' && nipStr[3] != '2') {
		return false
	}

	// Check fifth to eighth digits (year)
	year, err := strconv.Atoi(nipStr[4:8])
	if err != nil || year < 2000 || year > time.Now().Year() {
		return false
	}

	// Check ninth and tenth digits (month)
	month, err := strconv.Atoi(nipStr[8:10])
	if err != nil || month < 1 || month > 12 {
		return false
	}

	// Check eleventh to thirteenth digits (random)
	random, err := strconv.Atoi(nipStr[10:])
	if err != nil || random < 0 || random > 99999 {
		return false
	}

	return true
}

func ValidateUsername(username string) bool {
	if len(username) < 5 || len(username) > 30 {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func ValidatePassword(password string) bool {
	if len(password) < 5 || len(password) > 30 {
		return false
	}
	return true
}

func ValidateURL(url string) bool {
	regex := `^(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,})(\/[a-zA-Z0-9-._~:/?#@!$&'()*+,;=%]*)?$`
	return regexp.MustCompile(regex).MatchString(url)
}

func ValidatePhoneNumber(number string) bool {
	if len(number) < 10 || len(number) > 15 {
		return false
	}
	regex := `^\+62[0-9]{7,12}$`
	return regexp.MustCompile(regex).MatchString(number)
}

func ValidateMerchantName(name string) bool {
	if len(name) < 2 || len(name) > 30 {
		return false
	}
	return true
}

func ValidateMerchantCategory(category string) bool {
	validCategories := []string{
		"SmallRestaurant",
		"MediumRestaurant",
		"LargeRestaurant",
		"MerchandiseRestaurant",
		"BoothKiosk",
		"ConvenienceStore",
	}

	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}

	return false
}

func ValidateMerchantItem(itemName string) bool {
	validCategories := []string{
		"Beverage",
		"Food",
		"Snack",
		"Condiments",
		"Additions",
	}

	for _, validCategory := range validCategories {
		if itemName == validCategory {
			return true
		}
	}

	return false
}

func ValidateLocation(lat, lon float64) bool {
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return false
	}
	return true
}

func ValidateAdmin(email string, id string) (bool, error) {
	conn := db.CreateConn()
	var exists bool
	err := conn.QueryRow("SELECT EXISTS(SELECT 1 FROM admin WHERE email=$1 AND id=$2)", email, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func ValidateUser(email string, id string) (bool, error) {
	conn := db.CreateConn()
	var exists bool
	err := conn.QueryRow("SELECT EXISTS(SELECT 1 FROM public.user WHERE email=$1 AND id=$2)", email, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
