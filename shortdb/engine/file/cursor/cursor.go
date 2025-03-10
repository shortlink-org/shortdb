package cursor

import (
	page "github.com/shortlink-org/shortdb/shortdb/domain/page/v1"
	"github.com/shortlink-org/shortdb/shortdb/pkg/safecast"
)

func (c *Cursor) Advance() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.RowId > 0 && c.RowId%c.Table.GetOption().GetPageSize() == 0 {
		c.PageId = safecast.IntToInt32(int(c.RowId / c.Table.GetOption().GetPageSize()))
	}

	if (c.Table.GetStats().GetRowsCount() - 1) == c.RowId {
		c.EndOfTable = true
	} else {
		c.RowId += 1
	}
}

func (c *Cursor) Value() (*page.Row, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Table.GetPages() == nil {
		return nil, ErrGetPage
	}

	p := c.Table.GetPages()[c.PageId]
	if len(p.GetRows()) == 0 {
		p.Rows = make([]*page.Row, c.Table.GetOption().GetPageSize())
	}

	rowNum := int(c.RowId) % len(p.GetRows())

	if p.GetRows()[rowNum] == nil {
		p.Rows[rowNum] = &page.Row{}
	}

	return p.GetRows()[rowNum], nil
}
