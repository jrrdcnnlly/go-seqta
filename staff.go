package seqta

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// staffRow is a reciever for database select.
type staffRow struct {
	StaffID            int    `db:"staff_id"`
	StaffCode          string `db:"staff_code"`
	StaffTitle         string `db:"staff_title"`
	StaffFirstName     string `db:"staff_firstname"`
	StaffPreferredName string `db:"staff_prefname"`
	StaffSurname       string `db:"staff_surname"`
	StaffGender        string `db:"staff_gender"`
	StaffEmail         string `db:"staff_email"`
	StaffUsername      string `db:"staff_username"`
	StaffGovernmentID  string `db:"staff_government_id"`
}

// Staff is an admin member, or teacher in charge of a class or roll group.
type Staff struct {
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

// newStaff allocates and returns a new [Staff].
func newStaff(row staffRow) Staff {
	return Staff{
		ID:            row.StaffID,
		Code:          row.StaffCode,
		Title:         row.StaffTitle,
		FirstName:     row.StaffFirstName,
		PreferredName: row.StaffPreferredName,
		Surname:       row.StaffSurname,
		Gender:        row.StaffGender,
		Email:         row.StaffEmail,
		Username:      row.StaffUsername,
		GovernmentID:  row.StaffGovernmentID,
	}
}

// String returns the string representation of the staff.
func (s Staff) String() string {
	sb := strings.Builder{}

	len := 0

	if s.Title != "" {
		len, _ = sb.WriteString(s.Title)
	}

	if s.PreferredName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.PreferredName)
	} else if s.FirstName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.FirstName)
	}

	if s.Surname != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.Surname)
	}

	return sb.String()
}

// getStaffSQL returns the SQL query for [GetStaff].
func getStaffSQL(ex goqu.Ex) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("staff").As("staff_id"),
			goqu.COALESCE(goqu.C("code").Table("staff"), "").As("staff_code"),
			goqu.COALESCE(goqu.C("title").Table("staff"), "").As("staff_title"),
			goqu.COALESCE(goqu.C("firstname").Table("staff"), "").As("staff_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("staff"), "").As("staff_prefname"),
			goqu.COALESCE(goqu.C("surname").Table("staff"), "").As("staff_surname"),
			goqu.C("gender").Table("staff").As("staff_gender"),
			goqu.COALESCE(goqu.C("email").Table("staff"), "").As("staff_email"),
			goqu.COALESCE(goqu.C("username").Table("staff"), "").As("staff_username"),
			goqu.COALESCE(goqu.C("government_id").Table("staff"), "").As("staff_government_id"),
		).
		From("staff").
		Where(ex).
		ToSQL()
}

// GetStaff returns the staff selected by the goqu expression.
//
// Available columns are:
//   - staff.id								[int]
//   - staff.code							[string]
//   - staff.title						[string]
//   - staff.firstname				[string]
//   - staff.prefname					[string]
//   - staff.surname					[string]
//   - staff.gender						[string]
//   - staff.email						[string]
//   - staff.username					[string]
//   - staff.government_id		[string]
func GetStaff(db *sqlx.DB, ex goqu.Ex) ([]Staff, error) {
	ex["staff.email"] = goqu.Op{"neq": nil}
	ex["staff.username"] = goqu.Op{"neq": nil}
	ex["staff.external_id"] = goqu.Op{"neq": nil}

	sql, params, err := getStaffSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get staff: %w", err)
	}

	var rows []staffRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get staff: %w", err)
	}

	staff := make([]Staff, len(rows))
	for i, r := range rows {
		staff[i] = newStaff(r)
	}

	return staff, nil
}
