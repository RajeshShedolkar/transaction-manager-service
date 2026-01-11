


curl -X POST http://localhost:8080/api/v1/transactions \
-H "Content-Type: application/json" \
-d '{
  "paymentType":"IMMEDIATE",
  "paymentMode":"IMPS",
  "amount":2000,
  "currency":"INR"
}'
