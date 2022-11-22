-- MySQL dump 10.13  Distrib 8.0.23, for macos10.15 (x86_64)
--
-- Host: 8.142.177.235    Database: seed_pre
-- ------------------------------------------------------
-- Server version	8.0.30-0ubuntu0.20.04.2

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `agora_callback_record`
--

DROP TABLE IF EXISTS `agora_callback_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `agora_callback_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `notice_id` varchar(255) NOT NULL DEFAULT '' COMMENT '通知 ID，标识来自业务服务器的一次事件通知',
  `product_id` int NOT NULL DEFAULT '0' COMMENT '业务Id,1rtc2旁路推流CDN3云端录制4Cloud Player5旁路推流- 服务端',
  `event_type` int NOT NULL DEFAULT '0' COMMENT '事件类型ID',
  `notify_ms` varchar(13) NOT NULL DEFAULT '' COMMENT '对方推送时间',
  `payload` text COMMENT '详细内容',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb3 COMMENT='声网回调记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `agora_cloud_record`
--

DROP TABLE IF EXISTS `agora_cloud_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `agora_cloud_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `uid` int NOT NULL DEFAULT '0' COMMENT '用户ID',
  `listener_agora_uid` int NOT NULL DEFAULT '0' COMMENT '监听者的UID',
  `channel_name` varchar(255) NOT NULL DEFAULT '' COMMENT '频道名称',
  `resource_id` text COMMENT '声网返回的rid',
  `session_id` varchar(255) NOT NULL DEFAULT '' COMMENT '声网返回的sid',
  `status` int NOT NULL DEFAULT '0' COMMENT '0未知1已申请rid2已开始3已结束',
  `server_status` int NOT NULL DEFAULT '0' COMMENT '后端状态1未处理2已收到声网回调,开始合并视频3视频处理成功4处理异常',
  `start_time` int NOT NULL DEFAULT '0' COMMENT '开始录制时间',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '结束录制时间时间',
  `config_info` text COMMENT '请求声网,开始录制时设置的配置信息',
  `acquire_config` varchar(255) NOT NULL DEFAULT '' COMMENT '获取RID时的配置信息',
  `stop_action` int NOT NULL DEFAULT '0' COMMENT '0-未知1-正常停止2-页面刷新时拦截3-重新加载页面触发 4-声网回调触发;',
  `stop_res_info` text COMMENT '请求声网,停止录制时返回的文件信息',
  `video_url` varchar(255) NOT NULL DEFAULT '' COMMENT '最终录制好的视频的URL地址',
  `err_log` text COMMENT '请求3方返回的错误信息',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb3 COMMENT='声网录屏';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cicd_publish`
--

DROP TABLE IF EXISTS `cicd_publish`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cicd_publish` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `service_id` int NOT NULL DEFAULT '0' COMMENT '服务ID',
  `server_id` int NOT NULL DEFAULT '0' COMMENT '服务器ID',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1待部署2待发布3已发布4发布失败',
  `deploy_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1部署中2失败3完成',
  `service_info` varchar(255) NOT NULL DEFAULT '' COMMENT '服务信息-备份',
  `server_info` varchar(255) NOT NULL DEFAULT '' COMMENT '服务器信息-备份',
  `deploy_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1本地部署2远程同步部署',
  `code_dir` varchar(255) NOT NULL DEFAULT '' COMMENT '项目代码目录名',
  `log` text COMMENT '日志',
  `err_info` varchar(255) NOT NULL DEFAULT '' COMMENT '错误日志',
  `exec_time` int NOT NULL DEFAULT '0' COMMENT '执行时间',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  `step` int DEFAULT NULL COMMENT '执行到第几步了',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb3 COMMENT='cicd发布记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `email_log`
--

