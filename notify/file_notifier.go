package notify

import (
	"fmt"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// FileNotifier implements basic Notifier using file storage
type FileNotifier struct {
	fileName string
}

// NewFileNotifier constructor
func NewFileNotifier() *FileNotifier {
	return &FileNotifier{}
}

// Notify ...
func (n *FileNotifier) Notify(fileName string) error {
	file, err := os.OpenFile(n.fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	now := time.Now().Format(time.RFC3339)
	if _, err := file.WriteString(fmt.Sprintf("%s\t%s\n", now, fileName)); err != nil {
		return err
	}
	return nil
}

// Init ...
func (n *FileNotifier) Init(data []byte) error {
	var config struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}
	n.fileName = config.Path
	return nil
}
