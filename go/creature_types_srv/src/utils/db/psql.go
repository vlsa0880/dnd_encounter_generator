package db

import (
	"fmt"

	"creature_types_srv/src/utils/env"
)

func DsnByEnv() *string {
	address := env.Get("DB_PSQL_ADDRESS", "localhost")
	name := env.Get("DB_PSQL_NAME", "postgres")
	user := env.Get("DB_PSQL_USER", "postgres")
	pass := env.Get("DB_PSQL_PASS", "test")
	port := env.Get("DB_PSQL_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s",
		*address,
		*user,
		*name,
		*pass,
		*port,
	)
	return &dsn
}