DROP TABLE IF EXISTS `email_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` int NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `rule_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `receiver` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者邮件地址',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '失效时间',
  `auth_code` varchar(50) NOT NULL DEFAULT '' COMMENT '验证码',
  `auth_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '验证码状态1未使用2已使用3已超时',
  `send_uid` int NOT NULL DEFAULT '0' COMMENT '发送者UID，管理员是9999，未知8888',
  `send_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者的IP',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1成功2失败3发送中4等待发送',
  `carbon_copy` varchar(255) NOT NULL DEFAULT '' COMMENT '抄送邮件地址',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='邮件发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `email_rule`
--

DROP TABLE IF EXISTS `email_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '模板内容,可变量替换',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '分类,1验证码2通知3营销',
  `day_times` int NOT NULL DEFAULT '0' COMMENT '一天最多发送次数',
  `period` int NOT NULL DEFAULT '0' COMMENT '周期时间-秒',
  `period_times` int NOT NULL DEFAULT '0' COMMENT '周期时间内-发送次数',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '验证码要有失效时间',
  `memo` varchar(255) NOT NULL DEFAULT '' COMMENT '描述，主要是给3方审核用',
  `purpose` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用途,参考代码常量',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='邮件发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `game_match_group`
--

DROP TABLE IF EXISTS `game_match_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_match_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `rule_id` int NOT NULL DEFAULT '0' COMMENT 'rule_id',
  `self_id` int NOT NULL DEFAULT '0' COMMENT '用redis自生成的自增ID',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '报名跟报名成功会各创建一条group记录，1：报名，2匹配成功',
  `person` tinyint(1) NOT NULL DEFAULT '0' COMMENT '小组总人数',
  `weight` varchar(50) NOT NULL DEFAULT '' COMMENT '小组权重',
  `match_times` tinyint(1) NOT NULL DEFAULT '0' COMMENT '已匹配过的次数，超过3次，证明该用户始终不能匹配成功，直接丢弃，不过没用到',
  `sign_timeout` int NOT NULL DEFAULT '0' COMMENT '多少秒后无人来取，即超时，更新用户状态，删除数据',
  `success_timeout` int NOT NULL DEFAULT '0' COMMENT '匹配成功后，无人来取，超时',
  `sign_time` int NOT NULL DEFAULT '0' COMMENT '报名时间',
  `success_time` int NOT NULL DEFAULT '0' COMMENT '匹配成功时间',
  `player_ids` varchar(100) NOT NULL DEFAULT '' COMMENT '用户列表',
  `addition` varchar(100) NOT NULL DEFAULT '' COMMENT '请求方附加属性值',
  `team_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '组队互相PK的时候，得成两个队伍',
  `out_group_id` int NOT NULL DEFAULT '0' COMMENT '报名时，客户端请求时，自带的一个ID',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-小组信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `game_match_push`
--

