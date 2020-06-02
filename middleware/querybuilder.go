package middleware

import (
	"fmt"

	"GODAPP/models"
)

func GetQueryForTableList() string {
	return `select tablename from pg_catalog.pg_tables where schemaname !='pg_catalog' and schemaname != 'information_schema'`
}

func GetJSONKeysQuery(tableName string) string {
	return fmt.Sprintf(`
	with recursive extract_all as
	(
		select 
			key as path, 
			value
		from %s
		cross join lateral jsonb_each(data)
	union all
		select
			path || '.' || coalesce(obj_key, '[]'),
			coalesce(obj_value, arr_value)
		from extract_all
		left join lateral 
			jsonb_each(case jsonb_typeof(value) when 'object' then value end) 
			as o(obj_key, obj_value) 
			on jsonb_typeof(value) = 'object'
		left join lateral 
			jsonb_array_elements(case jsonb_typeof(value) when 'array' then value end) 
			with ordinality as a(arr_value, arr_key)
			on jsonb_typeof(value) = 'array'
		where obj_key is not null or arr_key is not null
	)
	select distinct path, jsonb_typeof(value)
	from extract_all where jsonb_typeof(value)<>'null' 
	order by path;`, tableName)
}

func GetSelectQuery(data models.SelectModel, limit int, offset int) string {
	query := `SELECT row_to_json(d) FROM( `
	query += getColumns(data)

	if len(data.Conditions) > 0 {
		query += getFilters(data.Conditions)
	}
	query += `)d`
	if limit != 0 {
		query += getLimitOffset(limit, offset)
	}

	return query
}

func GetCountQuery(data models.SelectModel) string {
	return fmt.Sprintf("SELECT COUNT(dt) from (%s) dt  ", GetSelectQuery(data, 0, 0))
}

func GetQueryForCopying(tablename string, path string) string {
	return fmt.Sprintf(`copy %s (data) from '%s'`, tablename, path)
}

func GetQueryForParseJSONARRAY(tablename string, data string) string {
	return fmt.Sprintf(`insert into %s (data) select json_array_elements('%s'::json)`, tablename, data)
}

func GetQueryForParseJSON(tablename string, data string) string {
	return fmt.Sprintf(`insert into %s (data) values('%s'::jsonb)`, tablename, data)
}

func GetQueryForCreatingHash(tablename string) string {
	return fmt.Sprintf(`update %s set hash = md5(data::text);`, tablename)
}

func GetQueryForUpdateTable(tablename string, temptable string) string {
	return fmt.Sprintf(`insert into %s (data, hash) select data, hash from %s 
		where not exists(select 1 from %s where %s.hash = %s.hash);`, tablename, temptable, tablename, tablename, temptable)
}

func GetQueryForClearTable(tablename string) string {
	return fmt.Sprintf(`delete from %s`, tablename)
}

func GetQueryForDropTable(tablename string) string {
	return fmt.Sprintf(`drop table %s`, tablename)
}

func getColumns(data models.SelectModel) string {
	query := `SELECT `
	for _, s := range data.Columns {
		query += fmt.Sprintf(` jsonb_path_query(data, '%s') as "%s", `, FormatForJsonPathPath(s), s)
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
	case "<":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ < %s)') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case "<=":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ <= %s)') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case ">":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ > %s)') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case ">=":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ >= %s)') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case "=":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ == %s)') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case "like":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ like_regex "%s" flag "i")') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case "!=":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ != "%s")') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	case "have key":
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (exists(@.%s'))') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	default:
		return fmt.Sprintf(` jsonb_path_exists(data, '%s ? (@ == "%s")') `,
			FormatForJsonPathPath(data.ColumnPath), data.Value)
	}
}

func getLimitOffset(limit int, offset int) string {
	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

func getLinkForPagination(host string, limit int, offset int) string {
	return fmt.Sprintf(`http://%s/api/get_data&limit=%d&offset=%d`, host, limit, offset) //TO DO: привязати ссилку глобально
}
