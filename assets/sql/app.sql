CREATE TABLE IF NOT EXISTS `system_config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
  `tenant` varchar(32) NOT NULL DEFAULT 'MAIN_SITE' COMMENT '租户名称',
  `env` varchar(50) NOT NULL COMMENT '环境',
  `type` varchar(32) NOT NULL COMMENT '配置类型',
  `config` mediumtext DEFAULT NULL COMMENT '配置内容',
  `description` varchar(256) DEFAULT NULL COMMENT '描述',
  `creator` varchar(32) DEFAULT NULL COMMENT '创建人',
  `modifier` varchar(32) DEFAULT NULL COMMENT '修改人',
  `deleted_at` timestamp(3) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_system_config_deleted_at` (`deleted_at`)
) AUTO_INCREMENT = 1400002 DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '系统配置表';
