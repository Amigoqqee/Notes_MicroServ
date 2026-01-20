#!/bin/bash

BASE_URL="http://localhost:8103/notes"
AUTH_BASE_URL="http://localhost:8101/auth"

echo "🚀 Полное тестирование Notes API"
echo "==============================="


echo ""
echo "🔍 Вход в систему (получение JWT токена)"
echo "Запрос: POST $AUTH_BASE_URL/login"
echo "Ответ:"
LOGIN_RESPONSE=$(curl -X "POST" "$AUTH_BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser","password":"password123"}' \
     -s)

echo "$LOGIN_RESPONSE"

LOGIN_STATUS=$(curl -X "POST" "$AUTH_BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser","password":"password123"}' \
     -w "%{http_code}" \
     -s -o /dev/null)

echo "📊 HTTP Статус: $LOGIN_STATUS"

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
TOKEN=$(echo "$TOKEN" | tr -d '\n\r ' | xargs)
echo "Извлеченный токен: $TOKEN"
echo "-------------------------------------------"

sleep 3

echo ""
echo "🔍 Создание новой заметки"
echo "Запрос: POST $BASE_URL/note"
echo "Ответ:"
CREATE_RESPONSE=$(curl -X "POST" "$BASE_URL/note" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Note","content":"Test Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n")
     
ID_NOTE=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "ID созданной заметки: $ID_NOTE"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Получение списка всех заметок"
echo "Запрос: GET $BASE_URL/notes"
echo "Ответ:"
curl -X "GET" "$BASE_URL/notes" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Получение заметки по ID"
echo "Запрос: GET $BASE_URL/note/$ID_NOTE"
echo "Ответ:"
curl -X "GET" "$BASE_URL/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Редактирование заметки по ID"
echo "Запрос: PUT $BASE_URL/note/$ID_NOTE"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Updated Note","content":"Updated Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Удаление заметки по ID"
echo "Запрос: DELETE $BASE_URL/note/$ID_NOTE"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

echo "✅ Все тесты завершены!"
echo "==============================="