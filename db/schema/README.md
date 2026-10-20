# db/schema

Актуальная схема базы данных PasteBox: `schema.sql`.

Файл получен штатной утилитой PostgreSQL `pg_dump --schema-only` (скрипт `scripts/dump-schema.sh`) из тестового стенда после применения всех миграций из каталога `migrations/`. При каждом изменении схемы файл выгружается заново и коммитится вместе с миграцией.

Выгрузка (в корне проекта, в каталоге с `docker-compose.yaml`):

```sh
sh scripts/dump-schema.sh > db/schema/schema.sql
```