DROP TABLE IF EXISTS `game_match_push`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_match_push` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `rule_id` int NOT NULL DEFAULT '0' COMMENT 'rule_id',
  `self_id` int NOT NULL DEFAULT '0' COMMENT '用redis自生成的自增ID',
  `a_time` int NOT NULL DEFAULT '0' COMMENT '添加时间',
  `link_id` int NOT NULL DEFAULT '0' COMMENT '小组总人数',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：1未推送2推送失败，等待重试3推送成功4推送失败，不再重试',
  `times` int NOT NULL DEFAULT '0' COMMENT '已推送次数',
  `category` int NOT NULL DEFAULT '0' COMMENT '1',
  `payload` text COMMENT '自定义的载体',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-推送消息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `game_match_rule`
--

DROP TABLE IF EXISTS `game_match_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_match_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `game_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '游戏关联ID',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1正常2关闭',
  `match_timeout` tinyint(1) NOT NULL DEFAULT '0' COMMENT '匹配超时时间(秒)',
  `success_timeout` tinyint(1) NOT NULL DEFAULT '0' COMMENT '匹配成功后，对方未接收超时时间(秒)',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1.N(TEAM)VS N(TEAM)2.N人够了就行(吃鸡模式)',
  `team_max_people` tinyint(1) NOT NULL DEFAULT '0' COMMENT '队伍最大人数',
  `condition_people` int NOT NULL DEFAULT '0' COMMENT '多少人，可开始一局游戏',
  `formula` varchar(100) NOT NULL DEFAULT '' COMMENT '权限计算公式',
  `weight_team_aggregation` varchar(50) NOT NULL DEFAULT '' COMMENT '每个小组的最终权重计算聚合方法 sum min max average',
  `weight_score_min` int NOT NULL DEFAULT '0' COMMENT '权重最小值',
  `weight_score_max` int NOT NULL DEFAULT '0' COMMENT '权重最大值',
  `weight_auto_assign` tinyint(1) NOT NULL DEFAULT '0' COMMENT '权重自动匹配',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  `fps` tinyint(1) NOT NULL DEFAULT '0' COMMENT '帧同步速率',
  `ready_timeout` int NOT NULL COMMENT '进入房间准备超时间',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `game_match_success`
--

DROP TABLE IF EXISTS `game_match_success`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_match_success` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `rule_id` int NOT NULL DEFAULT '0' COMMENT 'rule_id',
  `self_id` int NOT NULL DEFAULT '0' COMMENT '用redis自生成的自增ID',
  `a_time` int NOT NULL DEFAULT '0' COMMENT '匹配成功的时间',
  `timeout` int NOT NULL DEFAULT '0' COMMENT '多少秒后无人来取后超时，更新用户状态，删除数据',
  `teams` varchar(50) NOT NULL DEFAULT '' COMMENT '该结果，有几个 队伍，因为每个队伍是一个集合，要用来索引',
  `player_ids` varchar(100) NOT NULL DEFAULT '' COMMENT '玩家ID列表',
  `group_ids` varchar(100) NOT NULL DEFAULT '' COMMENT '小组ID列表',
  `push_self_id` int NOT NULL DEFAULT '0' COMMENT '推送ID,用redis生成的自增ID',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-成功';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `instance`
--

DROP TABLE IF EXISTS `instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `instance` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `platform` int NOT NULL DEFAULT '0' COMMENT '平台类型1自有2阿里3腾讯4华为5声网',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `host` varchar(255) NOT NULL DEFAULT '' COMMENT '主机地址',
  `port` varchar(50) NOT NULL DEFAULT '' COMMENT '主机端口号',
  `env` int NOT NULL DEFAULT '0' COMMENT '环境变量,1本地2开发3测试4预发布5线上',
  `user` varchar(100) NOT NULL DEFAULT '' COMMENT '用户名',
  `ps` varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
  `ext` varchar(255) NOT NULL DEFAULT '' COMMENT '自定义配置信息',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1正常2关闭3异常',
  `charge_user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '负责人姓名',
  `start_time` int NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '结束时间',
  `price` int NOT NULL DEFAULT '0' COMMENT '价格',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb3 COMMENT='服务-实例';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mail_group`
--

DROP TABLE IF EXISTS `mail_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mail_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `rule_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `people_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '接收人群，1单发2群发3指定group4指定tag5指定UIDS',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '模板内容,可变量替换',
  `receiver` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者，groupId，tagId , all ',
  `send_uid` int NOT NULL DEFAULT '0' COMMENT '发送者UID，管理员是9999，未知8888',
  `send_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者的IP',
  `send_time` int NOT NULL DEFAULT '0' COMMENT '发送时间',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='站内信 - 群发记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mail_log`
--

DROP TABLE IF EXISTS `mail_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mail_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` int NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `rule_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `receiver` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者uid',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '失效时间',
  `auth_code` varchar(50) NOT NULL DEFAULT '' COMMENT '验证码',
  `auth_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1未使用2已使用3已超时',
  `send_uid` int NOT NULL DEFAULT '0' COMMENT '发送者UID，管理员是9999，未知8888',
  `send_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者的IP',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1成功2失败3发送中4等待发送',
  `receiver_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '接收者已读',
  `receiver_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '接收者已删除',
  `send_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '发送者已删除',
  `mail_group_id` int NOT NULL DEFAULT '0' COMMENT '群发的ID',
  `send_time` int NOT NULL DEFAULT '0' COMMENT '发送时间',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='站内信-日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mail_rule`
--

