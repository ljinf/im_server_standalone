CREATE
DATABASE `im_chat` ;

USE
`im_chat`;

DROP TABLE IF EXISTS `register`;
CREATE TABLE `register`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    bigint(20) unsigned NOT NULL COMMENT '用户id',
    `phone`      varchar(11) DEFAULT NULL COMMENT '手机号',
    `email`      varchar(64) DEFAULT NULL COMMENT '邮箱',
    `password`   varchar(64) DEFAULT NULL COMMENT '密码',
    `created_at` datetime    DEFAULT NULL,
    `updated_at` datetime    DEFAULT NULL,
    `deleted_at` datetime    DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user` (`user_id`),
    KEY          user_phone_email_pass(`user_id`,`phone`,`email`,`password`),
    KEY          `deleted_idx` (`deleted_at`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT '注册信息表';


DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id`           bigint(20) unsigned NOT NULL COMMENT '用户id',
    `nick_name`         varchar(100) DEFAULT NULL COMMENT '昵称',
    `avatar`            varchar(255) DEFAULT NULL COMMENT '头像',
    `gender`            tinyint(2) DEFAULT 3 COMMENT '性别 1男 2女 3未知',
    `birth_day`         varchar(50)  DEFAULT NULL COMMENT '生日',
    `self_signature`    varchar(255) DEFAULT NULL COMMENT '个性签名',
    `friend_allow_type` int(10) NOT NULL DEFAULT '1' COMMENT '加好友验证类型（Friend_AllowType） 1无需验证 2需要验证',
    `silent_flag`       int(10) NOT NULL DEFAULT '0' COMMENT '禁言标识 1禁言',
    `status`            int(20) NOT NULL DEFAULT '1' COMMENT '用户状态  0:异常  1:正常',
    `created_at`        datetime     DEFAULT NULL,
    `updated_at`        datetime     DEFAULT NULL,
    `deleted_at`        datetime     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user` (`user_id`),
    KEY                 user_phone_email_pass(`user_id`,`nick_name`,`avatar`,`gender`),
    KEY                 `deleted_idx` (`deleted_at`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT '用户信息表';


DROP TABLE IF EXISTS `relationship_list`;
CREATE TABLE `relationship_list`
(
    `id`                BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`           BIGINT(20) UNSIGNED NOT NULL COMMENT '用户id 拥有者',
    `target_id`         BIGINT(20) UNSIGNED NOT NULL COMMENT '用户id 对方',
    `remark`            VARCHAR(64)  DEFAULT '' COMMENT '对方的别名备注',
    `relationship_type` TINYINT(2) DEFAULT '1' COMMENT '关系类型  1好友 2关注',
    `status`            TINYINT(2) DEFAULT '1' COMMENT '状态 1正常 2拉黑 3删除',
    `extra`             VARCHAR(256) DEFAULT '' COMMENT '其他信息',
    `created_at`        DATETIME     DEFAULT NULL,
    `updated_at`        DATETIME     DEFAULT NULL,
    `deleted_at`        DATETIME     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_target_relation_idx` (`user_id`,`target_id`,`relationship_type`),
    KEY                 `updated_idx`(`updated_at`),
    KEY                 `deleted_idx` (`deleted_at`)
) ENGINE=INNODB  DEFAULT CHARSET=utf8mb4 COMMENT '关系信息表';

DROP TABLE IF EXISTS `apply_friendship_list`;
CREATE TABLE `apply_friendship_list`
(
    `id`          BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT '用户id 拥有者',
    `target_id`   BIGINT(20) UNSIGNED NOT NULL COMMENT '用户id 对方',
    `remark`      VARCHAR(64)  DEFAULT '' COMMENT '对方的别名备注',
    `status`      TINYINT(2) DEFAULT '1' COMMENT '状态 1申请中 2通过 3被拒绝',
    `description` VARCHAR(256) DEFAULT '' COMMENT '申请描述',
    `created_at`  DATETIME     DEFAULT NULL,
    `updated_at`  DATETIME     DEFAULT NULL,
    `deleted_at`  DATETIME     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_target_idx` (`user_id`,`target_id`),
    KEY           `updated_idx`(`updated_at`),
    KEY           `deleted_idx` (`deleted_at`)
) ENGINE=INNODB  DEFAULT CHARSET=utf8mb4 COMMENT '好友申请记录表';


