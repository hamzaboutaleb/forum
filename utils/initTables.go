package utils

import "forum/database"

func InitTables() error {
	err := database.CreateUserTable()
	if err != nil {
		return err
	}
	return nil
}
