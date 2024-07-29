go 1.19

curl -H "Content-type: application/json" -X POST -d '{"name":"lsc","telephone":"17611129667","password":"123456"}' http://127.0.0.1:8081/api/auth/register

curl -H "Content-type: application/json" -X POST -d '{"telephone":"17611129667","password":"123456"}' http://127.0.0.1:8081/api/auth/login