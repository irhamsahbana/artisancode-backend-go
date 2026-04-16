package errmsg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

const (
	pgCodeForeignKeyViolation = "23503"
	pgCodeUniqueViolation     = "23505"
	pgCodeNotNullViolation    = "23502"
)

func errorPgHandler(lang Language, errPg *pgconn.PgError) (int, map[string][]string) {
	var (
		errors    = make(map[string][]string)
		code      = 500
		column    string
		columnMsg string
	)

	log.Debug().Str("code", errPg.Code).Str("constraint", errPg.ConstraintName).Msg("Postgres error")
	log.Debug().Msgf("postgres error detail: %s", errPg.Detail)

	switch errPg.Code {
	case pgCodeForeignKeyViolation:
		regex := regexp.MustCompile(`Key \(([^)]+)\)`)
		match := regex.FindStringSubmatch(errPg.Detail)

		if len(match) > 1 {
			column = match[1]
			columnMsg = strings.ReplaceAll(column, "_", " ")
		}

		if lang == LanguageEnglish {
			errors[column] = append(errors[column], columnMsg+" is invalid")
		} else {
			errors[column] = append(errors[column], columnMsg+" tidak valid")
		}
		code = 500
	case pgCodeUniqueViolation:
		code = 409
		regex := regexp.MustCompile(`Key \(([^)]+)\)`)
		match := regex.FindStringSubmatch(errPg.Detail)

		if len(match) > 1 {
			column = match[1]
		}

		if strings.Contains(column, ",") { // checking for unique_violation is compound key
			sliceOfColumns := strings.Split(column, ", ")
			columns := strings.Join(sliceOfColumns, "_dan_")
			column = columns
			if lang == LanguageEnglish {
				columnMsg = "combination of " + strings.ReplaceAll(columns, "_", " ")
				errors[column] = append(errors[column], fmt.Sprintf("%s already exists", columnMsg))
			} else {
				columnMsg = "kombinasi " + strings.ReplaceAll(columns, "_", " ")
				errors[column] = append(errors[column], fmt.Sprintf("%s sudah ada", columnMsg))
			}
		} else { // unique_violation is not compound key
			columnMsg = strings.ReplaceAll(column, "_", " ")
			msg := fmt.Sprintf("%s sudah ada", columnMsg)
			if lang == LanguageEnglish {
				msg = fmt.Sprintf("%s already exists", columnMsg)
			}
			if column == "email" {
				if lang == LanguageEnglish {
					msg = "email already registered"
				} else {
					msg = "email sudah terdaftar"
				}
			}
			errors[column] = append(errors[column], msg)
		}
	case pgCodeNotNullViolation: // null value in column violates not-null constraint
		// postgres: null value in column "product_id" of relation "product_inquiries" violates not-null constraint
		regex := regexp.MustCompile(`column \"(.+?)\" of relation \"(.+?)\"`)
		matches := regex.FindStringSubmatch(errPg.Error())
		if len(matches) >= 3 {
			column = matches[1]
			columnNameMsg := strings.ReplaceAll(column, "_", " ")
			if lang == LanguageEnglish {
				errors[column] = append(errors[column], fmt.Sprintf("%s must not be empty", columnNameMsg))
			} else {
				errors[column] = append(errors[column], fmt.Sprintf("%s tidak boleh kosong", columnNameMsg))
			}

		}
	}

	return code, errors
}
