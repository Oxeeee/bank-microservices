# Services
- Billing
- exchange
- notification
- transaction
- auth

# structs

1. Пользователь

```go
type User struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    Balance   float64   `json:"balance"`
    Currency  string    `json:"currency"`
    CreatedAt time.Time `json:"created_at"`
}
```

2. Транзакция
```go
type Transaction struct {
    ID            uuid.UUID `json:"id"`
    SenderID      uuid.UUID `json:"sender_id"`
    ReceiverID    uuid.UUID `json:"receiver_id,omitempty"`
    Amount        float64   `json:"amount"`
    Currency      string    `json:"currency"`
    ExchangeRate  float64   `json:"exchange_rate,omitempty"`
    Status        string    `json:"status"`
    CreatedAt     time.Time `json:"created_at"`
}
```

3. Валютные курсы
```go
type ExchangeRate struct {
    BaseCurrency  string    `json:"base_currency"`
    TargetCurrency string  	`json:"target_currency"`
    Rate          float64   `json:"rate"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

4. Платежи
```go
type BillPayment struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Provider  string          `json:"provider"`
	Amount    float64         `json:"amount"`
	Status    string          `json:"status"`
	Details   json.RawMessage `json:"details,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
```

# Models

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance NUMERIC(18,2) DEFAULT 0.00,
    currency TEXT NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
```

```sql
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('email', 'phone'))
);
```

---

```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES users(id),
    receiver_id UUID REFERENCES users(id),
    amount NUMERIC(18,2) NOT NULL,
    currency TEXT NOT NULL,
    exchange_rate NUMERIC(18,6),
    status TEXT NOT NULL CHECK (status IN ('pending', 'success', 'failed', 'canceled')),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

```sql
CREATE TABLE transaction_limits (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    daily_limit NUMERIC(18,2) NOT NULL,
    monthly_limit NUMERIC(18,2) NOT NULL
);
```

---

```sql
CREATE TABLE bill_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID REFERENCES users (id),
    provider TEXT NOT NULL,
    amount NUMERIC(18, 2) NOT NULL,
    currency TEXT NOT NULL,
    details JSONB,
    status TEXT NOT NULL CHECK (
        status IN (
            'pending',
            'success',
            'failed',
            'canceled',
            'reversed'
        )
    ),
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now ()
);

```























# description
Pet-проект: FinGo – Финансовый сервис на Go

Цель проекта: Создать высоконагруженный микросервисный REST API для финансовых операций с транзакциями, аутентификацией и интеграцией валютных переводов.

Архитектура проекта

Проект будет состоять из нескольких микросервисов:
	1.	Auth Service (Аутентификация)
	2.	Transaction Service (Транзакции)
	3.	Currency Exchange Service (Конвертация валют)
	4.	Billing Service (Оплата счетов)
	5.	Notification Service (Уведомления)
	6.	Audit & Logging Service (Аудит и логирование)

Общение между микросервисами будет через gRPC + REST API, а Kafka обеспечит асинхронную обработку.

Микросервисы

1. Auth Service (Авторизация и регистрация)
	•	Регистрация пользователей
	•	Аутентификация (JWT, OAuth2)
	•	Управление сессиями
	•	Подтверждение по почте (Resend)
	•	Хранение токенов обновления

Кеширование в Redis:
	•	Сессии пользователей (JWT)
	•	Коды подтверждения

2. Transaction Service (Переводы и пополнения)
	•	Пополнение баланса
	•	Переводы внутри страны (одна валюта)
	•	SWIFT-переводы (разные валюты)
	•	Контроль дубликатов транзакций
	•	Фрод-мониторинг

Кеширование в Redis:
	•	Курс валют (обновление раз в 10 мин)
	•	Лимиты пользователя (дневные, месячные)

Отправка в Kafka:
	•	Уведомления о переводах
	•	Аудит транзакций
	•	Фрод-события

3. Currency Exchange Service (Конвертация валют)
	•	Интеграция с провайдером курсов валют (например, openexchangerates.org)
	•	Поддержка маржинального курса
	•	Автообновление курсов

Кеширование в Redis:
	•	Курсы валют

4. Billing Service (Оплата счетов)
	•	Список доступных провайдеров
	•	Оплата коммуналки, интернета, мобильной связи
	•	Генерация платежных квитанций
	•	Статусы платежей

Отправка в Kafka:
	•	Уведомления пользователям о статусах платежей
	•	Аудит платежей

5. Notification Service (Уведомления)
	•	Email, SMS, WebSockets
	•	Подписка на события
	•	Очередь уведомлений

Использование Kafka:
	•	Очередь сообщений о переводах и платежах

6. Audit & Logging Service
	•	Логирование действий пользователей
	•	Хранение данных для отчетов и разборов инцидентов

Использование Kafka:
	•	Асинхронная обработка логов

Технологический стек
	•	Golang (Fiber/Gin + gRPC)
	•	PostgreSQL (Основная БД)
	•	Redis (Кеширование)
	•	Kafka (Очередь событий)
	•	Docker + Kubernetes (Продовая оркестрация)
	•	Prometheus + Grafana (Мониторинг)
	•	Jaeger (Трассировка запросов)
	•	gRPC (Межсервисное взаимодействие)

Структуры данных

1. Пользователь

```go
type User struct {
    ID        uuid.UUID json:"id"
    Email     string    json:"email"
    Password  string    json:"-"
    Balance   float64   json:"balance"
    Currency  string    json:"currency"
    CreatedAt time.Time json:"created_at"
}
```

2. Транзакция
```go
type Transaction struct {
    ID            uuid.UUID json:"id"
    SenderID      uuid.UUID json:"sender_id"
    ReceiverID    uuid.UUID json:"receiver_id,omitempty"
    Amount        float64   json:"amount"
    Currency      string    json:"currency"
    ExchangeRate  float64   json:"exchange_rate,omitempty"
    Status        string    json:"status"
    CreatedAt     time.Time json:"created_at"
}
```

3. Валютные курсы
```go
type ExchangeRate struct {
    BaseCurrency  string    json:"base_currency"
    TargetCurrency string   json:"target_currency"
    Rate          float64   json:"rate"
    UpdatedAt     time.Time json:"updated_at"
}
```

4. Платежи
```go
type BillPayment struct {
    ID         uuid.UUID json:"id"
    UserID     uuid.UUID json:"user_id"
    Provider   string    json:"provider"
    Amount     float64   json:"amount"
    Status     string    json:"status"
    CreatedAt  time.Time json:"created_at"
}
```

Заключение

Проект FinGo даст отличный опыт работы с:
	•	Аутентификацией и безопасностью (JWT, OAuth2)
	•	Финансовыми транзакциями (idempotency, дедупликация, фрод)
	•	Интеграцией с внешними сервисами (валютные API, SWIFT)
	•	Асинхронными процессами (Kafka, кеш Redis)
	•	Мониторингом и логированием (Prometheus, Grafana, Jaeger)

Готовый код можно будет оборачивать в OpenAPI и документировать для будущего продакшн-развертывания.
