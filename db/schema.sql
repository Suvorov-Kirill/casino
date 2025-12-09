CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     username TEXT UNIQUE,
                                     password TEXT,
                                     coins INTEGER,
                                     user_role TEXT NOT NULL DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS bets (
                                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                                    user_id INTEGER,
                                    amount INTEGER,
                                    game TEXT,
                                    result BOOLEAN,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

UPDATE users SET user_role = 'admin' WHERE id = 1;
