package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver       string        `mapstructure:"DB_DRIVER"`
	DBSource       string        `mapstructure:"DB_SOURCE"`
	ListenAddress  string        `mapstructure:"LISTEN_ADDRESS"`
	SecretKey      string        `mapstructure:"SECRET_KEY"`
	Duration       time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RequireAuth    bool          `mapstructure:"REQUIRE_AUTH"`
	RazorpayKey    string        `mapstructure:"RAZORPAY_KEY"`
	RazorpaySecret string        `mapstructure:"RAZORPAY_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // override with env vars if set
	// viper.SetEnvPrefix("")    // optional: set a prefix if needed (e.g. "MYAPP")
	viper.AllowEmptyEnv(true) // allow empty vars
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	fmt.Println("DB_DRIVER:", viper.Get("DB_DRIVER"))
	fmt.Println("DB_SOURCE:", viper.Get("DB_SOURCE"))
	fmt.Println("LISTEN_ADDRESS:", viper.Get("LISTEN_ADDRESS"))
	fmt.Println("SECRET_KEY:", viper.Get("SECRET_KEY"))
	fmt.Println("ACCESS_TOKEN_DURATION:", viper.Get("ACCESS_TOKEN_DURATION"))
	fmt.Println("REQUIRE_AUTH:", viper.Get("REQUIRE_AUTH"))
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err // other unexpected error
		}
		fmt.Println("No app.env found, continuing with environment variables...")
	}
	keys := []string{
		"DB_DRIVER", "DB_SOURCE", "LISTEN_ADDRESS",
		"SECRET_KEY", "ACCESS_TOKEN_DURATION", "REQUIRE_AUTH",
	}
	for _, key := range keys {
		if err := viper.BindEnv(key); err != nil {
			return config, err
		}
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	fmt.Printf("Loaded config: %+v\n", config)
	return
}
