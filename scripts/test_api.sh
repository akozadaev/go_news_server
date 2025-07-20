#!/bin/bash

# Скрипт для тестирования API эндпоинтов
# Убедитесь, что сервер запущен на localhost:8080

BASE_URL="http://localhost:8080"

echo "🧪 Тестирование API эндпоинтов"
echo "================================"

# Тест 1: Получение списка новостей
echo "1. Тестирование GET /list"
echo "------------------------"
response=$(curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/list")
echo "$response"
echo ""

# Тест 2: Обновление новости
echo "2. Тестирование POST /edit/1"
echo "----------------------------"
payload='{
  "Id": 1,
  "Title": "Updated Test Title",
  "Content": "Updated test content with <b>HTML</b>",
  "Categories": [1, 2, 3]
}'

response=$(curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X POST \
  -H "Content-Type: application/json" \
  -d "$payload" \
  "$BASE_URL/edit/1")
echo "$response"
echo ""

# Тест 3: Получение списка новостей после обновления
echo "3. Получение списка новостей после обновления"
echo "---------------------------------------------"
response=$(curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/list")
echo "$response"
echo ""

# Тест 4: Тест с пагинацией
echo "4. Тестирование пагинации GET /list?limit=5&offset=0"
echo "------------------------------------------------"
response=$(curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/list?limit=5&offset=0")
echo "$response"
echo ""

# Тест 5: Тест с неверным ID
echo "5. Тестирование с неверным ID POST /edit/invalid"
echo "-----------------------------------------------"
response=$(curl -s -w "\nHTTP Status: %{http_code}\n" \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"Id": "invalid"}' \
  "$BASE_URL/edit/invalid")
echo "$response"
echo ""

echo "✅ Тестирование завершено!" 