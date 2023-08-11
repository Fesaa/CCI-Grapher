CREATE TABLE IF NOT EXISTS usernames(
    user_id BIGINT NOT NULL,
    username VARCHAR(32) NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id BIGINT NOT NULL,
    channel_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    roles VARCHAR NOT NULL,
    time TIMESTAMP NOT NULL,
    PRIMARY KEY (message_id)
);
