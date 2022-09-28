package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type DB struct {
	DB     *sql.DB
	Config SQLConfig
}

func connectDB(con SQLConfig, rdb RDB) (DB, error) {
	conStr := fmt.Sprintf("host='%s' port=%d user='%s' "+
		"password='%s' dbname='%s' sslmode=disable",
		con.Host, con.Port, con.Username, con.Password, con.DataBase)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return DB{}, err
	}

	var rows *sql.Rows
	rows, err = db.Query(fmt.Sprintf("SELECT * FROM %s.%s;", con.Schema, con.Table))
	if err == nil {
		for rows.Next() {
			break
			//todo
		}
	} else {
		switch err.(*pq.Error).Code {
		case "3D000":
			err = createDataBase(con)
			if err != nil {
				return DB{}, err
			}
			err = createTable(db, con)
			if err != nil {
				return DB{}, err
			}
			break
		case "42P01":
			err = createTable(db, con)
			if err != nil {
				return DB{}, err
			}
			break
		default:
			return DB{}, err
		}
	}

	return DB{DB: db, Config: con}, nil
}

func createDataBase(con SQLConfig) error {
	conStr := fmt.Sprintf("host='%s' port=%d user='%s' "+
		"password='%s' sslmode=disable",
		con.Host, con.Port, con.Username, con.Password)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s\nWITH\nOWNER = '%s'\n"+
		"ENCODING = 'UTF-8'\nLC_COLLATE = 'Persian_Iran.1256'\n"+
		"LC_CTYPE = 'Persian_Iran.1256'\nTABLESPACE = pg_default\n"+
		"CONNECTION LIMIT = -1\nIS_TEMPLATE = False", con.DataBase, con.Username))
	_ = db.Close()
	return err
}

func createTable(db *sql.DB, con SQLConfig) error {
	row := db.QueryRow("SELECT exists(select schema_name "+
		"FROM information_schema.schemata WHERE schema_name = $1);", con.Schema)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s\nAUTHORIZATION %s",
			con.Schema, con.Username))
		if err != nil {
			return err
		}
	}
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s.%s\n(\n"+
		"\"ID\" integer NOT NULL,\n\"Routes\" json,\n\"Airlines\" json,\n"+
		"\"Suppliers\" json,\n\"Type\" \"char\"[],\n\"Value\" integer,\n"+
		"PRIMARY KEY (\"ID\")\n);\nALTER TABLE IF EXISTS %s.%s\nOWNER to %s;",
		con.Schema, con.Table, con.Schema, con.Table, con.Username))
	return err
}

func (db DB) addRule(rule Rule) {

}