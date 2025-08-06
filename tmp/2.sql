-- 创建数据库
CREATE DATABASE IF NOT EXISTS inner_join_demo;
USE inner_join_demo;

-- 创建部门表
CREATE TABLE departments (
    dept_id INT PRIMARY KEY AUTO_INCREMENT,
    dept_name VARCHAR(50) NOT NULL
) ENGINE=InnoDB;

-- 创建员工表（包含外键）
CREATE TABLE employees (
    emp_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    salary DECIMAL(10,2),
    dept_id INT,
    hire_date DATE,
    FOREIGN KEY (dept_id) REFERENCES departments(dept_id)
) ENGINE=InnoDB;

-- 插入部门数据
INSERT INTO departments (dept_name) VALUES
('研发部'),
('市场部'),
('财务部'),
('人事部');

-- 插入员工数据
INSERT INTO employees (name, salary, dept_id, hire_date) VALUES
('张三', 15000.00, 1, '2020-03-15'),
('李四', 12000.50, 2, '2021-07-22'),
('王五', 18000.00, 1, '2019-11-30'),
('赵六', 9000.00, 3, '2022-02-18'),
('钱七', 13500.00, NULL, '2021-09-05'),  -- 未分配部门
('孙八', 16000.00, 4, '2020-12-10');