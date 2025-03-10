package v1

import (
	"errors"
	"fmt"
)

// Common errors -------------------------------------------------------------------------------------------------------
var (
	ErrIncorrectSQLExpression = errors.New("incorrect sql-expression")
	ErrQueryTypeCannotBeEmpty = errors.New("query type cannot be empty")
	ErrTableNameCannotBeEmpty = errors.New("table name cannot be empty")
)

type ParserError struct {
	Err string
}

func (e *ParserError) Error() string {
	return e.Err
}

// Errors for SELECT ---------------------------------------------------------------------------------------------------
var (
	ErrExpectedFieldToSelect   = errors.New("at SELECT: expected field to SELECT")
	ErrExpectedCommaOrFrom     = errors.New("at SELECT: expected comma or FROM")
	ErrExpectedFrom            = errors.New("at SELECT: expected FROM")
	ErrExpectedQuotedTableName = errors.New("at SELECT: expected quoted table name")
)

type ExpectedFieldAliasToSelectError struct {
	Identifier string
}

func (e *ExpectedFieldAliasToSelectError) Error() string {
	return fmt.Sprintf("at SELECT: expected field alias for \"%s as\" to SELECT", e.Identifier)
}

// Errors for WHERE ----------------------------------------------------------------------------------------------------
var (
	ErrExpectedField                      = errors.New("at WHERE: expected field")
	ErrExpectedOperator                   = errors.New("at WHERE: expected operator")
	ErrExpectedQuotedValue                = errors.New("at WHERE: expected quoted value")
	ErrExpectedAnd                        = errors.New("at WHERE: expected AND")
	ErrEmptyWhereClause                   = errors.New("at WHERE: empty WHERE clause")
	ErrWhereClauseIsMandatory             = errors.New("at WHERE: WHERE clause is mandatory for UPDATE & DELETE")
	ErrConditionWithoutOperator           = errors.New("at WHERE: condition without operator")
	ErrConditionWithEmptyRightSideOperand = errors.New("at WHERE: condition with empty right side operand")
	ErrConditionWithEmptyLeftSideOperand  = errors.New("at WHERE: condition with empty left side operand")
)

// Errors for INSERT INTO ----------------------------------------------------------------------------------------------
var (
	ErrExpectedQuotedFieldName          = errors.New("at INSERT INTO: expected quoted field name")
	ErrNeedAtLeastOneRowToInsert        = errors.New("at INSERT INTO: need at least one row to insert")
	ErrValueCountDoesntMatchFieldCount  = errors.New("at INSERT INTO: value count doesn't match field count")
	ErrExpectedQuotedFieldNameToInsert  = errors.New("at INSERT INTO: expected quoted field name")
	ErrExpectedLessThanOneFieldToInsert = errors.New("at INSERT INTO: expected at least one field to insert")
	ErrExpectedValues                   = errors.New("at INSERT INTO: expected 'VALUES'")
	ErrExpectedOpeningParens            = errors.New("at INSERT INTO: expected opening parens")
	ErrNotMatchedFieldAndValueCount     = errors.New("at INSERT INTO: value count doesn't match field count")
	ErrExpectedCommaToInsert            = errors.New("at INSERT INTO: expected comma")
)

// Errors for DELETE FROM ----------------------------------------------------------------------------------------------
var (
	ErrExpectedQuotedTableNameToDelete = errors.New("at DELETE FROM: expected quoted table name")
	ErrExpectedWhere                   = errors.New("at DELETE FROM: expected WHERE")
)

// Errors for UPDATE ---------------------------------------------------------------------------------------------------
var (
	ErrExpectedQuotedTableNameToUpdate = errors.New("at UPDATE: expected quoted table name")
	ErrExpectedSet                     = errors.New("at UPDATE: expected 'SET'")
	ErrExpectedQuotedFieldNameToUpdate = errors.New("at UPDATE: expected quoted field name to update")
	ErrEcpectedEqualSign               = errors.New("at UPDATE: expected '='")
	ErrExpectedQuotedValueToUpdate     = errors.New("at UPDATE: expected quoted value")
	ErrExpectedComma                   = errors.New("at UPDATE: expected ','")
)

// Errors for CREATE TABLE ---------------------------------------------------------------------------------------------
var (
	ErrCreateTableTableNameCannotBeEmpty       = errors.New("at CREATE TABLE: table name cannot be empty")
	ErrCreateTableExpectedOpeningParens        = errors.New("at CREATE TABLE: expected opening parens")
	ErrCreateTableExpectedLessThanOneField     = errors.New("at CREATE TABLE: expected at least one field to create table")
	ErrCreateTableExpectedQuotedFieldName      = errors.New("at CREATE TABLE: expected quoted field name")
	ErrCreateTableUnsupportedTypeOfField       = errors.New("at CREATE TABLE: unsupported type of field")
	ErrCreateTableExpectedCommaOrClosingParens = errors.New("at CREATE TABLE: expected comma or closing parens")
)

// Errors for LIMIT ----------------------------------------------------------------------------------------------------
var (
	ErrEmptyLimitClause = errors.New("at LIMIT: empty LIMIT clause")
	ErrExpectedNumber   = errors.New("at LIMIT: required number")
)

// Errors for JOIN -----------------------------------------------------------------------------------------------------
var (
	ErrExpectedOperatorToJoin                    = errors.New("at ON: expected operator")
	ErrExpectedQuotedTableNameAndFieldNameToJoin = errors.New("at ON: expected <tablename>.<fieldname>")
)

// Errors for ORDER BY -------------------------------------------------------------------------------------------------
var (
	ErrExpectedOrder        = errors.New("expected ORDER")
	ErrExpectedFieldToOrder = errors.New("at ORDER BY: expected field to ORDER")
)

// Errors for INDEX ----------------------------------------------------------------------------------------------------
var (
	ErrIncorrectSQLExpressionForIndex  = errors.New("at INDEX: incorrect sql-expression")
	ErrExpectedQuotedIndexNameToDelete = errors.New("at DELETE INDEX: expected quoted index name")
)

// IncorrectTypeOfIndexError is an error for incorrect type of index
type IncorrectTypeOfIndexError struct {
	Type string
}

func (e *IncorrectTypeOfIndexError) Error() string {
	return "at INDEX: incorrect type of index - " + e.Type
}
