CREATE TABLE IF NOT EXISTS assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255),
    description TEXT,
    files VARCHAR[],
    deadline TIMESTAMP,
    promo VARCHAR,
    groups UUID[],
    teacher_id UUID,
    module_id UUID
);


CREATE TABLE IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    file VARCHAR(255),
    grade FLOAT,
    feedback TEXT,
    assignment_id UUID,
    student_id UUID,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);


COMMIT;