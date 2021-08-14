CREATE DATABASE IF NOT EXISTS sampleDB;

USE sampleDB;

CREATE TABLE IF NOT EXISTS sampleDB.people (
    name        VARCHAR(100),
    title       VARCHAR(10),
    description VARCHAR(100),
    PRIMARY KEY (name)
    );

DELETE FROM sampleDB.people;

INSERT INTO sampleDB.people VALUES ('EQ', 'Tech', 'Where are the cakes?');
INSERT INTO sampleDB.people VALUES ('Farhad', 'Dr.', 'Why ... why are you so nice?');
INSERT INTO sampleDB.people VALUES ('Sonos', 'Mr.', 'you are so loud!');
INSERT INTO sampleDB.people VALUES ('Margo', 'Ms.', 'Privet!');
INSERT INTO sampleDB.people VALUES ('Trace', 'Mr.', 'This is so cool!');