DROP TABLE IF EXISTS `mail_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `mail_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '模板内容,可变量替换',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '分类,1验证码2通知3营销',
  `day_times` int NOT NULL DEFAULT '0' COMMENT '一天最多发送次数',
  `period` int NOT NULL DEFAULT '0' COMMENT '周期时间-秒',
  `period_times` int NOT NULL DEFAULT '0' COMMENT '周期时间内-发送次数',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '验证码要有失效时间',
  `memo` varchar(255) NOT NULL DEFAULT '' COMMENT '描述，主要是给3方审核用',
  `purpose` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用途,参考代码常量',
  `people_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '接收人群，1单发2群发3指定group4指定tag5指定UIDS',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='站内信 - 发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `operation_record`
--

DROP TABLE IF EXISTS `operation_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `operation_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'ip',
  `method` varchar(50) NOT NULL DEFAULT '' COMMENT 'get|post|put|delete',
  `path` varchar(50) NOT NULL DEFAULT '' COMMENT 'uri请求路径',
  `status` int NOT NULL DEFAULT '0' COMMENT '请求状态',
  `latency` int NOT NULL DEFAULT '0' COMMENT '延迟',
  `agent` text COMMENT 'useragent',
  `error_message` varchar(255) NOT NULL DEFAULT '' COMMENT '错误信息',
  `body` text COMMENT '请求内容',
  `resp` text COMMENT '返回结果',
  `uid` int NOT NULL DEFAULT '0' COMMENT '用户Id',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='请求日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `project`
--

DROP TABLE IF EXISTS `project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型,1service 2frontend 3backend 4app',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '描述信息',
  `secret_key` varchar(100) NOT NULL DEFAULT '' COMMENT '密钥',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1正常2关闭',
  `access` varchar(255) NOT NULL DEFAULT '' COMMENT 'baseAuth 认证KEY',
  `lang` tinyint(1) NOT NULL DEFAULT '0' COMMENT '实现语言1php2go3java4js',
  `git` varchar(255) NOT NULL DEFAULT '' COMMENT 'git仓地址',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb3 COMMENT='服务/项目';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `server`
--

DROP TABLE IF EXISTS `server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `server` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `platform` int NOT NULL DEFAULT '0' COMMENT '平台类型1自有2阿里3腾讯4华为',
  `out_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '外网IP',
  `inner_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '内网IP',
  `env` int NOT NULL DEFAULT '0' COMMENT '环境变量,1本地2开发3测试4预发布5线上',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1正常2关闭',
  `ext` varchar(255) NOT NULL DEFAULT '' COMMENT '自定义配置信息',
  `charge_user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '负责人姓名',
  `start_time` int NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '结束时间',
  `price` int NOT NULL DEFAULT '0' COMMENT '价格',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3 COMMENT='服务器';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sms_log`
--

DROP TABLE IF EXISTS `sms_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sms_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` int NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `rule_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `receiver` varchar(50) NOT NULL DEFAULT '' COMMENT '接收者邮件地址',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '失效时间',
  `auth_code` varchar(50) NOT NULL DEFAULT '' COMMENT '验证码',
  `auth_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1未使用2已使用3已超时',
  `send_uid` int NOT NULL DEFAULT '0' COMMENT '发送者UID，管理员是9999，未知8888',
  `send_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '发送者的IP',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1成功2失败3发送中4等待发送',
  `out_no` varchar(50) NOT NULL DEFAULT '' COMMENT '3方ID',
  `channel` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1阿里2腾讯',
  `third_back_info` varchar(255) NOT NULL DEFAULT '' COMMENT '请示3方返回结果集',
  `third_callback_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '3方状态',
  `third_callback_info` varchar(255) NOT NULL DEFAULT '' COMMENT '3方回执-信息',
  `third_callback_time` varchar(255) NOT NULL DEFAULT '' COMMENT '3方回执-时间',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='短信发送日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sms_rule`
--

