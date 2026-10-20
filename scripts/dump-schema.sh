#!/bin/sh
# Выгружает схему БД (без данных, владельцев и прав доступа) из сервиса db текущего compose-проекта.
# Из вывода убираются комментарии и служебные строки \restrict, чтобы дампы разных стендов можно было сравнивать.
set -eu
docker compose exec -T db sh -c 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --schema-only --no-owner --no-privileges' \
  | grep -Ev '^(--|\\restrict|\\unrestrict)' | cat -s
