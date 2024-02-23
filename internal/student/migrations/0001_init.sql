-- +goose Up 
CREATE TABLE IF NOT EXISTS student_events (
	id SERIAL PRIMARY KEY,
	type VARCHAR(255) NOT NULL,
	data BYTEA NOT NULL,
	version INT NOT NULL,
	timestamp INT NOT NULL,
	aggregate_id VARCHAR(255) NOT NULL,
	UNIQUE (aggregate_id, version)
);

CREATE TABLE IF NOT EXISTS student_projections (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL, 
    last_name TEXT NOT NULL,
    school_id TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    date_of_enrollment DATE,
    version INT NOT NULL,
    active BOOLEAN NOT NULL, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down 
DROP TABLE IF EXISTS student_events;
DROP TABLE IF EXISTS student_projections;
