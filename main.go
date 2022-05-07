package main

import (
	"embed"
	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"timenoteWeb/auth"
	"timenoteWeb/config"
	"timenoteWeb/routes"
	"timenoteWeb/utils"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var templates embed.FS

var logger *logrus.Logger

func main() {

	// setup logger
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&formatter.Formatter{
		HideKeys:        false,
		TimestampFormat: "[2006-01-02 15:04:05]",
	})
	logger.Out = os.Stdout

	// If no data folder, create one
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.Mkdir("./data", 0777)
		if err != nil {
			return
		}
	}

	// Load config
	appConfig, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}

	// init gin
	r := gin.New()

	// setup recovery
	r.Use(gin.Recovery())

	// load templates
	templates, err := template.ParseFS(templates, "templates/*.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(templates)

	// setup static files
	r.Use(utils.StaticServer("/static", http.FS(static)))

	// setup logger
	r.Use(utils.LoggerMiddleware(logger))

	// setup webdav
	r.Use(utils.DavServer(
		"/dav",
		"./data",
		func(c *gin.Context) bool {
			return auth.BasicAuth(c, appConfig)
		},
		func(req *http.Request, err error) {
			logger.Error(err)
		}),
	)

	if gin.Mode() == gin.DebugMode {
		routes.DebugRoute(r, appConfig, logger)
	}

	// run
	srv := &http.Server{
		Addr:    appConfig.Server.Listen + ":" + strconv.Itoa(appConfig.Server.Port),
		Handler: r,
	}
	srv.SetKeepAlivesEnabled(true)
	logger.Info("Listening on " + appConfig.Server.Listen + ":" + strconv.Itoa(appConfig.Server.Port))
	logger.Fatal(srv.ListenAndServe())
}
