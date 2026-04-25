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
	"Please verify your email to continue": {
		ID: "Silakan verifikasi email Anda untuk melanjutkan",
		EN: "Please verify your email to continue",
	},
	"Email is not verified": {
		ID: "Email belum diverifikasi",
		EN: "Email is not verified",
	},
	"Email is already verified": {
		ID: "Email sudah diverifikasi",
		EN: "Email is already verified",
	},
	"Email verified successfully": {
		ID: "Email berhasil diverifikasi",
		EN: "Email verified successfully",
	},
	"Verification email has been resent": {
		ID: "Email verifikasi berhasil dikirim ulang",
		EN: "Verification email has been resent",
	},
	"If the email is registered, a password reset link has been sent": {
		ID: "Jika email terdaftar, tautan reset password sudah dikirim",
		EN: "If the email is registered, a password reset link has been sent",
	},
	"Password has been reset successfully": {
		ID: "Password berhasil direset",
		EN: "Password has been reset successfully",
	},
	"Invalid or expired email verification token": {
		ID: "Token verifikasi email tidak valid atau sudah kedaluwarsa",
		EN: "Invalid or expired email verification token",
	},
	"Invalid or expired password reset token": {
		ID: "Token reset password tidak valid atau sudah kedaluwarsa",
		EN: "Invalid or expired password reset token",
	},
	"Too many verification email requests. Please wait before trying again": {
		ID: "Terlalu banyak permintaan email verifikasi. Silakan tunggu sebentar sebelum mencoba lagi",
		EN: "Too many verification email requests. Please wait before trying again",
	},
	"Too many password reset requests. Please wait before trying again": {
		ID: "Terlalu banyak permintaan reset password. Silakan tunggu sebentar sebelum mencoba lagi",
		EN: "Too many password reset requests. Please wait before trying again",
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
	"Work shift is required": {
		ID: "Shift kerja wajib diisi",
		EN: "Work shift is required",
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
	"Failed to create user": {
		ID: "Gagal membuat pengguna",
		EN: "Failed to create user",
	},
	"Failed to update user": {
		ID: "Gagal memperbarui pengguna",
		EN: "Failed to update user",
	},
	"Work location name already exists": {
		ID: "Nama lokasi kerja sudah ada",
		EN: "Work location name already exists",
	},
	"Attendance log not found": {
		ID: "Log absensi tidak ditemukan",
		EN: "Attendance log not found",
	},
	"Attendance recorded successfully": {
		ID: "Absensi berhasil dicatat",
		EN: "Attendance recorded successfully",
	},
	"Check in already recorded for today": {
		ID: "Check in untuk hari ini sudah tercatat",
		EN: "Check in already recorded for today",
	},
	"Check out already recorded for today": {
		ID: "Check out untuk hari ini sudah tercatat",
		EN: "Check out already recorded for today",
	},
	"Check in must be recorded before check out": {
		ID: "Check in harus dicatat sebelum check out",
		EN: "Check in must be recorded before check out",
	},
	"Company policy not found": {
		ID: "Kebijakan perusahaan tidak ditemukan",
		EN: "Company policy not found",
	},
	"Employee profile not found": {
		ID: "Profil karyawan tidak ditemukan",
		EN: "Employee profile not found",
	},
	"Export job message bus is not configured": {
		ID: "Message bus export job belum dikonfigurasi",
		EN: "Export job message bus is not configured",
	},
	"Email message bus is not configured": {
		ID: "Message bus email belum dikonfigurasi",
		EN: "Email message bus is not configured",
	},
	"Export job not found": {
		ID: "Export job tidak ditemukan",
		EN: "Export job not found",
	},
	"Export job resource metadata is required": {
		ID: "Metadata resource export job wajib diisi",
		EN: "Export job resource metadata is required",
	},
	"File filter is required": {
		ID: "Filter file wajib diisi",
		EN: "File filter is required",
	},
	"File not found": {
		ID: "File tidak ditemukan",
		EN: "File not found",
	},
	"Invalid category for the selected parent": {
		ID: "Kategori tidak valid untuk parent yang dipilih",
		EN: "Invalid category for the selected parent",
	},
	"Invalid logged_at format": {
		ID: "Format logged_at tidak valid",
		EN: "Invalid logged_at format",
	},
	"Invalid selfie file": {
		ID: "File selfie tidak valid",
		EN: "Invalid selfie file",
	},
	"Job position not found": {
		ID: "Jabatan tidak ditemukan",
		EN: "Job position not found",
	},
	"Organization unit code already exists": {
		ID: "Kode unit organisasi sudah ada",
		EN: "Organization unit code already exists",
	},
	"Route not found": {
		ID: "Rute tidak ditemukan",
		EN: "Route not found",
	},
	"Selfie file is no longer available": {
		ID: "File selfie sudah tidak tersedia",
		EN: "Selfie file is no longer available",
	},
	"Selfie file not found": {
		ID: "File selfie tidak ditemukan",
		EN: "Selfie file not found",
	},
	"User role not found": {
		ID: "Peran pengguna tidak ditemukan",
		EN: "User role not found",
	},
	"You are not allowed to access export jobs": {
		ID: "Anda tidak diizinkan mengakses export job",
		EN: "You are not allowed to access export jobs",
	},
	"You are not allowed to access this resource": {
		ID: "Anda tidak diizinkan mengakses resource ini",
		EN: "You are not allowed to access this resource",
	},
	"You are not allowed to export attendance reports": {
		ID: "Anda tidak diizinkan mengekspor laporan absensi",
		EN: "You are not allowed to export attendance reports",
	},
	"You don't have access to this company": {
		ID: "Anda tidak memiliki akses ke perusahaan ini",
		EN: "You don't have access to this company",
	},
	"product not found": {
		ID: "produk tidak ditemukan",
		EN: "product not found",
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

	const timeFormatSuffix = " format must be HH:mm"
	if strings.HasSuffix(message, timeFormatSuffix) {
		fieldName := strings.TrimSuffix(message, timeFormatSuffix)
		if lang == LanguageEnglish {
			return message
		}
		return LocalizeFieldName(lang, fieldName) + " harus berformat HH:mm"
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
