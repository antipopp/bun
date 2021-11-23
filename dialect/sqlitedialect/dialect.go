package sqlitedialect

import (
	"database/sql"
	"encoding/hex"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/dialect/sqltype"
	"github.com/uptrace/bun/schema"
)

type Dialect struct {
	schema.BaseDialect

	tables   *schema.Tables
	features feature.Feature
}

func New() *Dialect {
	d := new(Dialect)
	d.tables = schema.NewTables(d)
	d.features = feature.CTE |
		feature.Returning |
		feature.InsertTableAlias |
		feature.DeleteTableAlias |
		feature.InsertOnConflict
	return d
}

func (d *Dialect) Init(*sql.DB) {}

func (d *Dialect) Name() dialect.Name {
	return dialect.SQLite
}

func (d *Dialect) Features() feature.Feature {
	return d.features
}

func (d *Dialect) Tables() *schema.Tables {
	return d.tables
}

func (d *Dialect) OnTable(table *schema.Table) {
	for _, field := range table.FieldMap {
		d.onField(field)
	}
}

func (d *Dialect) onField(field *schema.Field) {
	field.DiscoveredSQLType = fieldSQLType(field)
}

func (d *Dialect) IdentQuote() byte {
	return '"'
}

func (d *Dialect) AppendBytes(b []byte, bs []byte) []byte {
	if bs == nil {
		return dialect.AppendNull(b)
	}

	b = append(b, `X'`...)

	s := len(b)
	b = append(b, make([]byte, hex.EncodedLen(len(bs)))...)
	hex.Encode(b[s:], bs)

	b = append(b, '\'')

	return b
}

func fieldSQLType(field *schema.Field) string {
	switch field.DiscoveredSQLType {
	case sqltype.SmallInt, sqltype.BigInt:
		// INTEGER PRIMARY KEY is an alias for the ROWID.
		// It is safe to convert all ints to INTEGER, because SQLite types don't have size.
		return sqltype.Integer
	default:
		return field.DiscoveredSQLType
	}
}
