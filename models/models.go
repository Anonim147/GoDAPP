package models

type SelectModel struct {
	TableName  string            //`json:tablename`
	Columns    []string          //`json:columns`
	Conditions []SelectCondition //`json:conditions`
	Limit      int
}

type SelectCondition struct {
	ColumnPath     string //`json:columnpath`
	LogicalType    string
	ComparisonType string //`json:comparisontype`
	Value          string //`json:value`
}