DROP TABLE IF EXISTS `sms_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sms_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '项目ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '模板内容,可变量替换',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '分类,1验证码2通知3营销',
  `day_times` int NOT NULL DEFAULT '0' COMMENT '一天最多发送次数',
  `period` int NOT NULL DEFAULT '0' COMMENT '周期时间-秒',
  `period_times` int NOT NULL DEFAULT '0' COMMENT '周期时间内-发送次数',
  `expire_time` int NOT NULL DEFAULT '0' COMMENT '验证码要有失效时间',
  `memo` varchar(255) NOT NULL DEFAULT '' COMMENT '描述，主要是给3方审核用',
  `purpose` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用途,参考代码常量',
  `channel` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1阿里2腾讯',
  `third_back_info` varchar(255) NOT NULL DEFAULT '' COMMENT '请示3方返回结果集',
  `third_template_id` varchar(100) NOT NULL DEFAULT '' COMMENT '3方模板ID',
  `third_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '3方状态',
  `third_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '3方模板审核失败，理由信息',
  `third_callback_info` varchar(255) NOT NULL DEFAULT '' COMMENT '3方回执-信息',
  `third_callback_time` varchar(255) NOT NULL DEFAULT '' COMMENT '3方回执-时间',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='短信发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statistics_log`
--

DROP TABLE IF EXISTS `statistics_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `statistics_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `header_common` text COMMENT 'http公共请求头信息',
  `header_base` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT 'http请求头客户端基础信息',
  `project_id` int NOT NULL DEFAULT '0' COMMENT '项目ID',
  `category` int NOT NULL DEFAULT '0' COMMENT '分类，暂未使用',
  `action` varchar(255) NOT NULL DEFAULT '' COMMENT '动作标识',
  `uid` int NOT NULL DEFAULT '0' COMMENT '用户ID',
  `msg` varchar(255) NOT NULL DEFAULT '' COMMENT '自定义消息体',
  `sn` varchar(100) NOT NULL COMMENT '设备-序列号',
  `system_version` varchar(100) NOT NULL COMMENT '设备-版本号',
  `record_time` int NOT NULL COMMENT '记录时间',
  `package_name` varchar(100) NOT NULL COMMENT '包名',
  `app_version` varchar(100) NOT NULL COMMENT '应用版本号',
  `app_name` varchar(100) NOT NULL COMMENT '应用名',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=65 DEFAULT CHARSET=utf8mb3 COMMENT='接收前端推送的统计日志';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `twin_agora_room`
--

DROP TABLE IF EXISTS `twin_agora_room`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `twin_agora_room` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `room_id` varchar(255) NOT NULL DEFAULT '' COMMENT '房间ID',
  `channel` varchar(255) NOT NULL DEFAULT '' COMMENT '频道',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态,1发起呼叫，2正常通话中，3已结束',
  `end_status` int NOT NULL DEFAULT '0' COMMENT '结束的状态：(1)超时，(2)某一方退出,(3)某一方拒绝(4)发起方主动取消呼叫',
  `call_uid` varchar(255) NOT NULL DEFAULT '' COMMENT '发起呼叫者',
  `receive_uids` varchar(255) NOT NULL DEFAULT '' COMMENT '接收呼叫者消息',
  `receive_uids_accept` varchar(13) NOT NULL DEFAULT '' COMMENT '被呼叫的用户IDS，接收了此次呼叫',
  `receive_uids_deny` text COMMENT '被呼叫的用户IDS，拒绝了此次呼叫',
  `uids` text COMMENT '方便调试(ReceiveUidsAccept+CallUid)',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=652 DEFAULT CHARSET=utf8mb3 COMMENT='AR远程呼叫,房间记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `uuid` varchar(50) NOT NULL DEFAULT '' COMMENT 'UID字条串化',
  `project_id` tinyint(1) NOT NULL DEFAULT '0' COMMENT '项目ID',
  `sex` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别1男2女',
  `birthday` int NOT NULL DEFAULT '0' COMMENT '出生日期,unix时间戳',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户登录名',
  `password` varchar(50) NOT NULL DEFAULT '' COMMENT '用户登录密码',
  `pay_ps` varchar(50) NOT NULL DEFAULT '' COMMENT '用户支付密码',
  `nick_name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `robot` tinyint(1) NOT NULL DEFAULT '0' COMMENT '机器人',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态1正常2禁用',
  `guest` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否游客,1是2否',
  `test` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否测试,1是2否',
  `recommend` varchar(50) NOT NULL DEFAULT '' COMMENT '推荐人',
  `header_img` varchar(50) NOT NULL DEFAULT '' COMMENT '头像url地址',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  `channel_name` varchar(255) NOT NULL DEFAULT '' COMMENT '声网频道名',
  `role` int NOT NULL DEFAULT '0' COMMENT '角色1护士2专家',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `uuid_2` (`uuid`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb3 COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_login`
