package crdb

import (
	"fmt"
)

// More detail at
// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
type Config struct {
	User           string `json:"user,omitempty" yaml:"user,omitempty"`
	Password       string `json:"password,omitempty" yaml:"password,omitempty"`
	Host           string `json:"host,omitempty" yaml:"host,omitempty"`
	Port           int    `json:"port,omitempty" yaml:"port,omitempty"`
	DBName         string `json:"dbName,omitempty" yaml:"dbName,omitempty"`
	ConnectTimeout int    `json:"connectTimeout,omitempty" yaml:"connectTimeout,omitempty"`
	SSLMode        string `json:"sslMode,omitempty" yaml:"sslMode,omitempty"`
	SSLCert        string `json:"sslCert,omitempty" yaml:"sslCert,omitempty"`
	SSLKey         string `json:"sslKey,omitempty" yaml:"sslKey,omitempty"`
	SSLRootCert    string `json:"sslRootCert,omitempty" yaml:"sslRootCert,omitempty"`
}

// FormatDSN formats Config values into a connection as per
// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
func (d Config) FormatDSN() string {
	return fmt.Sprintf("user='%s' password='%s'"+
		" host='%s' port=%d dbname='%s'"+
		" connect_timeout=%d"+
		" sslmode='%s' sslcert='%s' sslkey='%s' sslrootcert='%s'",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.ConnectTimeout,
		d.SSLMode,
		d.SSLCert,
		d.SSLKey,
		d.SSLRootCert,
	)
}
