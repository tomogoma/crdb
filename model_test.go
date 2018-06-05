package crdb_test

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/tomogoma/crdb"
	"gopkg.in/yaml.v2"
	"github.com/pborman/uuid"
	"strings"
	"runtime"
	"fmt"
	"path/filepath"
)

var confFile = flag.String("conf", "dsn.conf.yml", "/path/to/testconf.yml")

func init() {
	flag.Parse()
}

func setUp(t *testing.T) crdb.Config {

	confB, err := ioutil.ReadFile(*confFile)
	if err != nil {
		t.Fatalf("Error setting up: read conf file: %v", err)
	}
	conf := crdb.Config{}
	if err := yaml.Unmarshal(confB, &conf); err != nil {
		t.Fatalf("Error setting up: unmarshal conf file content (%s): %v",
			*confFile, err)
	}
	return conf
}

func TestDBConn(t *testing.T) {
	validConf := setUp(t)
	tt := []struct {
		name   string
		dsn    string
		expErr bool
	}{
		{name: "valid dsn", dsn: validConf.FormatDSN(), expErr: false},
		{name: "invalid dsn", dsn: "an-invalid-dsn-string", expErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, err := crdb.DBConn(tc.dsn)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if db == nil {
				t.Fatalf("Got nil db")
			}
		})
	}
}

func TestTryConnect(t *testing.T) {

	conf := setUp(t)
	db := newConn(t, conf.FormatDSN())
	defer db.Close()

	tt := []struct {
		name   string
		dsn    string
		db     *sql.DB
		expErr bool
	}{
		{name: "valid first conn", dsn: conf.FormatDSN(), db: nil, expErr: false},
		{name: "valid already conn", dsn: conf.FormatDSN(), db: db, expErr: false},
		{name: "invalid dsn already conn", dsn: "invalid-dsn", db: db, expErr: false},
		{name: "invalid dsn first conn", dsn: "invalid-dsn", db: nil, expErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, err := crdb.TryConnect(tc.dsn, tc.db)
			if db != nil {
				defer db.Close()
			}
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if db == nil {
				t.Fatal("DB was nil")
			}
		})
	}

}

func TestInstantiateDB(t *testing.T) {

	t.Parallel()

	tts := []struct {
		name         string
		dbNameSuffix string
		tableDescs   []string
		expErr       bool
	}{
		{
			name:         "Ok",
			dbNameSuffix: "",
			tableDescs: []string{
				"CREATE TABLE foos (name VARCHAR(25))",
				"CREATE TABLE bars (name VARCHAR(25))",
			},
			expErr: false,
		},
		{
			name:         "BadDBName",
			dbNameSuffix: "-abc", // Hyphens(-) not allowed in db name in CockroachDB
			tableDescs: []string{
				"CREATE TABLE foos (name VARCHAR(25))",
				"CREATE TABLE bars (name VARCHAR(25))",
			},
			expErr: true,
		},
		{
			name:         "BadTableDesc",
			dbNameSuffix: "",
			tableDescs: []string{
				"CREATE TABLE foos (name VARCHAR(25))",
				"CREATE TABLE bars (name VARCHARS(25))", // VARCHARS is not valid SQL keyword
			},
			expErr: true,
		},
	}
	conf := setUp(t)
	var currConf crdb.Config

	for _, tc := range tts {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			currConf = conf // No deep copy as long as crdb.Config.DBName remains a none-pointer.
			currConf.DBName = randDBName(currConf.DBName) + tc.dbNameSuffix
			db := newConn(t, currConf.FormatDSN())
			defer tearDown(t, db, currConf.DBName)

			err := crdb.InstantiateDB(db, currConf.DBName, tc.tableDescs...)
			AssertNillable(t, !tc.expErr, err)
		})
	}
}

func AssertNillable(tb testing.TB, expNil bool, err error) {
	if expNil == (err == nil) {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	var msg string
	if expNil {
		msg = fmt.Sprintf("got error: %v", err.Error())
	} else {
		msg = "Expected an error, got nil"
	}
	fmt.Printf("\033[31m%s:%d: %s \033[39m\n\n", filepath.Base(file), line, msg)
	tb.FailNow()
}

func randDBName(prefix string) string {
	// Hyphens (-) are not acceptable db names in CockroachDB.
	suffixs := strings.Split(uuid.New(), "-")
	name := prefix
	for _, suffix := range suffixs {
		name = name + "_" + suffix
	}
	return name
}

func newConn(t *testing.T, DSN string) *sql.DB {
	db, err := crdb.DBConn(DSN)
	if err != nil {
		t.Fatalf("Error setting up: crdb.DBConn(): %v", err)
	}
	return db
}

func tearDown(t *testing.T, db *sql.DB, dbName string) {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbName + " CASCADE")
	if err != nil {
		t.Fatalf("Failed teardown: %v", err)
	}
}
