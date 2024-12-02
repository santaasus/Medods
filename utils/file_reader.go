package utils

import (
	// "fmt"
	"os"
	// "strings"
)

func FileByName(name string) ([]byte, error) {
	// path, _ := os.Getwd()
	// rootPath := strings.Split(path, "/Medods")[0]
	// if len(rootPath) == 0 {
	// 	err := fmt.Errorf("error for get root path for startServer")
	// 	return nil, err
	// }

	// In tests cases change rootPath + "/Medods" + "/" + name
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return data, nil
}
