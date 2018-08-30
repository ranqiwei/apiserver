package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Name string
}

//viper初始化配置
func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("apiserver")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

/*
//日志框架logrus初始化---未使用
func (c *Config) initLogrus() {
	l := logrus.New()

	if f := viper.GetString("log.JSONFormatter"); f == "true" || f == "True" {
		formatter1 := &logrus.JSONFormatter{}
		l.Formatter = new(&formatter1)
	}

	if viper.GetString("log.Logfile") == "true" {
		file, err := os.OpenFile("apiserver.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			l.Out = file
		} else {
			l.Info("Failed to log to file, using default stderr")
		}
	}
	l.Out = os.Stdout

	switch level := viper.GetString("log.Loglevel"); level {
	case "debug":
		l.Level = logrus.DebugLevel
	case "info":
		l.Level = logrus.InfoLevel
	default:
		l.Level = logrus.ErrorLevel
	}

}
*/

//日志框架初始化
func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}

	log.InitWithConfig(&passLagerCfg)
}

//监控配置文件并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

func Init(cfg string) error {
	c := Config{Name: cfg}
	//初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}
	//日志初始化
	c.initLog()
	//热加载
	c.watchConfig()
	return nil
}
