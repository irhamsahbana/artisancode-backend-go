package errmsg

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestErrorsWithLanguageLocalizesCustomError(t *testing.T) {
	t.Parallel()

	err := NewCustomErrors(404).
		SetMessage(MessageCompanyNotFound).
		Add("file", MessageFileIsRequired)

	code, errorsValue := ErrorsWithLanguage[error](LanguageIndonesian, err)
	if code != 404 {
		t.Fatalf("code = %d, want 404", code)
	}

	customError, ok := errorsValue.(*CustomError)
	if !ok {
		t.Fatalf("errors type = %T, want *CustomError", errorsValue)
	}

	if customError.Msg != "Perusahaan tidak ditemukan" {
		t.Fatalf("Msg = %q, want %q", customError.Msg, "Perusahaan tidak ditemukan")
	}

	fileErrors := customError.Errors["file"]
	if len(fileErrors) != 1 || fileErrors[0] != "file wajib diisi" {
		t.Fatalf("file errors = %#v, want localized file required error", fileErrors)
	}
}

func TestErrorsWithLanguageLocalizesValidationError(t *testing.T) {
	t.Parallel()

	type request struct {
		Email string `validate:"required,email"`
	}

	req := request{}
	validate := validator.New()
	err := validate.Struct(req)

	code, errorsValue := ErrorsWithLanguage(LanguageIndonesian, err, &req)
	if code != 400 {
		t.Fatalf("code = %d, want 400", code)
	}

	errorsMap, ok := errorsValue.(map[string][]string)
	if !ok {
		t.Fatalf("errors type = %T, want map[string][]string", errorsValue)
	}

	emailErrors := errorsMap["Email"]
	if len(emailErrors) != 1 || emailErrors[0] != "Email wajib diisi" {
		t.Fatalf("Email errors = %#v, want localized required error", emailErrors)
	}
}
