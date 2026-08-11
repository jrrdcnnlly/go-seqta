package seqta

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// headOfHouseRow is a reciever for database select.
type headOfHouseRow struct {
	HoHID            int    `db:"hoh_id"`
	HoHCode          string `db:"hoh_code"`
	HoHTitle         string `db:"hoh_title"`
	HoHFirstName     string `db:"hoh_firstname"`
	HoHPreferredName string `db:"hoh_prefname"`
	HoHSurname       string `db:"hoh_surname"`
	HoHGender        string `db:"hoh_gender"`
	HoHEmail         string `db:"hoh_email"`
	HoHUsername      string `db:"hoh_username"`
	HoHGovernmentID  string `db:"hoh_government_id"`
}

// HeadOfHouse is a head of house.
type HeadOfHouse struct {
	ID            int    `json:"id"`
	Code          string `json:"code"`
	Title         string `json:"title"`
	FirstName     string `json:"firstName"`
	PreferredName string `json:"preferredName"`
	Surname       string `json:"surname"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	GovernmentID  string `json:"governmentID"`
}

// newHeadOfHouse allocates and returns a new [HeadOfHouse].
func newHeadOfHouse(row headOfHouseRow) HeadOfHouse {
	return HeadOfHouse{
		ID:            row.HoHID,
		Code:          row.HoHCode,
		Title:         row.HoHTitle,
		FirstName:     row.HoHFirstName,
		PreferredName: row.HoHPreferredName,
		Surname:       row.HoHSurname,
		Gender:        row.HoHGender,
		Email:         row.HoHEmail,
		Username:      row.HoHUsername,
		GovernmentID:  row.HoHGovernmentID,
	}
}

// String returns the string representation of the head of house.
func (h HeadOfHouse) String() string {
	sb := strings.Builder{}

	len := 0

	if h.Title != "" {
		len, _ = sb.WriteString(h.Title)
	}

	if h.PreferredName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(h.PreferredName)
	} else if h.FirstName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(h.FirstName)
	}

	if h.Surname != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(h.Surname)
	}

	return sb.String()
}

// getHeadOfHouseSQL returns the SQL query for [GetHoH].
func getHeadOfHouseSQL(ex goqu.Expression) (sql string, paramns []any, err error) {
	return goqu.Dialect("postgres").
		Select(
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

// GetHeadOfHouse returns the directors of students selected by the goqu expression.
//
// Available columns are:
//   - house.id						[int]
//   - house.code					[string]
//   - house.name					[string]
//   - dos.id							[int]
//   - dos.code						[string]
//   - dos.title					[string]
//   - dos.firstname			[string]
//   - dos.prefname				[string]
//   - dos.surname				[string]
//   - dos.gender					[string]
//   - dos.email					[string]
//   - dos.username				[string]
//   - dos.government_id	[string]
func GetHeadOfHouse(db *sqlx.DB, ex goqu.Expression) ([]HeadOfHouse, error) {
	ex = goqu.And(
		ex,
		goqu.Ex{
			"house.external_id": goqu.Op{"neq": nil},
			"hoh.email":         goqu.Op{"neq": nil},
			"hoh.username":      goqu.Op{"neq": nil},
			"hoh.external_id":   goqu.Op{"neq": nil},
		},
	)

	sql, params, err := getHeadOfHouseSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get head of house: %w", err)
	}

	var rows []headOfHouseRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get head of house: %w", err)
	}

	dos := make([]HeadOfHouse, len(rows))
	for i, r := range rows {
		dos[i] = newHeadOfHouse(r)
	}

	return dos, nil
}
