CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    files VARCHAR[],
    deadline TIMESTAMP,
    year VARCHAR,
    groups string[],
    teacher_id UUID,
    module_id UUID
);


CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file VARCHAR(255),
    grade FLOAT,
    feedback TEXT,
    assignment_id UUID,
    student_id UUID,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);


COMMIT;