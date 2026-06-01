CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    coins INTEGER NOT NULL DEFAULT 100,
    user_role TEXT NOT NULL DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS bets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    game TEXT NOT NULL,
    result BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Таблица для сохранения состояния игры Блэкджек
CREATE TABLE IF NOT EXISTS blackjack_games (
                                               user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    bet_amount INT NOT NULL,
    player_hand TEXT NOT NULL, -- храним карты как строку, например: "H10,S5"
    dealer_hand TEXT NOT NULL,
    deck TEXT NOT NULL,        -- оставшаяся колода
    status VARCHAR(20) NOT NULL DEFAULT 'playing' -- 'playing', 'won', 'lost', 'push'
    );


UPDATE users SET user_role = 'admin' WHERE id = 1;

