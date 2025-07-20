# Go News Server

JSON REST сервер для управления новостями с двумя основными эндпоинтами:
- `POST /edit/:Id` - изменение новости по Id
- `GET /list` - список новостей

## Технологии

- **Web Framework**: Fiber
- **Database**: PostgreSQL
- **ORM**: Reform
- **Configuration**: Viper
- **Logging**: Zap
- **Dependency Injection**: Uber FX

## Структура базы данных

```sql
-- Таблица новостей
CREATE TABLE "News" (
  "Id" BIGSERIAL PRIMARY KEY,
  "Title" VARCHAR(255) NOT NULL,
  "Content" TEXT NOT NULL
);

-- Таблица категорий
CREATE TABLE "Categories" (
  "Id" BIGSERIAL PRIMARY KEY,
  "Name" VARCHAR(100) NOT NULL UNIQUE,
  "Description" TEXT,
  "CreatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "UpdatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица связи новостей и категорий
CREATE TABLE "NewsCategories" (
  "NewsId" BIGINT NOT NULL,
  "CategoryId" BIGINT NOT NULL,
  PRIMARY KEY ("NewsId", "CategoryId"),
  FOREIGN KEY ("NewsId") REFERENCES "News"("Id") ON DELETE CASCADE,
  FOREIGN KEY ("CategoryId") REFERENCES "Categories"("Id") ON DELETE CASCADE
);
```

## API Endpoints

### Новости

#### GET /list
Получение списка новостей с пагинацией.

**Query параметры:**
- `limit` (опционально) - количество записей (по умолчанию 10)
- `offset` (опционально) - смещение (по умолчанию 0)

**Пример ответа:**
```json
{
    "Success": true,
    "News": [
        {
            "Id": 64,
            "Title": "Lorem ipsum",
            "Content": "Dolor sit amet <b>foo</b>",
            "Categories": [1,2,3]
        }
    ]
}
```

#### POST /edit/:Id
Изменение новости по ID.

**Параметры пути:**
- `Id` - ID новости для редактирования

**Тело запроса:**
```json
{
    "Id": 64,
    "Title": "Lorem ipsum",
    "Content": "Dolor sit amet <b>foo</b>",
    "Categories": [1,2,3]
}
```

**Примечание:** Если какое-то из полей не задано - это поле не будет обновлено.

**Пример ответа:**
```json
{
    "Success": true
}
```

### Категории

#### GET /categories
Получение списка всех категорий с пагинацией.

**Query параметры:**
- `limit` (опционально) - количество записей (по умолчанию 10)
- `offset` (опционально) - смещение (по умолчанию 0)

**Пример ответа:**
```json
{
    "success": true,
    "categories": [
        {
            "id": 1,
            "name": "Technology",
            "description": "Technology related news",
            "created_at": "2025-07-20T09:56:38.619586Z",
            "updated_at": "2025-07-20T09:56:38.619586Z"
        }
    ],
    "total": 5
}
```

#### GET /categories/:id
Получение категории по ID.

**Параметры пути:**
- `id` - ID категории

**Пример ответа:**
```json
{
    "success": true,
    "category": {
        "id": 1,
        "name": "Technology",
        "description": "Technology related news",
        "created_at": "2025-07-20T09:56:38.619586Z",
        "updated_at": "2025-07-20T09:56:38.619586Z"
    }
}
```

#### POST /categories
Создание новой категории.

**Тело запроса:**
```json
{
    "name": "Science",
    "description": "Scientific discoveries and research"
}
```

**Пример ответа:**
```json
{
    "success": true,
    "category": {
        "id": 6,
        "name": "Science",
        "description": "Scientific discoveries and research",
        "created_at": "2025-07-20T13:01:40.913797Z",
        "updated_at": "2025-07-20T13:01:40.913797Z"
    }
}
```

#### PUT /categories/:id
Обновление категории.

**Параметры пути:**
- `id` - ID категории для обновления

**Тело запроса:**
```json
{
    "name": "Technology Updated",
    "description": "Updated technology description"
}
```

**Пример ответа:**
```json
{
    "success": true,
    "category": {
        "id": 1,
        "name": "Technology Updated",
        "description": "Updated technology description",
        "created_at": "2025-07-20T09:56:38.619586Z",
        "updated_at": "2025-07-20T13:01:40.937882Z"
    }
}
```

#### DELETE /categories/:id
Удаление категории.

**Параметры пути:**
- `id` - ID категории для удаления

**Пример ответа:**
```json
{
    "success": true,
    "message": "Category deleted successfully"
}
```

#### GET /news/:id/categories
Получение категорий для конкретной новости.

**Параметры пути:**
- `id` - ID новости

**Пример ответа:**
```json
{
    "success": true,
    "categories": [
        {
            "id": 1,
            "name": "Technology",
            "description": "Technology related news",
            "created_at": "2025-07-20T09:56:38.619586Z",
            "updated_at": "2025-07-20T09:56:38.619586Z"
        }
    ],
    "total": 1
}
```

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd go_news_server
```

### 2. Настройка базы данных
```bash
# Создайте базу данных PostgreSQL
createdb news_db

# Примените схему
psql -d news_db -f database/schema.sql
```

### 3. Настройка переменных окружения
```bash
# Скопируйте пример файла
cp env.example .env

