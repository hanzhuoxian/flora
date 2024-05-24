package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel            = "log.level"
	flagFormat           = "log.format"
	flagEnableColor      = "log.enable-color"
	flagEnableCaller     = "log.enable-caller"
	flagOutputPaths      = "log.output-paths"
	flagErrorOutputPaths = "log.error-output-paths"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// Options contains the configuration for the logger.
type Options struct {
	Level            string   `json:"level"              mapstructure:"level"`
	Format           string   `json:"format"             mapstructure:"format"`
	EnableColor      bool     `json:"enable-color"       mapstructure:"enable-color"`
	EnableCaller     bool     `json:"enable-caller"      mapstructure:"enable-caller"`
	OutputPaths      []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
}

// NewOptions returns a new Options instance.
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel.String(),
		Format:           consoleFormat,
		EnableColor:      false,
		EnableCaller:     false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Validate validates the Options.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds the flags for Options.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.BoolVar(&o.EnableCaller, flagEnableCaller, o.EnableCaller, "Enable adding caller info in the logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
}

func (o *Options) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return ""
	}

	return string(data)
}
