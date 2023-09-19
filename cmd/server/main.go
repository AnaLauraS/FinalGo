package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"os"

	"finalgo/pkg/middleware"
	"finalgo/cmd/server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	puerto = ":8080"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	// Usa la librería de env.
	errE := godotenv.Load()
	if errE != nil {
		log.Fatal(errE)
	}

	// Inicia el router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// Conecta a la base de datos
	db := connectDB()

	// Ejecuta la aplicación
	runApp(db, router)

	// cierro BD
	defer db.Close()
}

func runApp(db *sql.DB, engine *gin.Engine) {
	// Ejecuta la aplicación.
	router := routes.NewRouter(engine, db)
	// Mapea todas las rutas.
	router.MapRoutes()
	if err := engine.Run(puerto); err != nil {
		log.Fatalf("Error al ejecutar la aplicación: %v", err)
	}
}

func connectDB() *sql.DB {
	var (
		dbUsername = "root"
		dbPassword = "1234"
		dbHost     = "localhost"
		dbPort     = "3306"
		dbName     = "my_db"
	)

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
	}

	err = db.Ping()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	return db
}
