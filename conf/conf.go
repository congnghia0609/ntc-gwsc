/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper
var env string

// Init wdir
func Init(wdir string) {
	//// init Configuration
	environment := flag.String("e", "development", "run project with mode [-e development | test | production]")
	flag.Usage = func() {
		fmt.Println("Usage: ./[appname] -e development | test | production")
		os.Exit(1)
	}
	flag.Parse()
	log.Printf("============== NConf environment: %s", *environment)
	InitEnv(wdir, *environment)
}

// InitEnv is an exported method that takes the environment starts the viper (external lib) and
// returns the configuration struct.
func InitEnv(wdir string, environment string) {
	log.Printf("============== NConfig Init Environment: %s ==============", environment)
	var err error
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(environment)
	v.AddConfigPath(filepath.Join(wdir, "conf"))
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	config = v
	env = environment
}

// JoinPath join directory
func JoinPath(basedir string, path string) string {
	return filepath.Join(basedir, path)
}

// GetConfig get config
func GetConfig() *viper.Viper {
	return config
}

// GetEnv get env
func GetEnv() string {
	return env
}
