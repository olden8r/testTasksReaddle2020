package main

import (
	"fmt"
	"io/ioutil"
	"database/sql"
    _"github.com/go-sql-driver/mysql" //go get -u github.com/go-sql-driver/mysql
)

func main() {
	//connectionInformation := "user:password@/dbname"
	connectionInformation := "root:Vaslyaeva,27,12@tcp(127.0.0.1:3306)/employees"

	//opens a registered database driver
	myDataBase, err := sql.Open("mysql", connectionInformation)
    checkError(err)
    defer myDataBase.Close()

    //queryString := readSqlFile("./first_test2.sql")
    //queryString := readSqlFile("./second_test2.sql")
	queryString := readSqlFile("./third_test2.sql")
	showResultData(queryString, myDataBase)
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}

//read Sql file and return queries as a string
func readSqlFile(fileLocation string) string{
	query, err := ioutil.ReadFile(fileLocation)
	checkError(err)
	return string(query)
}

//print result of a query as a table
func showResultData(queryString string, db *sql.DB) {
	rows, err := db.Query(queryString)
	checkError(err)

    // Get column names
	columns, err := rows.Columns()
	checkError(err)

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the references into a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	fmt.Println("=======================================")
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		checkError(err)

		//print each column as a string
		var value string
		for i, col := range values {
			// check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	err = rows.Err()
	checkError(err)
}
