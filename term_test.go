package seqta

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var termData = []struct {
	rowData    []driver.Value
	row        termRow
	term       Term
	termString string
}{
	{
		[]driver.Value{1, "2022", "2022", time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		termRow{1, "2022", "2022", time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		Term{1, "2022", 2022, time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		"2022",
	},
	{
		[]driver.Value{2, "2023", "2023", time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		termRow{2, "2023", "2023", time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		Term{2, "2023", 2023, time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		"2023",
	},
	{
		[]driver.Value{3, "2024", "2024", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		termRow{3, "2024", "2024", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		Term{3, "2024", 2024, time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		"2024",
	},
	{
		[]driver.Value{4, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		termRow{4, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		Term{4, "2025", 2025, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		"2025",
	},
}

func TestTermString(t *testing.T) {
	testCases := []struct {
		term       Term
		termString string
	}{
		{
			termData[0].term,
			termData[0].termString,
		},
		{
			termData[1].term,
			termData[1].termString,
		},
		{
			termData[2].term,
			termData[2].termString,
		},
		{
			termData[3].term,
			termData[3].termString,
		},
	}

	for _, tc := range testCases {
		termString := tc.term.String()
		assert.Equal(t, tc.termString, termString)
	}
}

func TestGetTerm(t *testing.T) {
	columns := []string{
		"term_id",
		"term_code",
		"term_year",
		"term_start",
		"term_end",
	}

	testCases := []struct {
		ex      goqu.Expression
		sql     string
		rowData [][]driver.Value
		terms   []Term
	}{
		{
			goqu.Ex{},
			`SELECT "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "term"`,
			[][]driver.Value{termData[0].rowData, termData[1].rowData, termData[2].rowData, termData[3].rowData},
			[]Term{termData[0].term, termData[1].term, termData[2].term, termData[3].term},
		},
		{
			goqu.Ex{"term.id": 1},
			`SELECT "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "term" WHERE ("term"."id" = 1)`,
			[][]driver.Value{termData[0].rowData},
			[]Term{termData[0].term},
		},
		{
			goqu.Ex{"term.year": "2025"},
			`SELECT "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "term" WHERE ("term"."year" = '2025')`,
			[][]driver.Value{termData[3].rowData},
			[]Term{termData[3].term},
		},
	}

	loc, err := time.LoadLocation("UTC")
	assert.NoError(t, err)

	// Create mock database connection.
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Initialize sqlx with mock database connection.
	db := sqlx.NewDb(mockDB, "sqlmock")

	for _, tc := range testCases {
		rows := sqlmock.NewRows(columns)
		for _, r := range tc.rowData {
			rows.AddRow(r...)
		}

		mock.ExpectQuery(regexp.QuoteMeta(tc.sql)).WillReturnRows(rows)

		terms, err := GetTerm(db, loc, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.terms, terms)
	}
}
