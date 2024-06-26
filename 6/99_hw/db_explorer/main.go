// тут лежит тестовый код
// менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// DSN это соединение с базой
// вы можете изменить этот на тот который вам нужен
// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
// var DSN = "root@tcp(localhost:3306)/golang2017?charset=utf8"
// var DSN = "root@tcp(172.17.0.1:3306)/golang?charset=utf8"

var DSN = "root:1234@tcp(127.0.0.1:3306)/golang?charset=utf8"

// var DSN = "coursera:5QPbAUufx7@tcp(localhost:3306)/coursera?charset=utf8"
// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
// docker run -p 3307:3306 -v $(pwd):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
func main() {
	db, err := sql.Open("mysql", DSN)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	handler, err := NewDbExplorer(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", handler)
}
