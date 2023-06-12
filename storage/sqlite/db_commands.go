package sqlite

const CreateTables string = `
CREATE TABLE IF NOT EXISTS queues(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message_id INTEGER NOT NULL,
	is_ended INTEGER DEFAULT 0,
	description TEXT DEFAULT NULL
	);
CREATE UNIQUE INDEX IF NOT EXISTS idx_queues_message_id ON queues(message_id);

CREATE TABLE IF NOT EXISTS participants(
    message_id INTEGER NOT NULL,
    user_id BIGINT NOT NULL,
    user_name VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_prt_message_id ON participants(message_id);
CREATE INDEX IF NOT EXISTS idx_prt_user_id ON participants(user_id);

CREATE TABLE IF NOT EXISTS users(
    id BIGINT PRIMARY KEY,
    current_step INTEGER DEFAULT 0
);
`

const CreateQueue = `INSERT INTO queues (message_id, description) VALUES (?, ?);`
const CreateUser = `INSERT OR IGNORE INTO users VALUES (?, 0);`
const SetUserCurrentStep = `UPDATE users SET current_step = ? WHERE id = ?`

const AddUserToQueue = `INSERT INTO participants VALUES (?, ?, ?)`
const RemoveUserFromQueue = `DELETE FROM participants WHERE user_id = ? AND message_id = ?`

const GetUserCurrentStep = `SELECT current_step FROM users WHERE id = ?`
const GetUsersInQueue = `SELECT user_id, user_name FROM participants WHERE message_id = ?`
const GetDescriptionOfQueue = `SELECT description FROM queues WHERE message_id = ?`

const CountMatchesInParticipants = `SELECT COUNT(*) FROM participants WHERE user_id = ? AND message_id = ?;`
