package test

import (
	"fmt"

	"github.com/wallace5303/ee-go/eruntime"
)

func Info() {

	// base
	fmt.Println("ENV:", eruntime.ENV)
	fmt.Println("Port:", eruntime.Port)
	fmt.Println("SSL:", eruntime.SSL)
	fmt.Println("Debug:", eruntime.Debug)

	// base dir
	fmt.Println("BaseDir:", eruntime.BaseDir)
	fmt.Println("PublicDir:", eruntime.PublicDir)

	// user dir
	fmt.Println("UserHomeDir:", eruntime.UserHomeDir)
	fmt.Println("UserHomeConfDir:", eruntime.UserHomeConfDir)
	fmt.Println("WorkDir:", eruntime.WorkDir)
	fmt.Println("DataDir:", eruntime.DataDir)
	fmt.Println("TmpDir:", eruntime.TmpDir)
}
