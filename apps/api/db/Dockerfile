FROM postgis/postgis:16-3.4-alpine

COPY apps/api/migrations/001_init_discovery_schema.sql /docker-entrypoint-initdb.d/001_init_discovery_schema.sql
COPY apps/api/seeds/002_seed_discovery_data.sql /docker-entrypoint-initdb.d/002_seed_discovery_data.sql
