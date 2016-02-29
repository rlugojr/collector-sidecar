package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/pborman/uuid"
)

func GetRootPath() (string, error) {
	return filepath.Abs("/")
}

func GetSystemName() string {
	os := runtime.GOOS
	osRunes := []rune(os)
	osRunes[0] = unicode.ToUpper(osRunes[0])
	return string(osRunes)
}

func GetCollectorId(collectorId string) string {
	id := collectorId
	if strings.HasPrefix(collectorId, "file:") {
		filePath := strings.SplitAfter(collectorId, ":")[1]
		err := FileExists(filePath)
		if err != nil {
			log.Info("collector-id file doesn't exist, generating a new one")
			CreatePathToFile(filePath)
			ioutil.WriteFile(filePath, []byte(RandomUuid()), 0644)
		}
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("Can not read collector-id file: ", err)
		}
		id = string(file)
	}

	if id != "" {
		log.Info("Using collector-id: ", id)
	}
	return id
}

func RandomUuid() string {
	return uuid.NewRandom().String()
}

func FileExists(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return err
	}
	return nil

}

func IsDir(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Error(err)
		return false
	}
	if fi.Mode().IsDir() {
		return true
	}
	return false

}

func AppendIfDir(dir string, appendix string) (string, error) {
	file, err := os.Open(dir)
	if err != nil {
		log.Error("Can not access ", dir)
		return dir, err
	}

	fileInfo, err := file.Stat()
	switch {
	case err != nil:
		return "", err
	case fileInfo.IsDir():
		appended := filepath.Join(dir, appendix)
		return appended, nil
	default:
		return dir, nil
	}
}

func CreatePathToFile(filepath string) error {
	dir := path.Dir(filepath)
	_, err := os.Open(dir)
	if err != nil {
		log.Info("Trying to create directory for: ", filepath)
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			log.Error("Not able to create directory path: ", dir)
			return err
		}
	}
	return nil
}

func SplitCommaList(list string) []string {
	if list == "" {
		return make([]string, 0)
	}
	return strings.Split(list, ",")
}
