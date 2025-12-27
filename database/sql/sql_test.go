package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	b := New()
	assert.NotNil(t, b)
	assert.Empty(t, b.selectFields)
	assert.Empty(t, b.fromTable)
}

func TestBuilderSelect(t *testing.T) {
	b := New()
	b.Select("id", "name", "email")

	assert.Equal(t, []string{"id", "name", "email"}, b.selectFields)
}

func TestBuilderFrom(t *testing.T) {
	b := New()
	b.From("users")

	assert.Equal(t, "users", b.fromTable)
}

func TestBuilderJoin(t *testing.T) {
	b := New()
	b.Join("orders", "users.id = orders.user_id")
	b.LeftJoin("profiles", "users.id = profiles.user_id")
	b.RightJoin("payments", "users.id = payments.user_id")

	assert.Len(t, b.joins, 3)
	assert.Contains(t, b.joins[0], "JOIN")
	assert.Contains(t, b.joins[1], "LEFT JOIN")
	assert.Contains(t, b.joins[2], "RIGHT JOIN")
}

func TestBuilderWhere(t *testing.T) {
	b := New()
	b.Where("id = ?", 1)
	b.Where("name = ?", "test")

	assert.Len(t, b.whereClause, 2)
	assert.Len(t, b.args, 2)
}

func TestBuilderWhereOr(t *testing.T) {
	b := New()
	b.Where("id = ?", 1)
	b.WhereOr("name = ?", "test")

	assert.Len(t, b.whereClause, 1)
	assert.Contains(t, b.whereClause[0], "OR")
}

func TestBuilderGroupBy(t *testing.T) {
	b := New()
	b.GroupBy("category", "type")

	assert.Equal(t, []string{"category", "type"}, b.groupBy)
}

func TestBuilderHaving(t *testing.T) {
	b := New()
	b.Having("count > ?", 10)

	assert.Len(t, b.having, 1)
	assert.Len(t, b.args, 1)
}

func TestBuilderOrderBy(t *testing.T) {
	b := New()
	b.OrderBy("created_at DESC")
	b.OrderBy("id ASC")

	assert.Len(t, b.orderBy, 2)
}

func TestBuilderLimit(t *testing.T) {
	b := New()
	b.Limit(10)

	assert.Equal(t, 10, b.limit)
}

func TestBuilderOffset(t *testing.T) {
	b := New()
	b.Offset(5)

	assert.Equal(t, 5, b.offset)
}

func TestNewPager(t *testing.T) {
	t.Run("create pager with valid values", func(t *testing.T) {
		p := NewPager(1, 10)
		assert.Equal(t, 1, p.Page)
		assert.Equal(t, 10, p.PageSize)
	})

	t.Run("page less than 1 defaults to 1", func(t *testing.T) {
		p := NewPager(0, 10)
		assert.Equal(t, 1, p.Page)
	})

	t.Run("page size less than 1 defaults to 10", func(t *testing.T) {
		p := NewPager(1, 0)
		assert.Equal(t, 10, p.PageSize)
	})

	t.Run("page size capped at 100", func(t *testing.T) {
		p := NewPager(1, 200)
		assert.Equal(t, 100, p.PageSize)
	})
}

func TestPagerSetTotal(t *testing.T) {
	p := NewPager(1, 10)
	p.SetTotal(100)

	assert.Equal(t, int64(100), p.Total)
	assert.Equal(t, 10, p.TotalPage)
}

func TestPagerOffset(t *testing.T) {
	t.Run("offset for first page", func(t *testing.T) {
		p := NewPager(1, 10)
		assert.Equal(t, 0, p.Offset())
	})

	t.Run("offset for second page", func(t *testing.T) {
		p := NewPager(2, 10)
		assert.Equal(t, 10, p.Offset())
	})
}

func TestPagerLimit(t *testing.T) {
	p := NewPager(1, 25)
	assert.Equal(t, 25, p.Limit())
}

func TestPagerHasNext(t *testing.T) {
	t.Run("has next page", func(t *testing.T) {
		p := NewPager(1, 10)
		p.SetTotal(100)
		assert.True(t, p.HasNext())
	})

	t.Run("no next page", func(t *testing.T) {
		p := NewPager(10, 10)
		p.SetTotal(100)
		assert.False(t, p.HasNext())
	})
}

func TestPagerHasPrev(t *testing.T) {
	t.Run("has previous page", func(t *testing.T) {
		p := NewPager(2, 10)
		assert.True(t, p.HasPrev())
	})

	t.Run("no previous page", func(t *testing.T) {
		p := NewPager(1, 10)
		assert.False(t, p.HasPrev())
	})
}

func TestPagerPageInfo(t *testing.T) {
	p := NewPager(1, 10)
	p.SetTotal(100)

	info := p.PageInfo()

	assert.Equal(t, 1, info["page"])
	assert.Equal(t, 10, info["page_size"])
	assert.Equal(t, int64(100), info["total"])
	assert.Equal(t, 10, info["total_page"])
}

func TestPagerString(t *testing.T) {
	p := NewPager(2, 10)
	p.SetTotal(100)

	str := p.String()
	assert.Contains(t, str, "Page: 2/10")
	assert.Contains(t, str, "PageSize: 10")
	assert.Contains(t, str, "Total: 100")
}

func TestNewPaginateResult(t *testing.T) {
	items := []int{1, 2, 3}
	p := NewPager(1, 10)

	result := NewPaginateResult(items, p)

	assert.Equal(t, items, result.Items)
	assert.Equal(t, p, result.Pager)
}

func TestPagerEdgeCases(t *testing.T) {
	t.Run("zero total", func(t *testing.T) {
		p := NewPager(1, 10)
		p.SetTotal(0)
		assert.Equal(t, 0, p.TotalPage)
	})

	t.Run("exactly one page", func(t *testing.T) {
		p := NewPager(1, 10)
		p.SetTotal(10)
		assert.Equal(t, 1, p.TotalPage)
	})

	t.Run("one extra item makes second page", func(t *testing.T) {
		p := NewPager(1, 10)
		p.SetTotal(11)
		assert.Equal(t, 2, p.TotalPage)
	})
}

func TestBuilderChain(t *testing.T) {
	t.Run("build query with chain", func(t *testing.T) {
		b := New()
		b.Select("id", "name").
			From("users").
			Where("status = ?", "active").
			OrderBy("created_at DESC").
			Limit(10)

		assert.Equal(t, 2, len(b.selectFields))
		assert.Equal(t, "users", b.fromTable)
		assert.Equal(t, 1, len(b.whereClause))
		assert.Equal(t, 1, len(b.orderBy))
		assert.Equal(t, 10, b.limit)
	})
}
