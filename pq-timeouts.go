/*
Package pqtimeouts is a Postgres driver for Go that wraps lib/pq to provide read and write timeouts.


pq-timeouts adds two new connection string parameters: read_timeout and write_timeout. Otherwise, usage is the nearly
the same as lib/pq through the database/sql package:


	import (
		"database/sql"

		_ "github.com/lib/pq"
	)

	func main() {
		db, err := sql.Open("pq-timeouts", "user=pqtest dbname=pqtest read_timeout=500 write_timeout=1000")
		if err != nil {
			log.Fatal(err)
		}

		age := 21
		rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	}


Connections using a URL work as well:


	db,err := sql.Open("pq-timeouts", "postgres://pqtest:password@localhost/pqtest?read_timeout=500&write_timeout=1000")


read_timeout and write_timeout are specified in milliseconds. If read_timeout or write_timeout are not specified or set to 0,
no timeout is set and the driver behaves as standard lib/pq. For other connection options, check out the documentation for lib/pq:
https://godoc.org/github.com/lib/pq
*/
package pqtimeouts
