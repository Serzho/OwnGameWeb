#!/bin/bash
set -e

# Directory containing SQL files
SQL_DIR="/docker-entrypoint-initdb.d"

# Loop through all SQL files in the directory and execute them
for sql_file in "$SQL_DIR"/*.sql; do
    if [ -f "$sql_file" ]; then
        echo "Executing $sql_file..."
        psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$sql_file"
    fi
done

echo "All SQL files executed successfully."
