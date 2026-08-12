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

var hohData = []struct {
	rowData   []driver.Value
	row       headOfHouseRow
	hoh       HeadOfHouse
	hohString string
}{
	{
		[]driver.Value{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		headOfHouseRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		HeadOfHouse{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		"Mr John Smith",
	},
	{
		[]driver.Value{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		headOfHouseRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		HeadOfHouse{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		"Mr Bob Nguyen",
	},
	{
		[]driver.Value{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		headOfHouseRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		HeadOfHouse{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		"Mrs Jane Smith",
	},
	{
		[]driver.Value{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		headOfHouseRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		HeadOfHouse{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		"Ms Janice Jones",
	},
}

func TestHeadOfHouseString(t *testing.T) {
	testCases := []struct {
		hoh       HeadOfHouse
		hohString string
	}{
		{
			hohData[0].hoh,
			hohData[0].hohString,
		},
		{
			hohData[1].hoh,
			hohData[1].hohString,
		},
		{
			hohData[2].hoh,
			hohData[2].hohString,
		},
		{
			hohData[3].hoh,
			hohData[3].hohString,
		},
	}

	for _, tc := range testCases {
		hohString := tc.hoh.String()
		assert.Equal(t, tc.hohString, hohString)
	}
}

func TestGetHeadOfHouse(t *testing.T) {
	columns := []string{
		"hoh_id",
		"hoh_code",
		"hoh_title",
		"hoh_firstname",
		"hoh_prefname",
		"hoh_surname",
		"hoh_gender",
		"hoh_email",
		"hoh_username",
		"hoh_government_id",
	}

	testCases := []struct {
		ex      goqu.Expression
		sql     string
		rowData [][]driver.Value
		hoh     []HeadOfHouse
	}{
		{
			goqu.Ex{},
			`SELECT "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("hoh"."email" IS NOT NULL) AND ("hoh"."external_id" IS NOT NULL) AND ("hoh"."username" IS NOT NULL) AND ("house"."external_id" IS NOT NULL))`,
			[][]driver.Value{hohData[0].rowData, hohData[1].rowData, hohData[2].rowData, hohData[3].rowData},
			[]HeadOfHouse{hohData[0].hoh, hohData[1].hoh, hohData[2].hoh, hohData[3].hoh},
		},
		{
			goqu.Ex{"hoh.gender": "m"},
			`SELECT "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("hoh"."gender" = 'm') AND (("hoh"."email" IS NOT NULL) AND ("hoh"."external_id" IS NOT NULL) AND ("hoh"."username" IS NOT NULL) AND ("house"."external_id" IS NOT NULL)))`,
			[][]driver.Value{hohData[0].rowData, hohData[1].rowData},
			[]HeadOfHouse{hohData[0].hoh, hohData[1].hoh},
		},
		{
			goqu.Ex{"hoh.surname": "Smith"},
			`SELECT "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("hoh"."surname" = 'Smith') AND (("hoh"."email" IS NOT NULL) AND ("hoh"."external_id" IS NOT NULL) AND ("hoh"."username" IS NOT NULL) AND ("house"."external_id" IS NOT NULL)))`,
			[][]driver.Value{hohData[0].rowData, hohData[2].rowData},
			[]HeadOfHouse{hohData[0].hoh, hohData[2].hoh},
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

		hoh, err := GetHeadOfHouse(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.hoh, hoh)
	}
}