--

DROP TABLE IF EXISTS `user_login`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_login` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `source_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '来源类型',
  `uid` int NOT NULL DEFAULT '0' COMMENT 'uid',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型 1email2name3mobile3third4guest',
  `third_type` varchar(50) NOT NULL DEFAULT '' COMMENT '三方平台类型,参数常量USER_TYPE_THIRD',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT '请求方传输IP',
  `auto_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '程序自己计算的IP',
  `province` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `city` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `county` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `town` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `area_detail` varchar(255) NOT NULL DEFAULT '' COMMENT '页面来源',
  `app_version` varchar(50) NOT NULL DEFAULT '' COMMENT 'APP版本',
  `os` varchar(50) NOT NULL DEFAULT '' COMMENT '操作系统',
  `os_version` varchar(50) NOT NULL DEFAULT '' COMMENT '操作系统版本',
  `device` varchar(50) NOT NULL DEFAULT '' COMMENT '设备名称',
  `device_version` varchar(50) NOT NULL DEFAULT '' COMMENT '设备版本',
  `lat` varchar(50) NOT NULL DEFAULT '' COMMENT '伟度',
  `lon` varchar(50) NOT NULL DEFAULT '' COMMENT '经度',
  `device_id` varchar(50) NOT NULL DEFAULT '' COMMENT '设备ID',
  `dpi` varchar(50) NOT NULL DEFAULT '' COMMENT '分辨率',
  `referer` varchar(255) NOT NULL DEFAULT '' COMMENT '页面来源',
  `jwt` text COMMENT '登陆成功后的jwt',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=1713 DEFAULT CHARSET=utf8mb3 COMMENT='用户登陆记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_reg`
--

DROP TABLE IF EXISTS `user_reg`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_reg` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `project_id` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `source_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '来源类型',
  `uid` int NOT NULL DEFAULT '0' COMMENT 'uid',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型 1email2name3mobile3third4guest',
  `third_type` varchar(50) NOT NULL DEFAULT '' COMMENT '三方平台类型,参数常量USER_TYPE_THIRD',
  `channel` tinyint(1) NOT NULL DEFAULT '0' COMMENT '推广渠道1平台自己',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT '请求方传输IP',
  `auto_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '程序自己计算的IP',
  `province` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `city` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `county` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `town` int NOT NULL DEFAULT '0' COMMENT 'project_id',
  `area_detail` varchar(255) NOT NULL COMMENT '页面来源',
  `app_version` varchar(50) NOT NULL DEFAULT '' COMMENT 'APP版本',
  `os` tinyint(1) NOT NULL DEFAULT '0' COMMENT '操作系统',
  `os_version` varchar(50) NOT NULL DEFAULT '' COMMENT '操作系统版本',
  `device` varchar(50) NOT NULL DEFAULT '' COMMENT '设备名称',
  `device_version` varchar(50) NOT NULL DEFAULT '' COMMENT '设备版本',
  `lat` varchar(50) NOT NULL DEFAULT '' COMMENT '伟度',
  `lon` varchar(50) NOT NULL DEFAULT '' COMMENT '经度',
  `device_id` varchar(50) NOT NULL DEFAULT '' COMMENT '设备ID',
  `dpi` varchar(50) NOT NULL DEFAULT '' COMMENT '分辨率',
  `referer` varchar(255) NOT NULL DEFAULT '' COMMENT '页面来源',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='用户注册信息';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-11-22 20:06:22
