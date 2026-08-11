package seqta

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// houseRw is a reciever for database select.
type houseRow struct {
	HouseID   int    `db:"house_id"`
	HouseCode string `db:"house_code"`
	HouseName string `db:"house_name"`
	headOfHouseRow
}

// House is a school house.
type House struct {
	ID          int         `json:"id"`
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Coordinator HeadOfHouse `json:"coordinator"`
}

// newhouse allocates and returns a new [House].
func newHouse(row houseRow) House {
	return House{
		ID:          row.HouseID,
		Code:        row.HouseCode,
		Name:        row.HouseName,
		Coordinator: newHeadOfHouse(row.headOfHouseRow),
	}
}

// String returns the string representation of the house.
func (h House) String() string {
	return h.Name
}

// getHouseSQL returns the SQL query for [GetHouse].
func getHouseSQL(ex goqu.Expression) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("house").As("house_id"),
			goqu.C("code").Table("house").As("house_code"),
			goqu.C("name").Table("house").As("house_name"),
			goqu.C("id").Table("hoh").As("hoh_id"),
			goqu.COALESCE(goqu.C("code").Table("hoh"), "").As("hoh_code"),
			goqu.COALESCE(goqu.C("title").Table("hoh"), "").As("hoh_title"),
			goqu.COALESCE(goqu.C("firstname").Table("hoh"), "").As("hoh_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("hoh"), "").As("hoh_prefname"),
			goqu.COALESCE(goqu.C("surname").Table("hoh"), "").As("hoh_surname"),
			goqu.C("gender").Table("hoh").As("hoh_gender"),
			goqu.COALESCE(goqu.C("email").Table("hoh"), "").As("hoh_email"),
			goqu.COALESCE(goqu.C("username").Table("hoh"), "").As("hoh_username"),
			goqu.COALESCE(goqu.C("government_id").Table("hoh"), "").As("hoh_government_id"),
		).
		From("house").
		Join(
			goqu.T("house_coordinator"),
			goqu.On(goqu.Ex{"house.id": goqu.I("house_coordinator.house")}),
		).
		Join(
			goqu.T("staff").As("hoh"),
			goqu.On(goqu.Ex{"house_coordinator.staff": goqu.I("hoh.id")}),
		).
		Where(ex).
		ToSQL()
}

// GetHouse returns the houses selected by the goqu expression.
//
// Available columns are:
//   - house.id						[int]
//   - house.code					[string]
//   - house.name					[string]
//   - hoh.id							[int]
//   - hoh.name						[string]
//   - hoh.title					[string]
//   - hoh.firstname			[string]
//   - hoh.prefname				[string]
//   - hoh.surname				[string]
//   - hoh.gender					[string]
//   - hoh.email					[string]
//   - hoh.username				[string]
//   - hoh.government_id	[string]
func GetHouse(db *sqlx.DB, ex goqu.Expression) ([]House, error) {
	ex = goqu.And(
		ex,
		goqu.Ex{"house.external_id": goqu.Op{"neq": nil}},
	)

	sql, params, err := getHouseSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get house: %w", err)
	}

	var rows []houseRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get house: %w", err)
	}

	houses := make([]House, len(rows))
	for i, r := range rows {
		houses[i] = newHouse(r)
	}

	return houses, nil
}
