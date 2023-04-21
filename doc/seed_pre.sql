

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
-- Dumping data for table `agora_callback_record`
--

LOCK TABLES `agora_callback_record` WRITE;
/*!40000 ALTER TABLE `agora_callback_record` DISABLE KEYS */;
INSERT INTO `agora_callback_record` VALUES (1,'5c57f76ad32f47a5a89674618f674b6c',1,101,'1660714248801','{\"channelName\":\"test_webhook\",\"platform\":0,\"reason\":0,\"ts\":1560396834,\"uid\":0}',1660719985,1660719985,NULL),(2,'',0,0,'0','{\"channelName\":\"\",\"platform\":0,\"reason\":0,\"ts\":0,\"uid\":0}',1660721687,1660721687,NULL),(3,'',0,0,'0','{\"channelName\":\"\",\"platform\":0,\"reason\":0,\"ts\":0,\"uid\":0}',1660721715,1660721715,NULL),(4,'',0,0,'0','{\"channelName\":\"\",\"platform\":0,\"reason\":0,\"ts\":0,\"uid\":0}',1660721762,1660721762,NULL),(5,'',0,0,'0','{\"channelName\":\"\",\"platform\":0,\"reason\":0,\"ts\":0,\"uid\":0}',1660722243,1660722243,NULL),(6,'111',4,23,'0','{\"channelName\":\"222\",\"platform\":44444,\"reason\":4,\"ts\":2,\"uid\":1}',1660722361,1660722361,NULL),(7,'22',777,11,'33','{\"cname\":\"44\",\"sendts\":222,\"sequence\":333,\"serviceType\":444,\"sid\":\"555\",\"uid\":\"666\",\"details\":{\"errorCode\":55,\"errorLevel\":66,\"errorMsg\":\"77\",\"module\":88,\"msgName\":\"99\",\"stat\":111}}',1660722410,1660722410,NULL);
/*!40000 ALTER TABLE `agora_callback_record` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `agora_cloud_record`
--

LOCK TABLES `agora_cloud_record` WRITE;
/*!40000 ALTER TABLE `agora_cloud_record` DISABLE KEYS */;
INSERT INTO `agora_cloud_record` VALUES (1,1,123123,'ckck','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3F7qTRllOlXJ13mKiaiMoMVby52Fp9kCWDjl6m9EG9mWHvJryhTwFDJCC2RIceqkpT4HbKGBu6v4i-A-iz2paFuy41g10d_bU7WIODP4lxOsCJMvOHoQ_sT4H0z0VXbums-l7LMPNhj5sl4Boh8Dpeu_rTyuvXKR5yr0NS2uAFYqia1CuBhnOVneFF-kC05hB05WdO1uCEO8MOSLeqEZA0S6cD6vPlQUr61tJXs1C6r1ZpBADjzAx96UcvLl2x8rDY','',1,1,0,0,'','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',0,'','','',1660891114,1660891114,NULL),(2,1,123123,'ckck','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3GnFQveQMgEGCDzjQhfa2VYeVs-C4-5_0A8PvBaT1hR5Alumya65bSLk-6tpLfqoUfnuHG69Cc_hyub0xPU-K_JRsGPiun1uIiWiSpa7WDAxwI3B5ObbZU-N9CgKYBDuPeUiqW7wXI8Ykpp9pN0fY41H3stG2cB9f2TkbX34J3kp9xAzXvHzIvX9by4pD3YLVAo9xEUGTdFs9rFBbF7PfHp_K_V_YlByBRRI_dSu63b85ym0nuU72qYVylvpyEFzPk','0e4aecef1e4c7c7257a7208a20c3a86e',2,1,1660891394,0,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck\",\"1660891394\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIABgLIopL4sAEQR9Xr9m36xehErtFJdOvBH3P00Yb56YhoQj5djCc5IwIgBkkMABfnUAYwQAAQB+dQBjAgB+dQBjAwB+dQBjBAB+dQBj\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',0,'','','',1660891117,1660891394,NULL),(3,10,11111,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3ERqcFk_PNH1uyENy7jLlOMLkev57L9XAxUMcNkjXC9BweDgYRcqhPDCmQ76gXktUYOHYFU92veg3wyM9HlVrkJtZ6mg0iOt42KjFP0RQg9ICjVXDOUfpApV-w9758GI-0LOTU5mXGFdFYlqoJjaDNFIOEYx1K-rp-nbQuBeoF5kyNCsFcZrznjFeL01srWyLeglpyh51ZEIdqXNOZyR6-IRNQ__LKU-PGbLnB4wQnzCXktkHGntEa4xrIWBVtXzuE','36c9d6bcfc40cf65665922b6d3617a23',3,1,1665730017,1666149167,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1665730017\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIAAYN1ZoLR5oj5OCZCNtaYvfZ9RQGlolJA9qMgRgwrFHVa1N0WnAcd6gIgALKy8EJDNJYwQAAQAkM0ljAgAkM0ljAwAkM0ljBAAkM0lj\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',3,'','','',1665730016,1666149167,NULL),(4,9,22222,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3EgHxIxZhsmZ4DSvRRvxZB-QMN1t3PSOk1Rc9v1qo6a3hzaSl087zN5mp6HWjMDfib5xmnVb9xC6wLxkn-lTXnoyEMhfCtwjPG-GYsCSgHEl1YlK_wH5Cg7m-N7BTyniHedhedQ_I-sedSbCic3zBVg3wrI5rc8vxoQtUyecKeFYeB49qNt8O5hRtdF7eXYJjLqe8LB6FLfdpVsW0aRL1NuWFl01Fo73I5ixik0OYpT9PK47YjIgenKZuRo3GRzrOg','8f5e9b993f4486da99685d8ceb81dc68',3,1,1666075600,1666149202,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666075600\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIABdNYtoA6ta3X/fM+Fp4FA9yXSZ3/GzzqPy9Z+6raqMga1N0WneGKlFIgBxMy8EqJdPYwQAAQCol09jAgCol09jAwCol09jBACol09j\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',3,'','','',1666075599,1666149202,NULL),(5,10,11111,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3EgHxIxZhsmZ4DSvRRvxZB-t3sg4HcIVvBBTd2JZNRzbEa3qi1G86RLFd6E6ik7CfTybkk2M99MSO7REm9BOuZUy68peQ9fU0vZvkgng6pY0-I61MI5ZnGxWy2xXpmDw9j9mL0yBj1Cyaxgk-KMGyONsOTGOuQnKLtQInUuMv_pHCc9jBPJWB9qHgL4N-RrkR7hDnGCkXqovOhVNSnuG6v97lJ_79ecPsJ1NIhYRotJQV-38VluIjBhAoWIMm_07G8','1ea163a1ea4a4cd9434f11987d6ff76c',3,1,1666149187,1666149215,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666149187\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIACKC+2YSHroDUxezYFQjxpbJ8HgIG+qui0CR2x/LC6Lx61N0WnAcd6gIgDg3jcBTKtPYwQAAQBMq09jAgBMq09jAwBMq09jBABMq09j\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',1,'{\"fileListMode\":\"json\",\"uploadingStatus\":\"uploaded\",\"fileList\":[{\"fileName\":\"agoraRecord/ckck1/1666149187/1ea163a1ea4a4cd9434f11987d6ff76c_ckck1__uid_s_22222__uid_e_av.mpd\",\"trackType\":\"audio_and_video\",\"uid\":0,\"mixedAllUser\":false,\"isPlayable\":true,\"sliceStartTime\":1666149188378}],\"command\":\"\",\"subscribeModeBitmask\":0,\"vid\":\"\",\"payload\":{\"message\":\"\"}}','','',1666149187,1666149215,NULL),(6,9,22222,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3EgHxIxZhsmZ4DSvRRvxZB-WcL-_Ma_wjvTNge835B4tS2aPlsR6Vhf-3XGxojw9rmucptOwIelRhFwjPDn4vcvH2yOo79-esC38_NtpHtuVrvYsozn6OU0ECJTpVfjK_oPcIewOnTp1fbUJ-1hHZGC-VaCwssMgmx5-LVcX4XCWWPJIcXzrp4FV_ObFY0pV5DLt81cILnDgIqC_sMqZD_5qyigRKY0f6HvzadXt_wT1Tk649s-5D7NDdEJyjhya-k','1175a297c7481acbf207af8e758275dd',3,1,1666149193,1666149202,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666149192\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIABdNYtoA6ta3X/fM+Fp4FA9yXSZ3/GzzqPy9Z+6raqMga1N0WneGKlFIgBxMy8EqJdPYwQAAQCol09jAgCol09jAwCol09jBACol09j\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',2,'{\"fileListMode\":\"json\",\"uploadingStatus\":\"uploaded\",\"fileList\":[{\"fileName\":\"agoraRecord/ckck1/1666149192/1175a297c7481acbf207af8e758275dd_ckck1__uid_s_11111__uid_e_av.mpd\",\"trackType\":\"audio_and_video\",\"uid\":0,\"mixedAllUser\":false,\"isPlayable\":true,\"sliceStartTime\":1666149193776}],\"command\":\"\",\"subscribeModeBitmask\":0,\"vid\":\"\",\"payload\":{\"message\":\"\"}}','','',1666149192,1666149202,NULL),(7,10,11111,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3EgHxIxZhsmZ4DSvRRvxZB-IJuDiK3Wr8qEm1exye8DuST8rA4MHz2gKpysNZkSXjoYnD0rPjhU4OfLhd28kKN_86R0zbol15euMipr8GwrixwhylfpyIlSS-iS5xu34dbJyh1M_oPuD8UE-A37LgjsVuqmk2QxJpe99Bv4jrIGLoqb2Uc3bTWXXimOpBT6LWAm0Kt1JO32R-4FM13L-vSjiSwH2b7s3tuU9L-gXqqBhyzRCACHEDU2YxGuI_eClXc','0b205cf4df460643bca1ddb9b4585c33',3,1,1666149218,1666149226,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666149217\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIACKC+2YSHroDUxezYFQjxpbJ8HgIG+qui0CR2x/LC6Lx61N0WnAcd6gIgDg3jcBTKtPYwQAAQBMq09jAgBMq09jAwBMq09jBABMq09j\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',3,'','','',1666149217,1666149226,NULL),(8,9,11111,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3FH1NxkDj9X8DLUeWfawBvPcT0zLjZMEdzQr8G6Q3NKrXCoTam_FxIzB62yLwQEDhU6En1XrFJb7_CY7itsjmVuXSRniQXsYOWV4Os3q8DPoqhqoQ-wsL_1sNgLSxQBKJoELzGmGmRbNOPT43RSktLMAaFay-aWURTKMmM0YTGaPQlY_RwFW2_LKFn7_yp4lPQCg7jWbUdWwB0QGdI6AtG9O9T-44s9v_hmA8gUbAHEOUj9pwmAvzM3DD2uNZp57KU','0baccfb3e64aa066c6a897bcc3fa975c',3,1,1666249123,1666249130,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666249123\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIACIttG/iKFvbEBKh9bGuevY0YEsdklF7bRrKlqZVyaoiK1N0WnAcd6gIgCAvhMAQz1SYwQAAQBDPVJjAgBDPVJjAwBDPVJjBABDPVJj\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',1,'{\"fileListMode\":\"json\",\"uploadingStatus\":\"uploaded\",\"fileList\":[{\"fileName\":\"agoraRecord/ckck1/1666249123/0baccfb3e64aa066c6a897bcc3fa975c_ckck1__uid_s_22222__uid_e_av.mpd\",\"trackType\":\"audio_and_video\",\"uid\":0,\"mixedAllUser\":false,\"isPlayable\":true,\"sliceStartTime\":1666249124518}],\"command\":\"\",\"subscribeModeBitmask\":0,\"vid\":\"\",\"payload\":{\"message\":\"\"}}','','',1666249123,1666249130,NULL),(9,9,11111,'ckck1','nUwUbQf9Zg6tsgtLslGnDg0lk8RYaUE09pqOuSIgwfyLwpL9dJfszQxuJ9vAQmkifudr1BNw5HR4RVXufzmoefS4KNzD1jqxANsIiz10b3EgHxIxZhsmZ4DSvRRvxZB-b-36P7k2bO6mC4WD8O3cHyGnCWq9d9kc_R3k2gWfy2-R0LxLmQMTESV_0g_pYAqYruwzmbrfl8UpiKM2rtNVcZIZYdL-fIlhay9U1oHzgwnnDt1pda4PQBYbfxctDlGEsvNkykDRXh_O_Up5V477s8Wh8tqZTzH5oIx2Uxf1lU0FD3ZJD2yXBQANL_nz4NIfJuV9Zx0o108oWGu2LuAPcVBkzO7yxTkzVQ3JWjq4aqw','fc009dc0094b71b34d0c5ab0df6a9099',3,1,1666249133,1666249138,'{\"recordingConfig\":{\"maxIdleTime\":300,\"streamMode\":\"standard\",\"streamTypes\":2,\"channelType\":0,\"subscribeUidGroup\":5,\"videoStreamType\":0,\"subscribeVideoUids\":[\"#allstream#\"],\"subscribeAudioUids\":[\"#allstream#\"]},\"storageConfig\":{\"accessKey\":\"LTAI5tJbjZiWQ9Xn9N2brRFD\",\"region\":3,\"bucket\":\"servicebase\",\"secretKey\":\"GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm\",\"vendor\":2,\"fileNamePrefix\":[\"agoraRecord\",\"ckck1\",\"1666249133\"]},\"token\":\"006bf7dd6c465424e36835d9aa0ff664a7aIACIttG/iKFvbEBKh9bGuevY0YEsdklF7bRrKlqZVyaoiK1N0WnAcd6gIgCAvhMAQz1SYwQAAQBDPVJjAgBDPVJjAwBDPVJjBABDPVJj\"}','{\"region\":\"CN\",\"resourceExpiredHour\":72,\"scene\":0}',2,'{\"fileListMode\":\"json\",\"uploadingStatus\":\"uploaded\",\"fileList\":[{\"fileName\":\"agoraRecord/ckck1/1666249133/fc009dc0094b71b34d0c5ab0df6a9099_ckck1__uid_s_22222__uid_e_av.mpd\",\"trackType\":\"audio_and_video\",\"uid\":0,\"mixedAllUser\":false,\"isPlayable\":true,\"sliceStartTime\":1666249133938}],\"command\":\"\",\"subscribeModeBitmask\":0,\"vid\":\"\",\"payload\":{\"message\":\"\"}}','','',1666249132,1666249138,NULL);
/*!40000 ALTER TABLE `agora_cloud_record` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 COMMENT='邮件发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_log`
--

LOCK TABLES `email_log` WRITE;
/*!40000 ALTER TABLE `email_log` DISABLE KEYS */;
INSERT INTO `email_log` VALUES (1,6,'报警','报警，程序出错。级别：warning，项目ID:6，内容：商品库存不足，请及时补充货源',4,'mqzhifu@qq.com',0,'',0,9999,'192.168.22.173',0,'',1678960035,1678960035,NULL),(2,6,'报警','报警，程序出错。级别：warning，项目ID:6，内容：商品库存不足，请及时补充货源',4,'mqzhifu@sina.com',0,'',0,9999,'192.168.22.173',0,'',1678960103,1678960103,NULL);
/*!40000 ALTER TABLE `email_log` ENABLE KEYS */;
UNLOCK TABLES;

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
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COMMENT='邮件发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_rule`
--

LOCK TABLES `email_rule` WRITE;
/*!40000 ALTER TABLE `email_rule` DISABLE KEYS */;
INSERT INTO `email_rule` VALUES (1,6,'短信注册','{nickname},您好：欢迎注册本网站，验证码为：{auth_code},{auth_expire_time}秒后将失效，勿告诉他人，防止被骗',1,10,60,1,300,'',0,0,NULL),(2,6,'短信登陆','{nickname},您好：登陆验证码为：{auth_code},{auth_expire_time} 秒后将失效，勿告诉他人，防止被骗。',1,10,60,1,300,'',0,0,NULL),(3,6,'找加密码','找回密码',1,10,60,1,300,'0',0,0,NULL),(4,6,'报警','报警，程序出错。级别：{level}，项目ID:{project_id}，内容：{content}',2,10,60,1,300,'',0,0,NULL);
/*!40000 ALTER TABLE `email_rule` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=488 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-小组信息';
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
) ENGINE=InnoDB AUTO_INCREMENT=164 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-推送消息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_match_push`
--

LOCK TABLES `game_match_push` WRITE;
/*!40000 ALTER TABLE `game_match_push` DISABLE KEYS */;
INSERT INTO `game_match_push` VALUES (1,1,1,1667375885,2222,1,0,1,'2222%2%87.000%0%1667375885%0%1667375875%0%20,21%diiiiii%0%2222%',1667375885,1667375885,NULL),(2,1,2,1667376264,1,1,0,2,'1%1%1667376264%1667376324%1%20,21,10,11%2222,1111%0%',1667376264,1667376264,NULL),(3,1,1,1667529696,1111,1,0,1,'1111%2%86.000%0%1667381807%0%1667381797%0%10,11%diiiiii%0%1111%',1667529696,1667529696,NULL),(4,1,2,1667540734,1111,1,0,1,'1111%2%86.000%0%1667540734%0%1667540724%0%10,11%diiiiii%0%1111%',1667540734,1667540734,NULL),(5,1,3,1667542880,1111,1,0,1,'1111%2%86.000%0%1667542880%0%1667542870%0%10,11%diiiiii%0%1111%',1667542880,1667542880,NULL),(6,2,1,1668666855,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666855,1668666855,NULL),(7,2,2,1668666855,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666855,1668666855,NULL),(8,2,3,1668666855,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666855,1668666855,NULL),(9,2,4,1668666856,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666856,1668666856,NULL),(10,2,5,1668666856,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666856,1668666856,NULL),(11,2,6,1668666857,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666857,1668666857,NULL),(12,2,7,1668666857,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666857,1668666857,NULL),(13,2,8,1668666857,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666857,1668666857,NULL),(14,2,9,1668666858,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666858,1668666858,NULL),(15,2,10,1668666858,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666858,1668666858,NULL),(16,2,11,1668666858,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666858,1668666858,NULL),(17,2,12,1668666859,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666859,1668666859,NULL),(18,2,13,1668666859,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666859,1668666859,NULL),(19,2,14,1668666859,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666859,1668666859,NULL),(20,2,15,1668666860,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666860,1668666860,NULL),(21,2,16,1668666860,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666860,1668666860,NULL),(22,2,17,1668666861,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666861,1668666861,NULL),(23,2,18,1668666861,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666861,1668666861,NULL),(24,2,19,1668666861,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666861,1668666861,NULL),(25,2,20,1668666862,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666862,1668666862,NULL),(26,2,21,1668666862,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666862,1668666862,NULL),(27,2,22,1668666862,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666862,1668666862,NULL),(28,2,23,1668666863,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666863,1668666863,NULL),(29,2,24,1668666863,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666863,1668666863,NULL),(30,2,25,1668666863,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666863,1668666863,NULL),(31,2,26,1668666864,0,1,0,1,'0%1%0.000%0%1668666855%0%1668666845%0%1111%giveMeFive%0%0%',1668666864,1668666864,NULL),(32,2,1,1668668520,2,1,0,1,'2%1%0.000%0%1668668520%0%1668668510%0%1111%givemefive%0%2%',1668668520,1668668520,NULL),(33,2,2,1668671905,1,1,0,2,'1%2%1668671904%1668671964%1%1111,2222%4,3%0%',1668671905,1668671905,NULL),(34,2,3,1668680124,2,1,0,2,'2%2%1668680124%1668680184%1%1133,1122%2,1%0%',1668680124,1668680124,NULL),(35,2,4,1668680950,3,1,0,2,'3%2%1668680949%1668681009%1%1133,1122%4,3%0%',1668680950,1668680950,NULL),(36,2,3,1668681204,2,1,0,2,'2%2%1668681204%1668681264%1%1133,1122%4,3%0%',1668681204,1668681204,NULL),(37,2,1,1669810936,1,1,0,1,'1%1%0.000%0%1669810936%0%1669810926%0%1%html_test_frame_sync%0%1%',1669810936,1669810936,NULL),(38,2,2,1669866600,1,1,0,2,'1%2%1669866600%1669866660%1%2,1%3,2%0%',1669866600,1669866600,NULL),(39,2,3,1669880386,4,1,0,1,'4%1%0.000%0%1669880386%0%1669880376%0%1%html_test_frame_sync%0%4%',1669880386,1669880386,NULL),(40,2,4,1669889623,2,1,0,2,'2%2%1669889623%1669889683%1%2,1%6,5%0%',1669889623,1669889623,NULL),(41,2,1,1669889881,1,1,0,2,'1%2%1669889881%1669889941%1%2,1%2,1%0%',1669889881,1669889881,NULL),(42,2,2,1669890412,2,1,0,2,'2%2%1669890412%1669890472%1%2,1%4,3%0%',1669890412,1669890412,NULL),(43,2,3,1669890927,3,1,0,2,'3%2%1669890926%1669890986%1%1,2%5,6%0%',1669890927,1669890927,NULL),(44,2,4,1669893470,4,1,0,2,'4%2%1669893470%1669893530%1%2,1%8,7%0%',1669893470,1669893470,NULL),(45,2,5,1669894729,9,1,0,1,'9%1%0.000%0%1669894729%0%1669894719%0%1%html_test_frame_sync%0%9%',1669894729,1669894729,NULL),(46,2,6,1669895305,10,1,0,1,'10%1%0.000%0%1669895305%0%1669895295%0%1%html_test_frame_sync%0%10%',1669895305,1669895305,NULL),(47,2,7,1669895423,11,1,0,1,'11%1%0.000%0%1669895423%0%1669895413%0%1%html_test_frame_sync%0%11%',1669895423,1669895423,NULL),(48,2,8,1669895572,12,1,0,1,'12%1%0.000%0%1669895572%0%1669895562%0%1%html_test_frame_sync%0%12%',1669895572,1669895572,NULL),(49,2,9,1669895918,13,1,0,1,'13%1%0.000%0%1669895918%0%1669895908%0%1%html_test_frame_sync%0%13%',1669895918,1669895918,NULL),(50,2,10,1669895967,14,1,0,1,'14%1%0.000%0%1669895967%0%1669895957%0%1%html_test_frame_sync%0%14%',1669895967,1669895967,NULL),(51,2,11,1669895979,15,1,0,1,'15%1%0.000%0%1669895979%0%1669895969%0%1%html_test_frame_sync%0%15%',1669895979,1669895979,NULL),(52,2,12,1669896754,16,1,0,1,'16%1%0.000%0%1669896754%0%1669896744%0%1%html_test_frame_sync%0%16%',1669896754,1669896754,NULL),(53,2,13,1669896886,17,1,0,1,'17%1%0.000%0%1669896886%0%1669896876%0%1%html_test_frame_sync%0%17%',1669896886,1669896886,NULL),(54,2,14,1669904698,3232,1,0,1,'3232%1%0.000%0%1669904698%0%1669904688%0%1%html_test_frame_sync%0%3232%',1669904698,1669904698,NULL),(55,2,15,1669904876,4966,1,0,1,'4966%1%0.000%0%1669904876%0%1669904866%0%1%html_test_frame_sync%0%4966%',1669904876,1669904876,NULL),(56,2,16,1669904894,1473,1,0,1,'1473%1%0.000%0%1669904894%0%1669904884%0%1%html_test_frame_sync%0%1473%',1669904894,1669904894,NULL),(57,2,17,1669906618,5,1,0,2,'5%2%1669906618%1669906678%1%1,2%7420,1936%0%',1669906619,1669906619,NULL),(58,2,1,1669906686,1,1,0,2,'1%2%1669906686%1669906746%1%2,1%9932,8327%0%',1669906686,1669906686,NULL),(59,2,2,1669906714,2,1,0,2,'2%2%1669906714%1669906774%1%2,1%8027,7302%0%',1669906714,1669906714,NULL),(60,2,3,1669906818,3,1,0,2,'3%2%1669906818%1669906878%1%2,1%6400,2884%0%',1669906818,1669906818,NULL),(61,2,4,1669910402,4,1,0,2,'4%2%1669910402%1669910462%1%2,1%3531,3228%0%',1669910402,1669910402,NULL),(62,2,5,1669910634,5,1,0,2,'5%2%1669910634%1669910694%1%1,2%4940,9609%0%',1669910634,1669910634,NULL),(63,2,6,1669911159,6,1,0,2,'6%2%1669911158%1669911218%1%2,1%3715,6887%0%',1669911159,1669911159,NULL),(64,2,7,1669911246,7,1,0,2,'7%2%1669911246%1669911306%1%2,1%2699,1620%0%',1669911247,1669911247,NULL),(65,2,8,1669911520,8,1,0,2,'8%2%1669911520%1669911580%1%2,1%9925,3077%0%',1669911520,1669911520,NULL),(66,2,9,1669911571,9,1,0,2,'9%2%1669911571%1669911631%1%1,2%5579,3155%0%',1669911571,1669911571,NULL),(67,2,10,1669911737,10,1,0,2,'10%2%1669911736%1669911796%1%1,2%9205,6614%0%',1669911737,1669911737,NULL),(68,2,11,1669912013,11,1,0,2,'11%2%1669912012%1669912072%1%2,1%8004,2886%0%',1669912013,1669912013,NULL),(69,2,12,1669912258,12,1,0,2,'12%2%1669912258%1669912318%1%1,2%7056,1484%0%',1669912258,1669912258,NULL),(70,2,13,1669912403,13,1,0,2,'13%2%1669912403%1669912463%1%1,2%7107,6332%0%',1669912403,1669912403,NULL),(71,2,14,1669912448,14,1,0,2,'14%2%1669912447%1669912507%1%1,2%2311,1500%0%',1669912448,1669912448,NULL),(72,2,15,1669912728,15,1,0,2,'15%2%1669912727%1669912787%1%1,2%8324,3811%0%',1669912728,1669912728,NULL),(73,2,16,1669949649,16,1,0,2,'16%2%1669949649%1669949709%1%1,2%7042,6003%0%',1669949649,1669949649,NULL),(74,2,17,1669949727,17,1,0,2,'17%2%1669949727%1669949787%1%1,2%9928,1447%0%',1669949727,1669949727,NULL),(75,2,18,1669949761,18,1,0,2,'18%2%1669949761%1669949821%1%2,1%5649,1263%0%',1669949762,1669949762,NULL),(76,2,19,1669949963,19,1,0,2,'19%2%1669949962%1669950022%1%1,2%7977,7300%0%',1669949963,1669949963,NULL),(77,2,20,1669950050,20,1,0,2,'20%2%1669950049%1669950109%1%1,2%9566,4124%0%',1669950050,1669950050,NULL),(78,2,21,1669950208,21,1,0,2,'21%2%1669950208%1669950268%1%2,1%4576,7359%0%',1669950208,1669950208,NULL),(79,2,22,1669951751,22,1,0,2,'22%2%1669951751%1669951811%1%1,2%8940,5840%0%',1669951751,1669951751,NULL),(80,2,23,1669951811,23,1,0,2,'23%2%1669951811%1669951871%1%2,1%2557,1602%0%',1669951811,1669951811,NULL),(81,2,24,1669954700,24,1,0,2,'24%2%1669954700%1669954760%1%2,1%9421,8793%0%',1669954700,1669954700,NULL),(82,2,25,1669956235,25,1,0,2,'25%2%1669956234%1669956294%1%1,2%6359,6123%0%',1669956235,1669956235,NULL),(83,2,26,1669956287,26,1,0,2,'26%2%1669956286%1669956346%1%1,2%1698,8490%0%',1669956287,1669956287,NULL),(84,2,27,1669956494,27,1,0,2,'27%2%1669956493%1669956553%1%2,1%8482,8331%0%',1669956494,1669956494,NULL),(85,2,28,1670306368,28,1,0,2,'28%2%1670306368%1670306428%1%2,1%8740,8452%0%',1670306368,1670306368,NULL),(86,2,1,1670306460,1,1,0,2,'1%2%1670306460%1670306520%1%2,1%4575,2423%0%',1670306460,1670306460,NULL),(87,2,1,1670306499,8548,1,0,1,'8548%1%0.000%0%1670306499%0%1670306489%0%1%html_test_frame_sync%0%8548%',1670306499,1670306499,NULL),(88,2,2,1670306779,6608,1,0,1,'6608%1%0.000%0%1670306779%0%1670306769%0%1%html_test_frame_sync%0%6608%',1670306779,1670306779,NULL),(89,2,3,1670306792,8829,1,0,1,'8829%1%0.000%0%1670306792%0%1670306782%0%1%html_test_frame_sync%0%8829%',1670306792,1670306792,NULL),(90,2,4,1670307008,1,1,0,2,'1%2%1670307008%1670307068%1%2,1%8603,9900%0%',1670307008,1670307008,NULL),(91,2,5,1670307027,2,1,0,2,'2%2%1670307027%1670307087%1%1,2%9648,2666%0%',1670307027,1670307027,NULL),(92,2,6,1670307347,3,1,0,2,'3%2%1670307347%1670307407%1%1,2%1563,2346%0%',1670307348,1670307348,NULL),(93,2,7,1670307913,4,1,0,2,'4%2%1670307913%1670307973%1%2,1%9759,1653%0%',1670307913,1670307913,NULL),(94,2,8,1670308245,5,1,0,2,'5%2%1670308244%1670308304%1%2,1%9855,9673%0%',1670308245,1670308245,NULL),(95,2,9,1670311449,6,1,0,2,'6%2%1670311448%1670311508%1%1,2%8215,3265%0%',1670311449,1670311449,NULL),(96,2,10,1670311523,7,1,0,2,'7%2%1670311523%1670311583%1%2,1%4904,2073%0%',1670311523,1670311523,NULL),(97,2,11,1670312009,8,1,0,2,'8%2%1670312009%1670312069%1%2,1%7755,2913%0%',1670312009,1670312009,NULL),(98,2,12,1670312523,9,1,0,2,'9%2%1670312522%1670312582%1%2,1%2346,8391%0%',1670312523,1670312523,NULL),(99,2,13,1670312905,10,1,0,2,'10%2%1670312905%1670312965%1%2,1%9544,1075%0%',1670312905,1670312905,NULL),(100,2,14,1670313315,11,1,0,2,'11%2%1670313315%1670313375%1%2,1%4843,1452%0%',1670313315,1670313315,NULL),(101,2,15,1670313336,12,1,0,2,'12%2%1670313336%1670313396%1%2,1%8354,1401%0%',1670313336,1670313336,NULL),(102,2,16,1670313776,13,1,0,2,'13%2%1670313776%1670313836%1%1,2%3636,9527%0%',1670313776,1670313776,NULL),(103,2,17,1670314315,14,1,0,2,'14%2%1670314315%1670314375%1%2,1%1093,6772%0%',1670314315,1670314315,NULL),(104,2,18,1670315852,15,1,0,2,'15%2%1670315852%1670315912%1%1,2%8793,2357%0%',1670315853,1670315853,NULL),(105,2,19,1670316245,16,1,0,2,'16%2%1670316245%1670316305%1%1,2%4543,7810%0%',1670316245,1670316245,NULL),(106,2,20,1670319896,17,1,0,2,'17%2%1670319896%1670319956%1%2,1%5576,4282%0%',1670319896,1670319896,NULL),(107,2,21,1670320332,18,1,0,2,'18%2%1670320332%1670320392%1%2,1%8039,3534%0%',1670320332,1670320332,NULL),(108,2,22,1670321273,19,1,0,2,'19%2%1670321273%1670321333%1%2,1%1626,9897%0%',1670321273,1670321273,NULL),(109,2,23,1670322559,20,1,0,2,'20%2%1670322559%1670322619%1%2,1%9291,2480%0%',1670322559,1670322559,NULL),(110,2,24,1670322694,21,1,0,2,'21%2%1670322694%1670322754%1%1,2%2399,8552%0%',1670322694,1670322694,NULL),(111,2,25,1670322810,22,1,0,2,'22%2%1670322810%1670322870%1%1,2%4979,2829%0%',1670322810,1670322810,NULL),(112,2,26,1670322832,23,1,0,2,'23%2%1670322831%1670322891%1%1,2%1573,4651%0%',1670322832,1670322832,NULL),(113,2,27,1670323397,24,1,0,2,'24%2%1670323397%1670323457%1%2,1%8753,1308%0%',1670323397,1670323397,NULL),(114,2,28,1670325371,25,1,0,2,'25%2%1670325371%1670325431%1%2,1%7327,3519%0%',1670325371,1670325371,NULL),(115,2,29,1670325920,26,1,0,2,'26%2%1670325920%1670325980%1%1,2%9444,7948%0%',1670325920,1670325920,NULL),(116,2,30,1670326028,27,1,0,2,'27%2%1670326028%1670326088%1%2,1%7787,6267%0%',1670326028,1670326028,NULL),(117,2,31,1670326533,28,1,0,2,'28%2%1670326532%1670326592%1%1,2%5017,3483%0%',1670326533,1670326533,NULL),(118,2,32,1670326970,29,1,0,2,'29%2%1670326970%1670327030%1%2,1%6998,5528%0%',1670326970,1670326970,NULL),(119,2,33,1670327033,30,1,0,2,'30%2%1670327032%1670327092%1%2,1%7424,5693%0%',1670327033,1670327033,NULL),(120,2,34,1670327223,31,1,0,2,'31%2%1670327223%1670327283%1%1,2%4756,4564%0%',1670327223,1670327223,NULL),(121,2,35,1670327358,32,1,0,2,'32%2%1670327358%1670327418%1%1,2%5587,2597%0%',1670327358,1670327358,NULL),(122,2,36,1670327935,33,1,0,2,'33%2%1670327935%1670327995%1%2,1%8398,3223%0%',1670327935,1670327935,NULL),(123,2,37,1670330912,34,1,0,2,'34%2%1670330912%1670330972%1%2,1%1990,1262%0%',1670330912,1670330912,NULL),(124,2,38,1670331003,35,1,0,2,'35%2%1670331003%1670331063%1%2,1%3269,9181%0%',1670331003,1670331003,NULL),(125,2,39,1670331204,36,1,0,2,'36%2%1670331204%1670331264%1%2,1%6999,2690%0%',1670331204,1670331204,NULL),(126,2,40,1670331490,37,1,0,2,'37%2%1670331490%1670331550%1%1,2%1295,1287%0%',1670331490,1670331490,NULL),(127,2,41,1670331771,38,1,0,2,'38%2%1670331771%1670331831%1%2,1%2928,3983%0%',1670331771,1670331771,NULL),(128,2,42,1670332245,39,1,0,2,'39%2%1670332245%1670332305%1%2,1%8466,5595%0%',1670332245,1670332245,NULL),(129,2,43,1670332767,40,1,0,2,'40%2%1670332767%1670332827%1%2,1%6088,4068%0%',1670332767,1670332767,NULL),(130,2,44,1670333001,41,1,0,2,'41%2%1670333001%1670333061%1%2,1%7685,2616%0%',1670333002,1670333002,NULL),(131,2,45,1670333238,42,1,0,2,'42%2%1670333238%1670333298%1%2,1%3483,8547%0%',1670333238,1670333238,NULL),(132,2,46,1670333361,43,1,0,2,'43%2%1670333361%1670333421%1%2,1%9490,5288%0%',1670333361,1670333361,NULL),(133,2,47,1670334476,44,1,0,2,'44%2%1670334476%1670334536%1%2,1%6608,3401%0%',1670334477,1670334477,NULL),(134,2,48,1670335146,45,1,0,2,'45%2%1670335145%1670335205%1%1,2%1324,4731%0%',1670335146,1670335146,NULL),(135,2,49,1670335328,46,1,0,2,'46%2%1670335328%1670335388%1%2,1%9002,5384%0%',1670335328,1670335328,NULL),(136,2,50,1670335483,47,1,0,2,'47%2%1670335483%1670335543%1%1,2%7037,6136%0%',1670335483,1670335483,NULL),(137,2,51,1670335605,48,1,0,2,'48%2%1670335604%1670335664%1%1,2%6848,1737%0%',1670335605,1670335605,NULL),(138,2,52,1670335707,49,1,0,2,'49%2%1670335707%1670335767%1%2,1%7526,7503%0%',1670335707,1670335707,NULL),(139,2,53,1670335814,50,1,0,2,'50%2%1670335814%1670335874%1%1,2%8350,7229%0%',1670335815,1670335815,NULL),(140,2,54,1670335980,51,1,0,2,'51%2%1670335980%1670336040%1%1,2%3082,7537%0%',1670335980,1670335980,NULL),(141,2,55,1670337051,52,1,0,2,'52%2%1670337050%1670337110%1%2,1%9845,5554%0%',1670337051,1670337051,NULL),(142,2,56,1670337473,53,1,0,2,'53%2%1670337473%1670337533%1%2,1%5760,8503%0%',1670337473,1670337473,NULL),(143,2,57,1670337655,54,1,0,2,'54%2%1670337655%1670337715%1%1,2%6474,3238%0%',1670337655,1670337655,NULL),(144,2,58,1670338010,55,1,0,2,'55%2%1670338010%1670338070%1%2,1%8920,7693%0%',1670338010,1670338010,NULL),(145,2,59,1670824367,6895,1,0,1,'6895%1%0.000%0%1670567540%0%1670567530%0%1%html_test_frame_sync%0%6895%',1670824367,1670824367,NULL),(146,2,60,1671701210,6018,1,0,1,'6018%1%0.000%0%1671701210%0%1671701200%0%1%html_test_frame_sync%0%6018%',1671701210,1671701210,NULL),(147,2,61,1671703141,4589,1,0,1,'4589%1%0.000%0%1671703027%0%1671703017%0%1%html_test_frame_sync%0%4589%',1671703141,1671703141,NULL),(148,2,62,1671703160,56,1,0,2,'56%2%1671703160%1671703220%1%2,1%3827,4193%0%',1671703160,1671703160,NULL),(149,2,64,1671706848,57,1,0,2,'57%2%1671706848%1671706908%1%2,1%5470,1181%0%',1671706848,1671706848,NULL),(150,2,66,1671707273,58,1,0,2,'58%2%1671707273%1671707333%1%1,2%2077,9037%0%',1671707273,1671707273,NULL),(151,2,68,1671707936,59,1,0,2,'59%2%1671707935%1671707995%1%1,2%4888,5087%0%',1671707936,1671707936,NULL),(152,2,69,1671765653,60,1,0,2,'60%2%1671765653%1671765713%1%2,1%7632,7265%0%',1671765653,1671765653,NULL),(153,2,70,1671765727,61,1,0,2,'61%2%1671765727%1671765787%1%1,2%2960,1510%0%',1671765727,1671765727,NULL),(154,2,71,1671766337,62,1,0,2,'62%2%1671766337%1671766397%1%1,2%7424,5830%0%',1671766337,1671766337,NULL),(155,2,72,1671766355,63,1,0,2,'63%2%1671766355%1671766415%1%2,1%9433,1915%0%',1671766355,1671766355,NULL),(156,2,73,1671766548,64,1,0,2,'64%2%1671766548%1671766608%1%1,2%7103,6543%0%',1671766548,1671766548,NULL),(157,2,74,1671767190,65,1,0,2,'65%2%1671767190%1671767250%1%1,2%5882,1032%0%',1671767190,1671767190,NULL),(158,2,75,1671767372,66,1,0,2,'66%2%1671767372%1671767432%1%1,2%4716,2450%0%',1671767372,1671767372,NULL),(159,2,76,1671767421,67,1,0,2,'67%2%1671767421%1671767481%1%2,1%4589,3832%0%',1671767421,1671767421,NULL),(160,2,77,1671772692,68,1,0,2,'68%2%1671772692%1671772752%1%2,1%9747,4728%0%',1671772692,1671772692,NULL),(161,2,78,1671776023,69,1,0,2,'69%2%1671776023%1671776083%1%2,1%5321,3107%0%',1671776023,1671776023,NULL),(162,2,79,1671779743,70,1,0,2,'70%2%1671779742%1671779802%1%2,1%2365,8380%0%',1671779743,1671779743,NULL),(163,2,80,1671781035,71,1,0,2,'71%2%1671781035%1671781095%1%2,1%3146,1610%0%',1671781035,1671781035,NULL);
/*!40000 ALTER TABLE `game_match_push` ENABLE KEYS */;
UNLOCK TABLES;

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
  `off_line_wait_time` int NOT NULL DEFAULT '0' COMMENT '某玩家掉线等待时长(秒)',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_match_rule`
--

LOCK TABLES `game_match_rule` WRITE;
/*!40000 ALTER TABLE `game_match_rule` DISABLE KEYS */;
INSERT INTO `game_match_rule` VALUES (1,'测试5V5加权重公式',1,1,10,60,1,2,4,'( <age> * 2 ) + ( <level> * 5)','',0,0,0,0,0,NULL,20,0,10),(2,'测试吃鸡',1,1,10,60,2,2,2,'','',0,0,0,0,0,NULL,10,0,10);
/*!40000 ALTER TABLE `game_match_rule` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=114 DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-成功';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_match_success`
--



--
-- Table structure for table `game_sync_room`
--

DROP TABLE IF EXISTS `game_sync_room`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_sync_room` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
  `rule_id` int NOT NULL DEFAULT '0' COMMENT 'rule_id',
  `add_time` int NOT NULL DEFAULT '0' COMMENT '添加时间',
  `start_time` int NOT NULL DEFAULT '0' COMMENT '开始游戏时间',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '游戏结束时间',
  `ready_timeout` tinyint(1) NOT NULL DEFAULT '0' COMMENT '准备超时时间',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态',
  `sequence_number` int NOT NULL DEFAULT '0' COMMENT '匹配成功后，无人来取，超时',
  `rand_seek` int NOT NULL DEFAULT '0' COMMENT '当前逻辑帧号',
  `wait_player_offline` int NOT NULL DEFAULT '0' COMMENT '玩家掉线等待时间',
  `player_ids` varchar(100) NOT NULL DEFAULT '' COMMENT '玩家列表',
  `players_ack_list` varchar(100) NOT NULL DEFAULT '' COMMENT '最后一帧的确认情况',
  `end_total` varchar(255) NOT NULL DEFAULT '' COMMENT '结算信息',
  `logic_frame_history` text COMMENT '玩家的历史所有记录',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='游戏匹配-小组信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_sync_room`
--

LOCK TABLES `game_sync_room` WRITE;
/*!40000 ALTER TABLE `game_sync_room` DISABLE KEYS */;
/*!40000 ALTER TABLE `game_sync_room` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `instance`
--

LOCK TABLES `instance` WRITE;
/*!40000 ALTER TABLE `instance` DISABLE KEYS */;
INSERT INTO `instance` VALUES (1,1,'mysql','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(2,1,'redis','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(3,1,'email','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(4,1,'etcd','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(5,1,'alert','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(6,1,'cdn','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(7,1,'sms','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(8,1,'http','127.0.0.1','1111',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(9,1,'domain','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(10,1,'oss','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(11,1,'grpc','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(12,1,'gateway','0.0.0.0,127.0.0.1','1122,2233',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(13,1,'agora','127.0.0.1','',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(14,1,'super_visor','127.0.0.1','9001',1,'','','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(15,1,'mysql','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(16,1,'redis','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(17,1,'email','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(18,1,'etcd','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(19,1,'alert','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(20,1,'cdn','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(21,1,'sms','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(22,1,'http','192.168.1.21','2222',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(23,1,'domain','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(24,1,'oss','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(25,1,'grpc','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(26,1,'gateway','0.0.0.0,127.0.0.1','1122,2233',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(27,1,'agora','192.168.1.21','',2,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(28,1,'super_visor','192.168.1.21','9002',2,'user','123','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(29,1,'mysql','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(30,1,'redis','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(31,1,'email','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(32,1,'etcd','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(33,1,'alert','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(34,1,'cdn','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(35,1,'sms','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(36,1,'http','8.142.177.235','4444',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(37,1,'domain','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(38,1,'oss','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(39,1,'grpc','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(40,1,'gateway','0.0.0.0,127.0.0.1','1122,2233',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(41,1,'agora','8.142.177.235','',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(42,1,'super_visor','8.142.177.235','9988',4,'user','123','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(43,1,'redis','127.0.0.1','6375',5,'','ckckarar','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(44,2,'mysql','rm-8vb10pi2gz9rma8p8.mysql.zhangbei.rds.aliyuncs.com','3306',5,'seed','willbeOK618','',1,'小z',0,0,0,0,0,NULL),(45,1,'http','127.0.0.1','5555',5,'','','',1,'小z',0,0,0,0,0,NULL),(46,2,'oss','oss-cn-beijing.aliyuncs.com','servicebase',5,'LTAI5tJbjZiWQ9Xn9N2brRFD','GcVCuaZA7KWxV0o7UyzzSXhg9zCQfm','',1,'小z',0,0,0,0,0,NULL),(47,1,'etcd','127.0.0.1','2379',5,'','','',1,'小z',0,0,0,0,0,NULL),(48,1,'grpc','127.0.0.1','5656',5,'pbservice','','',1,'小z',0,0,0,0,0,NULL),(49,3,'email','smtp.exmail.qq.com','2EGdKudfF6KvdosN',5,'xxxxxx@seedreality.com','xxxxx','',1,'小z',0,0,0,0,0,NULL),(50,1,'gateway','0.0.0.0,127.0.0.1','1122,2233',5,'','','',1,'小z',0,0,0,0,0,NULL),(51,5,'agora','','',5,'8ff429463a234c7bae327d74941a5956','b58033d109354bce9205d5f2458900c9','',1,'小z',0,0,0,0,0,NULL),(52,2,'domain','static.seedreality.com','',5,'','','',1,'小z',0,0,0,0,0,NULL),(53,2,'ali_email','','',5,'','','',1,'小z',0,0,0,0,0,NULL),(54,2,'cdn','','',5,'','','',1,'小z',0,0,0,0,0,NULL),(55,1,'alert','','',5,'','','',1,'小z',0,0,0,0,0,NULL),(56,1,'super_visor','127.0.0.1','9988',5,'ckadmin','ckckarar','',1,'小z',0,0,0,0,0,NULL),(57,2,'sms','','',5,'','','',1,'小z',0,0,0,0,0,NULL);
/*!40000 ALTER TABLE `instance` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `mail_group`
--

LOCK TABLES `mail_group` WRITE;
/*!40000 ALTER TABLE `mail_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `mail_group` ENABLE KEYS */;
UNLOCK TABLES;

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
  `people_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '接收人群，1单发2群发3指定group4指定tag5指定UIDS',
  `created_at` bigint NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` bigint NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `deleted_at` bigint DEFAULT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`),
  KEY `deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COMMENT='站内信 - 发送规则配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--



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
-- Dumping data for table `operation_record`
--

--
-- Table structure for table `project`
--
INSERT INTO `project` VALUES (14,'120',1,'120眼镜','ck120T789!@#',1,'im120doctor',7,'',1650001049,1650001049,NULL);
INSERT INTO `project` VALUES (15,'120',1,'120WEB','ck120Ta!3D$)',1,'im120User',4,'',1650001049,1650001049,NULL);

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
-- Dumping data for table `project`
--

LOCK TABLES `project` WRITE;
/*!40000 ALTER TABLE `project` DISABLE KEYS */;
INSERT INTO `project` VALUES (1,'GameMatch',1,'小游戏-玩家匹配机制','ckgamematch',1,'imgamematch',2,'git://github.com/mqzhifu/gamematch.git',1650001049,1650001049,NULL),(2,'FrameSync',1,'游戏-帧同步','ckframesync',1,'imframesync',2,'git://github.com/mqzhifu/frame_sync.git',1650001049,1650001049,NULL),(6,'Zgoframe',1,'go框架测试','ckZgoframe',1,'imzgoframe',2,'git@github.com:mqzhifu/zgoframe.git',1650001049,1650001049,NULL),(9,'Gateway',1,'公共网关','ckgateway',2,'imgateway',2,'git://github.com/mqzhifu/gateway.git',1650001049,1650001049,NULL),(10,'Zwebuigo',1,'后台管理系统','ckZwebuigo',1,'imzwebuigo',2,'https://github.com/mqzhifu/zwebuigo.git',1650001049,1650001049,NULL),(11,'Zwebuivue',2,'后台管理系统-VUE','ckZwebuivue',1,'imzwebuivue',4,'https://github.com/mqzhifu/zwebuivue.git',1650001049,1650001049,NULL),(12,'TwinAgora',2,'数据孪生-专家指导(声网)','ckTwinAgora',1,'imtwinagora',4,'https://github.com/mqzhifu/twin_agora.git',1650001049,1650001049,NULL),(13,'AgoraUnity',5,'数据孪生-UNITY端','ckAgoraUnity',1,'imagoraunity',7,'http://192.168.1.22:40080/jiaxing.zhu/Agora.git',1650001049,1650001049,NULL);
/*!40000 ALTER TABLE `project` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `server`
--

LOCK TABLES `server` WRITE;
/*!40000 ALTER TABLE `server` DISABLE KEYS */;
INSERT INTO `server` VALUES (1,'本地',1,'127.0.0.1','127.0.0.1',1,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(2,'开发',1,'192.168.1.21','192.168.1.21',2,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(3,'测试',1,'2.2.2.2','127.0.0.1',3,2,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(4,'预发布',1,'8.142.177.235','172.27.198.210',4,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(5,'线上',1,'8.142.161.156','172.27.218.143',5,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL);
/*!40000 ALTER TABLE `server` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3 COMMENT='短信发送日志';

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
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COMMENT='短信发送规则配置';

LOCK TABLES `sms_rule` WRITE;
/*!40000 ALTER TABLE `sms_rule` DISABLE KEYS */;
INSERT INTO `sms_rule` VALUES (1,6,'短信注册','{nickname},您好：欢迎注册本网站，验证码为：{auth_code},{auth_expire_time}秒后将失效，勿告诉他人，防止被骗',1,10,60,1,300,'0',1,'','1',0,'','','1',0,0,NULL),(2,6,'短信登陆','{nickname},您好：登陆验证码为：{auth_code},{auth_expire_time} 秒后将失效，勿告诉他人，防止被骗。',1,10,60,1,300,'0',1,'','1',0,'','','1',0,0,NULL),(3,6,'找加密码','找回密码',1,10,60,1,300,'0',1,'','',1,'','','',1,0,NULL),(4,6,'报警','报警，程序出错。级别：{level}，项目ID:{project_id}，内容：{content}',2,10,60,1,0,'0',1,'','SMS_273495087',1,'','','1',300,0,NULL);
/*!40000 ALTER TABLE `sms_rule` ENABLE KEYS */;
UNLOCK TABLES;

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
-- Dumping data for table `statistics_log`
--



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
) ENGINE=InnoDB AUTO_INCREMENT=668 DEFAULT CHARSET=utf8mb3 COMMENT='AR远程呼叫,房间记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `twin_agora_room`
--


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
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'2d879cfe-d900-45ae-a3e5-af3517eb8d02',6,1,0,'frame_sync_1','e10adc3949ba59abbe56e057f20f883e','e10adc3949ba59abbe56e057f20f883e','sync_1','','',1,1,2,2,'','',1658995531,1658995531,NULL,'ckck',0),(2,'4d69dee4-38f3-47ed-8dee-c4792df2e2c6',6,2,0,'frame_sync_2','e10adc3949ba59abbe56e057f20f883e','','sync_2','','',1,1,2,2,'','',1658995531,1658995531,NULL,'ckck',0),(3,'4d69dee4-38f3-47ed-8dee-c4792df2e2c3',6,1,0,'frame_sync_3','e10adc3949ba59abbe56e057f20f883e','','sync_3','','',1,1,2,2,'','',1658995531,1658995531,NULL,'ckck',0),(4,'4d69dee4-38f3-47ed-8dee-c4792df2e2c1',6,2,0,'frame_sync_4','e10adc3949ba59abbe56e057f20f883e','','sync_4','','',1,1,2,2,'','',1658995531,1658995531,NULL,'ckck',0),(9,'111111',6,2,0,'calluser','e10adc3949ba59abbe56e057f20f883e','','calluser','','',1,1,2,2,'','',1658995531,1658995531,NULL,'',1),(10,'2222',6,1,0,'doctor','e10adc3949ba59abbe56e057f20f883e','','doctor','','',1,1,2,2,'','',1658995531,1658995531,NULL,'',2),(11,'3333',6,1,0,'vruser','e10adc3949ba59abbe56e057f20f883e','','vr','','',1,1,2,2,'','',1658995531,0,NULL,'',1);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=3262 DEFAULT CHARSET=utf8mb3 COMMENT='用户登陆记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_login`
--


/*!40000 ALTER TABLE `user_login` DISABLE KEYS */;

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


