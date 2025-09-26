-- ============================================================================
-- 主配置表：整合项目、环境、模块信息的一体化配置存储表
-- 该表将项目基本信息、环境配置信息、模块信息全部整合在一张表中
-- 便于查询和管理，适合中小型项目的配置管理需求
-- ============================================================================
CREATE TABLE IF NOT EXISTS config_master (
    -- 主键ID，自动递增
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- 项目名称：用于显示和识别项目
    project_name VARCHAR(100),
    -- 项目编码：项目的唯一标识符，用于程序中引用
    project_code VARCHAR(50),
    -- 环境名称：用于显示环境信息（如：开发环境、测试环境等）
    env_name VARCHAR(50),
    -- 环境编码：环境的唯一标识符，用于程序中引用（如：dev、test、prod）
    env_code VARCHAR(20),
    -- 模块名称：配置模块的显示名称（如：Redis配置、数据库配置等）
    module_name VARCHAR(100),
    -- 模块编码：配置模块的唯一标识符（如：redis、database、kafka等）
    module_code VARCHAR(50),
    -- 配置键：具体的配置项键名（如：spring.redis.host、server.port等）
    config_key VARCHAR(200),
    -- 自动别名：配置项的自动生成别名
    auto_alias VARCHAR(50),
    -- 配置别名：用户自定义的配置项名称
    config_alias VARCHAR(50),
    -- 配置值：配置项的具体值（如：localhost、6379、8080等）
    config_value TEXT,
    -- 配置类型：标识配置值的数据类型，便于程序处理
    -- 可选值：string（字符串）、number（数字）、boolean（布尔值）、json（JSON对象）、yaml（YAML格式）
    config_type VARCHAR(20) DEFAULT 'string',
    -- 配置描述：对配置项的详细说明，便于理解和维护
    description TEXT,
    -- 是否加密：标识配置值是否经过加密处理
    -- 0表示未加密，1表示已加密
    is_encrypted INTEGER DEFAULT 0,
    -- 排序字段：用于在展示时控制配置项的显示顺序
    sort_order INTEGER DEFAULT 0,
    -- 创建时间：记录配置项的创建时间
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 更新时间：记录配置项的最后更新时间
    updated_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 删除标识：用于逻辑删除，0表示未删除，1表示已删除
    is_deleted INTEGER DEFAULT 0,
    -- 唯一约束：确保同一项目同一环境同一模块下配置键的唯一性
    UNIQUE (project_code, env_code, module_code, config_key)
);

-- ============================================================================
-- 创建索引以提高查询性能
-- ============================================================================
-- 项目和环境组合索引：用于快速查询特定项目特定环境的所有配置
CREATE INDEX IF NOT EXISTS idx_config_project_env ON config_master (project_code, env_code);

-- 模块编码索引：用于快速查询特定模块的所有配置
CREATE INDEX IF NOT EXISTS idx_config_module ON config_master (module_code);

-- 配置键索引：用于快速查询特定配置键的所有配置
CREATE INDEX IF NOT EXISTS idx_config_key ON config_master (config_key);

-- 更新时间索引：用于按时间排序查询配置变更
CREATE INDEX IF NOT EXISTS idx_config_updated_time ON config_master (updated_time);

