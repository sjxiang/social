

CREATE DATABASE IF NOT EXISTS `social` DEFAULT CHARACTER SET = 'utf8mb4';

-- 表 1
-- 用户名和邮箱, 唯一
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `email`               varchar(255)        NOT NULL,
    `username`            varchar(255)        NOT NULL,
    `password`            varchar(255)        NOT NULL DEFAULT '',
    `is_active`           tinyint(4)          NOT NULL DEFAULT '0' COMMENT '0 未激活, 1 激活',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `role_id`             bigint(20)          NOT NULL DEFAULT '1' COMMENT '角色id',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_email` (`email`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';


-- 表 2
-- 角色, 单一
DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name`           varchar(64)         NOT NULL,
    `level`          tinyint(4)          NOT NULL DEFAULT '0' COMMENT '0 游客, 1 用户, 2 版主, 3 管理员',
    `description`    varchar(1024)       NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表';

INSERT INTO `roles` (`name`, `level`, `description`) VALUES ('guest', '0', 'A guest can only read posts');

INSERT INTO `roles` (`name`, `level`, `description`) VALUES ('user', '1', 'A user can create posts and comments');

INSERT INTO `roles` (`name`, `level`, `description`) VALUES ('moderator', '2', 'A moderator can update other users posts');

INSERT INTO `roles` (`name`, `level`, `description`) VALUES ('admin', '3', 'An admin can update and delete other users posts');


-- 表 3
-- 多个帖子可以是同一个用户发的
DROP TABLE IF EXISTS `posts`;

CREATE TABLE `posts` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `title`               varchar(255)        NOT NULL COMMENT '帖子标题',
    `user_id`             bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `content`             varchar(1024)       NOT NULL DEFAULT '' COMMENT '帖子内容',
    `tags`                varchar(100)        NOT NULL DEFAULT '' COMMENT '标签',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `version`             tinyint(4)          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='帖子表';


-- 表 4
-- 一个帖子可以有多个评论
DROP TABLE IF EXISTS `comments`;

CREATE TABLE `comments` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `post_id`             bigint(20)          NOT NULL COMMENT '帖子id',
    `user_id`             bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `content`             varchar(1024)       NOT NULL COMMENT '评论内容',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='评论表';


-- 表 5
-- 一个阿婆主可以有多个粉丝
DROP TABLE IF EXISTS `followers`;

CREATE TABLE `followers` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`             bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `follower_uid`        bigint(20)          NOT NULL COMMENT '粉丝的用户id',
    `created_at`          timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='关注表';


-- 表 6
-- 一个用户可以有多个邀请码
DROP TABLE IF EXISTS `user_invitations`;

CREATE TABLE `user_invitations` (
    `id`                  bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`             bigint(20)          NOT NULL COMMENT '阿婆主的用户id',
    `token`               VARCHAR(64)         NOT NULL COMMENT '邀请码',
    `expiry`              timestamp           NOT NULL COMMENT '过期时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='邀请表';

