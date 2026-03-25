package validator

import (
	"codebase-app/internal/adapter"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	// "github.com/go-playground/locales/en"
	// ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	}
	if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
		log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	}
	if err := v.RegisterValidation("exist", isExist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering exist validator")
	}
	if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
		log.Fatal().Err(err).Msg("Error while registering unique validator")
	}
	if err := v.RegisterValidation("uuidv7", uuidv7); err != nil {
		log.Fatal().Err(err).Msg("Error while registering uuidv7 validator")
	}

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

// blacklist email validator
func isEmailBlacklist(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	disallowedDomains := []string{"gmail", "yahoo", "outlook", "hotmail", "aol", "live", "inbox", "icloud", "mail", "gmx", "yandex"}

	for _, domain := range disallowedDomains {
		if strings.Contains(email, domain) {
			return false
		}
	}

	return true
}

func isStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= 'a' && char <= 'z':
			hasLowercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUppercase && hasLowercase && hasNumber
}

var allowedExistsQueries = map[string]string{
	"users.id":     "SELECT id FROM users WHERE id = $1",
	"roles.id":     "SELECT id FROM roles WHERE id = $1",
	"org_units.id": "SELECT id FROM org_units WHERE id = $1",
}

func isExist(fl validator.FieldLevel) bool {
	db := adapter.Adapters.Postgres
	fieldValue := fl.Field().String()
	tagValue := fl.Param()

	query, ok := allowedExistsQueries[tagValue]
	if !ok {
		return false
	}

	result := make(map[string]any)
	err := db.QueryRowx(query, fieldValue).MapScan(result)
	if err != nil {
		log.Warn().Err(err).Any("query", query).Msg("Error while querying the database")
		return false
	}

	return true
}

func isUniqueInSlice(fl validator.FieldLevel) bool {
	// Get the slice from the FieldLevel interface
	val := fl.Field()

	// Ensure the field is a slice
	if val.Kind() != reflect.Slice {
		return false
	}

	// Use a map to check for duplicates
	elements := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if _, found := elements[elem]; found {
			return false // Duplicate found
		}
		elements[elem] = true
	}
	return true
}

func uuidv7(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	return isUUIDv7(uuid)
}

func isUUIDv7(uuid string) bool {
	uuid = strings.ToLower(uuid)

	if len(uuid) != 36 {
		return false
	}
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}
	// versi nibble
	if uuid[14] != '7' {
		return false
	}
	// variant nibble (should be 8, 9, a, or b)
	if uuid[19] != '8' && uuid[19] != '9' && uuid[19] != 'a' && uuid[19] != 'b' {
		return false
	}
	// cek semua karakter selain '-' adalah hex
	for i, c := range uuid {
		if c == '-' {
			continue
		}
		if !unicode.IsDigit(c) && (c < 'a' || c > 'f') {
			fmt.Printf("invalid char at %d: %c\n", i, c)
			return false
		}
	}
	return true
}
