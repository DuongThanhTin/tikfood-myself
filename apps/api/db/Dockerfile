FROM postgis/postgis:16-3.4-alpine

COPY apps/api/migrations/*.sql /docker-entrypoint-initdb.d/
COPY apps/api/seeds/*.sql /docker-entrypoint-initdb.d/
