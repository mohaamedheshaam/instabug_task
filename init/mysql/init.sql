CREATE DATABASE IF NOT EXISTS chat

USE chats;

CREATE TABLE IF NOT EXISTS applications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    chats_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE KEY unique_token (token)
);

CREATE TABLE IF NOT EXISTS chats (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id VARCHAR(255) NOT NULL,
    number INT NOT NULL,
    messages_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    UNIQUE KEY unique_app_number (application_id, number),
    FOREIGN KEY (application_id) REFERENCES applications(token)
);

CREATE TABLE IF NOT EXISTS messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chat_id BIGINT UNSIGNED NOT NULL,
    number INT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    UNIQUE KEY unique_chat_number (chat_id, number)
);