package eruntime

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	Version = "0.1.0"
	ENV     = "dev" // 'dev' 'prod'
	// progressBar  float64 // 0 ~ 100
	// progressDesc string  // description

	AppName   = ""
	Platform  = "pc" // pc | mobile | web
	IsExiting = false
	Debug     = false
)

var (
	BaseDir, _      = os.Getwd()
	PublicDir       string // electron-egg public directory
	UserHomeDir     string // OS user home directory
	UserHomeConfDir string // OS user home config directory
	UserHomeAppDir  string // OS user home app directory
	WorkDir         string // App working directory
	DataDir         string // data directory
	TmpDir          string // tmp directory
)

var (
	Port                = "0"
	SSL                 = false
	HttpServerIsRunning = false
)

func InitDir() {
	PublicDir = filepath.Join(BaseDir, "public")
}

// Pwd gets the path of current working directory.
func IsProd() bool {
	return (ENV == "prod")
}

func IsDev() bool {
	return (ENV == "dev")
}

func IsDebug() bool {
	return Debug
}

func Pwd() string {
	file, _ := exec.LookPath(os.Args[0])
	pwd, _ := filepath.Abs(file)

	return filepath.Dir(pwd)
}
