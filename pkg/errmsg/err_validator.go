package errmsg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func errorValidationHandler[T any](lang Language, err error, payload *T) (int, map[string][]string) {
	var (
		errorMessages = make(map[string][]string)
		code          = 400
	)

	for _, err := range err.(validator.ValidationErrors) {
		var (
			// Get the JSON tag name
			namespace  = err.Namespace()               // ex: UpdateInterestRequest.interest
			fieldParts = strings.Split(namespace, ".") // ex: [UpdateInterestRequest, interest]

			field      string
			fieldInMsg string
			message    string

			value     = err.Value()
			valueType = reflect.TypeOf(value)

			// Get the error message
		)
		lastField := fieldParts[len(fieldParts)-1]                                // get the last element
		fieldParts = fieldParts[1:]                                               // remove the first element
		field = strings.Join(fieldParts, ".")                                     // join the rest of the elements
		if strings.Contains(lastField, "_") && strings.Contains(lastField, "]") { // check if the last element contains "_" and "]", ex: interested_in[0]
			// fieldInMsg = field
			// remove characters between "[" and "]"
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
		} else {
			fieldInMsg = strings.ReplaceAll(lastField, "_", " ")
			if strings.Contains(fieldInMsg, "[") {
				fieldInMsg = fieldInMsg[:strings.Index(fieldInMsg, "[")] // remove characters after "[" ("interested_in[0]" => "interested_in")
			}
		}

		message = defaultValidationMessage(lang, field, fieldInMsg, err.Tag(), err.Param())

		// get validate tag that causes the error
		switch err.Tag() {
		case "required":
			message = validationMessageRequired(lang, fieldInMsg)
		case "email":
			message = validationMessageEmail(lang, fieldInMsg)
		case "email_blacklist":
			message = validationMessageEmailBlacklist(lang, value)
		case "strong_password":
			message = validationMessageStrongPassword(lang, fieldInMsg)
		case "exist":
			message = validationMessageResourceNotExist(lang)
		case "datetime":
			message = validationMessageDatetime(lang, fieldInMsg, err.Param())
		case "ulid":
			message = validationMessageSimpleInvalid(lang, fieldInMsg, "ULID")
		case "base64":
			message = validationMessageSimpleFormat(lang, fieldInMsg, "base64")
		case "base64url":
			message = validationMessageSimpleFormat(lang, fieldInMsg, "base64url")
		case "base64rawurl":
			message = validationMessageSimpleFormat(lang, fieldInMsg, "base64rawurl")
		case "min":
			// check if the field is a number or a string
			if valueType.Kind() == reflect.Int || valueType.Kind() == reflect.Int8 || valueType.Kind() == reflect.Int16 || valueType.Kind() == reflect.Int32 || valueType.Kind() == reflect.Int64 || valueType.Kind() == reflect.Float32 || valueType.Kind() == reflect.Float64 {
				message = validationMessageMinValue(lang, fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.String {
				message = validationMessageMinChars(lang, fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				message = validationMessageMinItems(lang, fieldInMsg, err.Param())
			}
		case "max":
			// check if the field is a number or a string
			if _, ok := value.(int); ok {
				message = validationMessageMaxValue(lang, fieldInMsg, err.Param())
			}
			if _, ok := value.(float64); ok {
				message = validationMessageMaxValue(lang, fieldInMsg, err.Param())
			}
			if _, ok := value.(string); ok {
				message = validationMessageMaxChars(lang, fieldInMsg, err.Param())
			}
			if valueType.Kind() == reflect.Slice {
				message = validationMessageMaxItems(lang, fieldInMsg, err.Param())
			}
		case "gt":
			message = validationMessageCompare(lang, fieldInMsg, "gt", err.Param())
		case "gte":
			message = validationMessageCompare(lang, fieldInMsg, "gte", err.Param())
		case "lt":
			message = validationMessageCompare(lang, fieldInMsg, "lt", err.Param())
		case "lte":
			message = validationMessageCompare(lang, fieldInMsg, "lte", err.Param())
		case "latitude":
			message = validationMessageSimpleInvalid(lang, fieldInMsg, "latitude")
		case "longitude":
			message = validationMessageSimpleInvalid(lang, fieldInMsg, "longitude")
		case "numeric":
			message = validationMessageNumeric(lang, fieldInMsg)
		case "timezone":
			message = validationMessageTimezone(lang, fieldInMsg)
		case "eqfield":
			eqField := err.Param()
			eqFieldName := ""
			eqFieldTag, _ := reflect.TypeOf(payload).Elem().FieldByName(eqField)
			eqFieldJSONTag := eqFieldTag.Tag.Get("json")
			eqFieldQueryTag := eqFieldTag.Tag.Get("query")
			eqFieldFormTag := eqFieldTag.Tag.Get("form")
			eqFieldParamsTag := eqFieldTag.Tag.Get("params")

			if eqFieldJSONTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldJSONTag, "_", " ")
			}
			if eqFieldQueryTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldQueryTag, "_", " ")
			}
			if eqFieldFormTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldFormTag, "_", " ")
			}
			if eqFieldParamsTag != "" {
				eqFieldName = strings.ReplaceAll(eqFieldParamsTag, "_", " ")
			}

			message = validationMessageEqualField(lang, fieldInMsg, eqFieldName)
		case "oneof":
			message = validationMessageOneOf(lang, fieldInMsg, err.Param())
		case "unique_in_slice":
			message = validationMessageUniqueInSlice(lang, fieldInMsg)
		case "url":
			message = validationMessageSimpleInvalid(lang, fieldInMsg, "URL")
		}

		errorMessages[field] = append(errorMessages[field], message)
	}

	return code, errorMessages
}

