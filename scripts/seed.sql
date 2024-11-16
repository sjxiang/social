

INSERT INTO `plans` (`plan_name`, `plan_amount`) VALUES ('free', '0');
INSERT INTO `plans` (`plan_name`, `plan_amount`) VALUES ('basic', '1000');
INSERT INTO `plans` (`plan_name`, `plan_amount`) VALUES ('pro', '5000');
INSERT INTO `plans` (`plan_name`, `plan_amount`) VALUES ('enterprise', '10000');

INSERT INTO `followers` (`user_id`, `follower_id`, `created_At`) VALUES (1, 2, UTC_TIMESTAMP());
DELETE FROM `followers` WHERE `user_id` = 1 AND `follower_id` = 2;


INSERT INTO `posts` (`title`, `content`, `user_id`, `tags`, `created_At`, `version`) VALUES ("vue", "more ...", 7, "fe, interview", UTC_TIMESTAMP(), 0);

UPDATE `posts` SET `title` = 'react', `content` = 'todo...', `version` = `version` + 1, `updated_at` = UTC_TIMESTAMP() WHERE `id` = 1 AND `version` = 0;
