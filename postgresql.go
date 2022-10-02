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
	rows, err = db.Query(fmt.Sprintf("SELECT * FROM %s.%s;", con.Schema, con.RuleTable))
	if err == nil {
		var rules map[int]RuleJ
		for rows.Next() {
			var id int
			var rule RuleJ

			err = rows.Scan(&id, &rule.Routes, &rule.Airlines, &rule.Agencies, &rule.Suppliers, &rule.Type, &rule.Value)
			if err != nil {
				return DB{}, err
			}
			rules[id] = rule
		}
		for id, rule := range rules {
			rdb.addRule(id, rule)
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
		"CONNECTION LIMIT = -1\nIS_TEMPLATE = False;", con.DataBase, con.Username))
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
		"\"ID\" integer NOT NULL,\n\"Routes\" json,\n\"Airlines\" json,\n\"Agencies\" json"+
		"\"Suppliers\" json,\n\"Type\" \"char\"[] NOT NULL,\n\"Value\" integer NOT NULL,\n"+
		"PRIMARY KEY (\"ID\")\n);\nALTER TABLE IF EXISTS %s.%s\nOWNER to %s;",
		con.Schema, con.RuleTable, con.Schema, con.RuleTable, con.Username))
	return err
}

func (db DB) addRule(rule RuleJ) int {
	var id32 sql.NullInt32
	var id int
	row := db.DB.QueryRow(fmt.Sprintf("SELECT MAX(\"ID\") FROM %s.%s;",
		db.Config.Schema, db.Config.RuleTable))
	err := row.Scan(&id32)
	if err != nil {
		return -1
	}
	if id32.Valid {
		id = int(id32.Int32) + 1
	} else {
		id = 0
	}

	queryStr := fmt.Sprintf("INSERT INTO %s.%s VALUES ($1, $2, $3, $4, $5, $6, $7);",
		db.Config.Schema, db.Config.RuleTable)

	_, err = db.DB.Exec(queryStr,
		id, rule.Routes, rule.Airlines, rule.Agencies, rule.Suppliers, rule.Type, rule.Value)

	if err != nil {
		return -1
	}

	return id
}

func (db DB) validateRule(rule Rule) bool {
	mainQuery := "SELECT EXISTS (SELECT * FROM " + db.Config.Schema + ".%s WHERE %s = $1);"
	var tmpVar bool
	if rule.Routes != nil {
		tempQuery := fmt.Sprintf(mainQuery, db.Config.CityTable, "\"code\"")
		for _, route := range rule.Routes {
			if route.Origin.Valid {
				db.DB.QueryRow(tempQuery, route.Origin.String).Scan(&tmpVar)
				if !tmpVar {
					return false
				}
			}
			if route.Destination.Valid {
				db.DB.QueryRow(tempQuery, route.Origin.String).Scan(&tmpVar)
				if !tmpVar {
					return false
				}
			}
		}
	}
	if rule.Airlines != nil {
		tempQuery := fmt.Sprintf(mainQuery, db.Config.AirlineTable, "\"code\"")
		for _, airline := range rule.Airlines {
			db.DB.QueryRow(tempQuery, airline).Scan(&tmpVar)
			if !tmpVar {
				return false
			}
		}
	}
	if rule.Agencies != nil {
		tempQuery := fmt.Sprintf(mainQuery, db.Config.AgencyTable, "\"name\"")
		for _, agency := range rule.Agencies {
			db.DB.QueryRow(tempQuery, agency).Scan(&tmpVar)
			if !tmpVar {
				return false
			}
		}
	}
	if rule.Suppliers != nil {
		tempQuery := fmt.Sprintf(mainQuery, db.Config.SupplierTable, "\"name\"")
		for _, supplier := range rule.Suppliers {
			db.DB.QueryRow(tempQuery, supplier).Scan(&tmpVar)
			if !tmpVar {
				return false
			}
		}
	}
	if rule.Type != "PERCENTAGE" && rule.Type != "FIXED" {
		return false
	}
	return true
}
