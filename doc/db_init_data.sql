INSERT INTO `project` (`id`, `name`, `type`, `desc`, `secret_key`, `status`, `access`, `lang`, `git`, `created_at`, `updated_at`, `deleted_at`) VALUES
        (1,  'GameMatch',       1, '游戏-玩家匹配机制',  'ckGamematch',       0, 'imgamematch',         2, 'git://github.com/mqzhifu/gamematch.git', 1650001049, 1650001049, NULL),
        (2,  'FrameSync',       1, '游戏-帧同步'     ,  'ckFramesync',       0, 'imframesync',         2, 'git://github.com/mqzhifu/frame_sync.git', 1650001049, 1650001049, NULL),
        (6,  'Zgoframe',        1, 'go-框架测试'    ,   'ckZgoframe',       1, 'imzgoframe',           2, 'https://github.com/mqzhifu/zgoframe.git', 1650001049, 1650001049, NULL),
        (9,  'Gateway',         1, '公共网关',          'ckgateway',        0, 'imgateway',            2, 'git://github.com/mqzhifu/gateway.git', 1650001049, 1650001049, NULL),
        (10, 'Zwebuigo',        1, '后台管理系统-API',   'ckZwebuigo',        1, 'imzwebuigo',          2, 'https://github.com/mqzhifu/zwebuigo.git', 1650001049, 1650001049, NULL),
        (11, 'Zwebuivue',       2, '后台管理系统-VUE',   'ckZwebuivue',        1, 'imzwebuivue',        4, 'https://github.com/mqzhifu/zwebuivue.git', 1650001049, 1650001049, NULL),
        (12, 'TwinAgora',       2, '数据孪生-专家端',    'ckTwinAgora',        1, 'imtwinagora',         4, 'https://github.com/mqzhifu/twin_agora.git', 1650001049, 1650001049, NULL),
        (13, 'AgoraUnity',      5, '数据孪生-UNITY端',  'ckAgoraUnity',       1, 'imagoraunity',        7, 'http://192.168.1.22:40080/jiaxing.zhu/Agora.git', 1650001049, 1650001049, NULL),
        (14, 'AR120',           5, '120-眼镜端',       'ckAR120',           1, 'imar120',              7, '',                                                 1650001049, 1650001049, NULL),
        (15, 'WEB120',          5, '120-WEB端',        'ckWEB120',          1, 'imweb120',            7, '',                                                 1650001049, 1650001049, NULL),
        (16, 'Platform_console',5, '平台(开发)小助手',   'ckPlatform_console', 1, 'imaPlatform_console', 7, '',                                                1650001049, 1650001049, NULL);



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
      (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2c1', 6, 2, 0, 'frame_sync_4', 'e10adc3949ba59abbe56e057f20f883e', '', 'sync_4', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL),
       (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2y1', 6, 2, 0, 'Platform_console', 'e10adc3949ba59abbe56e057f20f883e', '', 'Platform_console', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL),
      (null, '4d69dee4-38f3-47ed-8dee-c4792df2e2x1', 6, 2, 0, 'Zgoframe', 'e10adc3949ba59abbe56e057f20f883e', '', 'Zgoframe', '', '', 1, 1, 2, 2, '', '', 1658995531, 1658995531, NULL);


INSERT INTO `pay_category` (`id`, `name`, `sn`, `status`, `sort`, `remark`, `icon`,  `created_at`, `updated_at`, `deleted_at`, `min_amt`, `max_amt`) VALUES
 (1, '数字人民币', 'rmb', 1, 3, 'test', '', 0, 0, NULL, 1,2),
 (2, '数字人民币', 'rmb', 1, 3, 'test', '', 0, 0,  NULL, 1,2),
(3, '支付宝', 'alipay', 1, 3, 'test', '', 0, 0, NULL, 1,2),
(4, '银行卡', 'bankcard', 1, 3, 'test', '', 0, 0,  NULL, 1,2),
(5, '微信', 'wx', 1, 3, 'test', '', 0, 0,  NULL, 1,2),
(6, 'USDT', 'usdt', 1, 3, 'test', '', 0, 0,  NULL, 1,2),
(7, '人工充值', 'manual', 1, 3, 'test', '', 0,  0, NULL, 1,2),
(8, '在线充值', 'daifu', 1, 3, 'test', '', 0,  0, NULL, 1,2);


INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '中国工商银行', 'ICBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, '中国银行', 'BOC', NULL, 1, NULL, 1709620140, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, '中国交通银行测试', 'BCM', NULL, 1, NULL, 1724681226, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, '北京银行', 'BOB', NULL, 1, NULL, 1708519226, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, '中国建设银行', 'BBD', NULL, 1, NULL, 1722946155, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, '中国农业银行', 'ABCC', NULL, 1, NULL, 1725527484, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, '广发银行', 'CGB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, '中国民生银行', 'CMBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, '中国招商银行', 'CMB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, '兴业银行', 'CIB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, '浦发银行', 'SPD', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, '平安银行', 'PAB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (13, '中信银行', 'CCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, '邮政银行', 'PSBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, '中国光大银行', 'CEB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, '华夏银行', 'HXB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, '上海银行', 'BOS', NULL, 1, NULL, 1704198786, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (18, '广州银行', 'GCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (19, '平顶山银行', 'BOP', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (20, '东营市商业银行', 'DYCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (21, '承德银行', 'BOCD', NULL, 1, NULL, 1709972477, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (22, '北京银行', 'BJBANK', NULL, 1, NULL, 1722339489, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (23, '营口银行', 'YKCB', NULL, 1, NULL, 1703656957, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (24, '东亚银行', 'HKBEA', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (25, '长沙银行', 'CSCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (26, '广东南粤银行', 'NYNB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (27, '青岛银行', 'QDCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (28, '西安银行', 'XABANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (29, '江苏银行', 'JSBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (30, '湖南省农村信用社', 'HNRCC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (31, '兰州银行', 'LZYH', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (32, '邢台银行', 'XTB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (33, '上海银行', 'SHBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (34, '甘肃省农村信用', 'GSRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (35, '晋商银行', 'JSB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (36, '北京农村商业银行', 'BJRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (37, '衡水银行', 'HSBK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (38, '南京银行', 'NJCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (39, '桂林银行', 'GLBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (40, '温州银行', 'WZCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (41, '贵阳市商业银行', 'GYCB', NULL, 1, NULL, 1703656360, 1703656360);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (42, '三门峡银行', 'SMXB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (43, '云南省农村信用社', 'YNRCC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (44, '宁夏银行', 'NXBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (45, '徽商银行', 'HSBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (46, '贵州省农村信用社', 'GZRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (47, '丹东银行', 'BODD', NULL, 1, NULL, 1721125997, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (48, '济宁银行', 'JNBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (49, '湖北银行宜昌分行', 'HBYCBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (50, '广东省农村信用社联合社', 'GDRCC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (51, '上饶银行', 'SRBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (52, '莱商银行', 'LSBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (53, '泰安市商业银行', 'TACCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (54, '常熟农村商业银行', 'CSRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (55, '乌鲁木齐市商业银行', 'URMQCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (56, '杭州银行', 'HZCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (57, '华融湘江银行', 'HRXJB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (58, '齐鲁银行', 'QLBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (59, '洛阳银行', 'LYBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (60, '武汉农村商业银行', 'WHRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (61, '吉林农信', 'JLRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (62, '常州农村信用联社', 'CZRCB', NULL, 1, NULL, 1703827521, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (63, '石嘴山银行', 'SZSBK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (64, '中原银行', 'XCYH', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (65, '农信银清算中心', 'NHQS', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (66, '大连银行', 'DLB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (67, '南充市商业银行', 'CGNB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (68, '乐山市商业银行', 'LSCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (69, '城市商业银行资金清算中心', 'CBBQS', NULL, 1, NULL, 1708582556, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (70, '广西省农村信用', 'GXRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (71, '宁夏黄河农村商业银行', 'NXRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (72, '国家开发银行', 'CDB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (73, '江苏太仓农村商业银行', 'TCRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (74, '潍坊银行', 'BANKWF', NULL, 1, NULL, 1724830840, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (75, '华夏银行', 'HXBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (76, '绍兴银行', 'SXCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (77, '汉口银行', 'HKB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (78, '宁波银行', 'NBBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (79, '河南省农村信用', 'HNRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (80, '成都银行', 'CDCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (81, '山东农信', 'SDRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (82, '江苏省农村信用联合社', 'JSRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (83, '天津银行', 'TCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (84, '赣州银行', 'GZB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (85, '重庆银行', 'CQBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (86, '龙江银行', 'DAQINGB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (87, '福建海峡银行', 'FJHXBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (88, '广发银行', 'GDB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (89, '重庆农村商业银行', 'CRCBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (90, '重庆三峡银行', 'CCQTGB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (91, '东莞农村商业银行', 'DRCBCL', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (92, '陕西信合', 'SXRCCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (93, '南昌银行', 'NCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (94, '吴江农商银行', 'WJRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (95, '恒丰银行', 'EGBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (96, '浙江民泰商业银行', 'MTBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (97, '吉林银行', 'JLBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (98, '浙江稠州商业银行', 'CZCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (99, '中信银行', 'CITIC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (100, '金华银行', 'JHBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (101, '湖北银行', 'HBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (102, '昆山农村商业银行', 'KSRB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (103, '三门峡银行', 'SCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (104, '郑州银行', 'ZZBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (105, '昆仑银行', 'KLB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (106, '晋城银行', 'JINCHB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (107, '阳泉银行', 'YQCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (108, '辽阳市商业银行', 'LYCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (109, '九江银行', 'JJBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (110, '渤海银行', 'BOHAIB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (111, '韩亚银行', 'HANABANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (112, '浦发银行', 'SPDB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (113, '鞍山银行', 'ASCB', NULL, 1, NULL, 1712656514, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (114, '内蒙古银行', 'H3CB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (115, '天津农商银行', 'TRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (116, '张家港农村商业银行', 'ZRCBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (117, '遵义市商业银行', 'ZYCBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (118, '张家口市商业银行', 'ZJKCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (119, '苏州银行', 'BOSZ', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (120, '嘉兴银行', 'JXBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (121, '青海银行', 'BOQH', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (122, '浙江省农村信用社联合社', 'ZJNX', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (123, '安徽省农村信用社', 'ARCU', NULL, 1, NULL, 1722508458, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (124, '晋中银行', 'JZBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (125, '河北银行', 'BHB', NULL, 1, NULL, 1714994740, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (126, '富滇银行', 'FDB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (127, '廊坊银行', 'LANGFB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (128, '成都农商银行', 'CDRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (129, '开封市商业银行', 'CBKF', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (130, '河北省农村信用社', 'HBRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (131, '四川省农村信用', 'SCRCU', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (132, '盛京银行', 'SJBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (133, '自贡市商业银行', 'ZGCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (134, '临商银行', 'LSBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (135, '上海农村商业银行', 'SHRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (136, '浙江泰隆商业银行', 'ZJTLCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (137, '交通银行', 'COMM', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (138, '信阳银行', 'XYBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (139, '汇丰银行', 'HSBC', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (140, '宜宾市商业银行', 'YBCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (141, '阜新银行', 'FXCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (142, '台州银行', 'TZCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (143, '玉溪市商业银行', 'YXCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (144, '尧都农商行', 'YDRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (145, '包商银行', 'BSB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (146, '齐商银行', 'ZBCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (147, '德州银行', 'DZBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (148, '江苏江阴农村商业银行', 'JRCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (149, '德阳商业银行', 'DYCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (150, '中山小榄村镇银行', 'XLBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (151, '广州农商银行', 'GRCB', NULL, 1, NULL, 1705159006, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (152, '抚顺银行', 'FSCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (153, '鄂尔多斯银行', 'ORBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (154, '湖州市商业银行', 'HZCCB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (155, '锦州银行', 'BOJZ', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (156, '安阳银行', 'AYCB', NULL, 1, NULL, 1724075433, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (157, '南海农商银行', 'NHB', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (158, '新乡银行', 'XXBANK', NULL, 1, NULL, NULL, NULL);
INSERT INTO `banks` (`id`, `name`, `code`, `address`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (159, '鄞州银行', 'NBYZ', NULL, 1, NULL, 1703656310, 1703656310);

