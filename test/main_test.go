package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/DarrelA/starter-go-postgresql/app"
	"github.com/DarrelA/starter-go-postgresql/db"
	"github.com/gofiber/fiber/v2"
)

var (
	rdbmsInstance       db.RDBMS
	inMemoryDbInstance  db.InMemoryDB
	appInstance         *fiber.App
	authServiceInstance *fiber.App
)

func TestMain(m *testing.M) {
	// Change to the project root directory from `/test` directory
	err := os.Chdir("../")
	if err != nil {
		fmt.Println("error changing directory:", err)
		os.Exit(1)
	}

	rdbmsInstance, inMemoryDbInstance = app.CreateDBConnections()
	appInstance, authServiceInstance = app.ConfigureAppInstance()
	go app.StartServer()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestMainFunc(t *testing.T) {
	logFilePath := "./deployments/logs/app.log"

	t.Run("test if app.log is writable", func(t *testing.T) {
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("cannot get cwd: %s", err)
			}

			t.Logf("@cwd: %s", cwd)
			t.Fatal("expected log file to be created, but it does not exist")
		}

		file, err := os.OpenFile(logFilePath, os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("expected log file to be writable, but it is not: %v", err)
		}

		file.Close()
	})
}

func TestAuthService(t *testing.T) {

	t.Run("testing /register endpoint", TestRegisterEndpoint)
}
