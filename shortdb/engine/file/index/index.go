package index

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	index "github.com/shortlink-org/shortdb/shortdb/domain/index/v1"
	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	binary_tree "github.com/shortlink-org/shortdb/shortdb/engine/file/index/binary-tree"
)

var ErrTreeInsert = errors.New("failed to insert value into tree")

func New(in *index.Index, rows []*page.Row) (Index[any], error) {
	var tree Index[any]

	switch in.GetType() {
	case index.Type_TYPE_UNSPECIFIED:
		return nil, ErrUnemployment
	case index.Type_TYPE_BTREE:
		return nil, ErrUnemployment
	case index.Type_TYPE_HASH:
		return nil, ErrUnemployment
	case index.Type_TYPE_BINARY_SEARCH:
		tree = binary_tree.New(func(a, b any) int {
			switch x, y := reflect.TypeOf(a), reflect.TypeOf(b); {
			case x.String() == "int" && y.String() == "int":
				aInt, aOk := a.(int)
				bInt, bOk := b.(int)

				if aOk && bOk {
					return aInt - bInt
				}

				return 0
			default:
				return 0
			}
		})

		for i := range rows {
			v, err := strconv.Atoi(string(rows[i].GetValue()["id"]))
			if err != nil {
				return nil, fmt.Errorf("failed to convert row id to integer: %w", err)
			}

			err = tree.Insert(v)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrTreeInsert, err)
			}
		}
	}

	return tree, nil
}
