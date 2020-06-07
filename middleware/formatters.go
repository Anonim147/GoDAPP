package middleware

import (
	"GODAPP/models"
	"fmt"
	"regexp"
	"strings"
)

func FormatForJsonPath(path string) string {
	path = strings.ReplaceAll(path, `.`, `"."`)
	path = strings.ReplaceAll(path, `."[]"`, `[*]`)
	path = fmt.Sprintf(`$."%s"`, path)
	return path
}

func FormatForJSONColumn(tablename string, path string) string {
	path = strings.ReplaceAll(path, `.`, `,`)
	return fmt.Sprintf(` %s.data #> '{ %s }' `, tablename, path)
}

func FormatForJSONWhere(tablename string, columns []string) string {
	fmt.Println(len(columns))
	condition := ` where `
	for index, column := range columns {
		fmt.Println(index)
		condition += fmt.Sprintf(` %s = %s`, FormatForJSONColumn("c1", column), FormatForJSONColumn("c2", column))
		if index != (len(columns) - 1) {
			condition += ` OR `
			fmt.Println(condition)
		}
	}
	return condition
}

func TransformKeys(keys []models.TableKey) []models.TableKey {
	var newKeys []models.TableKey
	re := regexp.MustCompile(`\.\[\]$`)
	for index, key := range keys {
		if !re.MatchString(key.KeyName) {
			if key.KeyType == `array` {
				key = сheckIfHasChildren(index, keys)
			}
			newKeys = append(newKeys, key)
		}
	}
	return newKeys
}

func сheckIfHasChildren(index int, keys []models.TableKey) models.TableKey {
	item := keys[index]
	re := regexp.MustCompile(fmt.Sprintf(`^%s\.\[\]\.\w`, item.KeyName))
	if index < len(keys) {
		for i := index + 1; i < len(keys); i++ {
			if re.MatchString(keys[i].KeyName) {
				return models.TableKey{
					KeyName: item.KeyName,
					KeyType: "complex array",
				}
			}
		}
	}
	return item
}
