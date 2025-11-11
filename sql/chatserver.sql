CREATE TABLE `users`
(
    `id`            BIGINT       NOT NULL AUTO_INCREMENT COMMENT '用户唯一ID',
    `username`      VARCHAR(50)  NOT NULL COMMENT '用户名（唯一）',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `created_at`    DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_active`     TINYINT  DEFAULT 1 COMMENT '账号状态（1=正常，0=禁用）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB CHARSET=utf8mb4 COMMENT='用户表';


CREATE TABLE `chat_histories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '对话组唯一ID',
    `user_id` BIGINT NOT NULL COMMENT '关联用户ID',
    `username` VARCHAR(50) DEFAULT NULL COMMENT '对话标题',
    `filepath` VARCHAR(100) DEFAULT NULL COMMENT '文件地址',
    `start_time` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '对话开始时间',
    `end_time` DATETIME DEFAULT NULL COMMENT '对话结束时间',
    `is_deleted` TINYINT DEFAULT 0 COMMENT '软删除标识（1=已删除，0=正常）',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB CHARSET=utf8mb4 COMMENT='历史对话表';


CREATE TABLE `applications`
(
    `id`          BIGINT       NOT NULL AUTO_INCREMENT COMMENT '应用唯一ID',
    `name`        VARCHAR(250) NOT NULL COMMENT '应用名称',
    `app_key`     VARCHAR(255)  NOT NULL COMMENT '应用地址（代码调用用）',
    `description` TEXT         DEFAULT NULL COMMENT '应用描述',
    `is_enabled`  TINYINT      DEFAULT 1 COMMENT '启用状态（1=可调用，0=禁用）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB CHARSET=utf8mb4 COMMENT='应用信息表';