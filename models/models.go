package models


type SelectModel struct {
	TableName string `json:tablename`,
	ColumnNames []string `json:columnnames`,
	Conditions []SelectCondition `json:conditions`,
}

type SelectCondition struct {
	ColumnPath string `json:column_path`,
	ComparisonType string `json:comparison_type`, 
	Value string `json:value`,
}
