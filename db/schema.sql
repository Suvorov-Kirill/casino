CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     username TEXT UNIQUE,
                                     password TEXT,
                                     coins INTEGER
);

CREATE TABLE IF NOT EXISTS bets (
                                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                                    user_id INTEGER,
                                    amount INTEGER,
                                    game TEXT,
                                    result TEXT,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);