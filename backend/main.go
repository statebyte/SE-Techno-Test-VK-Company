package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

type ContainerInfo struct {
	ContainerID   string    `json:"container_id"`
	ContainerName string    `json:"container_name"`
	IPAddress     string    `json:"ip_address"`
	Status        string    `json:"status"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var db *sql.DB

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASSWORD", "pass")
	dbName := getEnv("DB_NAME", "containers_db")
	serverPort := getEnv("PORT", "8080")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка sql.Open: %v", err)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Не удалось подключиться к базе: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("goose.Up: %v", err)
	}

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/api/containers", getContainersHandler)
	r.POST("/api/containers", postContainersHandler)

	log.Printf("Backend запущен на порту %s", serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func getContainersHandler(c *gin.Context) {
	rows, err := db.Query(`SELECT container_id, container_name, ip_address, status, updated_at 
		FROM containers 
		ORDER BY updated_at DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []ContainerInfo
	for rows.Next() {
		var ci ContainerInfo
		if err := rows.Scan(&ci.ContainerID, &ci.ContainerName, &ci.IPAddress, &ci.Status, &ci.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = append(result, ci)
	}
	c.JSON(http.StatusOK, result)
}

func postContainersHandler(c *gin.Context) {
	var containers []ContainerInfo
	if err := c.ShouldBindJSON(&containers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Невалидный JSON"})
		return
	}

	for _, container := range containers {
		_, err := db.Exec(`
			INSERT INTO containers (container_id, container_name, ip_address, status, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (container_id) DO UPDATE 
			SET container_name = EXCLUDED.container_name,
			    ip_address = EXCLUDED.ip_address,
				status = EXCLUDED.status,
				updated_at = EXCLUDED.updated_at
		`,
			container.ContainerID,
			container.ContainerName,
			container.IPAddress,
			container.Status,
			time.Now(),
		)
		if err != nil {
			log.Printf("Ошибка при вставке/обновлении контейнера %s: %v", container.ContainerID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getEnv(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}
