curl -X POST "http://localhost:8101/auth/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "password123"
         }' \
     -w "\nStatus: %{http_code}\n"

curl -X POST "http://localhost:8101/auth/login" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "password123"
         }' \
     -w "\nStatus: %{http_code}\n"

curl -X GET "http://localhost:8101/auth/user" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"

curl -X PUT "http://localhost:8101/auth/user" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "updated_username",
           "email": "updated@example.com"
         }' \
     -w "\nStatus: %{http_code}\n"

curl -X DELETE "http://localhost:8101/auth/user" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"
