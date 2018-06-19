package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
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
		case *[]byte:
			*dv = reflect.ValueOf(v).Bytes()

		default:
			panic(fmt.Sprintf("testRS::Scan: type %T is not recognized\n", dv))
		}

	}
	rs.index++
	return nil
}

func TestCollectTables(t *testing.T) {
	rs := new(testRS)
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
		t.Error(err)
	}

	if !reflect.DeepEqual(result, tableTestData) {
		t.Error(errors.New("collection is wrong"))
	}

	return

}
