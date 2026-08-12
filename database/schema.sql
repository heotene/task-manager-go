-- ============================================
-- MYJOB DATABASE SCHEMA
-- ============================================

-- USERS
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    fullname TEXT NOT NULL,

    email TEXT,

    phone TEXT,

    password TEXT NOT NULL,

    -- Verification
    verification_method TEXT NOT NULL DEFAULT 'email',

    email_verified INTEGER DEFAULT 0,

    phone_verified INTEGER DEFAULT 0,

    verification_code TEXT,

    verification_expires DATETIME,

    -- Password reset
    reset_token TEXT,

    reset_expires DATETIME,

    -- Account timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TASKS
-- ============================================

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    title TEXT NOT NULL,

    description TEXT,

    priority TEXT DEFAULT 'Medium',

    category TEXT DEFAULT 'General',

    due_date TEXT,

    completed INTEGER DEFAULT 0,

    user_id INTEGER NOT NULL,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);