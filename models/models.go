package models

type SelectModel struct {
	TableName  string            `json:"tablename"`
	Columns    []string          `json:"columns"`
	Conditions []SelectCondition `json:"conditions"`
	Limit      int
}

type SelectCondition struct {
	ColumnPath     string `json:"columnpath"`
	LogicalType    string `json:"logicaltype"`
	ComparisonType string `json:"comparisontype"`
	Value          string `json:"value"`
}

type Pagination struct {
	PrevLink string `json:"prev_link"`
	SelfLink string `json:"self_link"`
	NextLink string `json:"next_link"`
}
