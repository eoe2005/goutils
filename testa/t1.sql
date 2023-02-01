 CREATE TABLE IF NOT EXISTS `tb_msg_chat_content` (
 	`id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
 	`msg_key` varchar(44) NOT NULL DEFAULT '' COMMENT '消息key',
 	`user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '发布人',
 	`target_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标用户级id',
 	`msg` varchar(1024) NOT NULL DEFAULT '' COMMENT '消息内容',
 	`create_at` datetime NOT NULL DEFAULT current_timestamp() COMMENT '创建时间',
 	`update_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '最后更新时间',
 	PRIMARY KEY (`id`)
   ) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COMMENT='im消息 内容';
