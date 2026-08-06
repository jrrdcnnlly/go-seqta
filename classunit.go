package seqta

import (
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// classUnitRow is a reciever for database select.
type classUnitRow struct {
	ClassUnitID          int    `db:"classunit_id"`
	ClassUnitCode        string `db:"classunit_code"`
	ClassUnitName        string `db:"classunit_name"`
	ClassUnitClassNumber string `db:"classunit_class_number"`
	subjectRow
	termRow
}

// ClassUnit is as SEQTA class.
type ClassUnit struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	ClassNumber string  `json:"classNumber"`
	Subject     Subject `json:"subject"`
	Term        Term    `json:"term"`
}

// newClassUnit allocates and returns a new [ClassUnit].
func newClassUnit(row classUnitRow, loc *time.Location) ClassUnit {
	return ClassUnit{
		ID:          row.ClassUnitID,
		Code:        row.ClassUnitCode,
		Name:        row.ClassUnitName,
		ClassNumber: row.ClassUnitClassNumber,
		Subject:     newSubject(row.subjectRow),
		Term:        newTerm(row.termRow, loc),
	}
}

// String returns the string representation of the class unit.
func (c ClassUnit) String() string {
	return fmt.Sprintf("%s#%s %s", c.Name, c.ClassNumber, c.Subject.Description)
}

// getClassUnitSQL returns the SQL query for [GetClassUnit.]
func getClassUnitSQL(loc *time.Location, ex goqu.Ex) (sql string, params []any, err error) {
	epoch := time.Date(1979, time.January, 1, 0, 0, 0, 0, loc)

	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("classunit").As("classunit_id"),
			goqu.COALESCE(goqu.C("code").Table("classunit"), "").As("classunit_code"),
			goqu.COALESCE(goqu.C("name").Table("classunit"), "").As("classunit_name"),
			goqu.COALESCE(goqu.C("class_number").Table("classunit"), "").As("classunit_class_number"),
			goqu.C("id").Table("subject").As("subject_id"),
			goqu.C("code").Table("subject").As("subject_code"),
			goqu.C("name").Table("subject").As("subject_name"),
			goqu.C("description").Table("subject").As("subject_description"),
			goqu.C("id").Table("term").As("term_id"),
			goqu.COALESCE(goqu.C("code").Table("term"), "").As("term_code"),
			goqu.COALESCE(goqu.C("year").Table("term"), "").As("term_year"),
			goqu.COALESCE(goqu.C("start").Table("term"), epoch).As("term_start"),
			goqu.COALESCE(goqu.C("end").Table("term"), epoch).As("term_end"),
		).
		From("classunit").
		Join(
			goqu.T("subject"),
			goqu.On(goqu.Ex{"classunit.subject": goqu.I("subject.id")}),
		).
		Join(
			goqu.T("term"),
			goqu.On(goqu.Ex{"classunit.term": goqu.I("term.id")}),
		).
		Where(ex).
		ToSQL()
}

// GetClassUnit returns the classunits selected by the goqu expression.
//
// Available columns are:
//   - classunit.id						[int]
//   - classunit.code					[string]
//   - classunit.name					[string]
//   - classunit.class_number	[string]
//   - subject.id							[int]
//   - subject.code						[string]
//   - subject.name						[string]
//   - subject.description		[string]
//   - term.id								[int]
//   - term.code							[string]
//   - term.year							[string]
//   - term.start							[time.Time]
//   - term.end								[time.Time]
func GetClassUnit(db *sqlx.DB, loc *time.Location, ex goqu.Ex) ([]ClassUnit, error) {
	sql, params, err := getClassUnitSQL(loc, ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get class unit: %w", err)
	}

	var rows []classUnitRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get class unit: %w", err)
	}

	classunits := make([]ClassUnit, len(rows))
	for i, r := range rows {
		classunits[i] = newClassUnit(r, loc)
	}

	return classunits, nil
}
