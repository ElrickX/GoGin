// Package panda_lib provides database utilities for syncing retailer data.
// Use to load and create config.txt
package panda_lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// =====================================================
// CONSTANTS
// =====================================================
const (
	DB_USER     = "elrick.bong"
	DB_PASS     = "PadnapA$slBong!2025"
	DB_PORT     = 3306
	DB_NAME     = "retail_hub"
	CONFIG_FILE = "config.txt"
)

const EmailTableStyle = `
<style>
	table { border-collapse: collapse; width: 100%; font-family: sans-serif; margin-bottom: 20px; }
	th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
	th { background-color: #f2f2f2; font-weight: bold; }
	tr:nth-child(even) { background-color: #fafafa; }
	tr:hover { background-color: #f1f1f1; }
	.header-title { font-size: 18px; font-weight: bold; margin-bottom: 10px; color: #333; }
	.warn-msg { color: #856404; background-color: #fff3cd; border: 1px solid #ffeeba; padding: 10px; margin-bottom: 10px; border-radius: 4px; }
</style>
`

// =====================================================
// LOAD CONFIG FILE
// =====================================================
func LoadOrCreateConfigKV(path string) (map[string]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultContent := "DB_HOST=localhost\n"
		if err := os.WriteFile(path, []byte(defaultContent), 0644); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line: %s", line)
		}
		cfg[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return cfg, nil
}

func ConfigDB() *sql.DB {
	configPath := resolveConfigPath()

	fmt.Println("Using config:", configPath)

	cfg, err := LoadOrCreateConfigKV(configPath)
	if err != nil {
		log.Fatal(err)
	}

	dbHost := cfg["DB_HOST"]
	if dbHost == "" {
		log.Fatal("DB_HOST not found in config.txt")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		DB_USER,
		DB_PASS,
		dbHost,
		DB_PORT,
		DB_NAME,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func resolveConfigPath() string {
	// 1️⃣ ENV override (MOST IMPORTANT)
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}

	// 2️⃣ Working directory (good for go run / go test in project root)
	if _, err := os.Stat(CONFIG_FILE); err == nil {
		return CONFIG_FILE
	}

	// 3️⃣ Executable directory (for compiled binary in production)
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		p := filepath.Join(exeDir, CONFIG_FILE)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// 4️⃣ Final fallback
	log.Fatalf("config file not found: %s", CONFIG_FILE)
	return ""
}
