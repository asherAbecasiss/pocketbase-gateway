package main

import (
	"apiGateway/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"gopkg.in/yaml.v2"
)

type Api struct {
	App       *pocketbase.PocketBase
	SecretKey string
	Config    types.Config
}

func CreateServer() *Api {

	app := Api{App: pocketbase.New()}
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &app.Config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	return &app

}

func defaultPublicDir() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// most likely ran with go run
		return "./pb_public"
	}

	return filepath.Join(os.Args[0], "./pb_public")
}

func main() {
	pbApp := CreateServer()

	var publicDirFlag string

	// add "--publicDir" option flag
	pbApp.App.RootCmd.PersistentFlags().StringVar(
		&publicDirFlag,
		"publicDir",
		defaultPublicDir(),
		"the directory to serve static files",
	)

	// load js files to allow loading external JavaScript migrations
	jsvm.MustRegister(pbApp.App, jsvm.Config{
		HooksWatch: true, // make this false for production
	})

	// register the `migrate` command
	migratecmd.MustRegister(pbApp.App, pbApp.App.RootCmd, migratecmd.Config{
		TemplateLang: migratecmd.TemplateLangJS, // or migratecmd.TemplateLangGo (default)
		Automigrate:  true,
	})

	// // call this only if you want to auditlog tables named in AUDITLOG env var
	// auditlog.Register(pbApp.App)

	// // call this only if you want to use the configurable "hooks" functionality
	// hooks.PocketBaseInit(pbApp.App)

	//routes

	pbApp.InitRouting(publicDirFlag)

	if err := pbApp.App.Start(); err != nil {
		log.Fatal(err)
	}
}
