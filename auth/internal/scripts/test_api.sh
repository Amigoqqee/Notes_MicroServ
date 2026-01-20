#!/bin/bash

BASE_URL="http://localhost:8101/auth"

echo "🚀 Полное тестирование Auth API"
echo "==============================="


echo ""
echo "🔍 Регистрация нового пользователя"
echo "Запрос: POST $BASE_URL/register"
curl -X POST "$BASE_URL/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "password123"
         }' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"

sleep 5


echo ""
echo "🔍 Вход в систему (получение JWT токена)"
echo "Запрос: POST $BASE_URL/login"
echo "Ответ:"
LOGIN_RESPONSE=$(curl -X "POST" "$BASE_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser","password":"password123"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s)

echo "$LOGIN_RESPONSE"

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

echo "Извлеченный токен: $TOKEN"
echo "-------------------------------------------"


sleep 5


echo ""
echo "🔍 Получение информации о пользователе"
echo "Запрос: GET $BASE_URL/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"


sleep 5


echo ""
echo "🔍 Обновление информации о пользователе"
echo "Запрос: PUT $BASE_URL/user"
echo "Ответ:"
curl -X "PUT" "$BASE_URL/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"username":"updated_testuser","email":"updated@example.com"}' \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"


sleep 5


echo ""
echo "🔍 Проверка обновленной информации о пользователе"
echo "Запрос: GET $BASE_URL/user"
echo "Ответ:"
curl -X "GET" "$BASE_URL/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"


sleep 5


echo ""
echo "🔍 Удаление пользователя"
echo "Запрос: DELETE $BASE_URL/user"
echo "Ответ:"
curl -X "DELETE" "$BASE_URL/user" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -w "\n📊 HTTP Статус: %{http_code}\n" \
     -s
echo "-------------------------------------------"


echo ""
echo "✅ Тестирование всех endpoints завершено!"