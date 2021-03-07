package main

import "argus/video/pkg/models"

func main() {

	_db := models.GetDB()
	_db.Debug().AutoMigrate(&models.Task{})
}
