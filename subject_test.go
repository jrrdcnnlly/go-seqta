package seqta

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var subjectData = []struct {
	rowData       []driver.Value
	row           subjectRow
	subject       Subject
	subjectString string
}{
	{
		[]driver.Value{1, "7ART.A", "7ART.A", "Art"},
		subjectRow{1, "7ART.A", "7ART.A", "Art"},
		Subject{1, "7ART.A", "7ART.A", "Art"},
		"Art",
	},
	{
		[]driver.Value{2, "ATBIO.1", "ATBIO.1", "ATAR Biology"},
		subjectRow{2, "ATBIO.1", "ATBIO.1", "ATAR Biology"},
		Subject{2, "ATBIO.1", "ATBIO.1", "ATAR Biology"},
		"ATAR Biology",
	},
	{
		[]driver.Value{3, "ATBIO.2", "ATBIO.2", "ATAR Biology"},
		subjectRow{3, "ATBIO.2", "ATBIO.2", "ATAR Biology"},
		Subject{3, "ATBIO.2", "ATBIO.2", "ATAR Biology"},
		"ATAR Biology",
	},
	{
		[]driver.Value{4, "9DES.C", "9DES.C", "Design"},
		subjectRow{4, "9DES.C", "9DES.C", "Design"},
		Subject{4, "9DES.C", "9DES.C", "Design"},
		"Design",
	},
}

func TestSubjectString(t *testing.T) {
	testCases := []struct {
		subject       Subject
		subjectString string
	}{
		{
			subjectData[0].subject,
			subjectData[0].subjectString,
		},
		{
			subjectData[1].subject,
			subjectData[1].subjectString,
		},
		{
			subjectData[2].subject,
			subjectData[2].subjectString,
		},
		{
			subjectData[3].subject,
			subjectData[3].subjectString,
		},
	}

	for _, tc := range testCases {
		subjectString := tc.subject.String()
		assert.Equal(t, tc.subjectString, subjectString)
	}
}

func TestGetSubject(t *testing.T) {
	columns := []string{
		"subject_id",
		"subject_code",
		"subject_name",
		"subject_description",
	}

	testCases := []struct {
		ex       goqu.Expression
		sql      string
		rowData  [][]driver.Value
		subjects []Subject
	}{
		{
			goqu.Ex{},
			`SELECT "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description" FROM "subject"`,
			[][]driver.Value{subjectData[0].rowData, subjectData[1].rowData, subjectData[2].rowData, subjectData[3].rowData},
			[]Subject{subjectData[0].subject, subjectData[1].subject, subjectData[2].subject, subjectData[3].subject},
		},
		{
			goqu.Ex{"subject.description": "ATAR Biology"},
			`SELECT "subject"."id" AS "subject_id", "subject"."code" AS "subject_code", "subject"."name" AS "subject_name", "subject"."description" AS "subject_description" FROM "subject" WHERE ("subject"."description" = 'ATAR Biology')`,
			[][]driver.Value{subjectData[1].rowData, subjectData[2].rowData},
			[]Subject{subjectData[1].subject, subjectData[2].subject},
		},
	}

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

		subjects, err := GetSubject(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.subjects, subjects)
	}
}
