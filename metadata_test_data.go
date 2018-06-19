package main

var tableTestData = []*table{
	&table{
		schemaName: "public",
		tableName:  "table_info",
		columns: []*column{
			&column{
				schemaName: "public",
				tableName:  "table_info",
				columnName: "id",
				typeName:   "bigint",
				position:   1,
				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   0,
			},
			&column{
				schemaName: "public",
				tableName:  "table_info",
				columnName: "name",
				typeName:   "varchar",
				position:   2,
				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   0,
			},
			&column{
				schemaName: "public",
				tableName:  "table_info",
				columnName: "dumped",
				typeName:   "boolean",
				position:   3,

				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   1,
			},
			&column{
				schemaName: "public",
				tableName:  "table_info",
				columnName: "std",
				typeName:   "float",
				position:   4,

				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   1,
			},
		},
	},
	&table{
		schemaName: "public",
		tableName:  "column_info",
		columns: []*column{
			&column{
				schemaName: "public",
				tableName:  "column_info",
				columnName: "id",
				typeName:   "bigint",
				position:   1,

				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   0,
			},
			&column{
				schemaName: "public",
				tableName:  "column_info",
				columnName: "name",
				typeName:   "varchar",
				position:   2,

				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   0,
			},
			&column{
				schemaName: "public",
				tableName:  "column_info",
				columnName: "indexed",
				typeName:   "boolean",
				position:   3,

				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   0,
			},
			&column{
				schemaName: "public",
				tableName:  "column_info",
				columnName: "hash_unique_count",
				typeName:   "bigint",
				position:   4,
				charLength: 50,
				numScale:   0,
				numPrec:    30,
				nullable:   1,
			},
		},
	},
}
