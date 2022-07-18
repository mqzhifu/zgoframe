INSERT INTO `project` (`id`, `name`, `secret_key`, `status`, `created_at`, `updated_at`, `deleted_at`, `type`, `desc`, `git`, `access`) VALUES
(1, 'GameMatch', 'aaaaaa', 1, 0, 0, NULL, 1, '游戏匹配', 'git://github.com/mqzhifu/gamematch.git', 'imgamematch'),
(2, 'FrameSync', 'bbbb', 1, 0, 0, NULL, 1, '游戏帧同步', 'git://github.com/mqzhifu/frame_sync.git', 'imframesync'),
(3, 'LogSlave', 'ccccc', 1, 0, 0, NULL, 1, '日志接收', 'git://github.com/mqzhifu/log_slave.git', 'imlogsalve'),
(6, 'Zgoframe', 'asdf  sdf', 1, 0, 0, NULL, 1, 'go框架测试', 'git://github.com/mqzhifu/zgoframe.git', 'imzgoframe'),
(9, 'Gateway', 'adfasdf', 1, 0, 0, NULL, 1, '', 'git://github.com/mqzhifu/gateway.git', 'imgateway');



INSERT INTO `project` (`id`, `name`, `type`, `desc`, `secret_key`, `status`, `access`, `lang`, `git`, `created_at`, `updated_at`, `deleted_at`) VALUES
                                                                                                                                                    (6, 'Zgoframe', 1, 'go框架测试', 'ckZgoframe', 1, 'imzgoframe', 2, 'https://github.com/mqzhifu/zgoframe.git', 1650001049, 0, NULL),
                                                                                                                                                    (10, 'Zwebuigo', 1, '后台管理系统', 'ckZwebuigo', 1, 'imzwebuigo', 2, 'https://github.com/mqzhifu/zwebuigo.git', 1650001049, 0, NULL),
                                                                                                                                                    (11, 'Zwebuivue', 2, '后台管理系统-VUE', 'ckZwebuivue', 1, 'imzwebuivue', 4, 'https://github.com/mqzhifu/zwebuivue.git', 1650001049, 0, NULL),
                                                                                                                                                    (12, 'TwinAgora', 2, '数据孪生-专家指导(声网)', 'ckTwinAgora', 1, 'imtwinagora', 4, 'https://github.com/mqzhifu/twin_agora.git', 1650001049, 0, NULL),
                                                                                                                                                    (13, 'AgoraUnity', 5, '数据孪生-UNITY端', 'ckAgoraUnity', 1, 'imagoraunity', 7, 'http://192.168.1.22:40080/jiaxing.zhu/Agora.git', 1650001049, 0, NULL);