CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    files VARCHAR[],
    deadline TIMESTAMP,
    year VARCHAR,
    groups VARCHAR[],
    teacher_id VARCHAR,
    module_id VARCHAR
);


CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file VARCHAR(255),
    grade FLOAT,
    feedback TEXT,
    assignment_id UUID,
    student_id VARCHAR,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);

COMMIT;