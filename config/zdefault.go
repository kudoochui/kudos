package config

import (
	"github.com/kudoochui/kudos/utils"
	"os"
	"path/filepath"
)

const (
	DEV  = "development"
	PROD = "production"
)

var (
	appConfigProvider = "json"
	AppPath string
	WorkPath string
	RunMode string
)

type AppConfig struct {
	innerConfig Configer
}

func init() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	WorkPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	RunMode = DEV
	if envRunMode := os.Getenv("RUNMODE"); envRunMode != "" {
		RunMode = envRunMode
	}
}

func NewAppConfig(filename string) (*AppConfig, error) {
	appConfigPath := filepath.Join(WorkPath, "bin", "conf", filename)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "bin","conf", filename)
		if !utils.FileExists(appConfigPath) {
			return &AppConfig{innerConfig: NewFakeConfig()}, nil
		}
	}

	ac, err := NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &AppConfig{ac}, nil
}

func (b *AppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *AppConfig) String(key string) string {
	if v := b.innerConfig.String(RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *AppConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(RunMode + "::" + key); len(v) > 0 {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *AppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *AppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *AppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *AppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *AppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *AppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *AppConfig) GetSection(section string) (map[string]string, error) {
	if v, err := b.innerConfig.GetSection(RunMode + "::" + section); err == nil {
		return v, nil
	}
	return b.innerConfig.GetSection(section)
}

func (b *AppConfig) GetMap(key string) (map[string]interface{}, error) {
	if v, err := b.innerConfig.GetMap(RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.GetMap(key)
}

func (b *AppConfig) GetEnvMap() (map[string]interface{}, error) {
	if v, err := b.innerConfig.GetMap(RunMode); err == nil {
		return v, nil
	}
	return b.innerConfig.GetMap(RunMode)
}

func (b *AppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
