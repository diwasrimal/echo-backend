CREATE TABLE IF NOT EXISTS users (
	id bigserial NOT NULL PRIMARY KEY,
	fullname text NOT NULL,
	username text NOT NULL UNIQUE,
	password_hash text NOT NULL,
	bio text DEFAULT ''
);

CREATE TABLE IF NOT EXISTS messages (
	id bigserial NOT NULL PRIMARY KEY,
	sender_id bigserial NOT NULL REFERENCES users(id),
	receiver_id bigserial NOT NULL REFERENCES users(id),
	text text NOT NULL,
	timestamp timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS user_sessions (
	user_id bigserial NOT NULL PRIMARY KEY REFERENCES users(id),
	session_id text NOT NULL
);

CREATE TABLE IF NOT EXISTS conversations (
	user1_id bigserial NOT NULL REFERENCES users(id),
	user2_id bigserial NOT NULL REFERENCES users(id),
	timestamp timestamptz NOT NULL
);
CREATE UNIQUE INDEX unique_conversation_pair ON conversations(LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));

CREATE TABLE IF NOT EXISTS friends (
	user1_id bigserial NOT NULL REFERENCES users(id),
	user2_id bigserial NOT NULL REFERENCES users(id)
);
CREATE UNIQUE INDEX unique_friend_pair ON friends(LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));

CREATE TABLE IF NOT EXISTS friend_requests (
    requestor_id bigserial NOT NULL REFERENCES users(id),
    receiver_id bigserial NOT NULL REFERENCES users(id),
    UNIQUE(requestor_id, receiver_id)
);

-- To search users
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
