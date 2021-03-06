package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//Database driver
	_ "github.com/mattn/go-sqlite3"
)

var appointmentdb *sql.DB

//Appointment is the basic type to store an apointment
type Appointment struct {
	Description string
	Deadline    time.Time
	Ty          string
	Ch          Channel
}

//Channel Simplification of discordgo.Channel
type Channel struct {
	ID   string
	Name string
}

//InitDB opens database from file
func InitDB(appath string, lopath string) {
	dbb, err := sql.Open("sqlite3", appath)
	if err != nil {
		log.Panic(err)
	}
	appointmentdb = dbb

	dbb, err = sql.Open("sqlite3", lopath)
	if err != nil {
		log.Panic(err)
	}
	lookupdb = dbb

	sqlStmt := `CREATE TABLE IF NOT EXISTS "Channels" (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		channelID TEXT NOT NULL,
		channelName TEXT NOT NULL
	);`
	_, err = lookupdb.Exec(sqlStmt)
	if err != nil {
		log.Panic(err)
	}

	sqlStmt = `CREATE TABLE IF NOT EXISTS "Users" (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		userID TEXT NOT NULL,
		userName TEXT NOT NULL
	);`
	_, err = lookupdb.Exec(sqlStmt)
	if err != nil {
		log.Panic(err)
	}
}

//GetAppointmentsFromDatabase recieves all appointments from a given channel
func GetAppointmentsFromDatabase(channelID string) ([]Appointment, error) {
	var aps []Appointment
	sqlStmt := `SELECT * FROM "` + channelID + `";`
	stmt, err := appointmentdb.Prepare(sqlStmt)
	if err != nil {
		return aps, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return aps, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var ap Appointment
		var deadline string
		err = rows.Scan(&id, &ap.Description, &deadline, &ap.Ty)
		if err != nil {
			return aps, err
		}
		ap.Deadline, err = time.Parse(time.UnixDate, deadline)
		if err != nil {
			return aps, err
		}
		aps = append(aps, ap)
	}
	return aps, nil
}

//WriteAppointmentToDatabse writes Appointment to Database
func WriteAppointmentToDatabse(channelID string, ap Appointment) error {
	sqlStmt := `INSERT INTO "` + channelID + `" (
		description,deadline,type
	) VALUES (
		?,?,?
	);`
	stmt, err := appointmentdb.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ap.Description, ap.Deadline.Format(time.UnixDate), ap.Ty)
	return err
}

//MakeNewChannelTable creates new table for new Channel
func MakeNewChannelTable(channelID string) error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS "` + channelID + `"(
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		description TEXT,
		deadline TEXT NOT NULL,
		type TEXT NOT NULL);`
	_, err := appointmentdb.Exec(sqlStmt)
	return err
}

//CheckAppointmentExists checks if an appointment exists in the db, if exists: true, else false
func CheckAppointmentExists(channelID string, ap Appointment) (bool, error) {
	sqlStmt := fmt.Sprintf(`SELECT id FROM "%s" WHERE description = "%s" AND deadline = "%s" AND type = "%s";`,
		channelID, ap.Description, ap.Deadline.Format(time.UnixDate), ap.Ty)
	var id int
	err := appointmentdb.QueryRow(sqlStmt).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil

}

//DeleteAppointment deletes an appointment with the values of ap
func DeleteAppointment(channelID string, ap Appointment) error {
	sqlStmt := fmt.Sprintf(`DELETE FROM "%s" WHERE description = "%s" AND deadline = "%s" AND type = "%s";`,
		channelID, ap.Description, ap.Deadline.Format(time.UnixDate), ap.Ty)
	_, err := appointmentdb.Exec(sqlStmt)
	return err
}
