package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/AnaLauraS/FinalGo/cmd/server/routes"
	"github.com/AnaLauraS/FinalGo/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/AnaLauraS/FinalGo/docs"
	_ "github.com/go-sql-driver/mysql"
)

const (
	puerto = ":9090"
)

func main() {
	// Recover from panic.
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	// Load the environment variables.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando variables de entorno: %v", err)
	}

	// Connect to the database.
	db := connectDB()

	// Create a new Gin engine with recovery middleware.
	engine := gin.Default()

	// Add the swagger handler.
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the application.
	runApp(db, engine)

	// Close the database connection.
	defer db.Close()
}

func runApp(db *sql.DB, engine *gin.Engine) {
	// Run the application.
	router := routes.NewRouter(engine, db)
	// Map all routes.
	router.MapRoutes()
	if err := engine.Run(puerto); err != nil {
		log.Fatalf("Error al ejecutar la aplicación: %v", err)
	}
}

func connectDB() *sql.DB {
	var (
		dbUsername = "root"
		dbPassword = ""
		dbHost     = "localhost"
		dbPort     = "3306"
		dbName     = "my_db"
	)

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	return db
}