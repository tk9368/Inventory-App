-- 1. Employees Table
CREATE TABLE IF NOT EXISTS users (
    emp_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 2. Car Parts Table
CREATE TABLE IF NOT EXISTS parts (
    part_id VARCHAR(50) PRIMARY KEY,
    part_name VARCHAR(100) NOT NULL,
    stock_count INT NOT NULL DEFAULT 0
);

-- 3. Inventory Logs Table (Stock IN/OUT tracking)
CREATE TABLE IF NOT EXISTS part_logs (
    log_id SERIAL PRIMARY KEY,
    emp_id VARCHAR(50) REFERENCES users(emp_id),
    part_id VARCHAR(50) REFERENCES parts(part_id),
    action_type VARCHAR(10) NOT NULL, -- Isme 'IN' ya 'OUT' save hoga
    quantity INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Testing ke liye pehle se 2 employees aur 2 parts daal dete hain
INSERT INTO users (emp_id, name) VALUES ('EMP101', 'Rahul Sharma'), ('EMP102', 'Amit Verma');
INSERT INTO parts (part_id, part_name, stock_count) VALUES ('PART-BRAKE-01', 'Brake Pads', 10), ('PART-OIL-02', 'Engine Oil 5L', 5);

