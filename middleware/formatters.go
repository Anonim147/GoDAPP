package middleware

import (
	"GODAPP/models"
	"fmt"
	"regexp"
	"strings"
)

func FormatForJsonPathPath(path string) string {
	path = strings.ReplaceAll(path, `.`, `"."`)
	path = strings.ReplaceAll(path, `."[]"`, `[*]`)
	path = fmt.Sprintf(`$."%s"`, path)
	return path
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
