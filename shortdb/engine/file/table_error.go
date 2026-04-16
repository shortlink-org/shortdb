package file

import "errors"

var (
	ErrExistTable        = errors.New("at CREATE TABLE: exist table")
	ErrReservedTableName = errors.New("at CREATE TABLE: names shortdb_tables and shortdbcatalog are reserved for catalog discovery (SELECT)")
)
