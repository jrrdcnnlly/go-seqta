package seqta

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// directorOfStudentsRow is a reciever for database select.
type directorOfStudentsRow struct {
	DoSID            int    `db:"dos_id"`
	DoSCode          string `db:"dos_code"`
	DoSTitle         string `db:"dos_title"`
	DoSFirstName     string `db:"dos_firstname"`
	DoSPreferredName string `db:"dos_prefname"`
	DoSSurname       string `db:"dos_surname"`
	DoSGender        string `db:"dos_gender"`
	DoSEmail         string `db:"dos_email"`
	DoSUsername      string `db:"dos_username"`
	DoSGovernmentID  string `db:"dos_government_id"`
}

// DirectorOfStudents is a director of students, ie. school year coordinator.
type DirectorOfStudents struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Title         string `json:"title"`
	FirstName     string `json:"firstName"`
	PreferredName string `json:"preferredName"`
	Surname       string `json:"surname"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	GovernmentID  string `json:"governmentID"`
}

// newDirectorOfStudents allocates and returns a new [DirectorOfStudents].
func newDirectorOfStudents(row directorOfStudentsRow) DirectorOfStudents {
	return DirectorOfStudents{
		ID:            row.DoSID,
		Code:          row.DoSCode,
		Title:         row.DoSTitle,
		FirstName:     row.DoSFirstName,
		PreferredName: row.DoSPreferredName,
		Surname:       row.DoSSurname,
		Gender:        row.DoSGender,
		Email:         row.DoSEmail,
		Username:      row.DoSUsername,
		GovernmentID:  row.DoSGovernmentID,
	}
}

// String returns the string representation of the director of students.
func (d DirectorOfStudents) String() string {
	sb := strings.Builder{}

	len := 0

	if d.Title != "" {
		len, _ = sb.WriteString(d.Title)
	}

	if d.PreferredName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(d.PreferredName)
	} else if d.FirstName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(d.FirstName)
	}

	if d.Surname != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(d.Surname)
	}

	return sb.String()
}

// getDirectorOfStudentsSQL returns the SQL query for [GetDirectorOfStudents].
func getDirectorOfStudentsSQL(ex goqu.Ex) (sql string, paramns []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("dos").As("dos_id"),
			goqu.COALESCE(goqu.C("code").Table("dos"), "").As("dos_code"),
			goqu.COALESCE(goqu.C("title").Table("dos"), "").As("dos_title"),
			goqu.COALESCE(goqu.C("firstname").Table("dos"), "").As("dos_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("dos"), "").As("dos_prefname"),
			goqu.COALESCE(goqu.C("surname").Table("dos"), "").As("dos_surname"),
			goqu.C("gender").Table("dos").As("dos_gender"),
			goqu.COALESCE(goqu.C("email").Table("dos"), "").As("dos_email"),
			goqu.COALESCE(goqu.C("username").Table("dos"), "").As("dos_username"),
			goqu.COALESCE(goqu.C("government_id").Table("dos"), "").As("dos_government_id"),
		).
		From("schoolyear").
		Join(
			goqu.T("schoolyear_coordinator"),
			goqu.On(goqu.Ex{"schoolyear.id": goqu.I("schoolyear_coordinator.schoolyear")}),
		).
		Join(
			goqu.T("staff").As("dos"),
			goqu.On(goqu.Ex{"schoolyear_coordinator.staff": goqu.I("dos.id")}),
		).
		Where(ex).
		ToSQL()
}

// GetDirectorOfStudents returns the directors of students selected by the goqu expression.
//
// Available columns are:
//   - schoolyear.id				[int]
//   - schoolyear.code			[string]
//   - schoolyear.name			[string]
//   - dos.id								[int]
//   - dos.code							[string]
//   - dos.title						[string]
//   - dos.firstname				[string]
//   - dos.prefname					[string]
//   - dos.surname					[string]
//   - dos.gender						[string]
//   - dos.email						[string]
//   - dos.username					[string]
//   - dos.government_id		[string]
func GetDirectorOfStudents(db *sqlx.DB, ex goqu.Ex) ([]DirectorOfStudents, error) {
	ex["schoolyear.external_id"] = goqu.Op{"neq": nil}
	ex["dos.email"] = goqu.Op{"neq": nil}
	ex["dos.username"] = goqu.Op{"neq": nil}
	ex["dos.external_id"] = goqu.Op{"neq": nil}

	sql, params, err := getDirectorOfStudentsSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get director of students: %w", err)
	}

	var rows []directorOfStudentsRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get director of students: %w", err)
	}

	dos := make([]DirectorOfStudents, len(rows))
	for i, r := range rows {
		dos[i] = newDirectorOfStudents(r)
	}

	return dos, nil
}
