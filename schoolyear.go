package seqta

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// schoolYearRow is a reciever for database select.
type schoolYearRow struct {
	SchoolYearID   int    `db:"schoolyear_id"`
	SchoolYearCode string `db:"schoolyear_code"`
	SchoolYearName string `db:"schoolyear_name"`
	directorOfStudentsRow
}

// SchoolYear is a school year level.
type SchoolYear struct {
	ID          int                `json:"id"`
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Coordinator DirectorOfStudents `json:"coordinator"`
}

// newSchoolYear allocates and returns a new [SchoolYear].
func newSchoolYear(row schoolYearRow) SchoolYear {
	return SchoolYear{
		ID:          row.SchoolYearID,
		Code:        row.SchoolYearCode,
		Name:        row.SchoolYearName,
		Coordinator: newDirectorOfStudents(row.directorOfStudentsRow),
	}
}

// String returns the string representation of the school year.
func (s SchoolYear) String() string {
	return s.Name
}

// getSchoolYearSQL returns the SQL query for [GetSchoolYear].
func getSchoolYearSQL(ex goqu.Expression) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("schoolyear").As("schoolyear_id"),
			goqu.COALESCE(goqu.C("code").Table("schoolyear"), "").As("schoolyear_code"),
			goqu.COALESCE(goqu.C("name").Table("schoolyear"), "").As("schoolyear_name"),
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

// GetSchoolYear returns the school years specified by the goqu expression.
//
// Available columns are:
//   - schoolyear.id			[int]
//   - schoolyear.code		[string]
//   - schoolyear.name		[string]
//   - dos.id							[int]
//   - dos.code						[string]
//   - dos.title					[string]
//   - dos.firstname			[string]
//   - dos.prefname				[string]
//   - dos.surname				[string]
//   - dos.gender					[string]
//   - dos.email					[string]
//   - dos.username				[string]
//   - dos.government_id	[string]
func GetSchoolYear(db *sqlx.DB, ex goqu.Expression) ([]SchoolYear, error) {
	ex = goqu.And(
		ex,
		goqu.Ex{"schoolyear.external_id": goqu.Op{"neq": nil}},
	)

	sql, params, err := getSchoolYearSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get school year: %w", err)
	}

	var rows []schoolYearRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get school year: %w", err)
	}

	schoolyears := make([]SchoolYear, len(rows))
	for i, r := range rows {
		schoolyears[i] = newSchoolYear(r)
	}

	return schoolyears, nil
}
