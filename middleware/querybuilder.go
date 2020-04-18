package middleware

import (
	"fmt"
	"strings"

	"GODAPP/models"

	_ "github.com/lib/pq" //TODO: change to another driver and use sqlx
)

func GetSelectQuery(data models.SelectModel, limit int, offset int) string {
	query := `SELECT row_to_json(d) FROM( `
	query += getColumns(data)

	if len(data.Conditions) > 0 {
		query += getFilters(data.Conditions)
	}
	if limit != 0 {
		query += getLimitOffset(limit, offset)
	}
	query += `)d`
	return query
}

func GetCountQuery(data models.SelectModel) string {
	query := fmt.Sprintf("SELECT COUNT(id) from %s ", data.TableName)
	query += getFilters(data.Conditions)
	return query
}

func GetMergeQuery(data models.MergeModel) string {
	query := fmt.Sprintf("INSERT INTO %s(data) SELECT row_to_json(d) FROM( ", data.TargetTable)
	query += getMergeColumns(data)

	if len(data.Conditions) > 0 {
		query += getFilters(data.Conditions)
	}
	query += `)d`
	return query
}

func getColumns(data models.SelectModel) string {
	query := `SELECT `
	for _, s := range data.Columns {
		//todo: do another symbol
		query += fmt.Sprintf(` data #> '{ %s}' as "%s", `, strings.Replace(s, ".", ",", -1), s)
	}
	query = query[:len(query)-2]
	query += ` FROM ` + data.TableName

	return query
}

func getFilters(data []models.SelectCondition) string {
	query := " WHERE "

	for i, condition := range data {
		if i != 0 {
			query += mapLogicalType(condition)
		}
		query += mapCondition(condition)
	}
	return query
}

func mapLogicalType(data models.SelectCondition) string {
	switch data.LogicalType {
	case "and":
		return " AND "
	default:
		return " OR "
	}
}

func mapCondition(data models.SelectCondition) string {
	switch data.ComparisonType {
	case "lt":
		return fmt.Sprintf(` data #> '{ %s }' < '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "let":
		return fmt.Sprintf(` data #> '{ %s }' <= '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "gt":
		return fmt.Sprintf(` data #> '{ %s }' > '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "get":
		return fmt.Sprintf(` data #> '{ %s }' >= '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "like":
		return fmt.Sprintf(` data #>> '{ %s }' like '%%%s%%' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	case "eq":
		return fmt.Sprintf(` data #>> '{ %s }' = '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	default:
		return fmt.Sprintf(` data #>> '{ %s }' = '%s' `, strings.Replace(data.ColumnPath, ".", ",", -1), data.Value)
	}
}

func getLimitOffset(limit int, offset int) string {
	return fmt.Sprintf(" ORDER BY id ASC LIMIT %d OFFSET %d", limit, offset)
}

func getMergeColumns(data models.MergeModel) string {
	query := `SELECT `
	for srcCol, targCol := range data.MergeColumns {
		//todo: do another symbol
		query += fmt.Sprintf(` data #> '{ %s}' as "%s", `, strings.Replace(srcCol, ".", ",", -1), targCol)
	}
	query = query[:len(query)-2]
	query += ` FROM ` + data.SourceTable

	return query
}