DROP TABLE IF EXISTS `user_conversation_list`;
CREATE TABLE `user_conversation_list`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id`         bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `conversation_id` varchar(64) NOT NULL COMMENT '会话ID',
    `last_read_seq`   bigint(20) unsigned DEFAULT 0 COMMENT '此会话用户已读的最后一条消息',
    `notify_type`     int(11) DEFAULT 0 COMMENT '会话收到消息的提醒类型，0未屏蔽，正常提醒 1屏蔽 2强提醒',
    `is_top`          tinyint(2) DEFAULT 0 COMMENT '会话是否被置顶展示',
    `created_at`      int(11) NOT NULL DEFAULT 0,
    `updated_at`      int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_conversation_idx` (`user_id`,`conversation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会话链';

DROP TABLE IF EXISTS `user_msg_list`;
CREATE TABLE `user_msg_list`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id`         bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `msg_id`          bigint(20) unsigned NOT NULL COMMENT '消息ID',
    `conversation_id` varchar(64) NOT NULL COMMENT '会话ID',
    `seq`             bigint(20) unsigned DEFAULT 0 COMMENT '消息在会话中的序列号，用于保证消息的顺序',
    `created_at`      int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY               `user_conversation_seq_msg_idx` (`user_id`,`conversation_id`,`seq`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户消息链';

DROP TABLE IF EXISTS `msg_list`;
CREATE TABLE `msg_list`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `user_id`         bigint(20) unsigned NOT NULL COMMENT '发送者ID',
    `msg_id`          bigint(20) unsigned NOT NULL COMMENT '消息ID',
    `conversation_id` varchar(64) NOT NULL COMMENT '会话ID',
    `content`         text        NOT NULL COMMENT '消息文本',
    `content_type`    int(8) NOT NULL DEFAULT '1' COMMENT '内容类型  1文本  2图片 3音频文件  4音频文件  5实时语音  6实时视频',
    `status`          int(11) NOT NULL DEFAULT '0' COMMENT '消息状态枚举，0可见 1屏蔽 2撤回',
    `send_time`       DATETIME    NOT NULL COMMENT '发送时间',
    `created_at`      int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY               `msg_idx` (`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

DROP TABLE IF EXISTS `conversation_list`;
CREATE TABLE `conversation_list`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `conversation_id` varchar(64) NOT NULL COMMENT '会话ID',
    `type`            int(11) NOT NULL DEFAULT '0' COMMENT '会话类型枚举，0单聊 1群聊',
    `member`          int(11) NOT NULL DEFAULT '0' COMMENT '与会话相关的用户数量',
    `avatar`          varchar(256) DEFAULT '' COMMENT '群组头像',
    `announcement`    text COMMENT '群公告',
    `recent_msg_time` DATETIME    NOT NULL COMMENT '此会话最新产生消息的时间',
    `created_at`      int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY               conversation_idx(`conversation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话表';

DROP TABLE IF EXISTS `conversation_msg_list`;
CREATE TABLE `conversation_msg_list`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `conversation_id` varchar(64) NOT NULL COMMENT '会话ID',
    `msg_id`          bigint(20) unsigned NOT NULL COMMENT '消息ID',
    `seq`             bigint(20) unsigned DEFAULT 0 COMMENT '消息在会话中的序列号，用于保证消息的顺序',
    `created_at`      int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY               conversation_seq_msg_idx(`conversation_id`,`seq`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会话消息链';



