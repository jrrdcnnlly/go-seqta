package seqta

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

// studentRow is a reciever for database select.
type studentRow struct {
	StudentID            int    `db:"student_id"`
	StudentCode          string `db:"student_code"`
	StudentFirstName     string `db:"student_firstname"`
	StudentPreferredName string `db:"student_prefname"`
	StudentMiddleName    string `db:"student_middlename"`
	StudentSurname       string `db:"student_surname"`
	StudentGender        string `db:"student_gender"`
	StudentEmail         string `db:"student_email"`
	StudentUsername      string `db:"student_username"`
	StudentGovernmentID  string `db:"student_government_id"`
	StudentStatus        string `db:"student_status"`
	rollGroupRow
	schoolYearRow
	houseRow
}

// Student is a school student who is/has been enrolled.
type Student struct {
	ID            int        `json:"id"`
	Code          string     `json:"code"`
	FirstName     string     `json:"firstName"`
	PreferredName string     `json:"preferredName"`
	MiddleName    string     `json:"middleName"`
	Surname       string     `json:"surname"`
	Gender        string     `json:"gender"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	GovernmentID  string     `json:"governmentID"`
	Status        string     `json:"status"`
	RollGroup     RollGroup  `json:"rollGroup"`
	SchoolYear    SchoolYear `json:"schoolYear"`
	House         House      `json:"house"`
}

// newStudent allocates and returns a new [Student].
func newStudent(row studentRow) Student {
	return Student{
		ID:            row.StudentID,
		Code:          row.StudentCode,
		FirstName:     row.StudentFirstName,
		PreferredName: row.StudentPreferredName,
		MiddleName:    row.StudentMiddleName,
		Surname:       row.StudentSurname,
		Gender:        row.StudentGender,
		Email:         row.StudentEmail,
		Username:      row.StudentUsername,
		GovernmentID:  row.StudentGovernmentID,
		Status:        row.StudentStatus,
		RollGroup:     newRollGroup(row.rollGroupRow),
		SchoolYear:    newSchoolYear(row.schoolYearRow),
		House:         newHouse(row.houseRow),
	}
}

// String returns the string representation of the student.
func (s Student) String() string {
	sb := strings.Builder{}

	len := 0

	if s.PreferredName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.PreferredName)
	} else if s.FirstName != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.FirstName)
	}

	if s.Surname != "" {
		if len > 0 {
			len, _ = sb.WriteString(" ")
		}
		len, _ = sb.WriteString(s.Surname)
	}

	return sb.String()
}

// getStudentSQL returns the SQL query for [GetStudent].
func getStudentSQL(ex goqu.Expression) (sql string, params []any, err error) {
	return goqu.Dialect("postgres").
		Select(
			goqu.C("id").Table("student").As("student_id"),
			goqu.COALESCE(goqu.C("code").Table("student"), "").As("student_code"),
			goqu.COALESCE(goqu.C("firstname").Table("student"), "").As("student_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("student"), "").As("student_prefname"),
			goqu.COALESCE(goqu.C("middlename").Table("student"), "").As("student_middlename"),
			goqu.COALESCE(goqu.C("surname").Table("student"), "").As("student_surname"),
			goqu.C("gender").Table("student").As("student_gender"),
			goqu.COALESCE(goqu.C("email").Table("student"), "").As("student_email"),
			goqu.COALESCE(goqu.C("username").Table("student"), "").As("student_username"),
			goqu.COALESCE(goqu.C("government_id").Table("student"), "").As("student_government_id"),
			goqu.COALESCE(goqu.C("status").Table("student"), "").As("student_status"),
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
			goqu.C("id").Table("schoolyear").As("schoolyear_id"),
			goqu.COALESCE(goqu.C("code").Table("schoolyear"), "").As("schoolyear_code"),
			goqu.COALESCE(goqu.C("name").Table("schoolyear"), "").As("schoolyear_name"),
			goqu.C("id").Table("dos").As("dos_id"),
			goqu.COALESCE(goqu.C("code").Table("dos"), "").As("dos_code"),
			goqu.COALESCE(goqu.C("title").Table("dos"), "").As("dos_title"),
			goqu.COALESCE(goqu.C("firstname").Table("dos"), "").As("dos_firstname"),
			goqu.COALESCE(goqu.C("prefname").Table("dos"), "").As("dos_prefname"),
			goqu.COALESCE(goqu.C("surname").Table("dos"), "").As("dos_surname"),
			goqu.C("gender").Table("dos").As("dos_gender"),
			goqu.COALESCE(goqu.C("email").Table("dos"), "").As("dos_email"),
			goqu.COALESCE(goqu.C("username").Table("dos"), "").As("dos_username"),
			goqu.COALESCE(goqu.C("government_id").Table("dos"), "").As("dos_government_id"),
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
		From("student").
		Join(
			goqu.T("rollgroup"),
			goqu.On(goqu.Ex{"student.rollgroup": goqu.I("rollgroup.id")}),
		).
		Join(
			goqu.T("rollgroup_coordinator"),
			goqu.On(goqu.Ex{"rollgroup.id": goqu.I("rollgroup_coordinator.rollgroup")}),
		).
		Join(
			goqu.T("staff"),
			goqu.On(goqu.Ex{"rollgroup_coordinator.staff": goqu.I("staff.id")}),
		).
		Join(
			goqu.T("schoolyear"),
			goqu.On(goqu.Ex{"student.schoolyear": goqu.I("schoolyear.id")}),
		).
		Join(
			goqu.T("schoolyear_coordinator"),
			goqu.On(goqu.Ex{"schoolyear.id": goqu.I("schoolyear_coordinator.schoolyear")}),
		).
		Join(
			goqu.T("staff").As("dos"),
			goqu.On(goqu.Ex{"schoolyear_coordinator.staff": goqu.I("dos.id")}),
		).
		Join(
			goqu.T("house"),
			goqu.On(goqu.Ex{"student.house": goqu.I("house.id")}),
		).
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

func GetStudent(db *sqlx.DB, ex goqu.Expression) ([]Student, error) {
	ex = goqu.And(
		ex,
		goqu.Ex{
			"student.email":       goqu.Op{"neq": nil},
			"student.username":    goqu.Op{"neq": nil},
			"student.external_id": goqu.Op{"neq": nil},
		},
	)

	sql, params, err := getStudentSQL(ex)
	if err != nil {
		return nil, fmt.Errorf("seqta: get student: %w", err)
	}

	var rows []studentRow
	if err := db.Select(&rows, sql, params...); err != nil {
		return nil, fmt.Errorf("seqta: get student: %w", err)
	}

	students := make([]Student, len(rows))
	for i, r := range rows {
		students[i] = newStudent(r)
	}

	return students, nil
}
