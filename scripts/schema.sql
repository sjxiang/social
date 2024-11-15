

CREATE DATABASE IF NOT EXISTS `social` DEFAULT CHARACTER SET = 'utf8mb4';

-- 表 1
-- 用户名和邮箱, 唯一
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `email`               varchar(255)        NOT NULL,
    `username`            varchar(255)        NOT NULL,
    `password`            varchar(255)        NOT NULL DEFAULT ''  COMMENT '密码',
    `is_active`           tinyint(4)          NOT NULL DEFAULT '0' COMMENT '0 未激活, 1 激活',
    `role`                tinyint(4)          NOT NULL DEFAULT '0' COMMENT '角色, 0 游客, 1 用户, 2 版主, 3 管理员',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_email` (`email`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';



-- 表 2

DROP TABLE IF EXISTS `plans`;

CREATE TABLE `plans` (
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `plan_name`      varchar(64)         NOT NULL DEFAULT ''  COMMENT '订阅计划名称',
    `plan_amount`    bigint              NOT NULL DEFAULT '0' COMMENT '订阅计划费用',
    `created_at`     timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_plan_name` (`plan_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订阅计划表';


-- 表 3
DROP TABLE IF EXISTS `user_plans`;

CREATE TABLE `user_plans` (
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`        bigint              NOT NULL DEFAULT '0' COMMENT '用户 id',
    `plan_id`        bigint              NOT NULL DEFAULT '0' COMMENT '订阅计划 id',
    `created_at`     timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户订阅表';



-- 表 4
DROP TABLE IF EXISTS `posts`;

CREATE TABLE `posts` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `title`               varchar(255)        NOT NULL COMMENT '帖子标题',
    `user_id`             bigint(20)          NOT NULL COMMENT '作者的用户id',
    `content`             varchar(1024)       NOT NULL DEFAULT '' COMMENT '帖子内容',
    `tags`                varchar(100)        NOT NULL DEFAULT '' COMMENT '标签',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version`             tinyint(4)          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子表';

-- 掘金小册专栏, 连载 


-- 表 5
DROP TABLE IF EXISTS `comments`;

CREATE TABLE `comments` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `post_id`             bigint(20)          NOT NULL COMMENT '帖子id',
    `user_id`             bigint(20)          NOT NULL COMMENT '评论者的用户id',
    `content`             varchar(1024)       NOT NULL COMMENT '评论内容',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表';

-- 针对特朗普当选的新闻稿, 一口气喷了 30 楼 


-- 表 6
DROP TABLE IF EXISTS `followers`;

CREATE TABLE `followers` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`             bigint(20)          NOT NULL COMMENT '用户id',
    `follower_id`        bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    CONSTRAINT unique_follow UNIQUE (user_id, follower_id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='粉丝表';

-- 例, 我关注了峰哥亡命天涯、马督工 ...



-- 表 7
DROP TABLE IF EXISTS `user_invitations`;

CREATE TABLE `user_invitations` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`             bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `token`               VARCHAR(64)         NOT NULL COMMENT '激活码',
    `expiry`              timestamp           NOT NULL COMMENT '过期时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户激活码';

-- 注册完之后, 会给你提交的邮箱发激活码
