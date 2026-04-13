package errmsg

import (
	"context"
	"strings"
)

type Language string

const (
	LanguageIndonesian Language = "id"
	LanguageEnglish    Language = "en"
)

const DefaultLanguage = LanguageIndonesian

type LocalizedText struct {
	ID string
	EN string
}

type MessageCatalog map[string]LocalizedText

type contextLanguageKey string

const languageContextKey contextLanguageKey = "request_language"

var defaultCatalog = MessageCatalog{
	"Your request has failed to process": {
		ID: "Permintaan Anda gagal diproses",
		EN: "Your request has failed to process",
	},
	"Your request has been failed to process": {
		ID: "Permintaan Anda gagal diproses",
		EN: "Your request has been failed to process",
	},
	"Your request has been successfully processed": {
		ID: "Permintaan Anda berhasil diproses",
		EN: "Your request has been successfully processed",
	},
	"Company not found": {
		ID: "Perusahaan tidak ditemukan",
		EN: "Company not found",
	},
	"Company code already exists": {
		ID: "Kode perusahaan sudah ada",
		EN: "Company code already exists",
	},
	"Attendance radius must be >= 0": {
		ID: "Radius absensi harus >= 0",
		EN: "Attendance radius must be >= 0",
	},
	"Leave allowance annual must be >= 0": {
		ID: "Jatah cuti tahunan harus >= 0",
		EN: "Leave allowance annual must be >= 0",
	},
	"Overtime rate multiplier must be >= 0": {
		ID: "Pengali tarif lembur harus >= 0",
		EN: "Overtime rate multiplier must be >= 0",
	},
	"Timezone is required": {
		ID: "Timezone wajib diisi",
		EN: "Timezone is required",
	},
	"Date format is required": {
		ID: "Format tanggal wajib diisi",
		EN: "Date format is required",
	},
	"Time format is required": {
		ID: "Format waktu wajib diisi",
		EN: "Time format is required",
	},
	"Preferred language is required": {
		ID: "Bahasa default wajib diisi",
		EN: "Preferred language is required",
	},
	"Supported languages is required": {
		ID: "Daftar bahasa yang didukung wajib diisi",
		EN: "Supported languages is required",
	},
	"Supported languages contains unsupported language": {
		ID: "Daftar bahasa yang didukung berisi bahasa yang tidak didukung",
		EN: "Supported languages contains unsupported language",
	},
	"Preferred language must exist in supported languages": {
		ID: "Bahasa default harus ada di daftar bahasa yang didukung",
		EN: "Preferred language must exist in supported languages",
	},
	"Invalid credentials": {
		ID: "Kredensial tidak valid",
		EN: "Invalid credentials",
	},
	"Failed to generate token": {
		ID: "Gagal membuat token",
		EN: "Failed to generate token",
	},
	"Email is already registered": {
		ID: "Email sudah terdaftar",
		EN: "Email is already registered",
	},
	"Tenant code is already registered": {
		ID: "Kode tenant sudah terdaftar",
		EN: "Tenant code is already registered",
	},
	"Role not found": {
		ID: "Peran tidak ditemukan",
		EN: "Role not found",
	},
	"User not found": {
		ID: "Pengguna tidak ditemukan",
		EN: "User not found",
	},
	"User has no roles": {
		ID: "Pengguna tidak memiliki peran",
		EN: "User has no roles",
	},
	"Invalid or expired refresh token": {
		ID: "Refresh token tidak valid atau sudah kedaluwarsa",
		EN: "Invalid or expired refresh token",
	},
	"Employee not found": {
		ID: "Karyawan tidak ditemukan",
		EN: "Employee not found",
	},
	"Work location not found": {
		ID: "Lokasi kerja tidak ditemukan",
		EN: "Work location not found",
	},
	"Work shift not found": {
		ID: "Shift kerja tidak ditemukan",
		EN: "Work shift not found",
	},
	"Org unit not found": {
		ID: "Unit organisasi tidak ditemukan",
		EN: "Org unit not found",
	},
	"Organization unit not found": {
		ID: "Unit organisasi tidak ditemukan",
		EN: "Organization unit not found",
	},
	"You cannot delete your own account": {
		ID: "Anda tidak dapat menghapus akun Anda sendiri",
		EN: "You cannot delete your own account",
	},
	"File can no longer be deleted": {
		ID: "File tidak dapat dihapus lagi",
		EN: "File can no longer be deleted",
	},
	"Invalid join date or join date timezone": {
		ID: "Tanggal bergabung atau timezone tanggal bergabung tidak valid",
		EN: "Invalid join date or join date timezone",
	},
	"Employee number already exists in this tenant": {
		ID: "Nomor karyawan sudah ada di tenant ini",
		EN: "Employee number already exists in this tenant",
	},
	"Failed to get employee role": {
		ID: "Gagal mengambil peran karyawan",
		EN: "Failed to get employee role",
	},
	"Failed to create user account": {
		ID: "Gagal membuat akun pengguna",
		EN: "Failed to create user account",
	},
	"Failed to update employee password": {
		ID: "Gagal memperbarui kata sandi karyawan",
		EN: "Failed to update employee password",
	},
	"Work location name already exists": {
		ID: "Nama lokasi kerja sudah ada",
		EN: "Work location name already exists",
	},
	"file is required": {
		ID: "file wajib diisi",
		EN: "file is required",
	},
	"body is required": {
		ID: "body wajib diisi",
		EN: "body is required",
	},
	"filename is required": {
		ID: "filename wajib diisi",
		EN: "filename is required",
	},
	"failed to get page of results": {
		ID: "gagal mengambil halaman hasil",
		EN: "failed to get page of results",
	},
}

func ResolveLanguage(value string) Language {
	normalized := strings.TrimSpace(strings.ToLower(value))
	switch {
	case strings.HasPrefix(normalized, "en"):
		return LanguageEnglish
	case strings.HasPrefix(normalized, "id"):
		return LanguageIndonesian
	default:
		return DefaultLanguage
	}
}

func LanguageFromContext(ctx context.Context) Language {
	if ctx == nil {
		return DefaultLanguage
	}

	if language, ok := ctx.Value(languageContextKey).(string); ok && language != "" {
		return ResolveLanguage(language)
	}

	return DefaultLanguage
}

func ContextWithLanguage(ctx context.Context, language string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, languageContextKey, string(ResolveLanguage(language)))
}

func (t LocalizedText) Localize(lang Language) string {
	switch lang {
	case LanguageEnglish:
		if t.EN != "" {
			return t.EN
		}
		return t.ID
	default:
		if t.ID != "" {
			return t.ID
		}
		return t.EN
	}
}

func TranslateText(lang Language, message string) string {
	if message == "" {
		return ""
	}

	if localizedText, ok := defaultCatalog[message]; ok {
		return localizedText.Localize(lang)
	}

	trimmedMessage := strings.TrimSuffix(message, ".")
	if trimmedMessage != message {
		if localizedText, ok := defaultCatalog[trimmedMessage]; ok {
			return localizedText.Localize(lang)
		}
	}

	return message
}

func LocalizeFieldName(lang Language, field string) string {
	if field == "" {
		return ""
	}

	field = strings.ReplaceAll(field, "_", " ")
	if lang == LanguageEnglish {
		return field
	}

	return field
}
