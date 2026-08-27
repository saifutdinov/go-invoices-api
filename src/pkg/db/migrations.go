package db

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type LastMigration struct {
	version string
}

func MigrateDB(db *sql.DB, databaseMigrationsPath string) {
	if databaseMigrationsPath == "" {
		databaseMigrationsPath = "./db/migrations"
	}

	files, err := os.ReadDir(databaseMigrationsPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Slice(files, func(i, z int) bool {
		splitI := strings.Split(files[i].Name(), "-")
		iInts := convertVersionToInts(splitI[0])
		splitZ := strings.Split(files[z].Name(), "-")
		zInts := convertVersionToInts(splitZ[0])

		return doesNeedApplyMigration(iInts, zInts)
	})

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrate_history (
			id SERIAL PRIMARY KEY,
			version varchar, 
			comment varchar, 
			created_at TIMESTAMP without time zone DEFAULT now()
		);`); err != nil {
		panic(err)
	}

	res, _ := db.Query("SELECT version FROM migrate_history ORDER BY id DESC LIMIT 1")
	defer res.Close()
	lastMigration := &LastMigration{}
	if !res.Next() {
		lastMigration.version = "0.0.0"
	} else if err = res.Scan(&lastMigration.version); err != nil {
		fmt.Println(err)
		return
	}

	lastMigrationParts := convertVersionToInts(lastMigration.version)

	for _, file := range files {
		if file.Name() == "main.go" {
			continue
		}

		split := strings.Split(file.Name(), "-")
		migrationVersion, migrationComment := convertVersionToInts(split[0]), split[1]
		migrationVersionString := split[0]
		if doesNeedApplyMigration(lastMigrationParts, migrationVersion) {
			sqlFiles, err := os.ReadDir(databaseMigrationsPath + "/" + migrationVersionString + "-" + migrationComment)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, sqlFile := range sqlFiles {
				sqlFileName := sqlFile.Name()
				if _, err := db.Exec(stringifySQL(databaseMigrationsPath + "/" + migrationVersionString + "-" + migrationComment + "/" + sqlFileName)); err != nil {
					fmt.Printf("Exec error for file %s/%s: %s\n", migrationVersionString, sqlFileName, err.Error())
					return
				}
			}

			if _, err := db.Exec("INSERT INTO migrate_history (version, comment) VALUES ($1, $2)", migrationVersionString, migrationComment); err != nil {
				// panic(err)
				fmt.Println(err)
				return
			}
		}
	}
}

func convertVersionToInts(version string) (ints [3]int) {
	versionParts := strings.Split(version, ".")
	if len(versionParts) != 3 {
		return ints
	}
	for i := 0; i < 3; i++ {
		ints[i], _ = strconv.Atoi(versionParts[i])
	}
	return ints
}

func doesNeedApplyMigration(lastMigration [3]int, checkingMigration [3]int) bool {
	if checkingMigration[0] > lastMigration[0] {
		return true
	}

	if checkingMigration[0] == lastMigration[0] {
		if checkingMigration[1] > lastMigration[1] {
			return true
		}
		if checkingMigration[1] == lastMigration[1] {
			if checkingMigration[2] > lastMigration[2] {
				return true
			}
		}
	}

	return false
}

func stringifySQL(pathToSql string) string {
	sqlQuery, err := os.ReadFile(pathToSql)

	if err == nil {
		return string(sqlQuery)
	}

	return err.Error()
}
