package main

import (
	"fmt"
	"imgconv/manager"
	"imgconv/router"
	"imgconv/storage"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

func setupRoutes(engine *gin.Engine, router *router.Router) {
	engine.GET("/version", router.OnHEIFVersionRequest)
	engine.POST("/convert", router.OnSingleFileConversionRequest)
	engine.GET("/status/:id", router.OnGetConversionStatusRequest)
	engine.GET("/download/:id", router.OnDownloadConvertedImageRequest)
}

func callback(c *cli.Context) error {
	ip := c.String("ip")
	port := c.String("port")
	var stge storage.StorageRepo
	if c.String("file-storage") == "" {
		stge = storage.NewMemoryStorage()
	} else {
		stge = storage.NewFileStorage(c.String("file-storage"))
	}
	log.Printf("Serving on %s:%s\n", ip, port)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	engine := gin.Default()
	man := manager.NewConversionManager(stge)
	man.AddConversionListener(&manager.LoggerListener{})
	router := router.NewRouter(man)
	setupRoutes(engine, router)
	return engine.Run(fmt.Sprintf("%s:%s", ip, port))
}

func main() {
	app := cli.NewApp()
	app.Name = "IMGCONV"
	app.Usage = "A web server using libheif to convert images to png."
	app.Authors = []cli.Author{
		{Name: "JaskaranSM"},
	}
	app.Action = callback
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "port",
			Value: "6090",
			Usage: "port to listen on",
		},
		&cli.StringFlag{
			Name:  "ip",
			Value: "",
			Usage: "ip to listen on, by default webserver listens on localhost",
		},
		&cli.StringFlag{
			Name:  "file-storage",
			Value: "",
			Usage: "use file storage instead of memory storage",
		},
	}
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
