package property

import (
	"aerospike-cli/src/logger"
	"bufio"
	"github.com/magiconair/properties"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ConfigDirPath        = "CONFIG_DIR_PATH"
	PropertiesConfigName = "property.conf"
)

var Props = properties.MustLoadFile(filepath.Join(os.Getenv(ConfigDirPath), PropertiesConfigName), properties.UTF8)

func GetEnvOptions() string {
	filePath := filepath.Join(os.Getenv(ConfigDirPath), PropertiesConfigName)
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error.Println("failed to open the config file, %s", filePath)
		os.Exit(-1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error.Println("failed to close the config file, %s", filePath)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	pat := regexp.MustCompile("(^aerospike\\.client\\.host\\.)(.*?)(=)")
	var res []string
	for scanner.Scan() {
		submatch := pat.FindStringSubmatch(scanner.Text())
		if len(submatch) > 0 {
			res = append(res, submatch[2])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return strings.Join(res[:], ", ")
}