-- ============================================================================
-- 配置变更记录表：记录所有配置项的增删改操作历史
-- 采用简化设计，每次操作只记录一条变更记录
-- 变更内容以JSON格式保存，便于解析和展示
-- ============================================================================
CREATE TABLE IF NOT EXISTS config_change_log (
    -- 主键ID，自动递增
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- 项目编码：标识变更涉及的项目
    project_code VARCHAR(50),
    -- 环境编码：标识变更涉及的环境
    env_code VARCHAR(20),
    -- 模块编码：标识变更涉及的模块
    module_code VARCHAR(50),
    -- 配置键：标识变更的具体配置项
    config_key VARCHAR(200),
    -- 操作类型：标识本次变更的操作类型
    -- 可选值：INSERT（新增）、UPDATE（修改）、DELETE（删除）
    operation_type VARCHAR(10),
    -- 旧值：修改前的配置值或删除的配置值
    -- 对于INSERT操作，此字段为空
    old_value TEXT,
    -- 新值：修改后的配置值或新增的配置值
    -- 对于DELETE操作，此字段为空
    new_value TEXT,
    -- 变更内容：完整的变更信息，以JSON格式存储
    -- 包含操作类型、项目编码、环境编码、模块编码、配置键、旧值、新值、时间戳等信息
    change_content TEXT,
    -- 变更人：记录执行变更操作的用户或系统
    changed_by VARCHAR(50),
    -- 变更原因：记录本次变更的原因说明
    change_reason TEXT,
    -- 创建时间：记录变更发生的时间
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 创建索引以提高变更记录查询性能
-- ============================================================================
-- 项目和环境组合索引：用于快速查询特定项目特定环境的变更记录
CREATE INDEX IF NOT EXISTS idx_change_log_project_env ON config_change_log (project_code, env_code);

-- 配置键索引：用于快速查询特定配置项的变更历史
CREATE INDEX IF NOT EXISTS idx_change_log_config_key ON config_change_log (config_key);

-- 操作类型索引：用于按操作类型筛选变更记录
CREATE INDEX IF NOT EXISTS idx_change_log_operation ON config_change_log (operation_type);

-- 时间索引：用于按时间范围查询变更记录
CREATE INDEX IF NOT EXISTS idx_change_log_time ON config_change_log (created_time);

-- ============================================================================
-- 配置变更自动记录触发器
-- 通过数据库触发器自动记录所有配置项的增删改操作
-- 确保变更记录的完整性和一致性
-- ============================================================================
-- ----------------------------------------------------------------------------
-- 插入操作触发器：当向config_master表插入新记录时自动触发
-- ----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS config_insert_trigger AFTER INSERT ON config_master FOR EACH ROW BEGIN
-- 将新增操作记录到变更日志表中
INSERT INTO
    config_change_log (
        project_code,
        env_code,
        module_code,
        config_key,
        operation_type,
        new_value,
        change_content,
        changed_by
    )
VALUES
    (
        NEW.project_code,
        NEW.env_code,
        NEW.module_code,
        NEW.config_key,
        'INSERT',
        NEW.config_value,
        -- 构造JSON格式的变更内容
        '{"operation":"INSERT","project_code":"' || NEW.project_code || '","env_code":"' || NEW.env_code || '","module_code":"' || NEW.module_code || '","config_key":"' || NEW.config_key || '","new_value":"' || IFNULL (NEW.config_value, '') || '","timestamp":"' || datetime ('now') || '"}',
        'system' -- 系统自动触发
    );

END;

-- ----------------------------------------------------------------------------
-- 更新操作触发器：当更新config_master表记录时自动触发
-- 只有当配置值或删除状态发生变化时才记录变更
-- ----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS config_update_trigger AFTER
UPDATE ON config_master FOR EACH ROW WHEN OLD.config_value != NEW.config_value
OR OLD.is_deleted != NEW.is_deleted BEGIN
-- 将更新操作记录到变更日志表中
INSERT INTO
    config_change_log (
        project_code,
        env_code,
        module_code,
        config_key,
        operation_type,
        old_value,
        new_value,
        change_content,
        changed_by
    )
VALUES
    (
        NEW.project_code,
        NEW.env_code,
        NEW.module_code,
        NEW.config_key,
        'UPDATE',
        OLD.config_value,
        NEW.config_value,
        -- 构造JSON格式的变更内容，包含旧值和新值
        '{"operation":"UPDATE","project_code":"' || NEW.project_code || '","env_code":"' || NEW.env_code || '","module_code":"' || NEW.module_code || '","config_key":"' || NEW.config_key || '","old_value":"' || IFNULL (OLD.config_value, '') || '","new_value":"' || IFNULL (NEW.config_value, '') || '","timestamp":"' || datetime ('now') || '"}',
        'system' -- 系统自动触发
    );

END;

-- ----------------------------------------------------------------------------
-- 删除操作触发器：当逻辑删除config_master表记录时自动触发
-- 只有当删除状态从0变为1时才记录变更
-- ----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS config_delete_trigger AFTER
UPDATE ON config_master FOR EACH ROW WHEN OLD.is_deleted = 0
AND NEW.is_deleted = 1 BEGIN
-- 将删除操作记录到变更日志表中
INSERT INTO
    config_change_log (
        project_code,
        env_code,
        module_code,
        config_key,
        operation_type,
        old_value,
        change_content,
        changed_by
    )
VALUES
    (
        NEW.project_code,
        NEW.env_code,
        NEW.module_code,
        NEW.config_key,
        'DELETE',
        OLD.config_value,
        -- 构造JSON格式的变更内容，记录被删除的值
        '{"operation":"DELETE","project_code":"' || NEW.project_code || '","env_code":"' || NEW.env_code || '","module_code":"' || NEW.module_code || '","config_key":"' || NEW.config_key || '","deleted_value":"' || IFNULL (OLD.config_value, '') || '","timestamp":"' || datetime ('now') || '"}',
        'system' -- 系统自动触发
    );

END;
