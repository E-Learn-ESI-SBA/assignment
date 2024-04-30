CREATE TABLE IF NOT EXISTS assignments (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    deadline TIMESTAMP,
    promo INT,
    groups INT[],
    teacher_id INT,
    module_id INT
);


CREATE TABLE IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    file VARCHAR(255),
    grade FLOAT,
    feedback TEXT,
    assignment_id INT,
    student_id INT,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);


COMMIT;