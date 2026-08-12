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

var rollGroupData = []struct {
	rowData         []driver.Value
	row             rollGroupRow
	rollGroup       RollGroup
	rollGroupString string
}{
	{
		[]driver.Value{1, "G1", "Green 1", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		rollGroupRow{1, "G1", "Green 1", staffRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		RollGroup{1, "G1", "Green 1", Staff{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		"Green 1",
	},
	{
		[]driver.Value{2, "G2", "Green 2", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		rollGroupRow{2, "G2", "Green 2", staffRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		RollGroup{2, "G2", "Green 2", Staff{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		"Green 2",
	},
	{
		[]driver.Value{3, "B1", "Blue 1", 3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		rollGroupRow{3, "B1", "Blue 1", staffRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		RollGroup{3, "B1", "Blue 1", Staff{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		"Blue 1",
	},
	{
		[]driver.Value{4, "B2", "Blue 2", 4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		rollGroupRow{4, "B2", "Blue 2", staffRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		RollGroup{4, "B2", "Blue 2", Staff{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		"Blue 2",
	},
}

func TestRollGroupString(t *testing.T) {
	testCases := []struct {
		rollGroup       RollGroup
		rollGroupString string
	}{
		{
			rollGroupData[0].rollGroup,
			rollGroupData[0].rollGroupString,
		},
		{
			rollGroupData[1].rollGroup,
			rollGroupData[1].rollGroupString,
		},
		{
			rollGroupData[2].rollGroup,
			rollGroupData[2].rollGroupString,
		},
		{
			rollGroupData[3].rollGroup,
			rollGroupData[3].rollGroupString,
		},
	}

	for _, tc := range testCases {
		rollGroupString := tc.rollGroup.String()
		assert.Equal(t, tc.rollGroupString, rollGroupString)
	}
}

func TestGetRollGroup(t *testing.T) {
	columns := []string{
		"rollgroup_id",
		"rollgroup_code",
		"rollgroup_name",
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
		ex         goqu.Expression
		sql        string
		rowData    [][]driver.Value
		rollGroups []RollGroup
	}{
		{
			goqu.Ex{},
			`SELECT "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "rollgroup" INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") WHERE ("external_id" IS NOT NULL)`,
			[][]driver.Value{rollGroupData[0].rowData, rollGroupData[1].rowData, rollGroupData[2].rowData, rollGroupData[3].rowData},
			[]RollGroup{rollGroupData[0].rollGroup, rollGroupData[1].rollGroup, rollGroupData[2].rollGroup, rollGroupData[3].rollGroup},
		},
		{
			goqu.Ex{"rollgroup.id": 1},
			`SELECT "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "rollgroup" INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") WHERE (("rollgroup"."id" = 1) AND ("external_id" IS NOT NULL))`,
			[][]driver.Value{rollGroupData[0].rowData},
			[]RollGroup{rollGroupData[0].rollGroup},
		},
		{
			goqu.Ex{"rollgroup.name": goqu.Op{"like": "Green%"}},
			`SELECT "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "rollgroup" INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") WHERE (("rollgroup"."name" LIKE 'Green%') AND ("external_id" IS NOT NULL))`,
			[][]driver.Value{rollGroupData[0].rowData, rollGroupData[1].rowData},
			[]RollGroup{rollGroupData[0].rollGroup, rollGroupData[1].rollGroup},
		},
		{
			goqu.Ex{"staff.surname": "Smith"},
			`SELECT "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id" FROM "rollgroup" INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") WHERE (("staff"."surname" = 'Smith') AND ("external_id" IS NOT NULL))`,
			[][]driver.Value{rollGroupData[0].rowData, rollGroupData[2].rowData},
			[]RollGroup{rollGroupData[0].rollGroup, rollGroupData[2].rollGroup},
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

		rollGroups, err := GetRollGroup(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.rollGroups, rollGroups)
	}
}
