INSERT INTO `project` (`id`, `name`, `type`, `desc`, `secret_key`, `status`, `access`, `lang`, `git`, `created_at`, `updated_at`, `deleted_at`) VALUES
        (1, 'GameMatch', 1, '小游戏-玩家匹配机制', 'ckgamematch', 0, 'imgamematch', 2, 'git://github.com/mqzhifu/gamematch.git', 1650001049, 1650001049, NULL),
        (2, 'FrameSync', 1, '游戏-帧同步', 'ckframesync', 0, 'imframesync', 2, 'git://github.com/mqzhifu/frame_sync.git', 1650001049, 1650001049, NULL),
        (6, 'Zgoframe', 1, 'go框架测试', 'ckZgoframe', 1, 'imzgoframe', 2, 'https://github.com/mqzhifu/zgoframe.git', 1650001049, 1650001049, NULL),
        (9, 'Gateway', 1, '公共网关', 'ckgateway', 0, 'imgateway', 2, 'git://github.com/mqzhifu/gateway.git', 1650001049, 1650001049, NULL),
        (10, 'Zwebuigo', 1, '后台管理系统', 'ckZwebuigo', 1, 'imzwebuigo', 2, 'https://github.com/mqzhifu/zwebuigo.git', 1650001049, 1650001049, NULL),
        (11, 'Zwebuivue', 2, '后台管理系统-VUE', 'ckZwebuivue', 1, 'imzwebuivue', 4, 'https://github.com/mqzhifu/zwebuivue.git', 1650001049, 1650001049, NULL),
        (12, 'TwinAgora', 2, '数据孪生-专家指导(声网)', 'ckTwinAgora', 1, 'imtwinagora', 4, 'https://github.com/mqzhifu/twin_agora.git', 1650001049, 1650001049, NULL),
        (13, 'AgoraUnity', 5, '数据孪生-UNITY端', 'ckAgoraUnity', 1, 'imagoraunity', 7, 'http://192.168.1.22:40080/jiaxing.zhu/Agora.git', 1650001049, 1650001049, NULL);



INSERT INTO `server` (`id`, `name`, `platform`, `out_ip`, `inner_ip`, `env`, `status`, `ext`, `charge_user_name`, `start_time`, `end_time`, `price`, `created_at`, `updated_at`, `deleted_at`) VALUES
       (1, '本地', 1, '127.0.0.1', '127.0.0.1', 1, 1, '', '小z', 1650006845, 1650006845, 100, 1650006845, 0, NULL),
       (2, '开发', 1, '192.168.1.21', '192.168.1.21', 2, 1, '', '小z', 1650006845, 1650006845, 100, 1650006845, 0, NULL),
       (3, '测试', 1, '2.2.2.2', '127.0.0.1', 3, 2, '', '小z', 1650006845, 1650006845, 100, 1650006845, 0, NULL),
       (4, '预发布', 1, '8.142.177.235', '172.27.198.210', 4, 1, '', '小z', 1650006845, 1650006845, 100, 1650006845, 0, NULL),
       (5, '线上', 1, '8.142.161.156', '172.27.218.143', 5, 1, '', '小z', 1650006845, 1650006845, 100, 1650006845, 0, NULL);


INSERT INTO `sms_rule` (`id`, `project_id`, `title`, `content`, `type`, `day_times`, `period`, `period_times`, `expire_time`, `memo`, `channel`, `third_back_info`, `third_template_id`, `third_status`, `third_reason`, `third_callback_info`, `third_callback_time`, `created_at`, `updated_at`, `deleted_at`) VALUES
        (1, 6, '短信注册', '{nickname},您好：欢迎注册本网站，验证码为：{auth_code},{auth_expire_time}秒后将失效，勿告诉他人，防止被骗', 1, 10, 60, 1, 300, '0', 1, '', '1', 0, '', '', '1', 0, 0, null),
        (2, 6, '短信登陆', '{nickname},您好：登陆验证码为：{auth_code},{auth_expire_time} 秒后将失效，勿告诉他人，防止被骗。', 1, 10, 60, 1, 300, '0', 1, '', '1', 0, '', '', '1', 0, 0, null),
        (3, 6, '找加密码', '找回密码', 1, 10, 60, 1, 300, '0', 1, '', '', '1', '', '', '', '1', 0, null),
        (4, 6, '报警', '报警，程序出错。级别：{level}，项目ID:{project_id}，内容：{content}', 2, 10, 60, 1, 0, '0', 1, '', 'SMS_273495087', '1',  '', '', '1', 300, 0, null);

INSERT INTO `email_rule` (`id`, `project_id`, `title`, `content`, `type`, `day_times`, `period`, `period_times`, `expire_time`, `memo`, `created_at`, `updated_at`, `deleted_at`) VALUES
         (1, 6, '短信注册', '{nickname},您好：欢迎注册本网站，验证码为：{auth_code},{auth_expire_time}秒后将失效，勿告诉他人，防止被骗',  1, 10, 60, 1, 300,"", 0, 0, null),
         (2, 6, '短信登陆', '{nickname},您好：登陆验证码为：{auth_code},{auth_expire_time} 秒后将失效，勿告诉他人，防止被骗。', 1, 10, 60, 1, 300, "",0, 0, null),
         (3, 6, '找加密码', '找回密码', 1, 10, 60, 1, 300, '0', 0 , 0, null),
         (4, 6, '报警', '报警，程序出错。级别：{level}，项目ID:{project_id}，内容：{content}', 2, 10, 60, 1, 300,"",0, 0, null);





INSERT INTO `user` (`id`, `uuid`, `project_id`, `sex`, `birthday`, `username`, `password`, `pay_ps`, `nick_name`, `mobile`, `email`, `robot`, `status`, `guest`, `test`, `recommend`, `header_img`, `created_at`, `updated_at`, `deleted_at`) VALUES
      (null, '2d879cfe-d900-45ae-a3e5-af3517eb8d02', 6, 1, 0, 'frame_sync_1', 'e10adc3949ba59abbe56e057f20f883e', 'e10adc3949ba59abbe56e057f20f883e', 'sync_1', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL),
      (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2c6', 6, 2, 0, 'frame_sync_2', 'e10adc3949ba59abbe56e057f20f883e', '', 'sync_2', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL),
      (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2c3', 6, 1, 0, 'frame_sync_3', 'e10adc3949ba59abbe56e057f20f883e', '', 'sync_3', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL),
      (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2c1', 6, 2, 0, 'frame_sync_4', 'e10adc3949ba59abbe56e057f20f883e', '', 'sync_4', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL);