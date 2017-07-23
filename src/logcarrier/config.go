package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	cron "gopkg.in/robfig/cron.v2"
	yaml "gopkg.in/yaml.v2"
)

// Config structure
type Config struct {
	Listen      string   `yaml:"listen"`
	ListenDebug string   `yaml:"listen_debug"`
	WaitTimeout Duration `yaml:"wait_timeout"`
	Key         string   `yaml:"key"`
	LogFile     string   `yaml:"logfile"`

	Compression struct {
		Method CompressionMethod `yaml:"method"`
		Level  uint              `yaml:"level"`
	} `yaml:"compression"`

	Buffers struct {
		Input   Size `yaml:"input"`
		Framing Size `yaml:"framing"`
		ZSTDict Size `yaml:"zstdict"`

		Connections int `yaml:"connections"`
		Dumps       int `yaml:"dumps"`
		Logrotates  int `yaml:"logrotates"`
	} `yaml:"buffers"`

	Workers struct {
		Router     int `yaml:"route"`
		Dumper     int `yaml:"dumper"`
		Logrotater int `yaml:"logrotater"`

		FlusherSleep Duration `yaml:"flusher_sleep"`
	} `yaml:"workers"`

	Files struct {
		Root     string      `yaml:"root"`
		RootMode os.FileMode `yaml:"root_mode"`
		Name     string      `yaml:"name"`
		Rotation string      `yaml:"rotation"`
	} `yaml:"files"`

	Links struct {
		enabled  bool
		Root     string      `yaml:"root"`
		RootMode os.FileMode `yaml:"root_mode"`
		Name     string      `yaml:"name"`
		Rotation string      `yaml:"rotation"`
	} `yaml:"links"`

	Logrotate struct {
		Method   LogrotateMethod `yaml:"method"`
		Schedule string          `yaml:"schedule"`
	} `yaml:"logrotate"`
}

// sensible defaults
func initConfig(config *Config) {
	config.Listen = "0.0.0.0:1466"
	config.ListenDebug = ""
	config.WaitTimeout = Duration(60 * time.Second)
	config.Key = "key"
	config.LogFile = ""

	config.Compression.Method = Raw
	config.Compression.Level = 0

	config.Buffers.Input = 128 * 1024
	config.Buffers.Framing = 256 * 1024
	config.Buffers.ZSTDict = 128 * 1024
	config.Buffers.Connections = 1024
	config.Buffers.Dumps = 512
	config.Buffers.Logrotates = 512

	config.Workers.Router = 1024
	config.Workers.Dumper = 24
	config.Workers.Logrotater = 48
	config.Workers.FlusherSleep = Duration(time.Second * 30)

	config.Files.Root = "./logs"
	config.Files.RootMode = 0755
	config.Files.Name = "${dir}?/${name}"
	config.Files.Rotation = "${dir}/${name}-${ time | %Y.%m.%d-%H }"

	config.Logrotate.Method = LogrotatePeriodic
	config.Logrotate.Schedule = "0 */1 * * *"
}

// LoadConfig loads config from given file
func LoadConfig(filePath string) (res Config) {
	var err error
	initConfig(&res)
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read configuration file `\033[1m%s\033[0m`: \033[31m%s\033[0m\n", filePath, err)
			os.Exit(1)
		}
	}()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(data, &res); err != nil {
		return
	}
	lengths := map[string]string{
		"root":     res.Links.Root,
		"name":     res.Links.Name,
		"rotation": res.Links.Rotation,
	}

	if len(res.Files.Root) == 0 {
		err = fmt.Errorf("files.root must not be empty")
		return
	}
	if len(res.Files.Rotation) == 0 {
		err = fmt.Errorf("files.rotation must not be empty")
		return
	}

	// Check schedule format
	c := cron.New()
	if _, err = c.AddFunc(res.Logrotate.Schedule, func() {}); err != nil {
		err = fmt.Errorf("Malformed schedule `%s`: %s", res.Logrotate.Schedule, err)
		return
	}

	//
	max := 0
	maxarg := ""
	maxv := ""
	min := len(res.Links.Root)
	minarg := "root"
	for k, v := range lengths {
		if len(v) > max {
			max = len(v)
			maxarg = k
			maxv = v
		}
		if len(v) < min {
			min = len(v)
			minarg = k
		}
	}
	if min == 0 && max == 0 || (min > 0 && max > 0) {
		res.Links.enabled = min > 0
		return
	}
	err = fmt.Errorf(
		"links.* must be either all empty or all full: links.%s is empty and links.%s is not (=%s)",
		minarg, maxarg, maxv,
	)
	return
}
