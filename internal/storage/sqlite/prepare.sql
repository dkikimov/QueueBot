-- Initialization
CREATE TABLE IF NOT EXISTS queues(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message_id INTEGER NOT NULL,
	is_ended INTEGER DEFAULT 0,
	description TEXT DEFAULT NULL
	);
CREATE UNIQUE INDEX idx_queues_message_id ON queues(message_id);

CREATE TABLE IF NOT EXISTS participants(
    message_id INTEGER NOT NULL,
    user_id BIGINT NOT NULL,
    user_name VARCHAR NOT NULL
);

CREATE INDEX idx_prt_message_id ON participants(message_id);
CREATE UNIQUE INDEX idx_prt_user_id ON participants(user_id);

-- Create queue
INSERT INTO queues (message_id) VALUES (1);

-- Add participants
INSERT OR IGNORE INTO participants VALUES (1, 1, 'Kikimov Daniil');
INSERT INTO participants VALUES (1, 2, 'Pechkin Vanya');

-- View all participants in queue
-- EXPLAIN QUERY PLAN
SELECT user_name FROM participants WHERE message_id = 1;

SELECT user_id, user_name FROM participants WHERE message_id = 1;

-- Delete participant from queue
-- EXPLAIN QUERY PLAN
DELETE FROM participants WHERE message_id = 1 AND user_id = 1;

-- Count
SELECT COUNT(*) FROM participants WHERE user_id = 1;