package crdb_test

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/tomogoma/crdb"
	"gopkg.in/yaml.v2"
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

func newConn(t *testing.T, DSN string) *sql.DB {
	db, err := crdb.DBConn(DSN)
	if err != nil {
		t.Fatalf("Error setting up: crdb.DBConn(): %v", err)
	}
	return db
}
