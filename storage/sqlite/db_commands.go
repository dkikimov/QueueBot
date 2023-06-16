package sqlite

const CreateTables string = `
CREATE TABLE IF NOT EXISTS queues(
	message_id VARCHAR PRIMARY KEY,
	description TEXT DEFAULT NULL,
	is_ended INTEGER DEFAULT 0,
	is_started INTEGER DEFAULT 0,
	current_person INTEGER DEFAULT 0
	);

CREATE TABLE IF NOT EXISTS participants(
    message_id VARCHAR NOT NULL,
    user_id BIGINT NOT NULL,
    user_name VARCHAR NOT NULL,
    order_number INTEGER DEFAULT(random())
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

const AddUserToQueue = `INSERT INTO participants(message_id, user_id, user_name) VALUES (?, ?, ?)`
const RemoveUserFromQueue = `DELETE FROM participants WHERE user_id = ? AND message_id = ?`

const GetUserCurrentStep = `SELECT current_step FROM users WHERE id = ?`

const GetUsersInQueue = `SELECT user_id, user_name FROM participants WHERE message_id = ?`
const GetDescriptionOfQueue = `SELECT description FROM queues WHERE message_id = ?`

const CountMatchesInParticipants = `SELECT COUNT(*) FROM participants WHERE user_id = ? AND message_id = ?;`
const StartQueue = `UPDATE queues SET is_started = 1 WHERE message_id = ? AND is_started = 0 RETURNING is_started`

const IncrementCurrentPerson = `UPDATE queues SET current_person = current_person + 1 WHERE message_id = ? RETURNING current_person`
const GoToMenu = `UPDATE queues SET current_person = 0, is_started = 0 WHERE message_id = ?; `
