CREATE TABLE IF NOT EXISTS messages
(
    id      SERIAL,
    sender  VARCHAR(20)  NOT NULL,
    text    VARCHAR(120) NOT NULL,
    send_at TIMESTAMP
)