package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
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

func ValidateIRMobile(mobile string) bool {
	// الگوی بررسی شماره موبایل ایران
	re := regexp.MustCompile(`^09\d{9}$`)
	return len(mobile) == 11 && re.MatchString(mobile)
}

func Random4Digit() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(9000))
	return n.Int64() + 1000
}

func PrettyJson(v any) {
	jsonData, err := json.MarshalIndent(v, "", "  ") // تبدیل به JSON زیبا
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
	} else {
		// چاپ اطلاعات به‌صورت JSON زیبا
		fmt.Printf("~~~~~~~~~~~~~~~~~~~~ [pretty json ] : data:\n%s\n ~~~~~~~~~~~~~~~~~~~~", string(jsonData))
	}
}

// PrettyPrice show prices with commas
func PrettyPrice(price int) string {
	formattedNumber := strconv.FormatInt(int64(price), 10)

	return formatWithCommas(formattedNumber)
}

func formatWithCommas(input string) string {
	n := len(input)
	if n <= 3 {
		return input
	}
	return formatWithCommas(input[:n-3]) + "," + input[n-3:]
}

func ConvertIfNotNil[T any, R any](input *T, convertFunc func(*T) R) R {
	if input != nil {
		return convertFunc(input)
	}
	var zero R
	return zero
}

func ConvertSliceIfNotNil[T any, R any](input []*T, convertFunc func(*T) R) []R {
	if input == nil {
		return nil
	}

	var result []R
	for _, item := range input {
		result = append(result, convertFunc(item))
	}
	return result
}

// StringToUint convert string to uint function which is used in front-end
func StringToUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

// Trace will log file, line, function name, and message (error or string)
func Trace(input interface{}) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[Message]: %v (Unknown Location)\n", input)
		return
	}

	//get method name
	funcName := runtime.FuncForPC(pc).Name()

	//declare input type
	switch v := input.(type) {
	case error:
		log.Printf("[File]: %s | [Line]: %d | [Function]: %s | [Error]: %v\n", file, line, funcName, v)
	case string:
		log.Printf("[File]: %s | [Line]: %d | [Function]: %s | [Message]: %s\n", file, line, funcName, v)
	default:
		log.Printf("[File]: %s | [Line]: %d | [Function]: %s | [Unknown Type]: %v\n", file, line, funcName, v)
	}
}

// GetProductStoragePath if key value set to active then address will be bucket address else address will be server address
func GetProductStoragePath() string {
	if os.Getenv("STORAGE_STATUS") == "active" {
		return fmt.Sprintf("https://%s.parspack.net/uploads/media/products/", os.Getenv("STORAGE_BUCKET_NAME"))
	}

	return "/uploads/media/products/"

}
