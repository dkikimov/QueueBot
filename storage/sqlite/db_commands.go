package sqlite

const CreateTables string = `
CREATE TABLE IF NOT EXISTS queues(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	query_id VARCHAR NOT NULL,
	is_ended INTEGER DEFAULT 0,
	description TEXT DEFAULT NULL
	);
CREATE UNIQUE INDEX IF NOT EXISTS idx_queues_message_id ON queues(query_id);

CREATE TABLE IF NOT EXISTS participants(
    query_id VARCHAR NOT NULL,
    user_id BIGINT NOT NULL,
    user_name VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_prt_message_id ON participants(query_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_prt_user_id ON participants(user_id);
`

const CreateQueue string = `
INSERT INTO queues (query_id, description) VALUES (?, ?);
`
