curl -X POST "http://localhost/notes/note" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgzMzU5NTYsImlhdCI6MTc2ODI0OTU1NiwiaWQiOjMsInR5cGUiOiJhY2Nlc3NUb2tlbiJ9.9F0ggVQzC7hyzvMRwsbWXDEEUtvVN1Ufj1LmiaTiDg0" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Note2","content":"Test Content2"}' \
     -w "\nStatus: %{http_code}\n"

curl -X GET "http://localhost/notes/notes" \
     -H "Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgzMzU5NTYsImlhdCI6MTc2ODI0OTU1NiwiaWQiOjMsInR5cGUiOiJhY2Nlc3NUb2tlbiJ9.9F0ggVQzC7hyzvMRwsbWXDEEUtvVN1Ufj1LmiaTiDg0" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n" 

curl -X GET "http://localhost/notes/note/686bcae11f8babdb67eb3356" \
     -H "Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgzMzU5NTYsImlhdCI6MTc2ODI0OTU1NiwiaWQiOjMsInR5cGUiOiJhY2Nlc3NUb2tlbiJ9.9F0ggVQzC7hyzvMRwsbWXDEEUtvVN1Ufj1LmiaTiDg0" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n"

curl -X PUT "http://localhost/notes/note/686e1ce598e92c19e5ec3fb4" \
     -H "Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgzMzU5NTYsImlhdCI6MTc2ODI0OTU1NiwiaWQiOjMsInR5cGUiOiJhY2Nlc3NUb2tlbiJ9.9F0ggVQzC7hyzvMRwsbWXDEEUtvVN1Ufj1LmiaTiDg0" \
     -H "Content-Type: application/json" \
     -d '{"name":"Обновленное имя","content":"Обновленный контент"}' \
     -w "\nStatus: %{http_code}\n"
     
curl -X DELETE "http://localhost/notes/note/686bcae11f8babdb67eb3356" \
     -H "Authorization: Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjgzMzU5NTYsImlhdCI6MTc2ODI0OTU1NiwiaWQiOjMsInR5cGUiOiJhY2Nlc3NUb2tlbiJ9.9F0ggVQzC7hyzvMRwsbWXDEEUtvVN1Ufj1LmiaTiDg0" \
     -H "Content-Type: application/json" \
     -w "\nStatus: %{http_code}\n" 