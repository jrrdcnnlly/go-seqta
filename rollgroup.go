package seqta

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// rollGroupRow is a reciever for database select.
type rollGroupRow struct {
	RollGroupID   int    `db:"rollgroup_id"`
	RollGroupCode string `db:"rollgroup_code"`
	RollGroupName string `db:"rollgroup_name"`
	staffRow
}

// RollGroup is a school home room.
type RollGroup struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Coordinator Staff  `json:"coordinator"`
}

// newRollGroup allocates and returns a new [RollGroup].
func newRollGroup(row rollGroupRow) RollGroup {
	return RollGroup{
		ID:          row.RollGroupID,
		Code:        row.RollGroupCode,
		Name:        row.RollGroupName,
		Coordinator: newStaff(row.staffRow),
	}
}

// String returns the string representaion of the roll group.
func (r RollGroup) String() string {
	return r.Name
}

// getRollGroupSQL returns the SQL query for [GetRollGroup].
func getRollGroupSQL(ex goqu.Ex) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("rollgroup").As("rollgroup_id"),
			goqu.COALESCE(goqu.C("code").Table("rollgroup"), "").As("rollgroup_code"),
			goqu.COALESCE(goqu.C("name").Table("rollgroup"), "").As("rollgroup_name"),
			goqu.C("id").Table("staff").As("staff_id"),
			goqu.COALESCE(goqu.C("code").Table("staff"), "").As("staff_code"),
			goqu.COALESCE(goqu.C("title").Table("staff"), "").As("staff_title"),
			goqu.COALESCE(goqu.C("firstname").Table("staff"), "").As("staff_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("staff"), "").As("staff_prefname"),
			goqu.COALESCE(goqu.C("surname").Table("staff"), "").As("staff_surname"),
			goqu.C("gender").Table("staff").As("staff_gender"),
			goqu.COALESCE(goqu.C("email").Table("staff"), "").As("staff_email"),
			goqu.COALESCE(goqu.C("username").Table("staff"), "").As("staff_username"),
			goqu.COALESCE(goqu.C("government_id").Table("staff"), "").As("staff_government_id"),
		).
		From("rollgroup").
		Join(
			goqu.T("rollgroup_coordinator"),
			goqu.On(goqu.Ex{"rollgroup.id": goqu.I("rollgroup_coordinator.rollgroup")}),
		).
		Join(
			goqu.T("staff"),
			goqu.On(goqu.Ex{"rollgroup_coordinator.staff": goqu.I("staff.id")}),
		).
		Where(ex).
		ToSQL()
}

// GetRollGroup returns the roll groups specified by the goqu expression.
//
// Available columns are:
//   - rollgroup.id					[int]
//   - rollgroup.code				[string]
//   - rollgroup.name				[string]
//   - staff.id							[int]
//   - staff.code						[string]
//   - staff.title					[string]
//   - staff.firstname			[string]
//   - staff.prefname				[string]
//   - staff.surname				[string]
//   - staff.gender					[string]
//   - staff.email					[string]
//   - staff.username				[string]
//   - staff.government_id	[string]
func GetRollGroup(db *sqlx.DB, ex goqu.Ex) ([]RollGroup, error) {
	ex["rollgroup.external_id"] = goqu.Op{"neq": nil}

	sql, params, err := getRollGroupSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get roll group: %w", err)
	}

	var rows []rollGroupRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get roll group: %w", err)
	}

	rollgroups := make([]RollGroup, len(rows))
	for i, r := range rows {
		rollgroups[i] = newRollGroup(r)
	}

	return rollgroups, nil
}
