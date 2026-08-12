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

var staffData = []struct {
	rowData     []driver.Value
	row         staffRow
	staff       Staff
	staffString string
}{
	{
		[]driver.Value{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		staffRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		Staff{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		"Mr John Smith",
	},
	{
		[]driver.Value{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		staffRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		Staff{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		"Mr Bob Nguyen",
	},
	{
		[]driver.Value{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		staffRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		Staff{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		"Mrs Jane Smith",
	},
	{
		[]driver.Value{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		staffRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		Staff{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		"Ms Janice Jones",
	},
}

func TestStaffString(t *testing.T) {
	testCases := []struct {
		staff       Staff
		staffString string
	}{
		{
			staffData[0].staff,
			staffData[0].staffString,
		},
		{
			staffData[1].staff,
			staffData[1].staffString,
		},
		{
			staffData[2].staff,
			staffData[2].staffString,
		},
		{
			staffData[3].staff,
			staffData[3].staffString,
		},
	}

	for _, tc := range testCases {
		staffString := tc.staff.String()
		assert.Equal(t, tc.staffString, staffString)
	}
}

func TestGetStaff(t *testing.T) {
	columns := []string{
		"staff_id",
		"staff_code",
		"staff_title",
		"staff_firstname",
		"staff_prefname",
		"staff_surname",
		"staff_gender",
		"staff_email",
		"staff_username",
		"staff_government_id",
	}

	testCases := []struct {
		ex      goqu.Expression
		sql     string
		rowData [][]driver.Value
		staff   []Staff
	}{
		{
			goqu.Ex{},
			`SELECT "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "staff" WHERE (("staff"."email" IS NOT NULL) AND ("staff"."external_id" IS NOT NULL) AND ("staff"."username" IS NOT NULL))`,
			[][]driver.Value{staffData[0].rowData, staffData[1].rowData, staffData[2].rowData, staffData[3].rowData},
			[]Staff{staffData[0].staff, staffData[1].staff, staffData[2].staff, staffData[3].staff},
		},
		{
			goqu.Ex{"staff.gender": "m"},
			`SELECT "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "staff" WHERE (("staff"."gender" = 'm') AND (("staff"."email" IS NOT NULL) AND ("staff"."external_id" IS NOT NULL) AND ("staff"."username" IS NOT NULL)))`,
			[][]driver.Value{staffData[0].rowData, staffData[1].rowData},
			[]Staff{staffData[0].staff, staffData[1].staff},
		},
		{
			goqu.Ex{"staff.surname": "Smith"},
			`SELECT "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "staff" WHERE (("staff"."surname" = 'Smith') AND (("staff"."email" IS NOT NULL) AND ("staff"."external_id" IS NOT NULL) AND ("staff"."username" IS NOT NULL)))`,
			[][]driver.Value{staffData[0].rowData, staffData[2].rowData},
			[]Staff{staffData[0].staff, staffData[2].staff},
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

		staff, err := GetStaff(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.staff, staff)
	}
}
