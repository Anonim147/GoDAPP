package middleware

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadFileFromLink(url string) (string, error) {
	filename := getRandomString() + `.dat`
	filePath := filepath.Join(os.TempDir(), filename)
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return filePath, nil
}

func getDataFromFile(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return removeSpacesAndNewLines(data), nil
}

func removeSpacesAndNewLines(data []byte) []byte {
	//space := regexp.MustCompile(`\s+`)
	//s := space.ReplaceAllString(string(data), " ")
	s := strings.ReplaceAll(string(data), `'`, `''`)
	return []byte(s)
}

func getRandomString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
