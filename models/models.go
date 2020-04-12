package models

type SelectModel struct {
	TableName  string            //`json:tablename`
	Columns    []string          //`json:columns`
	Conditions []SelectCondition //`json:conditions`
}

type SelectCondition struct {
	ColumnPath     string //`json:columnpath`
	ComparisonType string //`json:comparisontype`
	Value          string //`json:value`
}
