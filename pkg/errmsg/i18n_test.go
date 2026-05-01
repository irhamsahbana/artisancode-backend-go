package errmsg

import "testing"

func TestTranslateTextUsesEmbeddedCatalog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		lang    Language
		message string
		want    string
	}{
		{
			name:    "indonesian",
			lang:    LanguageIndonesian,
			message: MessageCompanyNotFound,
			want:    "Perusahaan tidak ditemukan",
		},
		{
			name:    "english",
			lang:    LanguageEnglish,
			message: MessageCompanyNotFound,
			want:    "Company not found",
		},
		{
			name:    "trailing period",
			lang:    LanguageIndonesian,
			message: "company_not_found.",
			want:    "Perusahaan tidak ditemukan",
		},
		{
			name:    "hhmm fallback",
			lang:    LanguageIndonesian,
			message: "start_time format must be HH:mm",
			want:    "start time harus berformat HH:mm",
		},
		{
			name:    "unknown message",
			lang:    LanguageIndonesian,
			message: "Uncatalogued message",
			want:    "Uncatalogued message",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := TranslateText(tt.lang, tt.message)
			if got != tt.want {
				t.Fatalf("TranslateText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContextLanguageResolution(t *testing.T) {
	t.Parallel()

	ctx := ContextWithLanguage(t.Context(), "en-US,en;q=0.9")
	if got := LanguageFromContext(ctx); got != LanguageEnglish {
		t.Fatalf("LanguageFromContext() = %q, want %q", got, LanguageEnglish)
	}
}
