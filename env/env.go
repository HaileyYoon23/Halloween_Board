package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Dict struct {
	AppPort		int
	Domain		string
}

var dictInst Dict

func init() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8000
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = fmt.Sprintf("http://localhost:%d/",port)
	}
	domain = strings.TrimSuffix(domain, "/") + "/"

	dictInst = Dict{
		AppPort: port,
		Domain: domain,
	}
}

func GetEnv() *Dict {
	return &dictInst
}