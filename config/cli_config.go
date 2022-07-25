package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/RobbieMcKinstry/colorformatter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	defaultTimeout             = 5 * time.Second
	TextLogger     LoggerType  = "text"
	JSONLogger     LoggerType  = "json"
	ColorLine      ColorOption = "line"
	ColorOff       ColorOption = "off"
	ColorOn        ColorOption = "on"
	ColorAuto      ColorOption = "auto"
)

func loggerTypeFromStr(v string) (LoggerType, error) {
	var val = LoggerType(strings.ToLower(v))
	switch val {
	case TextLogger:
		fallthrough
	case JSONLogger:
		return val, nil
	default:
		return TextLogger, fmt.Errorf("Provided logger type was not valid. Expected either \"text\" or \"json\", but found \"%s\"", v)
	}
}

func colorOptionFromStr(v string) (ColorOption, error) {
	var val = ColorOption(strings.ToLower(v))
	switch val {
	case ColorLine:
		fallthrough
	case ColorOff:
		fallthrough
	case ColorOn:
		fallthrough
	case ColorAuto:
		return val, nil
	default:
		return ColorOff, fmt.Errorf("Provided color option not valid. Expected one of \"line\", \"off\", \"on\", or \"auto\", but found \"%s\"", v)
	}
}

type (
	CLIConfig struct {
		Color     ColorOption
		InCluster bool
		LogLevel  logrus.Level
		LoggerType
		Namespace   string
		ServiceName string
		Timeout     time.Duration
	}
	LoggerType  string
	ColorOption string
)

func NewConfigFromCLI(ctx *cli.Context) (*CLIConfig, error) {
	var typ, err = loggerTypeFromStr(ctx.String("logger-type"))
	if err != nil {
		return nil, err
	}
	level, err := logrus.ParseLevel(ctx.String("log-level"))
	if err != nil {
		return nil, err
	}
	color, err := colorOptionFromStr(ctx.String("color"))

	var conf = &CLIConfig{
		Color:       color,
		InCluster:   ctx.Bool("in-cluster"),
		Namespace:   ctx.String("namespace"),
		LoggerType:  typ,
		LogLevel:    level,
		ServiceName: ctx.String("service-name"),
		Timeout:     defaultTimeout,
	}
	conf.setupLogger()
	return conf, nil
}

// setupLogger will configure Logrus logging using
// configuration from the environment.
func (conf *CLIConfig) setupLogger() {
	// To set up the logger, we decide whether we want
	// JSON logging or human-oriented logging.
	// We also set the log level.
	var formatter logrus.Formatter
	// Only use text-based logging when developing locally.
	// If this app is deployed to any other environment,
	// use JSON logging.
	if conf.LoggerType == TextLogger {
		switch conf.Color {
		case ColorLine:
			formatter = &colorformatter.Colored{}
		case ColorOn:
			formatter = &logrus.TextFormatter{
				ForceColors: true,
			}
		case ColorOff:
			formatter = &logrus.TextFormatter{
				DisableColors: true,
			}
		case ColorAuto:
			formatter = &logrus.TextFormatter{}
		}
	} else {
		formatter = &logrus.JSONFormatter{}
	}
	logrus.SetFormatter(formatter)
	logrus.SetLevel(conf.LogLevel)
}
