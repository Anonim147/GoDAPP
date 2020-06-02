package models

type SelectModel struct {
	TableName  string            `json:"tablename"`
	Columns    []string          `json:"columns"`
	Conditions []SelectCondition `json:"conditions"`
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

type InsertTableModel struct {
	TableName string `json:"tablename"`
	FilePath  string `json:"filepath"`
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Value   interface{} `json:"value"`
}

type TableKey struct {
	KeyName string `json:"keyname"`
	KeyType string `json:"keytype"`
}
