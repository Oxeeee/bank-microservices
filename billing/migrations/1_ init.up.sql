CREATE TABLE bill_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    provider TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    currency TEXT NOT NULL,
    payment_method TEXT NOT NULL CHECK (payment_method IN ('balance', 'credit_card', 'crypto')),
    status TEXT NOT NULL CHECK (status IN ('pending', 'success', 'failed', 'canceled', 'reversed')),
    details JSONB,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('mobile_operator', 'internet_provider', 'utility')),
    active BOOLEAN DEFAULT true
);
