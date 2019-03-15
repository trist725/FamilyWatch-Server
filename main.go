package main

import (
	"FamilyWatch/db/mongo"
)

func main() {
	mongo.Init()
	defer mongo.Dispose()

}
