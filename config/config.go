package config

import (
	"errors"
	"fmt"
	"log"
	"micro/config/env"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

var (
	// Global config
	confs = Config{}
	// lock  = sync.Mutex{}
)

var Env env.EnvMode

func init() {
	str := os.Getenv("env")
	switch str {
	case string(env.PRODUCTION), string(env.DEVELOPMENT), string(env.STAGE), string(env.TEST):
		Env = env.EnvMode(str)
	case "":
		Env = env.DEVELOPMENT
		//panic("env must be set")
	default:
		panic(fmt.Sprintf("invalid env value: [%s]", str))
	}
}

// Config is grpc of configs we need for project
type Config struct {
	CORS     string   `yaml:"cors" `
	Debug    bool     `yaml:"debug"`
	Service  Service  `yaml:"usecase" `
	Jaeger   Jaeger   `yaml:"jaeger" `
	Etcd     Etcd     `yaml:"etcd" `
	Redis    Redis    `yaml:"redis" `
	Postgres Database `yaml:"database"`
	Nats     NATS     `yaml:"nats" `
	JWT      JWT      `yaml:"jwt" json:"jwt" `
	Logstash Logstash `yaml:"logstash"`
}

func Validate(c any) error {
	errmsg := ""
	numFields := reflect.TypeOf(c).NumField()
	for i := 0; i < numFields; i++ {
		fieldType := reflect.TypeOf(c).Field(i)
		tagval, ok := fieldType.Tag.Lookup("required")
		isRequired := ok && tagval == "true"
		if !isRequired {
			continue
		}
		fieldVal := reflect.ValueOf(c).Field(i)
		if fieldVal.Kind() == reflect.Struct {
			if err := Validate(fieldVal.Interface()); err != nil {
				errmsg += fmt.Sprintf("%s > [%v], ", fieldType.Name, err)
			}
		} else {
			if fieldVal.IsZero() {
				errmsg += fmt.Sprintf("%s, ", fieldType.Name)
			}
		}
	}
	if errmsg == "" {
		return nil
	}
	return errors.New(errmsg)
}

func C() *Config {
	return &confs
}

// init configs
func init() {
	dir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	loadConfigs()
	// viper.OnConfigChange(func(in fsnotify.Event) {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// 	lastUpdate := viper.GetTime("fsnotify")
	// 	if time.Since(lastUpdate) < time.Second {
	// 		return
	// 	}
	// 	viper.Set("fsnotify", time.Now())
	// 	log.Println("config file changed. restarting...")
	// 	shutdowner.Shutdown()
	// })
	// viper.WatchConfig()
}

func loadConfigs() {
	must(viper.Unmarshal(&confs),
		"could not unmarshal config file")
	must(Validate(confs),
		"some required configurations are missing")
	log.Printf("configs loaded from file successfully \n")
}

func must(err error, logMsg string) {
	if err != nil {
		log.Println(logMsg)
		panic(err.Error())
	}
}
