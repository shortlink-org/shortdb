package engine

import (
	"context"
	"fmt"

	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	query "github.com/shortlink-org/shortdb/shortdb/domain/query/v1"
	"github.com/shortlink-org/shortdb/shortdb/engine/file"
	"github.com/shortlink-org/shortdb/shortdb/engine/options"
)

type Engine interface {
	Exec(in *query.Query) (any, error)
	Close() error

	// Table
	CreateTable(in *query.Query) error
	DropTable(name string) error

	// Index
	CreateIndex(in *query.Query) error
	DropIndex(name string) error

	// Commands
	Select(in *query.Query) ([]*page.Row, error)
	Update(in *query.Query) error
	Insert(in *query.Query) error
	Delete(in *query.Query) error
}

//nolint:ireturn,nolintlint // ignore
func New(ctx context.Context, name string, ops ...options.Option) (Engine, error) {
	var err error

	var engine Engine

	switch name {
	case "file":
		engine, err = file.New(ctx, ops...)
		if err != nil {
			return nil, fmt.Errorf("failed to create file engine: %w", err)
		}
	default:
		engine, err = file.New(ctx, ops...)
		if err != nil {
			return nil, fmt.Errorf("failed to create file engine: %w", err)
		}
	}

	return engine, nil
}
