package db

import "github.com/Awaken1119/assistant-voice/db/mysql"

func SaveApp(appName string, appAddr string) error {
	var appData mysql.Application

	db := mysql.InitDB()

	appData.Name = appName
	appData.AppKey = appAddr

	if err := db.Model(mysql.Application{}).Where("name = ?", appName).FirstOrCreate(&appData).Error; err != nil {
		return err
	}
	return nil
}
