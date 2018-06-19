package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type testRS struct {
	index int
	data  [][]interface{}
}

func (rs testRS) Next() bool {
	return rs.index < len(rs.data)
}
func (rs *testRS) Scan(dst ...interface{}) error {
	for i, v := range rs.data[rs.index] {
		switch dv := dst[i].(type) {
		case *string:
			*dv = reflect.ValueOf(v).String()
		case *int64:
			*dv = reflect.ValueOf(v).Int()
		case *int:
			*dv = int(reflect.ValueOf(v).Int())
		case *float64:
			*dv = reflect.ValueOf(v).Float()
		case *bool:
			*dv = reflect.ValueOf(v).Bool()
		case *[]byte:
			*dv = reflect.ValueOf(v).Bytes()

		default:
			panic(fmt.Sprintf("testRS::Scan: type %T is not recognized\n", dv))
		}

	}
	rs.index++
	return nil
}

/*
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
*/
func TestCollectTables(t *testing.T) {
	rs := new(testRS)
	for _, tb := range tableTestData {
		for _, cl := range tb.columns {
			var row []interface{}
			row = append(row, tb.schemaName)
			row = append(row, tb.tableName)
			row = append(row, cl.columnName)
			row = append(row, cl.position)
			row = append(row, cl.defValue)
			row = append(row, cl.nullable)
			row = append(row, cl.typeName)
			row = append(row, cl.charLength)
			row = append(row, cl.numPrec)
			row = append(row, cl.numScale)
			rs.data = append(rs.data, row)
		}
	}
	result, err := collectTables(rs)
	if err != nil {
		t.Error(errors.Wrap(err, "Error while scanning table/column values"))
	}
	if !reflect.DeepEqual(result, tableTestData) {
		t.Error(errors.New("table/column collection is wrong"))
	}
	return
}

/*
C.TABLE_SCHEMA
		,C.TABLE_NAME
		,C.INDEX_NAME
		,C.NON_UNIQUE
		,C.INDEX_TYPE_NAME
		,C.CONSTRAINT_NAME
		,C.COLUMN_NAME
		,C.ORDINAL_POSITION
		,C.ASC_OR_DESC
*/

func TestCollectIndexes(t *testing.T) {
	rs := new(testRS)

	for _, idx := range indexTestData {
		for _, cl := range idx.columns {
			var row []interface{}
			row = append(row, idx.schemaName)
			row = append(row, idx.tableName)
			row = append(row, idx.indexName)
			row = append(row, idx.nonUnique)
			row = append(row, idx.typeName)
			row = append(row, idx.consName)
			row = append(row, cl.columnName)
			row = append(row, cl.position)
			row = append(row, cl.asc)
			rs.data = append(rs.data, row)
		}
	}
	result, err := collectIndexes(rs)
	if err != nil {
		t.Error(errors.Wrap(err, "Error while scanning index/column values"))
	}

	if !reflect.DeepEqual(result, indexTestData) {
		t.Error(errors.New("index/column collection is wrong"))
	}
	return
}
