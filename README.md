# kvstore

Мини key-value хранилище на Go: TCP-сервер, in-memory мапа, простой текстовый протокол, TTL, конкурентный доступ.

Пет-проект для закрепления конкурентности и сетевого программирования на Go. Не замена Redis: только базовый минимум.

## Возможности

- In-memory хранилище ключ-значение
- TCP-сервер, горутина на каждого клиента
- Простой текстовый протокол (не RESP)
- Команды: `GET`, `SET`, `DELETE`, `TTL`
- TTL: `SET key value <seconds>`
- Конкурентный доступ через `sync.RWMutex`
- Graceful shutdown по SIGINT/SIGTERM
- Фоновый уборщик просроченных ключей (раз в секунду)
- Ленивое удаление просроченных ключей при `GET`

## Что НЕ входит

LRU-вытеснение, persistence, Pub/Sub, пайплайнинг, RESP-протокол, pprof, CI/CD, Docker, шардирование мапы.

## Стек

Только стандартная библиотека Go (1.21+). Внешних зависимостей нет.

- `net`, `bufio` — TCP-сервер и чтение/запись
- `sync.RWMutex`, `sync.WaitGroup` — конкурентный доступ и shutdown
- `time.Ticker` — чистка TTL
- `os/signal` — graceful shutdown
- `testing`, `b.RunParallel` — тесты и бенчмарки

## Запуск

```bash
go run ./cmd/server
# слушает :6379
```

## Пример сессии
$ nc localhost 6379
SET foo bar 10
+OK
GET foo
bar
TTL foo
9
(ждём 11 секунд)
GET foo
-ERR not found
TTL foo
-2
SET baz qux
+OK
TTL baz
-1
SET foo
-ERR SET requires 2 or 3 arguments
TTL
-ERR TTL requires 1 argument, got 0
UNKNOWN x
-ERR unknown command: UNKNOWN
