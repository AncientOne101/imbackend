/*如果数据库中不存在user_info表,则创建user_info表,如果有则覆盖*/
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`(
    /*用户id*/
    `id`               bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id,自增,唯一主键',
    /*用户昵称*/
    `name`             varchar(255)    NOT NULL DEFAULT '' COMMENT '用户昵称',
    /*用户密码*/
    `password`         varchar(255)    NOT NULL DEFAULT '' COMMENT '用户密码',
    /*用户邮箱*/
    `email`            varchar(255)    NOT NULL DEFAULT '' COMMENT '用户邮箱',
    /*用户性别*/
    `gender`           tinyint         NOT NULL DEFAULT 0 COMMENT '用户性别,0表示未知,1表示男,2表示女',
    /*创建时间*/
    `create_time`      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    /*修改时间*/
    `update_time`      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    /*逻辑删除,默认为0,表示未删除,1表示删除*/
    `is_deleted`       bigint(1)      NOT NULL DEFAULT 0 COMMENT '逻辑删除,默认为0,表示未删除,1表示删除',
    /*用户头像*/
    `avatar_url`       varchar(255)    NOT NULL DEFAULT '' COMMENT '用户头像',
    /*唯一主键*/
    PRIMARY KEY (`id`),
    /*唯一键*/
    UNIQUE KEY `idx_email_unique` (`email`),
    /*唯一键*/
    UNIQUE KEY `idx_name_unique` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;




