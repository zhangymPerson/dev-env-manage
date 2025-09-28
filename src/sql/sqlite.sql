-- ============================================================================
-- 主配置表：仅存储当前有效的最新配置项
-- 所有历史版本由 config_history 表管理
-- 彻底移除 is_deleted 字段，删除操作通过物理 DELETE 实现
-- ============================================================================
CREATE TABLE IF NOT EXISTS config_master (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键ID，自动递增
    project VARCHAR(100), -- 项目标识（统一原 project_name 和 project_code）
    env VARCHAR(20), -- 环境标识（如：dev、test、prod）
    module VARCHAR(50), -- 模块标识（如：redis、database）
    config_key VARCHAR(200), -- 配置项键名（如：spring.redis.host）
    auto_alias VARCHAR(50), -- 自动生成的配置别名
    config_alias VARCHAR(50), -- 用户自定义的配置别名
    config_value TEXT, -- 配置项的具体值
    config_type VARCHAR(20) DEFAULT 'string', -- 配置值类型（string/number/boolean/json/yaml）
    description TEXT, -- 配置项的详细说明
    is_encrypted INTEGER DEFAULT 0, -- 是否加密（0=未加密，1=已加密）
    sort_order INTEGER DEFAULT 0, -- 显示排序序号
    created_time DATETIME DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_time DATETIME DEFAULT CURRENT_TIMESTAMP, -- 记录最后更新时间
    -- 唯一约束：确保同一项目、环境、模块下的配置键唯一
    UNIQUE (project, env, module, config_key)
);

-- ============================================================================
-- 配置历史版本表：记录 config_master 中每一版的历史快照
-- 每次 UPDATE 或 DELETE 前，将旧记录完整存入此表
-- ============================================================================
CREATE TABLE IF NOT EXISTS config_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 主键ID，自动递增
    project VARCHAR(100), -- 项目标识
    env VARCHAR(20), -- 环境标识
    module VARCHAR(50), -- 模块标识
    config_key VARCHAR(200), -- 配置项键名
    auto_alias VARCHAR(50), -- 自动生成的配置别名
    config_alias VARCHAR(50), -- 用户自定义的配置别名
    config_value TEXT, -- 配置项的具体值
    config_type VARCHAR(20), -- 配置值类型
    description TEXT, -- 配置项的详细说明
    is_encrypted INTEGER, -- 是否加密
    sort_order INTEGER, -- 显示排序序号
    created_time DATETIME, -- 原始创建时间（来自主表）
    updated_time DATETIME, -- 最后更新时间（来自主表）
    version DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 版本时间戳，作为历史版本唯一标识
    changed_by VARCHAR(50) DEFAULT 'system' -- 变更操作者（默认为系统）
);

-- ============================================================================
-- 索引：提升主表和历史表的查询性能
-- ============================================================================
-- 主表：按项目+环境快速查询当前配置
CREATE INDEX IF NOT EXISTS idx_master_project_env ON config_master (project, env);

-- 主表：按配置键快速定位
CREATE INDEX IF NOT EXISTS idx_master_key ON config_master (config_key);

-- 历史表：按配置项+版本倒序查询变更历史
CREATE INDEX IF NOT EXISTS idx_history_key_version ON config_history (project, env, module, config_key, version DESC);

-- 历史表：按版本时间倒序查询全局变更流
CREATE INDEX IF NOT EXISTS idx_history_version ON config_history (version DESC);

-- ============================================================================
-- 触发器1：配置更新时，自动将旧记录存入历史表
-- 捕获所有有意义的字段变更（排除仅 updated_time 自动刷新）
-- ============================================================================
CREATE TRIGGER IF NOT EXISTS config_update_history_trigger AFTER
UPDATE ON config_master FOR EACH ROW WHEN OLD.config_value != NEW.config_value
OR IFNULL (OLD.auto_alias, '') != IFNULL (NEW.auto_alias, '')
OR IFNULL (OLD.config_alias, '') != IFNULL (NEW.config_alias, '')
OR OLD.config_type != NEW.config_type
OR IFNULL (OLD.description, '') != IFNULL (NEW.description, '')
OR OLD.is_encrypted != NEW.is_encrypted
OR OLD.sort_order != NEW.sort_order BEGIN
INSERT INTO
    config_history (
        project,
        env,
        module,
        config_key,
        auto_alias,
        config_alias,
        config_value,
        config_type,
        description,
        is_encrypted,
        sort_order,
        created_time,
        updated_time,
        version,
        changed_by
    )
VALUES
    (
        OLD.project,
        OLD.env,
        OLD.module,
        OLD.config_key,
        OLD.auto_alias,
        OLD.config_alias,
        OLD.config_value,
        OLD.config_type,
        OLD.description,
        OLD.is_encrypted,
        OLD.sort_order,
        OLD.created_time,
        OLD.updated_time,
        DATETIME ('now'),
        'system'
    );

END;

-- ============================================================================
-- 触发器2：配置删除时，自动将被删除的记录存入历史表
-- 物理 DELETE 操作触发，确保删除前状态被完整归档
-- ============================================================================
CREATE TRIGGER IF NOT EXISTS config_delete_history_trigger BEFORE DELETE ON config_master FOR EACH ROW BEGIN
INSERT INTO
    config_history (
        project,
        env,
        module,
        config_key,
        auto_alias,
        config_alias,
        config_value,
        config_type,
        description,
        is_encrypted,
        sort_order,
        created_time,
        updated_time,
        version,
        changed_by
    )
VALUES
    (
        OLD.project,
        OLD.env,
        OLD.module,
        OLD.config_key,
        OLD.auto_alias,
        OLD.config_alias,
        OLD.config_value,
        OLD.config_type,
        OLD.description,
        OLD.is_encrypted,
        OLD.sort_order,
        OLD.created_time,
        OLD.updated_time,
        DATETIME ('now'),
        'system'
    );

END;