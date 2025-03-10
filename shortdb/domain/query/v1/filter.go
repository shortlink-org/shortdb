package v1

type FilterType interface {
	int | int64 | uint64 | float64 | string
}

func Filter[V FilterType](lValue, rValue V, operator Operator) bool {
	switch operator {
	case Operator_OPERATOR_EQ:
		return lValue == rValue
	case Operator_OPERATOR_GT:
		return lValue > rValue
	case Operator_OPERATOR_GTE:
		return lValue >= rValue
	case Operator_OPERATOR_LT:
		return lValue < rValue
	case Operator_OPERATOR_LTE:
		return lValue <= rValue
	case Operator_OPERATOR_NE:
		return lValue != rValue
	default:
		return false
	}
}

func FilterBool(lValue, rValue bool, operator Operator) bool {
	switch operator {
	case Operator_OPERATOR_EQ:
		return lValue == rValue
	case Operator_OPERATOR_NE:
		return lValue != rValue
	default:
		return false
	}
}
