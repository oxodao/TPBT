package models

type DbModel interface {
	GetTableName() string
	GetCreationScript() string
}
