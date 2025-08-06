-- 外键特性演示 SQL 文件
-- 创建数据库
CREATE DATABASE IF NOT EXISTS fk_demo;
USE fk_demo;

-- =====================================================================
-- 1. 基础外键约束演示
-- =====================================================================
-- 创建主表：部门表
CREATE TABLE departments (
    dept_id INT AUTO_INCREMENT PRIMARY KEY,  -- 主键
    dept_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- 创建从表：员工表（基础外键约束）
CREATE TABLE employees (
    emp_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    dept_id INT,  -- 外键字段
    
    -- 基础外键约束：引用departments表的dept_id
    -- 默认行为：RESTRICT（拒绝违反约束的操作）
    FOREIGN KEY (dept_id) 
        REFERENCES departments(dept_id)
) ENGINE=InnoDB;

-- 插入部门数据
INSERT INTO departments (dept_name) VALUES
('研发部'),
('市场部'),
('财务部');

-- 插入有效员工数据（dept_id存在）
INSERT INTO employees (name, dept_id) VALUES
('张三', 1),  -- 研发部
('李四', 2); -- 市场部

-- 尝试插入无效员工（违反外键约束）
-- 错误：Cannot add or update a child row: a foreign key constraint fails
INSERT INTO employees (name, dept_id) VALUES ('王五', 99);

-- =====================================================================
-- 2. 级联删除 (ON DELETE CASCADE)
-- =====================================================================
-- 创建新员工表（带级联删除）
CREATE TABLE employees_cascade (
    emp_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    dept_id INT,
    
    -- 级联删除：当部门被删除时，自动删除该部门所有员工
    FOREIGN KEY (dept_id) 
        REFERENCES departments(dept_id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- 插入测试数据
INSERT INTO employees_cascade (name, dept_id) VALUES
('赵六', 1),  -- 研发部
('钱七', 3); -- 财务部

-- 验证级联删除
-- 删除研发部（dept_id=1）
DELETE FROM departments WHERE dept_id = 1;

-- 检查结果：employees_cascade表中dept_id=1的员工自动删除
SELECT * FROM employees_cascade;  -- 只剩钱七（dept_id=3）

-- =====================================================================
-- 3. 级联更新 (ON UPDATE CASCADE)
-- =====================================================================
-- 创建新员工表（带级联更新）
CREATE TABLE employees_update (
    emp_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    dept_id INT,
    
    -- 级联更新：当部门ID变更时，自动更新员工部门ID
    FOREIGN KEY (dept_id) 
        REFERENCES departments(dept_id)
        ON UPDATE CASCADE
) ENGINE=InnoDB;

-- 插入测试数据
INSERT INTO employees_update (name, dept_id) VALUES
('孙八', 2),  -- 市场部（当前dept_id=2）
('周九', 3); -- 财务部

-- 更新部门ID（将市场部dept_id从2改为200）
UPDATE departments SET dept_id = 200 WHERE dept_id = 2;

-- 检查结果：员工表中部门ID自动更新
SELECT * FROM employees_update;  -- 孙八的dept_id变为200

-- =====================================================================
-- 4. SET NULL 操作 (ON DELETE SET NULL)
-- =====================================================================
-- 创建新员工表（带SET NULL）
CREATE TABLE employees_setnull (
    emp_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    dept_id INT,
    
    -- 当部门删除时，员工部门ID设为NULL
    FOREIGN KEY (dept_id) 
        REFERENCES departments(dept_id)
        ON DELETE SET NULL
) ENGINE=InnoDB;

-- 插入测试数据
INSERT INTO employees_setnull (name, dept_id) VALUES
('吴十', 200),  -- 市场部（新ID）
('郑十一', 3); -- 财务部

-- 删除财务部
DELETE FROM departments WHERE dept_id = 3;

-- 检查结果：郑十一的dept_id变为NULL
SELECT * FROM employees_setnull;  -- 郑十一的dept_id为NULL

-- =====================================================================
-- 5. 复合外键演示
-- =====================================================================
-- 创建主表：项目表（复合主键）
CREATE TABLE projects (
    project_code VARCHAR(10),
    sub_id INT,
    project_name VARCHAR(50) NOT NULL,
    PRIMARY KEY (project_code, sub_id)  -- 复合主键
) ENGINE=InnoDB;

-- 创建从表：任务表（复合外键）
CREATE TABLE tasks (
    task_id INT AUTO_INCREMENT PRIMARY KEY,
    description TEXT,
    project_code VARCHAR(10),
    sub_id INT,
    
    -- 复合外键：引用projects表的复合主键
    FOREIGN KEY (project_code, sub_id)
        REFERENCES projects(project_code, sub_id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

-- 插入项目数据
INSERT INTO projects (project_code, sub_id, project_name) VALUES
('WEB', 1, '官网改版'),
('APP', 2, '移动端开发');

-- 插入有效任务
INSERT INTO tasks (description, project_code, sub_id) VALUES
('设计首页', 'WEB', 1),
('开发登录模块', 'APP', 2);

-- 尝试插入无效任务（违反复合外键）
INSERT INTO tasks (description, project_code, sub_id) VALUES
('测试功能', 'DATA', 1);  -- 项目不存在

-- =====================================================================
-- 6. 外键约束验证
-- =====================================================================
-- 尝试删除有外键引用的主表记录（默认RESTRICT行为）
-- 错误：Cannot delete or update a parent row: a foreign key constraint fails
DELETE FROM departments WHERE dept_id = 200;

-- 解决方案：先删除关联的子表记录
DELETE FROM employees_update WHERE dept_id = 200;
DELETE FROM departments WHERE dept_id = 200;  -- 现在可以成功删除

-- =====================================================================
-- 7. 外键管理操作
-- =====================================================================
-- 查看表的外键约束
SHOW CREATE TABLE employees_cascade;

-- 临时禁用外键检查（用于批量导入数据）
SET FOREIGN_KEY_CHECKS = 0;

-- 此时可以插入违反外键约束的数据（仅用于演示，实际慎用）
INSERT INTO employees (name, dept_id) VALUES ('临时员工', 99);

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 删除外键约束
ALTER TABLE employees DROP FOREIGN KEY employees_ibfk_1;

-- 添加新外键约束
ALTER TABLE employees 
ADD CONSTRAINT fk_dept 
FOREIGN KEY (dept_id) 
REFERENCES departments(dept_id)
ON DELETE SET NULL;

-- =====================================================================
-- 验证总结
-- =====================================================================
-- 最终数据状态查询
SELECT * FROM departments;
SELECT * FROM employees;
SELECT * FROM employees_cascade;
SELECT * FROM employees_update;
SELECT * FROM employees_setnull;
SELECT * FROM projects;
SELECT * FROM tasks;

-- 清理数据库（取消注释执行）
-- DROP DATABASE IF EXISTS fk_demo;