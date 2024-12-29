package eboot

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wallace5303/ee-go/eapp"
	"github.com/wallace5303/ee-go/econfig"
	"github.com/wallace5303/ee-go/eerror"
	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/ehttp"
	"github.com/wallace5303/ee-go/elog"
	"github.com/wallace5303/ee-go/eos"
	"github.com/wallace5303/ee-go/eruntime"
	"github.com/wallace5303/ee-go/estatic"
)

var (
// cmdENV  = "prod" // 'dev' 'prod'
// cmdPort = "0"
)

// ee-go instance type
type Ego struct {
}

// run program
func (ego *Ego) Run() {
	eapp.Run()
}

// create new ego
func New(staticFS embed.FS) *Ego {
	// args
	environment := flag.String("env", "prod", "dev/prod")
	baseDir := flag.String("basedir", "./", "base directory")
	port := flag.String("port", "0", "service port")
	ssl := flag.String("ssl", "false", "https/wss service")
	debug := flag.String("debug", "false", "debug")
	flag.Parse()

	// fmt.Println("cmdENV:", *environment)
	// fmt.Println("baseDir:", *baseDir)
	// fmt.Println("goport:", *port)
	// fmt.Println("ssl:", *ssl)

	eruntime.ENV = *environment
	eruntime.Debug, _ = strconv.ParseBool(*debug)
	eruntime.BaseDir = filepath.Join(eruntime.BaseDir, *baseDir)
	eruntime.Port = *port
	eruntime.SSL, _ = strconv.ParseBool(*ssl)

	// static "./public"
	estatic.StaticFS = staticFS

	initApp()

	// debug
	//test.Info()

	ego := &Ego{}
	return ego
}

// Initialize the app
func initApp() {

	// init dir
	eruntime.InitDir()

	// init config
	econfig.Init()

	pkg := eapp.ReadPackage()
	pkgName := pkg["name"].(string)
	if pkgName == "" {
		eerror.ThrowWithCode("The app name is required!", eerror.ExitAppNameIsEmpty)
	}
	eruntime.AppName = pkgName

	// init user dir
	initUserDir()

	// init logger
	coreLogCfg := econfig.GetCoreLogger()
	logCfg := econfig.GetLogger()
	elog.InitCoreLog(coreLogCfg)
	elog.InitLog(logCfg)

	// http server
	httpCfg := econfig.GetHttp()
	if httpCfg["enable"] == true {
		ehttp.CreateServer(httpCfg)
	}

}

// Initialize user directory
func initUserDir() {
	eruntime.UserHomeDir, _ = eos.GetUserHomeDir()
	eruntime.UserHomeConfDir = filepath.Join(eruntime.UserHomeDir, ".config", eruntime.AppName)
	if !ehelper.FileIsExist(eruntime.UserHomeConfDir) {
		if err := os.MkdirAll(eruntime.UserHomeConfDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create user home conf folder [%s] failed: %s", eruntime.UserHomeConfDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateUserHomeConfDir)
		}
	}
	hiddenAppName := fmt.Sprintf(".%s", eruntime.AppName)
	eruntime.UserHomeAppDir = filepath.Join(eruntime.UserHomeDir, hiddenAppName)
	if !ehelper.FileIsExist(eruntime.UserHomeAppDir) {
		if err := os.MkdirAll(eruntime.UserHomeAppDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create user home app folder [%s] failed: %s", eruntime.UserHomeAppDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateUserHomeAppDir)
		}
	}

	eruntime.WorkDir = eruntime.BaseDir
	if eruntime.IsProd() {
		// userhome/.appname
		eruntime.WorkDir = eruntime.UserHomeAppDir
	}
	if !ehelper.FileIsExist(eruntime.WorkDir) {
		if err := os.MkdirAll(eruntime.WorkDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create work folder [%s] failed: %s", eruntime.WorkDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateWorkDir)
		}
	}

	eruntime.DataDir = filepath.Join(eruntime.WorkDir, "data")
	if !ehelper.FileIsExist(eruntime.DataDir) {
		if err := os.MkdirAll(eruntime.DataDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create data folder [%s] failed: %s", eruntime.DataDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateDataDir)
		}
	}

	logDir := filepath.Join(eruntime.WorkDir, "logs")
	if !ehelper.FileIsExist(logDir) {
		if err := os.MkdirAll(logDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create logs folder [%s] failed: %s", logDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateLogDir)
		}
	}
	elog.SetLogDir(logDir)

	eruntime.TmpDir = filepath.Join(eruntime.DataDir, "tmp")
	os.RemoveAll(eruntime.TmpDir)
	if !ehelper.FileIsExist(eruntime.TmpDir) {
		if err := os.MkdirAll(eruntime.TmpDir, 0755); err != nil && !os.IsExist(err) {
			errMsg := fmt.Sprintf("create tmp folder [%s] failed: %s", eruntime.TmpDir, err)
			eerror.ThrowWithCode(errMsg, eerror.ExitCreateTmpDir)
		}
	}
	os.Setenv("TMPDIR", eruntime.TmpDir)
	os.Setenv("TEMP", eruntime.TmpDir)
	os.Setenv("TMP", eruntime.TmpDir)
}
