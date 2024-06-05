CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    file VARCHAR,
    deadline TIMESTAMP,
    year VARCHAR(10),
    teacher_id VARCHAR,
    module_id VARCHAR
);



CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file VARCHAR(255),
    grade FLOAT DEFAULT 0,
    feedback TEXT DEFAULT '',
    assignment_id UUID,
    student_id VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    EvaluatedAt TIMESTAMP,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);

COMMIT;