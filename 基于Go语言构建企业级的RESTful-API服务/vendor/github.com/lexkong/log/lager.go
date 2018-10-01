package log

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lexkong/log/lager"
	"gopkg.in/yaml.v2"
)

// constant values for logrotate parameters
const (
	RollingPolicySize = "size"
	LogRotateDate     = 1
	LogRotateSize     = 10
	LogBackupCount    = 7
)

// Lager struct for logger parameters
type Lager struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

//PassLagerCfg is the struct for lager information(passlager.yaml)
type PassLagerCfg struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

// Logger is the global variable for the object of lager.Logger
var Logger lager.Logger

// logFilePath log file path
var logFilePath string

// PassLagerDefinition is having the information about loging
var PassLagerDefinition *PassLagerCfg = DefaultLagerDefinition()

// Initialize Build constructs a *Lager.Logger with the configured parameters.
func Initialize(writers, loggerLevel, loggerFile, rollingPolicy string, logFormatText bool,
	LogRotateDate, LogRotateSize, LogBackupCount int) {
	lag := &Lager{
		Writers:        writers,
		LoggerLevel:    loggerLevel,
		LoggerFile:     loggerFile,
		LogFormatText:  logFormatText,
		RollingPolicy:  rollingPolicy,
		LogRotateDate:  LogRotateDate,
		LogRotateSize:  LogRotateSize,
		LogBackupCount: LogBackupCount,
	}

	Logger = newLog(lag)
	initLogRotate(logFilePath, lag)
	return
}

// newLog new log
func newLog(lag *Lager) lager.Logger {
	checkPassLagerDefinition(lag)

	if filepath.IsAbs(lag.LoggerFile) {
		createLogFile("", lag.LoggerFile)
		logFilePath = filepath.Join("", lag.LoggerFile)
	} else {
		createLogFile(os.Getenv("CHASSIS_HOME"), lag.LoggerFile)
		logFilePath = filepath.Join(os.Getenv("CHASSIS_HOME"), lag.LoggerFile)
	}
	writers := strings.Split(strings.TrimSpace(lag.Writers), ",")
	if len(strings.TrimSpace(lag.Writers)) == 0 {
		writers = []string{"stdout"}
	}
	LagerInit(Config{
		Writers:       writers,
		LoggerLevel:   lag.LoggerLevel,
		LoggerFile:    logFilePath,
		LogFormatText: lag.LogFormatText,
	})

	logger := NewLogger(lag.LoggerFile)
	return logger
}

// checkPassLagerDefinition check pass lager definition
func checkPassLagerDefinition(lag *Lager) {
	if lag.LoggerLevel == "" {
		lag.LoggerLevel = "DEBUG"
	}

	if lag.LoggerFile == "" {
		lag.LoggerFile = "log/chassis.log"
	}

	if lag.RollingPolicy == "" {
		log.Println("RollingPolicy is empty, use default policy[size]")
		lag.RollingPolicy = RollingPolicySize
	} else if lag.RollingPolicy != "daily" && lag.RollingPolicy != RollingPolicySize {
		log.Printf("RollingPolicy is error, RollingPolicy=%s, use default policy[size].", lag.RollingPolicy)
		lag.RollingPolicy = RollingPolicySize
	}

	if lag.LogRotateDate <= 0 || lag.LogRotateDate > 10 {
		lag.LogRotateDate = LogRotateDate
	}

	if lag.LogRotateSize <= 0 || lag.LogRotateSize > 50 {
		lag.LogRotateSize = LogRotateSize
	}

	if lag.LogBackupCount < 0 || lag.LogBackupCount > 100 {
		lag.LogBackupCount = LogBackupCount
	}
}

// createLogFile create log file
func createLogFile(localPath, outputpath string) {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, outputpath)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, outputpath)), "\\", "/", -1), os.ModePerm)
	} else if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, outputpath), "\\", "/", -1), os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

// readPassLagerConfigFile is unmarshal the paas lager configuration file(lager.yaml)
func InitWithFile(lagerFile string) error {
	if lagerFile == "" {
		log.Printf("log config file is empty, use default config: `%s`\n", marshalDefinition())
		return Init()
	}

	passLagerDef := PassLagerCfg{}
	yamlFile, err := ioutil.ReadFile(lagerFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v, use default config: `%s`\n", err, marshalDefinition())
		return Init()
	}

	err = yaml.Unmarshal(yamlFile, &passLagerDef)
	if err != nil {
		log.Printf("Unmarshal: %v, use default config: `%s`\n", err, marshalDefinition())
		return Init()
	}

	PassLagerDefinition = &passLagerDef
	return Init()
}

func InitWithConfig(passLagerDef *PassLagerCfg) error {
	PassLagerDefinition = passLagerDef
	return Init()
}

func DefaultLagerDefinition() *PassLagerCfg {
	cfg := PassLagerCfg{
		Writers:        "stdout,file",
		LoggerLevel:    "DEBUG",
		LoggerFile:     "logs/chassis.log",
		LogFormatText:  false,
		RollingPolicy:  RollingPolicySize,
		LogRotateDate:  1,
		LogRotateSize:  10,
		LogBackupCount: 7,
	}

	return &cfg
}

func Init() error {
	Initialize(PassLagerDefinition.Writers, PassLagerDefinition.LoggerLevel,
		PassLagerDefinition.LoggerFile, PassLagerDefinition.RollingPolicy,
		PassLagerDefinition.LogFormatText, PassLagerDefinition.LogRotateDate,
		PassLagerDefinition.LogRotateSize, PassLagerDefinition.LogBackupCount)

	return nil
}

func marshalDefinition() string {
	data, _ := json.Marshal(PassLagerDefinition)
	return string(data)
}
