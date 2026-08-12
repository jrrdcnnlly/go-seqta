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

var schoolYearData = []struct {
	rowData          []driver.Value
	row              schoolYearRow
	schoolYear       SchoolYear
	schoolYearString string
}{
	{
		[]driver.Value{1, "7", "Year 7", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		schoolYearRow{1, "7", "Year 7", directorOfStudentsRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		SchoolYear{1, "7", "Year 7", DirectorOfStudents{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		"Year 7",
	},
	{
		[]driver.Value{2, "8", "Year 8", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		schoolYearRow{2, "8", "Year 8", directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		SchoolYear{2, "8", "Year 8", DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		"Year 8",
	},
	{
		[]driver.Value{3, "9", "Year 9", 3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		schoolYearRow{3, "9", "Year 9", directorOfStudentsRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		SchoolYear{3, "9", "Year 9", DirectorOfStudents{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		"Year 9",
	},
	{
		[]driver.Value{4, "10", "Year 10", 4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		schoolYearRow{4, "10", "Year 10", directorOfStudentsRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		SchoolYear{4, "10", "Year 10", DirectorOfStudents{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		"Year 10",
	},
}

func TestSchoolYearString(t *testing.T) {
	testCases := []struct {
		schoolYear       SchoolYear
		schoolYearString string
	}{
		{
			schoolYearData[0].schoolYear,
			schoolYearData[0].schoolYearString,
		},
		{
			schoolYearData[1].schoolYear,
			schoolYearData[1].schoolYearString,
		},
		{
			schoolYearData[2].schoolYear,
			schoolYearData[2].schoolYearString,
		},
		{
			schoolYearData[3].schoolYear,
			schoolYearData[3].schoolYearString,
		},
	}

	for _, tc := range testCases {
		schoolYearString := tc.schoolYear.String()
		assert.Equal(t, tc.schoolYearString, schoolYearString)
	}
}

func TestGetSchoolYear(t *testing.T) {
	columns := []string{
		"schoolyear_id",
		"schoolyear_code",
		"schoolyear_name",
		"dos_id",
		"dos_code",
		"dos_title",
		"dos_firstname",
		"dos_prefname",
		"dos_surname",
		"dos_gender",
		"dos_email",
		"dos_username",
		"dos_government_id",
	}

	testCases := []struct {
		ex          goqu.Expression
		sql         string
		rowData     [][]driver.Value
		schoolYears []SchoolYear
	}{
		{
			goqu.Ex{},
			`SELECT "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE ("schoolyear"."external_id" IS NOT NULL)`,
			[][]driver.Value{schoolYearData[0].rowData, schoolYearData[1].rowData, schoolYearData[2].rowData, schoolYearData[3].rowData},
			[]SchoolYear{schoolYearData[0].schoolYear, schoolYearData[1].schoolYear, schoolYearData[2].schoolYear, schoolYearData[3].schoolYear},
		},
		{
			goqu.Ex{"schoolyear.id": 1},
			`SELECT "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("schoolyear"."id" = 1) AND ("schoolyear"."external_id" IS NOT NULL))`,
			[][]driver.Value{schoolYearData[0].rowData},
			[]SchoolYear{schoolYearData[0].schoolYear},
		},
		{
			goqu.Ex{"dos.gender": "m"},
			`SELECT "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("dos"."gender" = 'm') AND ("schoolyear"."external_id" IS NOT NULL))`,
			[][]driver.Value{schoolYearData[0].rowData, schoolYearData[1].rowData},
			[]SchoolYear{schoolYearData[0].schoolYear, schoolYearData[1].schoolYear},
		},
		{
			goqu.Ex{"dos.surname": "Smith"},
			`SELECT "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("dos"."surname" = 'Smith') AND ("schoolyear"."external_id" IS NOT NULL))`,
			[][]driver.Value{schoolYearData[0].rowData, schoolYearData[2].rowData},
			[]SchoolYear{schoolYearData[0].schoolYear, schoolYearData[2].schoolYear},
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

		schoolYears, err := GetSchoolYear(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.schoolYears, schoolYears)
	}
}
