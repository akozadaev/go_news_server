# Быстрый старт

## Вариант 1: Запуск с Docker Compose (рекомендуется)

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd go_news_server
   ```

2. **Запустите с помощью Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **Проверьте работу API:**
   ```bash
   # Получить список новостей
   curl http://localhost:8080/list
   
   # Обновить новость
   curl -X POST http://localhost:8080/edit/1 \
     -H "Content-Type: application/json" \
     -d '{
       "Id": 1,
       "Title": "Test Title",
       "Content": "Test content",
       "Categories": [1, 2, 3]
     }'
   ```

## Вариант 2: Локальный запуск

1. **Установите PostgreSQL:**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   
   # macOS
   brew install postgresql
   ```

2. **Создайте базу данных:**
   ```bash
   createdb news_db
   psql -d news_db -f database/schema.sql
   ```

3. **Настройте переменные окружения:**
   ```bash
   cp env.example .env
   # Отредактируйте .env файл под ваши настройки
   ```

4. **Запустите приложение:**
   ```bash
   go run cmd/go_news_server/main.go
   ```

## Тестирование

### Автоматическое тестирование
```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Ручное тестирование API
```bash
# Используйте готовый скрипт
./scripts/test_api.sh

# Или тестируйте вручную
curl http://localhost:8080/list
```

## Полезные команды

```bash
# Показать справку по командам
make help

# Установка зависимостей
make deps

# Сборка приложения
make build

# Запуск в режиме разработки
make dev

# Очистка
make clean
```

## Структура API

### GET /list
Получение списка новостей с пагинацией.

**Query параметры:**
- `limit` (опционально) - количество записей (по умолчанию 10)
- `offset` (опционально) - смещение (по умолчанию 0)

**Пример:**
```bash
curl "http://localhost:8080/list?limit=5&offset=0"
```

### POST /edit/:Id
Обновление новости по ID.

**Пример:**
```bash
curl -X POST http://localhost:8080/edit/1 \
  -H "Content-Type: application/json" \
  -d '{
    "Id": 1,
    "Title": "Updated Title",
    "Content": "Updated content",
    "Categories": [1, 2, 3]
  }'
```

## Логи

Логи приложения сохраняются в директории `logs/`:
- `application.log` - информационные логи
- `application_error.log` - логи ошибок

## Мониторинг

Приложение запускается на порту 8080 по умолчанию.

Для проверки статуса:
```bash
curl http://localhost:8080/list
```

## Устранение неполадок

1. **Ошибка подключения к базе данных:**
   - Проверьте настройки в `.env` файле
   - Убедитесь, что PostgreSQL запущен
   - Проверьте, что база данных создана

2. **Ошибка порта:**
   - Измените `SERVER_PORT` в `.env` файле
   - Убедитесь, что порт не занят другим процессом

3. **Ошибки в логах:**
   - Проверьте файлы в директории `logs/`
   - Убедитесь, что у приложения есть права на запись 