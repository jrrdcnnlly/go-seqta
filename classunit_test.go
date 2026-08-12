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

var classUnitData = []struct {
	rowData         []driver.Value
	row             classUnitRow
	classUnit       ClassUnit
	classUnitString string
}{
	{
		[]driver.Value{
			1, "20257ART.1", "20257ART", "1",
			1, "7ART", "7ART", "Art",
			1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC),
		},
		classUnitRow{
			1, "20257ART.1", "20257ART", "1",
			subjectRow{1, "7ART", "7ART", "Art"},
			termRow{1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		ClassUnit{
			1, "20257ART.1", "20257ART", "1",
			Subject{1, "7ART", "7ART", "Art"},
			Term{1, "2025", 2025, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		"20257ART#1 Art",
	},
	{
		[]driver.Value{
			2, "20257ART.2", "20257ART", "2",
			1, "7ART", "7ART", "Art",
			1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC),
		},
		classUnitRow{
			2, "20257ART.2", "20257ART", "2",
			subjectRow{1, "7ART", "7ART", "Art"},
			termRow{1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		ClassUnit{
			2, "20257ART.2", "20257ART", "2",
			Subject{1, "7ART", "7ART", "Art"},
			Term{1, "2025", 2025, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		"20257ART#2 Art",
	},
	{
		[]driver.Value{
			3, "2025ATBIO.1", "2025ATBIO", "1",
			2, "ATBIO", "ATBIO", "ATAR Biology",
			1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC),
		},
		classUnitRow{
			3, "2025ATBIO.1", "2025ATBIO", "1",
			subjectRow{2, "ATBIO", "ATBIO", "ATAR Biology"},
			termRow{1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		ClassUnit{
			3, "2025ATBIO.1", "2025ATBIO", "1",
			Subject{2, "ATBIO", "ATBIO", "ATAR Biology"},
			Term{1, "2025", 2025, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		"2025ATBIO#1 ATAR Biology",
	},
	{
		[]driver.Value{
			3, "2025ATBIO.2", "2025ATBIO", "2",
			2, "ATBIO", "ATBIO", "ATAR Biology",
			1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC),
		},
		classUnitRow{
			3, "2025ATBIO.2", "2025ATBIO", "2",
			subjectRow{2, "ATBIO", "ATBIO", "ATAR Biology"},
			termRow{1, "2025", "2025", time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		ClassUnit{
			3, "2025ATBIO.2", "2025ATBIO", "2",
			Subject{2, "ATBIO", "ATBIO", "ATAR Biology"},
			Term{1, "2025", 2025, time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, time.December, 31, 23, 59, 59, 999999, time.UTC)},
		},
		"2025ATBIO#2 ATAR Biology",
	},
}

func TestClassUnitString(t *testing.T) {
	testCases := []struct {
		classUnit       ClassUnit
		classUnitString string
	}{
		{
			classUnitData[0].classUnit,
			classUnitData[0].classUnitString,
		},
		{
			classUnitData[1].classUnit,
			classUnitData[1].classUnitString,
		},
		{
			classUnitData[2].classUnit,
			classUnitData[2].classUnitString,
		},
		{
			classUnitData[3].classUnit,
			classUnitData[3].classUnitString,
		},
	}

	for _, tc := range testCases {
		classUnitString := tc.classUnit.String()
		assert.Equal(t, tc.classUnitString, classUnitString)
	}
}

func TestGetClassUnit(t *testing.T) {
	columns := []string{
		"classunit_id",
		"classunit_code",
		"classunit_name",
		"classunit_class_number",
		"subject_id",
		"subject_code",
		"subject_name",
		"subject_description",
		"term_id",
		"term_code",
		"term_year",
		"term_start",
		"term_end",
	}

	testCases := []struct {
		ex         goqu.Expression
		sql        string
		rowData    [][]driver.Value
		classUnits []ClassUnit
	}{
		{
			goqu.Ex{},
			`SELECT "classunit"."id" AS "classunit_id", COALESCE("classunit"."code", '') AS "classunit_code", COALESCE("classunit"."name", '') AS "classunit_name", COALESCE("classunit"."class_number", '') AS "classunit_class_number", "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description", "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "classunit" INNER JOIN "subject" ON ("classunit"."subject" = "subject"."id") INNER JOIN "term" ON ("classunit"."term" = "term"."id")`,
			[][]driver.Value{classUnitData[0].rowData, classUnitData[1].rowData, classUnitData[2].rowData, classUnitData[3].rowData},
			[]ClassUnit{classUnitData[0].classUnit, classUnitData[1].classUnit, classUnitData[2].classUnit, classUnitData[3].classUnit},
		},
		{
			goqu.Ex{"classunit.id": 1},
			`SELECT "classunit"."id" AS "classunit_id", COALESCE("classunit"."code", '') AS "classunit_code", COALESCE("classunit"."name", '') AS "classunit_name", COALESCE("classunit"."class_number", '') AS "classunit_class_number", "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description", "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "classunit" INNER JOIN "subject" ON ("classunit"."subject" = "subject"."id") INNER JOIN "term" ON ("classunit"."term" = "term"."id") WHERE ("classunit"."id" = 1)`,
			[][]driver.Value{classUnitData[0].rowData},
			[]ClassUnit{classUnitData[0].classUnit},
		},
		{
			goqu.Ex{"classunit.code": "2025ATBIO.1"},
			`SELECT "classunit"."id" AS "classunit_id", COALESCE("classunit"."code", '') AS "classunit_code", COALESCE("classunit"."name", '') AS "classunit_name", COALESCE("classunit"."class_number", '') AS "classunit_class_number", "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description", "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "classunit" INNER JOIN "subject" ON ("classunit"."subject" = "subject"."id") INNER JOIN "term" ON ("classunit"."term" = "term"."id") WHERE ("classunit"."code" = '2025ATBIO.1')`,
			[][]driver.Value{classUnitData[2].rowData},
			[]ClassUnit{classUnitData[2].classUnit},
		},
		{
			goqu.Ex{"subject.description": "Art"},
			`SELECT "classunit"."id" AS "classunit_id", COALESCE("classunit"."code", '') AS "classunit_code", COALESCE("classunit"."name", '') AS "classunit_name", COALESCE("classunit"."class_number", '') AS "classunit_class_number", "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description", "term"."id" AS "term_id", COALESCE("term"."code", '') AS "term_code", COALESCE("term"."year", '') AS "term_year", COALESCE("term"."start", '1979-01-01T00:00:00Z') AS "term_start", COALESCE("term"."end", '1979-01-01T00:00:00Z') AS "term_end" FROM "classunit" INNER JOIN "subject" ON ("classunit"."subject" = "subject"."id") INNER JOIN "term" ON ("classunit"."term" = "term"."id") WHERE ("subject"."description" = 'Art')`,
			[][]driver.Value{classUnitData[0].rowData, classUnitData[1].rowData},
			[]ClassUnit{classUnitData[0].classUnit, classUnitData[1].classUnit},
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

		classUnits, err := GetClassUnit(db, loc, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.classUnits, classUnits)
	}
}
