package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type db struct {
	login    string
	password string
	path     string
	conn     *sql.DB
}

type schema struct {
	schemaName string
	tables     []*table
}

type table struct {
	schemaName string
	tableName  string
	columns    []*column
	indexes    []*index
	cons       []*constraint
}
type column struct {
	schemaName string
	tableName  string
	columnName string
	typeName   string
	position   int
	nullable   int
	charLength int
	defValue   string
	numScale   int
	numPrec    int
}

type IResultSet interface {
	Next() bool
	Scan(...interface{}) error
}

func (h2 db) grabSchema(schemaName string) (s *schema, err error) {
	s = &schema{
		schemaName: schemaName,
	}
	s.tables, err = h2.grabTables(schemaName)
	for _, t := range s.tables {
		t.cons, err = h2.grabConstraints(t)
		if err != nil {
			return
		}
		t.indexes, err = h2.grabIndexes(t)
		if err != nil {
			return
		}
	}
	return
}

func collectTables(rs IResultSet) (result []*table, err error) {
	var curTab *table
	for rs.Next() {
		var c = new(column)

		err = rs.Scan(

			&c.schemaName,
			&c.tableName,
			&c.columnName,
			&c.position,
			&c.defValue,
			&c.nullable,
			&c.typeName,
			&c.charLength,
			&c.numPrec,
			&c.numScale,
		)

		if err != nil {
			return
		}

		if curTab == nil || curTab.tableName != c.tableName || curTab.schemaName != c.schemaName {
			if curTab != nil {
				result = append(result, curTab)
			}
			curTab = &table{
				schemaName: c.schemaName,
				tableName:  c.tableName,
			}
		}
		curTab.columns = append(curTab.columns, c)
	}
	if curTab != nil {
		result = append(result, curTab)
	}
	return
}

func (h2 *db) grabTables(schemaName string) (result []*table, err error) {
	query := `
	SELECT 
        C.TABLE_SCHEMA
        ,C.TABLE_NAME
        ,C.COLUMN_NAME 
		,C.ORDINAL_POSITION 
		,C.COLUMN_DEFAULT 
		,C.NULLABLE 
		,C.TYPE_NAME 
		,C.CHARACTER_MAXIMUM_LENGTH 
		,C.NUMERIC_PRECISION 
		,C.NUMERIC_SCALE
	FROM INFORMATION_SCHEMA.COLUMNS C
  		INNER JOIN INFORMATION_SCHEMA.TABLES  T
		ON T.TABLE_CATALOG = C.TABLE_CATALOG
		AND T.TABLE_SCHEMA = C.TABLE_SCHEMA
		AND T.TABLE_NAME = C.TABLE_NAME
	WHERE T.TABLE_SCHEMA = '%s'
		AND T.TABLE_TYPE='TABLE' 
		AND T.STORAGE_TYPE = 'CACHED' 
	ORDER BY 
		C.TABLE_NAME ASC
		,C.ORDINAL_POSITION ASC
`
	query = fmt.Sprintf(query, schemaName)

	rs, err := h2.conn.Query(query)
	if err != nil {
		err = errors.Wrapf(err, "could not read H2 table metadata from schema %s", schemaName)
		return
	}
	defer rs.Close()
	result, err = collectTables(rs)
	if err != nil {
		errors.Wrapf(err, "could not read H2 table metadata from schema %s", schemaName)
		return
	}
	return
}

type constraint struct {
	consName   string
	references string
	columList  string
	typeName   string
	indexName  string
	checkExpr  string
}

func collectConstraints(rs IResultSet) (result []*constraint, err error) {
	for rs.Next() {
		var c = new(constraint)
		err = rs.Scan(
			&c.typeName,
			&c.consName,
			&c.columList,
			&c.indexName,
			&c.checkExpr,
		)
		if err != nil {
			return
		}
		result = append(result, c)
	}
	return
}

func (h2 *db) grabConstraints(t *table) (result []*constraint, err error) {

	query := `
	SELECT 
	C.CONSTRAINT_TYPE
	,C.CONSTRAINT_NAME
	,C.COLUMN_LIST
	,C.UNIQUE_INDEX_NAME
	,C.CHECK_EXPRESSION
   FROM INFORMATION_SCHEMA.CONSTRAINTS C
   INNER JOIN INFORMATION_SCHEMA.TABLES T
	ON T.TABLE_CATALOG = C.TABLE_CATALOG
   AND T.TABLE_SCHEMA = C.TABLE_SCHEMA
   AND T.TABLE_NAME = C.TABLE_NAME
   WHERE T.TABLE_SCHEMA = '%s'
   AND T.TABLE_TYPE='TABLE' 
   AND T.STORAGE_TYPE = 'CACHED' 
   AND T.TABLE_NAME = '%s'
   `
	query = fmt.Sprintf(query, t.schemaName, t.tableName)
	rs, err := h2.conn.Query(query)
	if err != nil {
		err = errors.Wrapf(err, "could not open H2 constraint metadata for %s.%s ", t.schemaName, t.tableName)
		return
	}
	defer rs.Close()
	result, err = collectConstraints(rs)
	if err != nil {
		err = errors.Wrapf(err, "could not read H2 constraint metadata for %s.%s ", t.schemaName, t.tableName)
		return
	}
	return
}

