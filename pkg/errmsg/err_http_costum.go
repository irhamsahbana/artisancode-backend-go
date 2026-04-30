package errmsg

type CustomError struct {
	Code      int
	ErrorCode string
	Errors    map[string][]string
	Msg       string
}

func (e *CustomError) Error() string {
	return e.Msg
}

func NewCustomErrors(errCode int, opts ...Option) *CustomError {
	err := &CustomError{
		Code:   errCode,
		Errors: make(map[string][]string),
		Msg:    "Your request has failed to process",
		// Msg:    "Permintaan Anda gagal diproses",
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func (e *CustomError) SetCode(code int) *CustomError {
	e.Code = code
	return e
}

func (e *CustomError) SetMessage(msg string) *CustomError {
	e.Msg = msg
	return e
}

func (e *CustomError) SetLocalizedMessage(msg LocalizedText) *CustomError {
	e.Msg = msg.Localize(DefaultLanguage)
	return e
}

func (e *CustomError) Add(field, msg string) *CustomError {
	e.Errors[field] = append(e.Errors[field], msg)
	return e
}

func (e *CustomError) AddLocalized(field string, msg LocalizedText) *CustomError {
	e.Errors[field] = append(e.Errors[field], msg.Localize(DefaultLanguage))
	return e
}

func (e *CustomError) HasErrors() bool {
	return len(e.Errors) > 0
}

type Option func(*CustomError)

func WithMessage(msg string) Option {
	return func(err *CustomError) {
		err.Msg = msg
	}
}

func errorCustomHandler(lang Language, err *CustomError) (int, *CustomError) {
	localized := &CustomError{
		Code:      err.Code,
		ErrorCode: err.ErrorCode,
		Errors:    make(map[string][]string, len(err.Errors)),
		Msg:       TranslateText(lang, err.Msg),
	}

	for field, messages := range err.Errors {
		localizedMessages := make([]string, 0, len(messages))
		for _, message := range messages {
			localizedMessages = append(localizedMessages, TranslateText(lang, message))
		}
		localized.Errors[field] = localizedMessages
	}

	return localized.Code, localized
}

func (e *CustomError) SetErrorCode(code string) *CustomError {
	e.ErrorCode = code
	return e
}
