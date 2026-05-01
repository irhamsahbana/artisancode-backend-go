package errmsg

import (
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
	data := map[string]string{
		"Field":      field,
		"FieldInMsg": fieldInMsg,
		"Tag":        tag,
		"Param":      param,
	}
	if param != "" {
		return TranslateTextWithData(lang, MessageValidationDefaultWithParam, data)
	}

	return TranslateTextWithData(lang, MessageValidationDefault, data)
}

func validationMessageRequired(lang Language, field string) string {
	return validationMessage(lang, MessageValidationRequired, field, nil)
}

func validationMessageEmail(lang Language, field string) string {
	return validationMessage(lang, MessageValidationEmail, field, nil)
}

func validationMessageEmailBlacklist(lang Language, value any) string {
	return TranslateTextWithData(lang, MessageValidationEmailBlacklist, map[string]any{
		"Value": value,
	})
}

func validationMessageStrongPassword(lang Language, field string) string {
	return validationMessage(lang, MessageValidationStrongPassword, field, nil)
}

func validationMessageResourceNotExist(lang Language) string {
	return TranslateText(lang, MessageValidationResourceNotExist)
}

func validationMessageDatetime(lang Language, field string, example string) string {
	return validationMessage(lang, MessageValidationDatetime, field, map[string]string{
		"Example": example,
	})
}

func validationMessageSimpleInvalid(lang Language, field string, kind string) string {
	return validationMessage(lang, MessageValidationSimpleInvalid, field, map[string]string{
		"Kind": kind,
	})
}

func validationMessageSimpleFormat(lang Language, field string, format string) string {
	return validationMessage(lang, MessageValidationSimpleFormat, field, map[string]string{
		"Format": format,
	})
}

func validationMessageMinValue(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMinValue, field, param)
}

func validationMessageMinChars(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMinChars, field, param)
}

func validationMessageMinItems(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMinItems, field, param)
}

func validationMessageMaxValue(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMaxValue, field, param)
}

func validationMessageMaxChars(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMaxChars, field, param)
}

func validationMessageMaxItems(lang Language, field string, param string) string {
	return validationMessageWithParam(lang, MessageValidationMaxItems, field, param)
}

func validationMessageCompare(lang Language, field string, operator string, param string) string {
	code := MessageValidationCompareLte
	switch operator {
	case "gt":
		code = MessageValidationCompareGt
	case "gte":
		code = MessageValidationCompareGte
	case "lt":
		code = MessageValidationCompareLt
	}

	return validationMessageWithParam(lang, code, field, param)
}

func validationMessageNumeric(lang Language, field string) string {
	return validationMessage(lang, MessageValidationNumeric, field, nil)
}

func validationMessageTimezone(lang Language, field string) string {
	return validationMessage(lang, MessageValidationTimezone, field, nil)
}

func validationMessageEqualField(lang Language, field string, eqFieldName string) string {
	return validationMessage(lang, MessageValidationEqfield, field, map[string]string{
		"EqField": eqFieldName,
	})
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
	return validationMessage(lang, MessageValidationOneof, field, map[string]string{
		"Values": oneOfValuesStr,
	})
}

func validationMessageUniqueInSlice(lang Language, field string) string {
	return validationMessage(lang, MessageValidationUniqueInSlice, field, nil)
}

func validationMessageWithParam(lang Language, code string, field string, param string) string {
	return validationMessage(lang, code, field, map[string]string{
		"Param": param,
	})
}

func validationMessage(lang Language, code string, field string, data map[string]string) string {
	templateData := make(map[string]string, len(data)+1)
	templateData["Field"] = field
	for key, value := range data {
		templateData[key] = value
	}

	return TranslateTextWithData(lang, code, templateData)
}
