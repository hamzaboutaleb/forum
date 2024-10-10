package utils

import "forum/database"

func InitTables() error {
	err := database.CreateUserTable()
	if err != nil {
		return err
	}
	err = database.CreateSessionTable()
	if err != nil {
		return err
	}
	return nil
}
