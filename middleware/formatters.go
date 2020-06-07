package middleware

import (
	"GODAPP/models"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func FormatCSVtoJSON(path string) ([]byte, error) {
	csvFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		return nil, errors.New("File is empty or length of the lines are not the same")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headersArr = append(headersArr, headE)
	}
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, d := range content {
		buffer.WriteString("{")
		for j, y := range d {
			buffer.WriteString(`"` + headersArr[j] + `":`)
			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			if j < len(d)-1 {
				buffer.WriteString(",")
			}

		}
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	rawMessage := json.RawMessage(buffer.String())
	x, err := json.MarshalIndent(rawMessage, "", "  ")
	if err != nil {
		return nil, err
	}
	return x, nil
}
