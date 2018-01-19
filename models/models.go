package models

import (
	"database/sql/driver"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

//DB connect
var (
	DB *sqlx.DB
)

type TransactionalInformation struct {
	ReturnStatus     bool
	ReturnMessage    []string
	ValidationErrors map[string]InterfaceArray
	TotalPages       int
	TotalRows        int
	PageSize         int
	IsAuthenticated  bool
}

// Int64Array is a type implementing the sql/driver/value interface
// This is due to the native driver not supporting arrays...
type Int64Array []int64

// Value returns the driver compatible value
func (a Int64Array) Value() (driver.Value, error) {
	var strs []string
	for _, i := range a {
		strs = append(strs, strconv.FormatInt(i, 10))
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}

// InterfaceArray is a type implementing the sql/driver/value interface
// This is due to the native driver not supporting arrays...
type InterfaceArray []interface{}

// Value returns the driver compatible value
func (a InterfaceArray) Value() (driver.Value, error) {
	var strs []string
	for _, i := range a {
		if str, ok := i.(string); ok {
			strs = append(strs, q(str))
		} else {
			strs = append(strs, "''")
		}
	}
	return "{" + strings.Join(strs, ",") + "}", nil
}

// q
func q(s string) string {
	re := strings.NewReplacer("'", "''")
	return "'" + re.Replace(s) + "'"
}
