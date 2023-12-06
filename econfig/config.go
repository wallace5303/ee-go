package econfig

import (
	"path/filepath"

	"github.com/wallace5303/ee-go/eerror"
	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/eruntime"
	"github.com/wallace5303/ee-go/estatic"
	"github.com/wallace5303/ee-go/eutil"

	"github.com/spf13/viper"
)

var (
	Vip *viper.Viper
)

// Initialize config
func Init() {
	var defaultCfg map[string]any
	var envCfg map[string]any
	defaultCfgName := "public/config/config.default.json"
	envCfgName := "public/config/config.prod.json"

	// dev
	if eruntime.IsDev() || eruntime.IsDebug() {
		// 优先读项目中的（构建后，项目中的是不存在的）
		defaultConfigPath := filepath.Join(eruntime.BaseDir, "go", "config", "config.default.json")
		devConfigPath := filepath.Join(eruntime.BaseDir, "go", "config", "config.local.json")
		if ehelper.FileIsExist(defaultConfigPath) && ehelper.FileIsExist(devConfigPath) {
			defaultCfg = eutil.ReadJsonStrict(defaultConfigPath)
			envCfg = eutil.ReadJsonStrict(devConfigPath)
		}
	}

	if len(defaultCfg) == 0 || len(envCfg) == 0 {
		// 读 嵌入的StaticFS
		if estatic.FileIsExist(defaultCfgName) && estatic.FileIsExist(envCfgName) {
			defaultCfg = estatic.ReadJsonStrict(defaultCfgName)
			envCfg = estatic.ReadJsonStrict(envCfgName)
		} else {
			// 读 外部的 （config 没有被嵌入）
			defaultConfigPath := filepath.Join(eruntime.BaseDir, defaultCfgName)
			devConfigPath := filepath.Join(eruntime.BaseDir, envCfgName)
			if ehelper.FileIsExist(defaultConfigPath) && ehelper.FileIsExist(devConfigPath) {
				defaultCfg = eutil.ReadJsonStrict(defaultConfigPath)
				envCfg = eutil.ReadJsonStrict(devConfigPath)
			}
		}
	}

	// 都没有，直接报错
	if len(defaultCfg) == 0 || len(envCfg) == 0 {
		eerror.ThrowWithCode("The config file does not exist !", eerror.ExitConfigFileNotExist)
	}

	// merge
	ehelper.Mapserge(envCfg, defaultCfg, nil)

	Vip = viper.New()
	for key, value := range defaultCfg {
		Vip.Set(key, value)
	}

	//fmt.Println("defaultCfg: ", Vip.AllSettings())
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
func Get(key string) any {
	return Vip.Get(key)
}

// Return all config as a map[string]any.
func GetAll() map[string]any {
	return Vip.AllSettings()
}

// Return logger config as a map[string]any.
func GetLogger() map[string]any {
	cfg := Vip.Get("logger")
	logCfg, ok := cfg.(map[string]any)
	if !ok {
		eerror.ThrowWithCode("Get logger config error !", eerror.ExitConfigLogErr)
	}
	return logCfg
}

// Return coreLogger config as a map[string]any.
func GetCoreLogger() map[string]any {
	cfg := Vip.Get("core_logger")
	logCfg, ok := cfg.(map[string]any)
	if !ok {
		eerror.ThrowWithCode("Get core_logger config error !", eerror.ExitConfigCoreLogErr)
	}
	return logCfg
}

// Return http config as a map[string]any.
func GetHttp() map[string]any {
	cfg := Vip.Get("http")
	httpCfg, ok := cfg.(map[string]any)
	if !ok {
		eerror.ThrowWithCode("Get http config error !", eerror.ExitConfigHttpErr)
	}

	return httpCfg
}

// Return static config as a map[string]any.
func GetStatic() map[string]any {
	cfg := Vip.Get("static")
	staticCfg, ok := cfg.(map[string]any)
	if !ok {
		eerror.ThrowWithCode("Get static config error !", eerror.ExitConfigStaticErr)
	}
	return staticCfg
}
