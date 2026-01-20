#!/bin/bash

BASE_URL="http://localhost"
SERVICE_NAME_AUTH="auth"
SERVICE_NAME_NOTES="notes"

RANDOM_USERNAME="testuser_$(date +%s)_$RANDOM"

echo "🚀 Полное тестирование Auth API"
echo "🎲 Используемый username: $RANDOM_USERNAME"
echo "==============================="

echo ""
echo "🔍 Регистрация нового пользователя"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_AUTH/register"
curl -X POST "$BASE_URL/$SERVICE_NAME_AUTH/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "'$RANDOM_USERNAME'",
           "password": "password123"
         }' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

sleep 1

echo ""
echo "🔍 Вход в систему (получение JWT токена)"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_AUTH/login"
echo "Ответ:"
LOGIN_RESPONSE=$(curl -X "POST" "$BASE_URL/$SERVICE_NAME_AUTH/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "'$RANDOM_USERNAME'","password":"password123"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s)

echo "$LOGIN_RESPONSE"

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

echo "Извлеченный токен: $TOKEN"
echo "-------------------------------------------"

sleep 1

echo ""
echo "🔍 Получение информации о пользователе"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

sleep 1

echo ""
echo "🔍 Обновление информации о пользователе"
echo "Запрос: PUT $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"username":"updated_'$RANDOM_USERNAME'","email":"updated@example.com"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

sleep 1

echo ""
echo "🔍 Проверка обновленной информации о пользователе"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

sleep 1

echo ""
echo "🔍 Удаление пользователя"
echo "Запрос: DELETE $BASE_URL/$SERVICE_NAME_AUTH/user"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/$SERVICE_NAME_AUTH/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s

echo ""
echo "✅ Тестирование сервиса AUTH  завершено!"
echo "-------------------------------------------"

echo ""
echo "🚀 Полное тестирование Notes API"
echo "==============================="


echo ""
echo "🔍 Создание новой заметки"
echo "Запрос: POST $BASE_URL/$SERVICE_NAME_NOTES/note"
echo "Ответ:"
CREATE_RESPONSE=$(curl -X "POST" "$BASE_URL/$SERVICE_NAME_NOTES/note" \
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
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/notes"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/notes" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Получение заметки по ID"
echo "Запрос: GET $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "GET" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Редактирование заметки по ID"
echo "Запрос: PUT $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"name":"Updated Note","content":"Updated Content"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

sleep 2

echo ""
echo "🔍 Удаление заметки по ID"
echo "Запрос: DELETE $BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/$SERVICE_NAME_NOTES/note/$ID_NOTE" \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -w "\n📊 HTTP Статус: %{http_code}\n"
echo "-------------------------------------------"

echo "✅ Все тесты завершены!"
echo "==============================="