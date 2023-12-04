package eapp

import (
	"path/filepath"

	"github.com/wallace5303/ee-go/econfig"
	"github.com/wallace5303/ee-go/eerror"
	"github.com/wallace5303/ee-go/ehelper"
	"github.com/wallace5303/ee-go/eruntime"
	"github.com/wallace5303/ee-go/estatic"
	"github.com/wallace5303/ee-go/eutil"
)

// read electron package.json
func ReadPackage() map[string]any {
	var ret map[string]any

	if eruntime.IsDev() {
		// 优先读项目中的 (构建后，不嵌入是没有的)
		pkgPath := filepath.Join(eruntime.BaseDir, "package.json")
		if ehelper.FileIsExist(pkgPath) {
			ret = eutil.ReadJsonStrict(pkgPath)
		}
	}

	if len(ret) == 0 {
		// 读嵌入的
		staticCfg := econfig.GetStatic()
		if staticCfg["enable"] == true {
			ret = estatic.ReadJsonStrict("public/package.json")
		} else {
			// 读外部的
			pkgPath := filepath.Join(eruntime.BaseDir, "public/package.json")
			if ehelper.FileIsExist(pkgPath) {
				ret = eutil.ReadJsonStrict(pkgPath)
			}
		}
	}

	if len(ret) == 0 {
		eerror.ThrowWithCode("The package.json does not exist!", eerror.ExitPackageFile)
	}

	return ret
}
