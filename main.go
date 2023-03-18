package main

import (
	"30httpHandlers/pkg/storage"
	"30httpHandlers/pkg/app"	
)

func main() {
	repository := storage.NewService()
	var app = app.NewApp(repository)
	app.Run()
}

/*
curl -X POST -d "{\"name\": \"Vasiliy\", \"age\": 20}" http://localhost:8080/create
curl -X POST -d "{\"name\": \"Ivan\", \"age\": 30}" http://localhost:8080/create
curl -X POST -d "{\"name\": \"Boss\", \"age\": 30}" http://localhost:8080/create
curl -X POST -d "{\"sourceId\": 1, \"targetId\": 2}" http://localhost:8080/make_friends
curl -X DELETE -d "{\"targetId\": 1}" http://localhost:8080/delete
curl -X PUT -d "{\"age\": 25}" http://localhost:8080/2

http://localhost:8080/get
http://localhost:8080/get/1
*/