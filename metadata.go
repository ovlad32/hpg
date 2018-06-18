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
	columnName string
	typeName   string
	position   int
	nullable   int
	charLength int
	defValue   string
	numScale   int
	numPrec    int
}
type index struct {
	schemaName string
	tableName  string
	indexName  string
	typeName   string
	consName   string
	unique     bool
	columns    []*indexColumn
}
type indexColumn struct {
	columnName string
	asc        string
}
type constraint struct {
	consName   string
	references string
	columList  string
	typeName   string
	indexName  string
	checkExpr  string
}

func (h2 *db) grabMetadataTable(schemaName string) (result *schema, err error) {
	query := `
	SELECT 
		C.TABLE_NAME
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
	var tableName, prevTable string
	var current *table
	for rs.Next() {
		c := &column{}
		err = rs.Scan(
			&tableName,
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
			errors.Wrapf(err, "could not read table metadata")
			return
		}
		if result == nil {
			result = &schema{
				schemaName: schemaName,
			}
		}
		if prevTable != tableName {
			current = &table{
				schemaName: schemaName,
				tableName:  tableName,
			}
			result.tables = append(result.tables, current)
			prevTable = tableName
		}

		current.columns = append(current.columns, c)
	}
	return
}
func (h2 *db) grabConstraints(t table) (result []*constraint, err error) {

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
	//  var prev,curr string
	for rs.Next() {
		c := &constraint{}

		err = rs.Scan(
			&c.typeName,
			&c.consName,
			&c.columList,
			&c.indexName,
			&c.checkExpr,
		)

		if err != nil {
			err = errors.Wrapf(err, "could not read H2 constraint metadata for %s.%s ", t.schemaName, t.tableName)
			return
		}
		result = append(result, c)
	}
	return
}

func (h2 *db) grabIndexes(t table) (result []*index, err error) {
	query := `
	SELECT 
		C.NON_UNIQUE
		,C.INDEX_NAME
		,C.INDEX_TYPE_NAME
		,C.CONSTRAINT_NAME
		,C.COLUMN_NAME
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
	var curr *index
	for rs.Next() {
		i := &index{
			schemaName: t.schemaName,
			tableName:  t.tableName,
		}
		ic := new(indexColumn)
		err = rs.Scan(
			&i.unique,
			&i.indexName,
			&i.typeName,
			&i.consName,
			&ic.columnName,
			&ic.asc,
		)
		if err != nil {
			err = errors.Wrapf(err, "could not read H2 constraint metadata for %s.%s ", t.schemaName, t.tableName)
			return
		}
		if curr == nil || curr.indexName != i.indexName {
			if curr != nil {
				result = append(result, curr)
			}
			curr = i
		}
		curr.columns = append(curr.columns, ic)
	}
	result = append(result, curr)

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
	result = fmt.Sprintf("drop tab;e if exists %s.%s;", t.schemaName, t.tableName)
	return
}

func (c *column) createDDL() (result string) {
	//TODO: accomplish it
	return
}
