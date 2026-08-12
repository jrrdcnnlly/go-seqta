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

var houseData = []struct {
	rowData     []driver.Value
	row         houseRow
	house       House
	houseString string
}{
	{
		[]driver.Value{1, "Green", "Green", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"},
		houseRow{1, "Green", "Green", headOfHouseRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		House{1, "Green", "Green", HeadOfHouse{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		"Green",
	},
	{
		[]driver.Value{2, "Red", "Red", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"},
		houseRow{2, "Red", "Red", headOfHouseRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		House{2, "Red", "Red", HeadOfHouse{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
		"Red",
	},
	{
		[]driver.Value{3, "Blue", "Blue", 3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"},
		houseRow{3, "Blue", "Blue", headOfHouseRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		House{3, "Blue", "Blue", HeadOfHouse{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		"Blue",
	},
	{
		[]driver.Value{4, "Orange", "Orange", 4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"},
		houseRow{4, "Orange", "Orange", headOfHouseRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		House{4, "Orange", "Orange", HeadOfHouse{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
		"Orange",
	},
}

func TestHouseString(t *testing.T) {
	testCases := []struct {
		house       House
		houseString string
	}{
		{
			houseData[0].house,
			houseData[0].houseString,
		},
		{
			houseData[1].house,
			houseData[1].houseString,
		},
		{
			houseData[2].house,
			houseData[2].houseString,
		},
		{
			houseData[3].house,
			houseData[3].houseString,
		},
	}

	for _, tc := range testCases {
		houseString := tc.house.String()
		assert.Equal(t, tc.houseString, houseString)
	}
}

func TestGetHouse(t *testing.T) {
	columns := []string{
		"house_id",
		"house_code",
		"house_name",
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
		houses  []House
	}{
		{
			goqu.Ex{},
			`SELECT "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE ("house"."external_id" IS NOT NULL)`,
			[][]driver.Value{houseData[0].rowData, houseData[1].rowData, houseData[2].rowData, houseData[3].rowData},
			[]House{houseData[0].house, houseData[1].house, houseData[2].house, houseData[3].house},
		},
		{
			goqu.Ex{"house.name": "Green"},
			`SELECT "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("house"."name" = 'Green') AND ("house"."external_id" IS NOT NULL))`,
			[][]driver.Value{houseData[0].rowData},
			[]House{houseData[0].house},
		},
		{
			goqu.Ex{"hoh.surname": "Smith"},
			`SELECT "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "house" INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("hoh"."surname" = 'Smith') AND ("house"."external_id" IS NOT NULL))`,
			[][]driver.Value{houseData[0].rowData, houseData[2].rowData},
			[]House{houseData[0].house, houseData[2].house},
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

		houses, err := GetHouse(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.houses, houses)
	}
}
