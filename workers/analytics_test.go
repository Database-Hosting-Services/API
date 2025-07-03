package workers_test

import (
	"DBHS/config"
	"DBHS/workers"
	"log"
	"os"
	"testing"

	// "github.com/stretchr/testify/assert"
)

func beforeTest() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	config.Init(infoLog, errorLog)
}

func afterTest() {
	config.CloseDB()
}

func TestGatherAnalytics(t *testing.T) {
	beforeTest()
	defer afterTest()
	workers.GatherAnalytics(config.App)
}
