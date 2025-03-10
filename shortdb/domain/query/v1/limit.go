package v1

import (
	"strconv"

	field "github.com/shortlink-org/shortdb/shortdb/domain/field/v1"
	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
)

func (q *Query) IsLimit() bool {
	return q.GetLimit() != 0
}

//nolint:revive // simple type assertion
func (q *Query) IsFilter(record *page.Row, fields map[string]field.Type) bool {
	for _, condition := range q.GetConditions() {
		var err error

		payload := record.GetValue()[condition.GetLValue()]

		switch fields[condition.GetLValue()] {
		case field.Type_TYPE_INTEGER:
			var (
				LValue int
				RValue int
			)

			LValue, err = strconv.Atoi(string(payload))
			if err != nil {
				return false
			}

			RValue, err = strconv.Atoi(condition.GetRValue())
			if err != nil {
				return false
			}

			return Filter(LValue, RValue, condition.GetOperator())
		case field.Type_TYPE_STRING:
			LValue := string(payload)

			return Filter(LValue, condition.GetRValue(), condition.GetOperator())
		case field.Type_TYPE_BOOLEAN:
			var (
				LValue bool
				RValue bool
			)

			LValue, err = strconv.ParseBool(string(payload))
			if err != nil {
				return false
			}

			RValue, err = strconv.ParseBool(condition.GetRValue())
			if err != nil {
				return false
			}

			return FilterBool(LValue, RValue, condition.GetOperator())
		case field.Type_TYPE_UNSPECIFIED:
			return false
		default:
			return false
		}
	}

	return true
}