# Отредактируйте .env файл под ваши настройки
```

### 4. Установка зависимостей
```bash
go mod download
```

### 5. Запуск сервера
```bash
go run cmd/go_news_server/main.go
```

## Конфигурация

Все настройки выполняются через переменные окружения:

### База данных
- `DB_HOST` - хост базы данных (по умолчанию localhost)
- `DB_PORT` - порт базы данных (по умолчанию 5432)
- `DB_USER` - пользователь базы данных
- `DB_PASSWORD` - пароль базы данных
- `DB_NAME` - имя базы данных
- `DB_MAX_OPEN_CONNS` - максимальное количество открытых соединений (по умолчанию 50)
- `DB_MAX_IDLE_CONNS` - максимальное количество неактивных соединений (по умолчанию 20)
- `DB_CONN_MAX_LIFETIME` - время жизни соединения в минутах (по умолчанию 5)

### Сервер
- `SERVER_HOST` - хост сервера (по умолчанию localhost)
- `SERVER_PORT` - порт сервера (по умолчанию 8080)
- `SERVER_READ_TIMEOUT` - таймаут чтения в секундах (по умолчанию 15)

### Логирование
- `level` - уровень логирования (по умолчанию -1)
- `encoding` - формат логирования (console/json)
- `info_filename` - файл для информационных логов
- `error_filename` - файл для логов ошибок

## Особенности реализации

1. **Connection Pool**: Настроен пул соединений с базой данных для оптимизации производительности
2. **Graceful Shutdown**: Корректное завершение работы сервера
3. **Dependency Injection**: Использование Uber FX для управления зависимостями
4. **Structured Logging**: Логирование с использованием Zap
5. **Error Handling**: Грамотная обработка ошибок
6. **Configuration Management**: Централизованное управление конфигурацией через Viper
7. **CRUD Operations**: Полный набор операций для управления категориями
8. **Validation**: Валидация входных данных и проверка уникальности

## Структура проекта

```
go_news_server/
├── cmd/go_news_server/     # Точка входа приложения
├── config/                 # Конфигурационные файлы
├── database/               # SQL скрипты
├── internal/               # Внутренний код приложения
│   ├── handlers/          # HTTP обработчики
│   ├── middleware/        # Промежуточное ПО
│   ├── models/            # Модели данных
│   ├── repository/        # Слой доступа к данным
│   ├── routes/            # Маршрутизация
│   ├── server/            # Настройки сервера
│   └── services/          # Бизнес-логика
├── pkg/                   # Переиспользуемые пакеты
│   ├── config/           # Конфигурация
│   └── logging/          # Логирование
└── logs/                 # Логи приложения
```

## Результаты тестирования

Проект полностью протестирован и работает корректно:

### ✅ Функциональность новостей
- **GET /list** - получение списка новостей работает корректно
- **POST /edit/:Id** - обновление новостей работает корректно
- **Частичное обновление** - если поле не задано, оно не обновляется
- **Пагинация** - параметры limit и offset работают корректно
- **Обработка ошибок** - неверные ID возвращают HTTP 400

### ✅ Функциональность категорий
- **GET /categories** - получение списка категорий работает корректно
- **GET /categories/:id** - получение категории по ID работает корректно
- **POST /categories** - создание категории работает корректно
- **PUT /categories/:id** - обновление категории работает корректно
- **DELETE /categories/:id** - удаление категории работает корректно
- **GET /news/:id/categories** - получение категорий для новости работает корректно
- **Валидация** - проверка уникальности имен категорий
- **Пагинация** - работает для списка категорий

### ✅ Технические аспекты
- **Подключение к PostgreSQL** - работает через Docker контейнер
- **Connection Pool** - настроен и работает
- **Логирование** - структурированные логи с Zap
- **Конфигурация** - через переменные окружения
- **Docker** - контейнеризация работает корректно
- **Dependency Injection** - Uber FX настроен корректно
- **Статический анализ кода** - настроен и используется golangci-lint

### 📊 Примеры ответов API

**Получение списка новостей:**
```json
{
  "Success": true,
  "News": [
    {
      "id": 1,
      "title": "Test News",
      "content": "Test content",
      "categories": [1, 2, 3]
    }
  ]
}
```

**Получение списка категорий:**
```json
{
  "success": true,
  "categories": [
    {
      "id": 1,
      "name": "Technology",
      "description": "Technology related news",
      "created_at": "2025-07-20T09:56:38.619586Z",
      "updated_at": "2025-07-20T09:56:38.619586Z"
    }
  ],
  "total": 5
}
```

**Обновление новости:**
```json
{
  "Success": true
}
```

**Создание категории:**
```json
{
  "success": true,
  "category": {
    "id": 6,
    "name": "Science",
    "description": "Scientific discoveries",
    "created_at": "2025-07-20T13:01:40.913797Z",
    "updated_at": "2025-07-20T13:01:40.913797Z"
  }
}
```

**Ошибка при неверном ID:**
```
HTTP 404 Not Found
{"message":"Category not found","success":false}
```

## Команды для разработки

```bash
# Установка и настройка
make install
make setup-db

# Разработка
make dev
make test

# Сборка и запуск
make build
make run

# Docker
make docker-build
make docker-run
``` 
