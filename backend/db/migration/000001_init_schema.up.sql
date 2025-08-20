-- ====================================
-- Users Table
-- ====================================
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR NOT NULL,
  hashed_password VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  full_name VARCHAR NOT NULL,
  role TEXT CHECK(role IN ('admin', 'user')) DEFAULT 'user',
  is_admin BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ====================================
-- Sarees Table
-- ====================================
CREATE TABLE sarees (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    stock INT DEFAULT 0,
    image_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- ====================================
-- Orders Table
-- ====================================
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    saree_id INT REFERENCES sarees(id),
    quantity INT DEFAULT 1,
    total_amount NUMERIC(10,2) NOT NULL,
    status TEXT CHECK(status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')) DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW()
);

-- ====================================
-- Payments (tracks Razorpay transactions)
-- ====================================
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    razorpay_order_id TEXT NOT NULL,   -- Razorpay order (created before checkout)
    razorpay_payment_id TEXT NOT NULL,          -- Filled after success
    razorpay_signature TEXT NOT NULL,           -- For verification
    amount NUMERIC(10,2) NOT NULL,
    currency TEXT DEFAULT 'INR',
    status TEXT CHECK(status IN (
        'created',   -- razorpay order created
        'success',   -- payment verified
        'failed'
    )) DEFAULT 'created',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ====================================
-- Trigger to auto-update `updated_at`
-- ====================================
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to relevant tables
-- Users
CREATE TRIGGER set_updated_at_users
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Apply trigger to relevant tables
-- Users
CREATE TRIGGER set_updated_at_orders
BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Apply trigger to relevant tables
-- sarees
CREATE TRIGGER set_updated_at_sarees
BEFORE UPDATE ON sarees
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();


-- Apply trigger to relevant tables
-- payments
CREATE TRIGGER set_updated_at_payments
BEFORE UPDATE ON payments
FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();