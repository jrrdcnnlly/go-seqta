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

var dosData = []struct {
	rowData   []driver.Value
	row       directorOfStudentsRow
	dos       DirectorOfStudents
	dosString string
}{
	{
		[]driver.Value{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		directorOfStudentsRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		DirectorOfStudents{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		"Mr John Smith",
	},
	{
		[]driver.Value{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		"Mr Bob Nguyen",
	},
	{
		[]driver.Value{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		directorOfStudentsRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		DirectorOfStudents{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		"Mrs Jane Smith",
	},
	{
		[]driver.Value{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		directorOfStudentsRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		DirectorOfStudents{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		"Ms Janice Jones",
	},
}

func TestDirectorOfStudentsString(t *testing.T) {
	testCases := []struct {
		dos       DirectorOfStudents
		dosString string
	}{
		{
			dosData[0].dos,
			dosData[0].dosString,
		},
		{
			dosData[1].dos,
			dosData[1].dosString,
		},
		{
			dosData[2].dos,
			dosData[2].dosString,
		},
		{
			dosData[3].dos,
			dosData[3].dosString,
		},
	}

	for _, tc := range testCases {
		dosString := tc.dos.String()
		assert.Equal(t, tc.dosString, dosString)
	}
}

func TestGetDirectorOfStudents(t *testing.T) {
	columns := []string{
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
		ex      goqu.Expression
		sql     string
		rowData [][]driver.Value
		dos     []DirectorOfStudents
	}{
		{
			goqu.Ex{},
			`SELECT "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("dos"."email" IS NOT NULL) AND ("dos"."external_id" IS NOT NULL) AND ("dos"."username" IS NOT NULL) AND ("schoolyear"."external_id" IS NOT NULL))`,
			[][]driver.Value{dosData[0].rowData, dosData[1].rowData, dosData[2].rowData, dosData[3].rowData},
			[]DirectorOfStudents{dosData[0].dos, dosData[1].dos, dosData[2].dos, dosData[3].dos},
		},
		{
			goqu.Ex{"dos.gender": "m"},
			`SELECT "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("dos"."gender" = 'm') AND (("dos"."email" IS NOT NULL) AND ("dos"."external_id" IS NOT NULL) AND ("dos"."username" IS NOT NULL) AND ("schoolyear"."external_id" IS NOT NULL)))`,
			[][]driver.Value{dosData[0].rowData, dosData[1].rowData},
			[]DirectorOfStudents{dosData[0].dos, dosData[1].dos},
		},
		{
			goqu.Ex{"dos.surname": "Smith"},
			`SELECT "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id" FROM "schoolyear" INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") WHERE (("dos"."surname" = 'Smith') AND (("dos"."email" IS NOT NULL) AND ("dos"."external_id" IS NOT NULL) AND ("dos"."username" IS NOT NULL) AND ("schoolyear"."external_id" IS NOT NULL)))`,
			[][]driver.Value{dosData[0].rowData, dosData[2].rowData},
			[]DirectorOfStudents{dosData[0].dos, dosData[2].dos},
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

		dos, err := GetDirectorOfStudents(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.dos, dos)
	}
}
