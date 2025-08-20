-- ====================================
-- Drop Triggers
-- ====================================
DROP TRIGGER IF EXISTS set_updated_at_users ON users;
DROP TRIGGER IF EXISTS set_updated_at_orders ON orders;
DROP TRIGGER IF EXISTS set_updated_at_payments ON sarees;
DROP TRIGGER IF EXISTS set_updated_at_genres ON payments;

-- ====================================
-- Drop Trigger Function
-- ====================================
DROP FUNCTION IF EXISTS trigger_set_updated_at;

-- ====================================
-- ===============================
-- DROP TABLES (in dependency order)
-- ===============================

-- payments → depends on orders
DROP TABLE IF EXISTS payments;

-- orders → depends on sarees,users
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS sarees;

DROP TABLE IF EXISTS users;

