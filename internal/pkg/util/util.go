package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateFilename year/month/day/uuid-string.extension =>e.g (2024/04/30/078998-897-546.jgp)
func GenerateFilename(originalFilename string) string {
	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())
	random := uuid.New().String()

	extension := filepath.Ext(originalFilename)

	return fmt.Sprintf("%s/%s/%s/%s%s", year, month, day, random, extension)
}

func checkNationalCode(code string) bool {

	reg, err := regexp.Compile("/[^0-9]/")

	if err != nil {
		log.Fatal(err)
	}

	code = reg.ReplaceAllString(code, "")
	if len(code) != 10 {
		return false
	}

	codes := strings.Split(code, "")
	last, err := strconv.Atoi(codes[9])

	i := 10
	sum := 0

	for in, el := range codes {
		temp, err := strconv.Atoi(el)

		if err != nil {
			log.Fatal(err)
		}

		if in == 9 {
			break
		}

		sum += temp * i
		i -= 1
	}

	mod := sum % 11
	if mod >= 2 {
		mod = 11 - mod
	}
	return mod == last
}

func GenerateCSRFToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(token)
}

func GeneratePageNumbers(currentPage int, totalPages int) []int {
	var pages []int

	// Calculate the number of pages to display
	numPagesToShow := int(math.Min(float64(totalPages), 4)) // Show up to 5 pages
	startPage := int(math.Max(float64(currentPage-2), 1))
	endPage := int(math.Min(float64(startPage+numPagesToShow-1), float64(totalPages)))

	// Generate page numbers
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	return pages
}

// StructToMap generate a map from a struct which key is struct field name
func StructToMap(yourStruct interface{}) map[string]interface{} {
	val := reflect.ValueOf(yourStruct)
	structMap := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		structMap[field.Name] = val.Field(i).Interface()
	}

	return structMap
}

func AllowImageExtensions() []string {
	return []string{".jpg", ".png", ".jpeg"}
}
