package main

import (
	"log"
	"net/http"
	"os"

	"finalgo/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	puerto = ":8080"
)

// @title           Swagger Clinica Odontologica API
// @version         1.0
// @description     Documentacion de la clinica odontologica
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	// uso la libreria de env.
	errE := godotenv.Load()
	if errE != nil {
		log.Fatal(errE)
	}

	// inicio el router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// API de prueba
	router.GET("/ping", middleware.Authenticate(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"mensaje": "pong",
		})
	})

	// arranco a funcionar las apis
	if err := router.Run(puerto); err != nil {
		panic(err)
	}
}

// func runApp(db *sql.DB, engine *gin.Engine) {
// 	// Run the application.
// 	router := routes.NewRouter(engine, db)
// 	// Map all routes.
// 	router.MapRoutes()
// 	if err := engine.Run(puerto); err != nil {
// 		log.Fatalf("Error al ejecutar la aplicación: %v", err)
// 	}
// }

// func connectDB() *sql.DB {
// 	var (
// 		dbUsername = "root"
// 		dbPassword = ""
// 		dbHost     = "localhost"
// 		dbPort     = "3306"
// 		dbName     = "my_db"
// 	)

// 	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

// 	db, err := sql.Open("mysql", dataSource)
// 	if err != nil {
// 		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatalf("Error al conectar con la base de datos: %v", err)
// 	}

// 	return db
// }
