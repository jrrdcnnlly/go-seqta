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

var studentData = []struct {
	rowData       []driver.Value
	row           studentRow
	student       Student
	studentString string
}{
	{
		[]driver.Value{
			1, "S0001", "George", "", "Michael", "Pell", "m", "george.pell@example.com", "george.pell@example.com", "10002000", "FULL",
			1, "G1", "Green 1", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100",
			2, "8", "Year 8", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200",
			1, "Green", "Green", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100",
		},
		studentRow{
			1, "S0001", "George", "", "Michael", "Pell", "m", "george.pell@example.com", "george.pell@example.com", "10002000", "FULL",
			rollGroupRow{1, "G1", "Green 1", staffRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
			schoolYearRow{2, "8", "Year 8", directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			houseRow{1, "Green", "Green", headOfHouseRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		},
		Student{
			1, "S0001", "George", "", "Michael", "Pell", "m", "george.pell@example.com", "george.pell@example.com", "10002000", "FULL",
			RollGroup{1, "G1", "Green 1", Staff{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
			SchoolYear{2, "8", "Year 8", DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			House{1, "Green", "Green", HeadOfHouse{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		},
		"George Pell",
	},
	{
		[]driver.Value{
			2, "S0002", "Michael", "Mike", "", "Jackson", "m", "michael.jackson@example.com", "michael.jackson@example.com", "10003000", "FULL",
			1, "G1", "Green 1", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100",
			2, "8", "Year 8", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200",
			1, "Green", "Green", 1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100",
		},
		studentRow{
			2, "S0002", "Michael", "Mike", "", "Jackson", "m", "michael.jackson@example.com", "michael.jackson@example.com", "10003000", "FULL",
			rollGroupRow{1, "G1", "Green 1", staffRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
			schoolYearRow{2, "8", "Year 8", directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			houseRow{1, "Green", "Green", headOfHouseRow{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		},
		Student{
			2, "S0002", "Michael", "Mike", "", "Jackson", "m", "michael.jackson@example.com", "michael.jackson@example.com", "10003000", "FULL",
			RollGroup{1, "G1", "Green 1", Staff{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
			SchoolYear{2, "8", "Year 8", DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			House{1, "Green", "Green", HeadOfHouse{1, "E0001", "Mr", "John", "", "Smith", "m", "john.smith@example.com", "john.smith@example.com", "3100100"}},
		},
		"Mike Jackson",
	},
	{
		[]driver.Value{
			3, "S0003", "Denise", "", "", "O'Collins", "f", "denise.ocollins@example.com", "denise.ocollins@example.com", "10004000", "FULL",
			4, "B2", "Blue 2", 4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400",
			2, "8", "Year 8", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200",
			3, "Blue", "Blue", 3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300",
		},
		studentRow{
			3, "S0003", "Denise", "", "", "O'Collins", "f", "denise.ocollins@example.com", "denise.ocollins@example.com", "10004000", "FULL",
			rollGroupRow{4, "B2", "Blue 2", staffRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
			schoolYearRow{2, "8", "Year 8", directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			houseRow{3, "Blue", "Blue", headOfHouseRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		},
		Student{
			3, "S0003", "Denise", "", "", "O'Collins", "f", "denise.ocollins@example.com", "denise.ocollins@example.com", "10004000", "FULL",
			RollGroup{4, "B2", "Blue 2", Staff{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
			SchoolYear{2, "8", "Year 8", DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			House{3, "Blue", "Blue", HeadOfHouse{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		},
		"Denise O'Collins",
	},
	{
		[]driver.Value{
			4, "S0004", "Alison", "Ali", "", "Tyson", "f", "alison.tyson@example.com", "alison.tyson@example.com", "10005000", "LEFT",
			4, "B2", "Blue 2", 4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400",
			2, "8", "Year 8", 2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200",
			3, "Blue", "Blue", 3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300",
		},
		studentRow{
			4, "S0004", "Alison", "Ali", "", "Tyson", "f", "alison.tyson@example.com", "alison.tyson@example.com", "10005000", "LEFT",
			rollGroupRow{4, "B2", "Blue 2", staffRow{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
			schoolYearRow{2, "8", "Year 8", directorOfStudentsRow{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			houseRow{3, "Blue", "Blue", headOfHouseRow{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		},
		Student{
			4, "S0004", "Alison", "Ali", "", "Tyson", "f", "alison.tyson@example.com", "alison.tyson@example.com", "10005000", "LEFT",
			RollGroup{4, "B2", "Blue 2", Staff{4, "E0004", "Ms", "Janice", "", "Jones", "f", "janice.jones@example.com", "janice.jones@example.com", "3200400"}},
			SchoolYear{2, "8", "Year 8", DirectorOfStudents{2, "E0002", "Mr", "Robert", "Bob", "Nguyen", "m", "robert.nguyen@example.com", "robert.nguyen@example.com", "3200200"}},
			House{3, "Blue", "Blue", HeadOfHouse{3, "E0003", "Mrs", "Jane", "", "Smith", "f", "jane.smith@example.com", "jane.smith@example.com", "3100300"}},
		},
		"Ali Tyson",
	},
}

func TestStudentString(t *testing.T) {
	testCases := []struct {
		student       Student
		studentString string
	}{
		{
			studentData[0].student,
			studentData[0].studentString,
		},
		{
			studentData[1].student,
			studentData[1].studentString,
		},
		{
			studentData[2].student,
			studentData[2].studentString,
		},
		{
			studentData[3].student,
			studentData[3].studentString,
		},
	}

	for _, tc := range testCases {
		studentString := tc.student.String()
		assert.Equal(t, tc.studentString, studentString)
	}
}

func TestGetStudent(t *testing.T) {
	columns := []string{
		"student_id",
		"student_code",
		"student_firstname",
		"student_prefname",
		"student_middlename",
		"student_surname",
		"student_gender",
		"student_email",
		"student_username",
		"student_government_id",
		"student_status",
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
		ex       goqu.Expression
		sql      string
		rowData  [][]driver.Value
		students []Student
	}{
		{
			goqu.Ex{},
			`SELECT "student"."id" AS "student_id", COALESCE("student"."code", '') AS "student_code", COALESCE("student"."firstname", '') AS "student_firstname", COALESCE("student"."prefname", '') AS "student_prefname", COALESCE("student"."middlename", '') AS "student_middlename", COALESCE("student"."surname", '') AS "student_surname", "student"."gender" AS "student_gender", COALESCE("student"."email", '') AS "student_email", COALESCE("student"."username", '') AS "student_username", COALESCE("student"."government_id", '') AS "student_government_id", COALESCE("student"."status", '') AS "student_status", "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id", "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id", "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "student" INNER JOIN "rollgroup" ON ("student"."rollgroup" = "rollgroup"."id") INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") INNER JOIN "schoolyear" ON ("student"."schoolyear" = "schoolyear"."id") INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") INNER JOIN "house" ON ("student"."house" = "house"."id") INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("student"."email" IS NOT NULL) AND ("student"."external_id" IS NOT NULL) AND ("student"."username" IS NOT NULL))`,
			[][]driver.Value{studentData[0].rowData, studentData[1].rowData, studentData[2].rowData, studentData[3].rowData},
			[]Student{studentData[0].student, studentData[1].student, studentData[2].student, studentData[3].student},
		},
		{
			goqu.Ex{"student.id": 1},
			`SELECT "student"."id" AS "student_id", COALESCE("student"."code", '') AS "student_code", COALESCE("student"."firstname", '') AS "student_firstname", COALESCE("student"."prefname", '') AS "student_prefname", COALESCE("student"."middlename", '') AS "student_middlename", COALESCE("student"."surname", '') AS "student_surname", "student"."gender" AS "student_gender", COALESCE("student"."email", '') AS "student_email", COALESCE("student"."username", '') AS "student_username", COALESCE("student"."government_id", '') AS "student_government_id", COALESCE("student"."status", '') AS "student_status", "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id", "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id", "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "student" INNER JOIN "rollgroup" ON ("student"."rollgroup" = "rollgroup"."id") INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") INNER JOIN "schoolyear" ON ("student"."schoolyear" = "schoolyear"."id") INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") INNER JOIN "house" ON ("student"."house" = "house"."id") INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("student"."id" = 1) AND (("student"."email" IS NOT NULL) AND ("student"."external_id" IS NOT NULL) AND ("student"."username" IS NOT NULL)))`,
			[][]driver.Value{studentData[0].rowData},
			[]Student{studentData[0].student},
		},
		{
			goqu.Ex{"student.status": "LEFT"},
			`SELECT "student"."id" AS "student_id", COALESCE("student"."code", '') AS "student_code", COALESCE("student"."firstname", '') AS "student_firstname", COALESCE("student"."prefname", '') AS "student_prefname", COALESCE("student"."middlename", '') AS "student_middlename", COALESCE("student"."surname", '') AS "student_surname", "student"."gender" AS "student_gender", COALESCE("student"."email", '') AS "student_email", COALESCE("student"."username", '') AS "student_username", COALESCE("student"."government_id", '') AS "student_government_id", COALESCE("student"."status", '') AS "student_status", "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id", "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id", "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "student" INNER JOIN "rollgroup" ON ("student"."rollgroup" = "rollgroup"."id") INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") INNER JOIN "schoolyear" ON ("student"."schoolyear" = "schoolyear"."id") INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") INNER JOIN "house" ON ("student"."house" = "house"."id") INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("student"."status" = 'LEFT') AND (("student"."email" IS NOT NULL) AND ("student"."external_id" IS NOT NULL) AND ("student"."username" IS NOT NULL)))`,
			[][]driver.Value{studentData[3].rowData},
			[]Student{studentData[3].student},
		},
		{
			goqu.Ex{"house.name": "Green"},
			`SELECT "student"."id" AS "student_id", COALESCE("student"."code", '') AS "student_code", COALESCE("student"."firstname", '') AS "student_firstname", COALESCE("student"."prefname", '') AS "student_prefname", COALESCE("student"."middlename", '') AS "student_middlename", COALESCE("student"."surname", '') AS "student_surname", "student"."gender" AS "student_gender", COALESCE("student"."email", '') AS "student_email", COALESCE("student"."username", '') AS "student_username", COALESCE("student"."government_id", '') AS "student_government_id", COALESCE("student"."status", '') AS "student_status", "rollgroup"."id" AS "rollgroup_id", COALESCE("rollgroup"."code", '') AS "rollgroup_code", COALESCE("rollgroup"."name", '') AS "rollgroup_name", "staff"."id" AS "staff_id", COALESCE("staff"."code", '') AS "staff_code", COALESCE("staff"."title", '') AS "staff_title", COALESCE("staff"."firstname", '') AS "staff_firstname", COALESCE("staff"."prefname", '') AS "staff_prefname", COALESCE("staff"."surname", '') AS "staff_surname", "staff"."gender" AS "staff_gender", COALESCE("staff"."email", '') AS "staff_email", COALESCE("staff"."username", '') AS "staff_username", COALESCE("staff"."government_id", '') AS "staff_government_id", "schoolyear"."id" AS "schoolyear_id", COALESCE("schoolyear"."code", '') AS "schoolyear_code", COALESCE("schoolyear"."name", '') AS "schoolyear_name", "dos"."id" AS "dos_id", COALESCE("dos"."code", '') AS "dos_code", COALESCE("dos"."title", '') AS "dos_title", COALESCE("dos"."firstname", '') AS "dos_firstname", COALESCE("dos"."prefname", '') AS "dos_prefname", COALESCE("dos"."surname", '') AS "dos_surname", "dos"."gender" AS "dos_gender", COALESCE("dos"."email", '') AS "dos_email", COALESCE("dos"."username", '') AS "dos_username", COALESCE("dos"."government_id", '') AS "dos_government_id", "house"."id" AS "house_id", "house"."code" AS "house_code", "house"."name" AS "house_name", "hoh"."id" AS "hoh_id", COALESCE("hoh"."code", '') AS "hoh_code", COALESCE("hoh"."title", '') AS "hoh_title", COALESCE("hoh"."firstname", '') AS "hoh_firstname", COALESCE("hoh"."prefname", '') AS "hoh_prefname", COALESCE("hoh"."surname", '') AS "hoh_surname", "hoh"."gender" AS "hoh_gender", COALESCE("hoh"."email", '') AS "hoh_email", COALESCE("hoh"."username", '') AS "hoh_username", COALESCE("hoh"."government_id", '') AS "hoh_government_id" FROM "student" INNER JOIN "rollgroup" ON ("student"."rollgroup" = "rollgroup"."id") INNER JOIN "rollgroup_coordinator" ON ("rollgroup"."id" = "rollgroup_coordinator"."rollgroup") INNER JOIN "staff" ON ("rollgroup_coordinator"."staff" = "staff"."id") INNER JOIN "schoolyear" ON ("student"."schoolyear" = "schoolyear"."id") INNER JOIN "schoolyear_coordinator" ON ("schoolyear"."id" = "schoolyear_coordinator"."schoolyear") INNER JOIN "staff" AS "dos" ON ("schoolyear_coordinator"."staff" = "dos"."id") INNER JOIN "house" ON ("student"."house" = "house"."id") INNER JOIN "house_coordinator" ON ("house"."id" = "house_coordinator"."house") INNER JOIN "staff" AS "hoh" ON ("house_coordinator"."staff" = "hoh"."id") WHERE (("house"."name" = 'Green') AND (("student"."email" IS NOT NULL) AND ("student"."external_id" IS NOT NULL) AND ("student"."username" IS NOT NULL)))`,
			[][]driver.Value{studentData[0].rowData, studentData[1].rowData},
			[]Student{studentData[0].student, studentData[1].student},
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

		students, err := GetStudent(db, tc.ex)
		assert.NoError(t, err)
		assert.Equal(t, tc.students, students)
	}
}
