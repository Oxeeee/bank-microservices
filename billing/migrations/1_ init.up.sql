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
