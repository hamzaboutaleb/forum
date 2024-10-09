package config

import (
	"database/sql"
)

const (
	DATABASE_NAME = "database.db"
	DRIVER_NAME   = "sqlite3"
	ADDRS         = ":8000"
	TEMPLATE_DIR  = "./template"
)

var (
	DB   *sql.DB          = nil
	TMPL *TemplateManager = nil
)
