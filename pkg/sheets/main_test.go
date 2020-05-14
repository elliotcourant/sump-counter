package sheets

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}
