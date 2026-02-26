SERVER_PORT	8080
SMTP_HOST	smtp.mail.ru
SMTP_PORT	587
SMTP_USERNAME	supdev@list.ru

SMTP_PASSWORD	pass
FROM_EMAIL	supdev@list.ru

KAFKA_BROKERS	kafka:29092
KAFKA_TOPIC	email-auth-codes
KAFKA_GROUP_ID	email-service

P.S. в godotenv.Load() можно указать путь к .env, но это нужно только для локальной разработки 
