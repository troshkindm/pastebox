#!/bin/sh
# Резервная копия БД (формат custom) из сервиса db; запускается в каталоге compose-проекта,
# копия сохраняется в соседний каталог ../backups.
set -eu
mkdir -p ../backups
file="../backups/pastebox-$(date +%Y%m%d-%H%M%S).dump"
docker compose exec -T db sh -c 'pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" -Fc' > "$file"
echo "Backup contents (first lines of the table of contents):"
docker compose exec -T db pg_restore --list < "$file" | grep -E 'TABLE|INDEX|CONSTRAINT' | head -n 6
ls -lh "$file"
