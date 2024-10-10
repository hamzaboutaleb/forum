package config

func IsAuth(id string) bool {
	_, err := SESSION.GetSession(id)
	return err == nil
}
