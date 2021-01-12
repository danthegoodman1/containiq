package setup
import (
	"io/ioutil"
	"path/filepath"
	"os"
	y "gopkg.in/yaml.v2"
	"github.com/sirupsen/logrus"
)

type GeneralConfig struct {
	Monitoring struct {
		Namespaces struct {
			Watch []string `yaml:"watch"`
		} `yaml:"namespaces"`
		Level struct {
			Watch []string `yaml:"watch"`
		} `yaml:"level"`
		Resource struct {
			Watch []string `yaml:"watch"`
		} `yaml:"resource"`
	} `yaml:"monitoring"`
}

type NotifyConfig struct {
	Source struct {
		Slack struct {
			Key     string `yaml:"Key"`
			Enabled bool   `yaml:"enabled"`
			Channel string `yaml:"channel"`
		} `yaml:"Slack"`
		Webhook struct {
			Enabled bool   `yaml:"enabled"`
			URL     string `yaml:"url"`
		} `yaml:"Webhook"`
	} `yaml:"Source"`
}
type Config struct {
	GeneralConfig
	NotifyConfig
}

func Setup() *Config {

	var generalConfig  GeneralConfig
	var notifyConfig NotifyConfig
	filePath := os.Getenv("CONFIG_FILE_LOCATION")
	filename, _ := filepath.Abs(filePath)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		logrus.Error(err)
	}

	err = y.Unmarshal(yamlFile, &generalConfig)
	if err != nil {
		logrus.Error(err)
	}

	filePathSecret := os.Getenv("NOTIFICATION_FILE_LOCATION")
	filenameSecret, _ := filepath.Abs(filePathSecret)
	yamlFileSecret, err := ioutil.ReadFile(filenameSecret)

	if err != nil {
		logrus.Error(err)
	}

	err = y.Unmarshal(yamlFileSecret, &notifyConfig)
	if err != nil {
		logrus.Error(err)
	}

	var c = &Config {
		GeneralConfig: generalConfig,
		NotifyConfig: notifyConfig,
	}
	return c

}
