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

var indexTestData = []*index{
	&index{
		schemaName: "public",
		tableName:  "table_info",
		indexName:  "table_info_pk",
		consName:   "table_info_pk",
		nonUnique:  false,
		typeName:   "PRIMARY KEY",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "id", position: 1, asc: "A",
			},
		},
	},
	&index{
		schemaName: "public",
		tableName:  "table_info",
		indexName:  "table_info_wfn",
		nonUnique:  true,
		typeName:   "INDEX",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "name", position: 1, asc: "A",
			},
			&indexColumn{
				columnName: "workflow_id", position: 2, asc: "A",
			},
		},
	},
	&index{
		schemaName: "public",
		tableName:  "column_info",
		indexName:  "column_info_pk",
		consName:   "column_info_pk",
		nonUnique:  false,
		typeName:   "PRIMARY KEY",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "id", position: 1, asc: "A",
			},
		},
	},
	&index{
		schemaName: "public",
		tableName:  "column_info",
		indexName:  "column_info_uq_tn",
		consName:   "column_info_uq_tn",
		nonUnique:  false,
		typeName:   "UNIQUE",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "table_info_id", position: 1, asc: "A",
			},
			&indexColumn{
				columnName: "name", position: 2, asc: "A",
			},
		},
	},
	&index{
		schemaName: "public",
		tableName:  "column_info",
		indexName:  "column_info_uq_tp",
		consName:   "column_info_uq_tp",
		nonUnique:  false,
		typeName:   "UNIQUE",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "table_info_id", position: 1, asc: "A",
			},
			&indexColumn{
				columnName: "position", position: 2, asc: "A",
			},
		},
	},
	&index{
		schemaName: "public",
		tableName:  "column_info",
		indexName:  "column_info_name",
		consName:   "column_info_name",
		nonUnique:  true,
		typeName:   "INDEX",
		columns: []*indexColumn{
			&indexColumn{
				columnName: "name", position: 1, asc: "A",
			},
			&indexColumn{
				columnName: "real_type", position: 2, asc: "D",
			},
			&indexColumn{
				columnName: "indexed", position: 3, asc: "D",
			},
		},
	},
}
