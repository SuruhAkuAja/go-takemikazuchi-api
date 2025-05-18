package helper

import (
	"fmt"
	"go-takemikazuchi-api/pkg/exception"
	"math/rand"
	"reflect"
	"strconv"
)

func CheckErrorOperation(indicatedError error, clientError *exception.ClientError) bool {
	if indicatedError != nil {
		panic(clientError)
		return true
	}
	return false
}

func ParseNullableValue(value interface{}) string {
	v := reflect.ValueOf(value)

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "" // Jika nil, kembalikan string kosong
		}
		v = v.Elem()
	}

	// Check for sql.Null* types
	if v.Kind() == reflect.Struct {
		validField := v.FieldByName("Valid")
		valueField := v.FieldByName("String") // Sesuaikan jika nullable bukan string

		if validField.IsValid() && validField.Kind() == reflect.Bool {
			if validField.Bool() {
				return valueField.String() // Ambil nilai sebagai string
			}
			return ""
		}
	}

	// Jika bukan tipe nullable, langsung ubah menjadi string
	if v.IsValid() && v.CanInterface() {
		return fmt.Sprintf("%v", v.Interface())
	}

	return "" // Jika tidak valid, kembalikan string kosong
}

func GenerateOneTimePasswordToken() string {
	num := rand.Intn(9000) + 1000
	return strconv.Itoa(num)
}

func ConvertStringToUint64(id string) (uint64, error) {
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to uint64: %w", err)
	}
	return idUint64, nil
}