func defaultValidationMessage(lang Language, field string, fieldInMsg string, tag string, param string) string {
	if param != "" {
		if lang == LanguageEnglish {
			return fmt.Sprintf("field validation for '%s' failed on the '%s' tag with param '%s'", field, tag, param)
		}
		return fmt.Sprintf("validasi untuk '%s' gagal pada tag '%s' dengan parameter '%s'", fieldInMsg, tag, param)
	}

	if lang == LanguageEnglish {
		return fmt.Sprintf("field validation for '%s' failed on the '%s' tag", field, tag)
	}
	return fmt.Sprintf("validasi untuk '%s' gagal pada tag '%s'", fieldInMsg, tag)
}

func validationMessageRequired(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s is required", field)
	}
	return fmt.Sprintf("%s wajib diisi", field)
}

func validationMessageEmail(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s is not a valid email address", field)
	}
	return fmt.Sprintf("%s bukan alamat email yang valid", field)
}

func validationMessageEmailBlacklist(lang Language, value any) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("email %v is not allowed", value)
	}
	return fmt.Sprintf("email %v tidak diizinkan", value)
}

func validationMessageStrongPassword(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be at least 8 characters and contain at least one uppercase letter, one lowercase letter, and one number", field)
	}
	return fmt.Sprintf("%s minimal 8 karakter dan harus mengandung setidaknya satu huruf besar, satu huruf kecil, dan satu angka", field)
}

func validationMessageResourceNotExist(lang Language) string {
	if lang == LanguageEnglish {
		return "resource does not exist."
	}
	return "sumber data tidak ditemukan."
}

func validationMessageDatetime(lang Language, field string, example string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s is not a valid datetime format (Ex: %s)", field, example)
	}
	return fmt.Sprintf("%s bukan format tanggal dan waktu yang valid (Contoh: %s)", field, example)
}

func validationMessageSimpleInvalid(lang Language, field string, kind string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s is not a valid %s", field, kind)
	}
	return fmt.Sprintf("%s bukan %s yang valid", field, kind)
}

func validationMessageSimpleFormat(lang Language, field string, format string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s is not a valid %s format", field, format)
	}
	return fmt.Sprintf("%s bukan format %s yang valid", field, format)
}

func validationMessageMinValue(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be at least %s", field, param)
	}
	return fmt.Sprintf("%s harus minimal %s", field, param)
}

func validationMessageMinChars(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	}
	return fmt.Sprintf("%s harus minimal %s karakter", field, param)
}

func validationMessageMinItems(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must have at least %s items", field, param)
	}
	return fmt.Sprintf("%s harus minimal %s item", field, param)
}

func validationMessageMaxValue(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must not be greater than %s", field, param)
	}
	return fmt.Sprintf("%s harus tidak lebih dari %s", field, param)
}

func validationMessageMaxChars(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must not be greater than %s characters", field, param)
	}
	return fmt.Sprintf("%s harus tidak lebih dari %s karakter", field, param)
}

func validationMessageMaxItems(lang Language, field string, param string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must not have more than %s items", field, param)
	}
	return fmt.Sprintf("%s harus tidak lebih dari %s item", field, param)
}

func validationMessageCompare(lang Language, field string, operator string, param string) string {
	if lang == LanguageEnglish {
		switch operator {
		case "gt":
			return fmt.Sprintf("%s must be greater than %s", field, param)
		case "gte":
			return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
		case "lt":
			return fmt.Sprintf("%s must be less than %s", field, param)
		default:
			return fmt.Sprintf("%s must be less than or equal to %s", field, param)
		}
	}

	switch operator {
	case "gt":
		return fmt.Sprintf("%s harus lebih dari %s", field, param)
	case "gte":
		return fmt.Sprintf("%s harus lebih dari atau sama dengan %s", field, param)
	case "lt":
		return fmt.Sprintf("%s harus kurang dari %s", field, param)
	default:
		return fmt.Sprintf("%s harus kurang dari atau sama dengan %s", field, param)
	}
}

func validationMessageNumeric(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be a number", field)
	}
	return fmt.Sprintf("%s harus angka", field)
}

func validationMessageTimezone(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be a valid timezone (Ex: Asia/Jakarta)", field)
	}
	return fmt.Sprintf("%s harus zona waktu yang valid (Contoh: Asia/Jakarta)", field)
}

func validationMessageEqualField(lang Language, field string, eqFieldName string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be equal to %s", field, eqFieldName)
	}
	return fmt.Sprintf("%s harus sama dengan %s", field, eqFieldName)
}

func validationMessageOneOf(lang Language, field string, param string) string {
	oneOfValues := strings.Split(param, " ")
	lastIndex := len(oneOfValues) - 1
	if lastIndex >= 0 {
		if lang == LanguageEnglish {
			oneOfValues[lastIndex] = "or " + oneOfValues[lastIndex]
		} else {
			oneOfValues[lastIndex] = "atau " + oneOfValues[lastIndex]
		}
	}
	oneOfValuesStr := strings.Join(oneOfValues, ", ")
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s must be one of %s", field, oneOfValuesStr)
	}
	return fmt.Sprintf("%s harus salah satu dari %s", field, oneOfValuesStr)
}

func validationMessageUniqueInSlice(lang Language, field string) string {
	if lang == LanguageEnglish {
		return fmt.Sprintf("%s elements must be unique", field)
	}
	return fmt.Sprintf("elemen %s harus unik", field)
}
