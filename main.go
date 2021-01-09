package main

import "os"

// DB_HOST=fullstack-mysql
// # DB_HOST=127.0.0.1                           # when running the app without docker
// DB_DRIVER=mysql
// DB_USER=tester
// DB_PASSWORD=testing
// DB_NAME=imagestore
// DB_PORT=3306

func main() {
	a := App{}
	//	a.Initialize("user", "password", "db", "db_mysql", 3306)
	a.Initialize(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		3306)

	a.Run(":8081")
}