type index struct {
	schemaName string
	tableName  string
	indexName  string
	typeName   string
	consName   string
	nonUnique  bool
	columns    []*indexColumn
}
type indexColumn struct {
	columnName string
	position   int
	asc        string
}

func collectIndexes(rs IResultSet) (result []*index, err error) {
	var curIndex *index
	for rs.Next() {
		var i = new(index)
		var ic = new(indexColumn)

		err = rs.Scan(
			&i.schemaName,
			&i.tableName,
			&i.indexName,
			&i.nonUnique,
			&i.typeName,
			&i.consName,
			&ic.columnName,
			&ic.position,
			&ic.asc,
		)
		if err != nil {
			return
		}
		if curIndex == nil || curIndex.indexName != i.indexName || curIndex.tableName != i.tableName || curIndex.schemaName != i.schemaName {
			if curIndex != nil {
				result = append(result, curIndex)
			}
			curIndex = i
		}
		curIndex.columns = append(curIndex.columns, ic)
	}
	if curIndex != nil {
		result = append(result, curIndex)
	}
	return

}

func (h2 *db) grabIndexes(t *table) (result []*index, err error) {
	query := `
	SELECT 
		C.TABLE_SCHEMA
		,C.TABLE_NAME
		,C.INDEX_NAME
		,C.NON_UNIQUE
		,C.INDEX_TYPE_NAME
		,C.CONSTRAINT_NAME
		,C.COLUMN_NAME
		,C.ORDINAL_POSITION
		,C.ASC_OR_DESC
	FROM INFORMATION_SCHEMA.INDEXES C
		INNER JOIN  INFORMATION_SCHEMA.TABLES T
			ON T.TABLE_CATALOG = C.TABLE_CATALOG
			AND T.TABLE_SCHEMA = C.TABLE_SCHEMA
			AND T.TABLE_NAME = C.TABLE_NAME
	WHERE T.TABLE_SCHEMA = '%s'
		AND T.TABLE_TYPE='TABLE' 
		AND T.STORAGE_TYPE = 'CACHED' 
		AND T.TABLE_NAME = '%s'
      ORDER BY C.INDEX_NAME,C.ORDINAL_POSITION asc
   `
	query = fmt.Sprintf(query, t.schemaName, t.tableName)
	rs, err := h2.conn.Query(query)
	if err != nil {
		err = errors.Wrapf(err, "could not read H2 index metadata for %s.%s ", t.schemaName, t.tableName)
		return
	}
	defer rs.Close()
	result, err = collectIndexes(rs)
	if err != nil {
		err = errors.Wrapf(err, "could not read H2 index metadata for %s.%s ", t.schemaName, t.tableName)
		return
	}
	return
}

func (i *index) createDDL() (result string) {
	var list []string
	for _, ic := range i.columns {
		list = append(list, ic.createDDL())
	}
	result = fmt.Sprintf("create index if not exists %s on %s.%s(%s);",
		i.indexName,
		i.schemaName,
		i.tableName,
		strings.Join(list, ","),
	)
	return
}
func (ic *indexColumn) createDDL() (result string) {
	asc := ic.asc
	if asc == "A" {
		asc = "ASC"
	} else if asc == "D" {
		asc = "DESC"
	}
	result = fmt.Sprintf("%s %s", ic.columnName, asc)
	return
}

func (i *index) dropDDL() (result string) {
	result = fmt.Sprintf("drop index if exists %s.%s;", i.schemaName, i.indexName)
	return
}

func (t *table) createDDL() (result string) {
	var list []string
	for _, ic := range t.columns {
		list = append(list, ic.createDDL())
	}
	result = fmt.Sprintf("create table if not exists %s.%s(%s);",
		t.schemaName,
		t.tableName,
		strings.Join(list, ","),
	)
	return
}

func (t *table) dropDDL() (result string) {
	result = fmt.Sprintf("drop table if exists %s.%s;", t.schemaName, t.tableName)
	return
}

func (c *column) createDDL() (result string) {
	//TODO: accomplish it
	pgColumn := c
	switch strings.ToUpper(c.typeName) {
	case "BIGINT", "BOOLEAN", "INTEGER", "SMALLINT":
		pgColumn.typeName = strings.ToLower(c.typeName)
	case "DECIMAL":
		pgColumn.typeName = "numeric"
	case "VARCHAR":
		pgColumn.typeName = "character varying"
	case "DOUBLE":
		pgColumn.typeName = "double precision"
	case "TIMESTAMP":
		pgColumn.typeName = "timestamp without time zone"
	case "CLOB", "BLOB", "VARBINARY":
		pgColumn.typeName = "bytea"
	}

	return
}
