# Changelog

## [1.0.0] - 2024-01-XX

### Добавлено
- ✅ Полная реализация REST API сервера согласно требованиям
- ✅ Два основных эндпоинта:
  - `POST /edit/:Id` - изменение новости по Id
  - `GET /list` - список новостей с пагинацией
- ✅ Поддержка PostgreSQL с использованием Reform ORM
- ✅ Connection pool для оптимизации производительности
- ✅ Конфигурация через переменные окружения с использованием Viper
- ✅ Структурированное логирование с Zap
- ✅ Dependency Injection с Uber FX
- ✅ Graceful shutdown для корректного завершения работы
- ✅ Обработка ошибок и валидация входных данных
- ✅ Пагинация для списка новостей
- ✅ Поддержка категорий новостей через связную таблицу

### Исправлено
- 🔧 Исправлены SQL запросы для PostgreSQL (правильные кавычки и синтаксис)
- 🔧 Исправлена обработка категорий в репозитории
- 🔧 Исправлено формирование DataSourceName в конфигурации
- 🔧 Исправлены форматы JSON ответов согласно требованиям
- 🔧 Добавлена поддержка частичного обновления полей (если поле не задано - не обновляется)

### Улучшено
- 📈 Добавлены настройки connection pool (MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
- 📈 Улучшена структура проекта с четким разделением слоев
- 📈 Добавлены интерфейсы для лучшей тестируемости
- 📈 Добавлены unit тесты для хендлеров
- 📈 Добавлена поддержка Docker и Docker Compose
- 📈 Создан Makefile с полезными командами
- 📈 Добавлена документация и инструкции по запуску

### Технические детали
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL с Reform ORM
- **Configuration**: Viper для управления конфигурацией
- **Logging**: Zap для структурированного логирования
- **Testing**: Testify для unit тестов
- **Containerization**: Docker с многоэтапной сборкой
- **Dependency Injection**: Uber FX

### Структура базы данных
```sql
-- Таблица новостей
CREATE TABLE "News" (
  "Id" BIGSERIAL PRIMARY KEY,
  "Title" VARCHAR(255) NOT NULL,
  "Content" TEXT NOT NULL
);

-- Таблица связи новостей и категорий
CREATE TABLE "NewsCategories" (
  "NewsId" BIGINT NOT NULL,
  "CategoryId" BIGINT NOT NULL,
  PRIMARY KEY ("NewsId", "CategoryId")
);
```

### API Endpoints

#### GET /list
Получение списка новостей с пагинацией.
- Query параметры: `limit`, `offset`
- Формат ответа: `{"Success": true, "News": [...]}`

#### POST /edit/:Id
Обновление новости по ID.
- Поддерживает частичное обновление (только заданные поля)
- Формат запроса: `{"Id": 1, "Title": "...", "Content": "...", "Categories": [...]}`
- Формат ответа: `{"Success": true}`

### Файлы конфигурации
- `env.example` - пример переменных окружения
- `database/schema.sql` - схема базы данных
- `docker-compose.yml` - конфигурация для Docker Compose
- `Dockerfile` - многоэтапная сборка Docker образа

### Команды для разработки
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

### Тестирование
- Unit тесты для хендлеров
- Скрипт для ручного тестирования API (`scripts/test_api.sh`)
- Поддержка тестирования с покрытием кода

### Документация
- `README.md` - полная документация проекта
- `QUICKSTART.md` - инструкции по быстрому запуску
- `CHANGELOG.md` - история изменений 