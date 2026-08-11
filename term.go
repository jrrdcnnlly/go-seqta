package seqta

import (
	"fmt"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jmoiron/sqlx"
)

// termRow is a reciever for database select.
type termRow struct {
	TermID    int       `db:"term_id"`
	TermCode  string    `db:"term_code"`
	TermYear  string    `db:"term_year"`
	TermStart time.Time `db:"term_start"`
	TermEnd   time.Time `db:"term_end"`
}

// Term is a school term.
type Term struct {
	ID    int       `json:"int"`
	Code  string    `json:"code"`
	Year  int       `json:"year"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// newTerm allocates and returns a new [Term].
func newTerm(row termRow, loc *time.Location) Term {
	year, err := strconv.ParseInt(row.TermYear, 10, 0)
	if err != nil {
		year = 0
	}

	return Term{
		ID:   row.TermID,
		Code: row.TermCode,
		Year: int(year),
		Start: time.Date(
			row.TermStart.Year(),
			row.TermStart.Month(),
			row.TermStart.Day(),
			row.TermStart.Hour(),
			row.TermStart.Minute(),
			row.TermStart.Second(),
			row.TermStart.Nanosecond(),
			loc,
		),
		End: time.Date(
			row.TermEnd.Year(),
			row.TermEnd.Month(),
			row.TermEnd.Day(),
			row.TermEnd.Hour(),
			row.TermEnd.Minute(),
			row.TermEnd.Second(),
			row.TermEnd.Nanosecond(),
			loc,
		),
	}
}

// String returns the string representation of the term.
func (t Term) String() string {
	return strconv.Itoa(t.Year)
}

// getTermSQL returns the SQL query for [GetTerm].
func getTermSQL(loc *time.Location, ex goqu.Expression) (sql string, params []any, err error) {
	epoch := time.Date(1979, time.January, 1, 0, 0, 0, 0, loc)

	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("term").As("term_id"),
			goqu.COALESCE(goqu.C("code").Table("term"), "").As("term_code"),
			goqu.COALESCE(goqu.C("year").Table("term"), "").As("term_year"),
			goqu.COALESCE(goqu.C("start").Table("term"), epoch).As("term_start"),
			goqu.COALESCE(goqu.C("end").Table("term"), epoch).As("term_end"),
		).
		From("term").
		Where(ex).
		ToSQL()
}

// GetTerm returns the terms selected by the goqu expression.
//
// Available columns are:
//   - term.id			[int]
//   - term.code		[string]
//   - term.year		[string]
//   - term.start		[time.Time]
//   - term.end			[time.Time]
func GetTerm(db *sqlx.DB, loc *time.Location, ex goqu.Expression) ([]Term, error) {
	sql, params, err := getTermSQL(loc, ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get term: %w", err)
	}

	var rows []termRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get term: %w", err)
	}

	terms := make([]Term, len(rows))
	for i, r := range rows {
		terms[i] = newTerm(r, loc)
	}

	return terms, nil
}
