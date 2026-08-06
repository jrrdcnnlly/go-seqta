package seqta

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// subjectRow is a reciever for database select.
type subjectRow struct {
	SubjectID          int    `db:"subject_id"`
	SubjectCode        string `db:"subject_code"`
	SubjectName        string `db:"subject_name"`
	SubjectDescription string `db:"subject_description"`
}

// Subject is a school subject.
type Subject struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// newSubject allocates and returns a new [Subject].
func newSubject(row subjectRow) Subject {
	return Subject{
		ID:          row.SubjectID,
		Code:        row.SubjectCode,
		Name:        row.SubjectName,
		Description: row.SubjectDescription,
	}
}

// String returns the string representation of the subject.
func (s Subject) String() string {
	return s.Description
}

// getSubjectSQL returns the SQL query for [GetSubject].
func getSubjectSQL(ex goqu.Ex) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("subject").As("subject_id"),
			goqu.C("code").Table("subject").As("subject_code"),
			goqu.C("name").Table("subject").As("subject_name"),
			goqu.C("description").Table("subject").As("subject_description"),
		).
		From("subject").
		Where(ex).
		ToSQL()
}

// GetSubject returns the subjects selected by the goqu expression.
//
// Available columns are:
//   - subject.id						[int]
//   - subject.code					[string]
//   - subject.name					[string]
//   - subject.description	[string]
func GetSubject(db *sqlx.DB, ex goqu.Ex) ([]Subject, error) {
	sql, params, err := getSubjectSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get subject: %w", err)
	}

	var rows []subjectRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get subject: %w", err)
	}

	subjects := make([]Subject, len(rows))
	for i, r := range rows {
		subjects[i] = newSubject(r)
	}

	return subjects, nil
}
