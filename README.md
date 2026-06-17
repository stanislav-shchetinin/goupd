# goupd

CLI, которая для указанного Git-репозитория выводит данные о Go-модуле и список зависимостей, доступных для обновления.

На вход подаётся адрес Git-репозитория (или путь к локальному каталогу), на выходе — имя модуля, версия Go и список зависимостей, которые можно обновить (включая мажорные версии v2+).

## Установка

Требуется установленный Go (>= 1.22) и `git` в `PATH`.

```bash
go build -o goupd ./cmd/goupd
```

## Использование

```bash
goupd [flags] <git-repo-url|local-path>
```

Примеры:

```bash
goupd https://github.com/foo

goupd --ref v1.4.0 https://github.com/foo

goupd --format json --direct-only https://github.com/foo

goupd ./path/to/module
```

### Флаги

| Флаг            | По умолчанию | Описание                                                        |
| --------------- | ------------ | --------------------------------------------------------------- |
| `--ref`         | (default)    | Ветка, тег или коммит для checkout                              |
| `--format`      | `text`       | Формат вывода: `text` или `json`                                |
| `--major`       | `true`       | Искать мажорные (v2+) обновления через Go module proxy          |
| `--direct-only` | `false`      | Показывать только прямые зависимости                            |
| `--timeout`     | `2m`         | Общий таймаут на сетевые операции                               |

## Как это работает

1. Репозиторий клонируется (`git clone --depth 1`) во временный каталог; локальный путь используется напрямую.
2. `go.mod` парсится через `golang.org/x/mod/modfile` — извлекаются имя модуля, версия Go и список зависимостей.
3. Обновления в рамках текущего мажора берутся из `go list -m -u -json all`.
4. Мажорные обновления (v2+) обнаруживаются отдельно: для каждой прямой зависимости опрашивается Go module proxy по путям `path/vN`, поскольку `go list` такие апдейты не возвращает (меняется путь модуля).
5. Версии сравниваются и классифицируются (`patch` / `minor` / `major`) через `golang.org/x/mod/semver`.

## Пример вывода

```
Module:  github.com/acme/foo
Go:      1.22

DEPENDENCY                   CURRENT   LATEST    TYPE
github.com/pkg/errors        v0.9.0    v0.9.1    patch
github.com/spf13/cobra       v1.7.0    v1.8.1    minor
github.com/redis/go-redis    v6.15.9   v9.5.1    major (-> github.com/redis/go-redis/v9)
```

## Тесты

```bash
go test ./...
```
