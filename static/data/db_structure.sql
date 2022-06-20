
create table user_login(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` int  not null  default   0  comment 'project_id'  , 
    `source_type` tinyint(1)  not null  default   0  comment '来源类型'  , 
    `uid` int  not null  default   0  comment 'uid'  , 
    `type` tinyint(1)  not null  default   0  comment '类型 1email2name3mobile3third4guest'  , 
    `third_type` varchar(50)  not null  default   ''  comment '三方平台类型,参数常量USER_TYPE_THIRD'  , 
    `ip` varchar(50)  not null  default   ''  comment '请求方传输IP'  , 
    `auto_ip` varchar(50)  not null  default   ''  comment '程序自己计算的IP'  , 
    `province` int  not null  default   0  comment 'project_id'  , 
    `city` int  not null  default   0  comment 'project_id'  , 
    `county` int  not null  default   0  comment 'project_id'  , 
    `town` int  not null  default   0  comment 'project_id'  , 
    `area_detail` varchar(255)  not null  comment '页面来源'  , 
    `app_version` varchar(50)  not null  default   ''  comment 'APP版本'  , 
    `os` tinyint(1)  not null  default   0  comment '操作系统'  , 
    `os_version` varchar(50)  not null  default   ''  comment '操作系统版本'  , 
    `device` varchar(50)  not null  default   ''  comment '设备名称'  , 
    `device_version` varchar(50)  not null  default   ''  comment '设备版本'  , 
    `lat` varchar(50)  not null  default   ''  comment '伟度'  , 
    `lon` varchar(50)  not null  default   ''  comment '经度'  , 
    `device_id` varchar(50)  not null  default   ''  comment '设备ID'  , 
    `dpi` varchar(50)  not null  default   ''  comment '分辨率'  , 
    `referer` varchar(255)  not null  default   ''  comment '页面来源'  , 
    `jwt` text  null  comment '登陆成功后的jwt'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='用户登陆记录'
 ;

create table instance(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `platform` int  not null  default   0  comment '平台类型1自有2阿里3腾讯4华为'  , 
    `name` varchar(50)  not null  default   ''  comment '名称'  , 
    `host` varchar(255)  not null  default   ''  comment '主机地址'  , 
    `port` varchar(50)  not null  default   ''  comment '主机端口号'  , 
    `env` int  not null  default   0  comment '环境变量,1本地2开发3测试4预发布5线上'  , 
    `user` varchar(100)  not null  default   ''  comment '用户名'  , 
    `ps` varchar(100)  not null  default   ''  comment '密码'  , 
    `ext` varchar(255)  not null  default   ''  comment '自定义配置信息'  , 
    `status` tinyint(1)  not null  default   0  comment '状态1正常2关闭3异常'  , 
    `charge_user_name` varchar(50)  not null  default   ''  comment '负责人姓名'  , 
    `start_time` int  not null  default   0  comment '开始时间'  , 
    `end_time` int  not null  default   0  comment '结束时间'  , 
    `price` int  not null  default   0  comment '价格'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='服务-实例'
 ;

create table mail_rule(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` tinyint(1)  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '模板内容,可变量替换'  , 
    `type` tinyint(1)  not null  default   0  comment '分类,1验证码2通知3营销'  , 
    `day_times` int  not null  default   0  comment '一天最多发送次数'  , 
    `period` int  not null  default   0  comment '周期时间-秒'  , 
    `period_times` int  not null  default   0  comment '周期时间内-发送次数'  , 
    `expire_time` int  not null  default   0  comment '验证码要有失效时间'  , 
    `memo` varchar(255)  not null  default   ''  comment '描述，主要是给3方审核用'  , 
    `purpose` tinyint(1)  not null  default   0  comment '用途,参考代码常量'  , 
    `people_type` tinyint(1)  not null  default   0  comment '接收人群，1单发2群发3指定group4指定tag5指定UIDS'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='站内信 - 发送规则配置'
 ;

create table mail_group(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `rule_id` tinyint(1)  not null  default   0  comment '规则ID'  , 
    `people_type` tinyint(1)  not null  default   0  comment '接收人群，1单发2群发3指定group4指定tag5指定UIDS'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '模板内容,可变量替换'  , 
    `receiver` varchar(50)  not null  default   ''  comment '接收者，groupId，tagId , all '  , 
    `send_uid` int  not null  default   0  comment '发送者UID，管理员是9999，未知8888'  , 
    `send_ip` varchar(50)  not null  default   ''  comment '发送者的IP'  , 
    `send_time` int  not null  default   0  comment '发送时间'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='站内信 - 群发记录'
 ;

create table mail_log(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` int  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '内容'  , 
    `rule_id` tinyint(1)  not null  default   0  comment '规则ID'  , 
    `receiver` varchar(50)  not null  default   ''  comment '接收者uid'  , 
    `expire_time` int  not null  default   0  comment '失效时间'  , 
    `auth_code` varchar(50)  not null  default   ''  comment '验证码'  , 
    `auth_status` tinyint(1)  not null  default   0  comment '1未使用2已使用3已超时'  , 
    `send_uid` int  not null  default   0  comment '发送者UID，管理员是9999，未知8888'  , 
    `send_ip` varchar(50)  not null  default   ''  comment '发送者的IP'  , 
    `status` tinyint(1)  not null  default   0  comment '1成功2失败3发送中4等待发送'  , 
    `receiver_read` tinyint(1)  not null  default   0  comment '接收者已读'  , 
    `receiver_del` tinyint(1)  not null  default   0  comment '接收者已删除'  , 
    `send_del` tinyint(1)  not null  default   0  comment '发送者已删除'  , 
    `mail_group_id` int  not null  default   0  comment '群发的ID'  , 
    `send_time` int  not null  default   0  comment '发送时间'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='站内信-日志'
 ;

create table project(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `name` varchar(50)  not null  default   ''  comment '名称'  , 
    `type` tinyint(1)  not null  default   0  comment '类型,1service 2frontend 3backend 4app'  , 
    `desc` varchar(255)  not null  default   ''  comment '描述信息'  , 
    `secret_key` varchar(100)  not null  default   ''  comment '密钥'  , 
    `status` tinyint(1)  not null  default   0  comment '状态1正常2关闭'  , 
    `access` varchar(255)  not null  default   ''  comment 'baseAuth 认证KEY'  , 
    `lang` tinyint(1)  not null  default   0  comment '实现语言1php2go3java4js'  , 
    `git` varchar(255)  not null  default   ''  comment 'git仓地址'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='服务/项目'
 ;

create table cicd_publish(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `service_id` int  not null  default   0  comment '服务ID'  , 
    `server_id` int  not null  default   0  comment '服务器ID'  , 
    `status` tinyint(1)  not null  default   0  comment '1待部署2待发布3已发布4发布失败'  , 
    `deploy_status` tinyint(1)  not null  default   0  comment '1部署中2失败3完成'  , 
    `service_info` varchar(255)  not null  default   ''  comment '服务信息-备份'  , 
    `server_info` varchar(255)  not null  default   ''  comment '服务器信息-备份'  , 
    `deploy_type` tinyint(1)  not null  default   0  comment '1本地部署2远程同步部署'  , 
    `code_dir` varchar(255)  not null  default   ''  comment '项目代码目录名'  , 
    `log` text  null  comment '日志'  , 
    `err_info` varchar(255)  not null  default   ''  comment '错误日志'  , 
    `exec_time` int  not null  default   0  comment '执行时间'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='cicd发布记录'
 ;

create table server(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `name` varchar(50)  not null  default   ''  comment '名称'  , 
    `platform` int  not null  default   0  comment '平台类型1自有2阿里3腾讯4华为'  , 
    `out_ip` varchar(15)  not null  default   ''  comment '外网IP'  , 
    `inner_ip` varchar(15)  not null  default   ''  comment '内网IP'  , 
    `env` int  not null  default   0  comment '环境变量,1本地2开发3测试4预发布5线上'  , 
    `status` tinyint(1)  not null  default   0  comment '状态1正常2关闭'  , 
    `ext` varchar(255)  not null  default   ''  comment '自定义配置信息'  , 
    `charge_user_name` varchar(50)  not null  default   ''  comment '负责人姓名'  , 
    `start_time` int  not null  default   0  comment '开始时间'  , 
    `end_time` int  not null  default   0  comment '结束时间'  , 
    `price` int  not null  default   0  comment '价格'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='服务器'
 ;

create table sms_rule(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` tinyint(1)  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '模板内容,可变量替换'  , 
    `type` tinyint(1)  not null  default   0  comment '分类,1验证码2通知3营销'  , 
    `day_times` int  not null  default   0  comment '一天最多发送次数'  , 
    `period` int  not null  default   0  comment '周期时间-秒'  , 
    `period_times` int  not null  default   0  comment '周期时间内-发送次数'  , 
    `expire_time` int  not null  default   0  comment '验证码要有失效时间'  , 
    `memo` varchar(255)  not null  default   ''  comment '描述，主要是给3方审核用'  , 
    `purpose` tinyint(1)  not null  default   0  comment '用途,参考代码常量'  , 
    `channel` tinyint(1)  not null  default   0  comment '1阿里2腾讯'  , 
    `third_back_info` varchar(255)  not null  default   ''  comment '请示3方返回结果集'  , 
    `third_template_id` varchar(100)  not null  default   ''  comment '3方模板ID'  , 
    `third_status` tinyint(1)  not null  default   0  comment '3方状态'  , 
    `third_reason` varchar(255)  not null  default   ''  comment '3方模板审核失败，理由信息'  , 
    `third_callback_info` varchar(255)  not null  default   ''  comment '3方回执-信息'  , 
    `third_callback_time` varchar(255)  not null  default   ''  comment '3方回执-时间'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='短信发送规则配置'
 ;

create table sms_log(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` int  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '内容'  , 
    `rule_id` tinyint(1)  not null  default   0  comment '规则ID'  , 
    `receiver` varchar(50)  not null  default   ''  comment '接收者邮件地址'  , 
    `expire_time` int  not null  default   0  comment '失效时间'  , 
    `auth_code` varchar(50)  not null  default   ''  comment '验证码'  , 
    `auth_status` tinyint(1)  not null  default   0  comment '1未使用2已使用3已超时'  , 
    `send_uid` int  not null  default   0  comment '发送者UID，管理员是9999，未知8888'  , 
    `send_ip` varchar(50)  not null  default   ''  comment '发送者的IP'  , 
    `status` tinyint(1)  not null  default   0  comment '1成功2失败3发送中4等待发送'  , 
    `out_no` varchar(50)  not null  default   ''  comment '3方ID'  , 
    `channel` tinyint(1)  not null  default   0  comment '1阿里2腾讯'  , 
    `third_back_info` varchar(255)  not null  default   ''  comment '请示3方返回结果集'  , 
    `third_callback_status` tinyint(1)  not null  default   0  comment '3方状态'  , 
    `third_callback_info` varchar(255)  not null  default   ''  comment '3方回执-信息'  , 
    `third_callback_time` varchar(255)  not null  default   ''  comment '3方回执-时间'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='短信发送日志'
 ;

create table email_rule(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` tinyint(1)  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '模板内容,可变量替换'  , 
    `type` tinyint(1)  not null  default   0  comment '分类,1验证码2通知3营销'  , 
    `day_times` int  not null  default   0  comment '一天最多发送次数'  , 
    `period` int  not null  default   0  comment '周期时间-秒'  , 
    `period_times` int  not null  default   0  comment '周期时间内-发送次数'  , 
    `expire_time` int  not null  default   0  comment '验证码要有失效时间'  , 
    `memo` varchar(255)  not null  default   ''  comment '描述，主要是给3方审核用'  , 
    `purpose` tinyint(1)  not null  default   0  comment '用途,参考代码常量'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='邮件发送规则配置'
 ;

create table email_log(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` int  not null  default   0  comment '项目ID'  , 
    `title` varchar(50)  not null  default   ''  comment '标题'  , 
    `content` varchar(255)  not null  default   ''  comment '内容'  , 
    `rule_id` tinyint(1)  not null  default   0  comment '规则ID'  , 
    `receiver` varchar(50)  not null  default   ''  comment '接收者邮件地址'  , 
    `expire_time` int  not null  default   0  comment '失效时间'  , 
    `auth_code` varchar(50)  not null  default   ''  comment '验证码'  , 
    `auth_status` tinyint(1)  not null  default   0  comment '验证码状态1未使用2已使用3已超时'  , 
    `send_uid` int  not null  default   0  comment '发送者UID，管理员是9999，未知8888'  , 
    `send_ip` varchar(50)  not null  default   ''  comment '发送者的IP'  , 
    `status` tinyint(1)  not null  default   0  comment '状态1成功2失败3发送中4等待发送'  , 
    `carbon_copy` varchar(255)  not null  default   ''  comment '抄送邮件地址'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='邮件发送规则配置'
 ;

create table user(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `uuid` varchar(50)  not null  default   ''  comment 'UID字条串化'  , 
    `project_id` tinyint(1)  not null  default   0  comment '项目ID'  , 
    `sex` tinyint(1)  not null  default   0  comment '性别1男2女'  , 
    `birthday` int  not null  default   0  comment '出生日期,unix时间戳'  , 
    `username` varchar(50)  not null  default   ''  comment '用户登录名'  , 
    `password` varchar(50)  not null  default   ''  comment '用户登录密码'  , 
    `pay_ps` varchar(50)  not null  default   ''  comment '用户支付密码'  , 
    `nick_name` varchar(50)  not null  default   ''  comment '用户昵称'  , 
    `mobile` varchar(50)  not null  default   ''  comment '手机号'  , 
    `email` varchar(50)  not null  default   ''  comment '邮箱'  , 
    `robot` tinyint(1)  not null  default   0  comment '机器人'  , 
    `status` tinyint(1)  not null  default   0  comment '状态1正常2禁用'  , 
    `guest` tinyint(1)  not null  default   0  comment '是否游客,1是2否'  , 
    `test` tinyint(1)  not null  default   0  comment '是否测试,1是2否'  , 
    `recommend` varchar(50)  not null  default   ''  comment '推荐人'  , 
    `header_img` varchar(50)  not null  default   ''  comment '头像url地址'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 Unique  (`uuid`)  , index  (`uuid`) , index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='用户表'
 ;

create table user_reg(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `project_id` int  not null  default   0  comment 'project_id'  , 
    `source_type` tinyint(1)  not null  default   0  comment '来源类型'  , 
    `uid` int  not null  default   0  comment 'uid'  , 
    `type` tinyint(1)  not null  default   0  comment '类型 1email2name3mobile3third4guest'  , 
    `third_type` varchar(50)  not null  default   ''  comment '三方平台类型,参数常量USER_TYPE_THIRD'  , 
    `channel` tinyint(1)  not null  default   0  comment '推广渠道1平台自己'  , 
    `ip` varchar(50)  not null  default   ''  comment '请求方传输IP'  , 
    `auto_ip` varchar(50)  not null  default   ''  comment '程序自己计算的IP'  , 
    `province` int  not null  default   0  comment 'project_id'  , 
    `city` int  not null  default   0  comment 'project_id'  , 
    `county` int  not null  default   0  comment 'project_id'  , 
    `town` int  not null  default   0  comment 'project_id'  , 
    `area_detail` varchar(255)  not null  comment '页面来源'  , 
    `app_version` varchar(50)  not null  default   ''  comment 'APP版本'  , 
    `os` tinyint(1)  not null  default   0  comment '操作系统'  , 
    `os_version` varchar(50)  not null  default   ''  comment '操作系统版本'  , 
    `device` varchar(50)  not null  default   ''  comment '设备名称'  , 
    `device_version` varchar(50)  not null  default   ''  comment '设备版本'  , 
    `lat` varchar(50)  not null  default   ''  comment '伟度'  , 
    `lon` varchar(50)  not null  default   ''  comment '经度'  , 
    `device_id` varchar(50)  not null  default   ''  comment '设备ID'  , 
    `dpi` varchar(50)  not null  default   ''  comment '分辨率'  , 
    `referer` varchar(255)  not null  default   ''  comment '页面来源'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='用户注册信息'
 ;

create table operation_record(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `ip` varchar(50)  not null  default   ''  comment 'ip'  , 
    `method` varchar(50)  not null  default   ''  comment 'get|post|put|delete'  , 
    `path` varchar(50)  not null  default   ''  comment 'uri请求路径'  , 
    `status` int  not null  default   0  comment '请求状态'  , 
    `latency` int  not null  default   0  comment '延迟'  , 
    `agent` text  null  comment 'useragent'  , 
    `error_message` varchar(255)  not null  default   ''  comment '错误信息'  , 
    `body` text  null  comment '请求内容'  , 
    `resp` text  null  comment '返回结果'  , 
    `uid` int  not null  default   0  comment '用户Id'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='请求日志'
 ;

create table statistics_log(
    `id` int  UNSIGNED  not null  auto_increment  comment '主键自增ID'  primary key  , 
    `header_common` text  null  comment 'http公共请求头信息'  , 
    `header_base` varchar(255)  not null  default   ''  comment 'http请求头客户端基础信息'  , 
    `project_id` int  not null  default   0  comment '项目ID'  , 
    `category` int  not null  default   0  comment '分类，暂未使用'  , 
    `action` varchar(255)  not null  default   ''  comment '动作标识'  , 
    `uid` int  not null  default   0  comment '用户ID'  , 
    `msg` varchar(255)  not null  default   ''  comment '自定义消息体'  , 
    `created_at` bigint  not null  default   0  comment '创建时间'  , 
    `updated_at` bigint  not null  default   0  comment '最后更新时间'  , 
    `deleted_at` bigint  default   null  comment '是否删除'  , 
 index  (`deleted_at`) )
 engine=innodb charset=utf8 comment='接收前端推送的统计日志'
 ;
