-- MySQL dump 10.13  Distrib 8.0.22, for macos10.15 (x86_64)
--
-- Host: 8.142.177.235    Database: test
-- ------------------------------------------------------
-- Server version	8.0.22

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
-- Dumping data for table `sms_rule`
--

LOCK TABLES `sms_rule` WRITE;
/*!40000 ALTER TABLE `sms_rule` DISABLE KEYS */;
INSERT INTO `sms_rule` VALUES (1,0,'',0,'',0,10,30,'',0,NULL,0,'短信注册',1,'','','{nickname},您好：欢迎注册本网站，验证码为：{auth_code},{auth_expire_time}秒后将失效，勿告诉他人，防止被骗',1,300,'',11),(2,0,'',0,'',0,10,60,'',0,NULL,6,'短信登陆',1,'','','{nickname},您好：登陆验证码为：{auth_code},{auth_expire_time} 秒后将失效，勿告诉他人，防止被骗。',1,300,'',12),(3,0,'找加密码',0,'',0,10,30,'',0,NULL,6,'找回密码',1,'','','找回密码',1,300,'',13);
/*!40000 ALTER TABLE `sms_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `mail_rule`
--

LOCK TABLES `mail_rule` WRITE;
/*!40000 ALTER TABLE `mail_rule` DISABLE KEYS */;
INSERT INTO `mail_rule` VALUES (1,6,'通知完成任务','您好，恭喜您完成了新手任务，我们将会给您大大的奖励哟~',2,10,60,1,0,'',1,1,0,0,NULL);
/*!40000 ALTER TABLE `mail_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `project`
--

LOCK TABLES `project` WRITE;
/*!40000 ALTER TABLE `project` DISABLE KEYS */;
INSERT INTO `project` VALUES (6,'Zgoframe',1,'go框架测试','dddd',1,'imzgoframe','git@github.com:mqzhifu/zgoframe.git',1649122850,1650001049,NULL,0),(10,'Zwebui',1,'后台管理系统','aaaaaa',1,'imzwebui','git@github.com:mqzhifu/zwebui.git',1649125912,1650001049,NULL,0);
/*!40000 ALTER TABLE `project` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `server`
--

LOCK TABLES `server` WRITE;
/*!40000 ALTER TABLE `server` DISABLE KEYS */;
INSERT INTO `server` VALUES (1,'本地',1,'127.0.0.1','127.0.0.1',1,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(2,'开发',1,'192.168.1.21','192.168.1.21',2,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(3,'测试',1,'2.2.2.2','127.0.0.1',3,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(4,'预发布',1,'8.142.177.235','172.27.198.210',4,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL),(5,'线上',1,'8.142.161.156','172.27.218.143',5,1,'','小z',1650006845,1650006845,100,1650006845,0,NULL);
/*!40000 ALTER TABLE `server` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `instance`
--

LOCK TABLES `instance` WRITE;
/*!40000 ALTER TABLE `instance` DISABLE KEYS */;
INSERT INTO `instance` VALUES (1,1,'mysql','127.0.0.1','3306',1,'root','mqzhifu','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(2,1,'redis','127.0.0.1','6370',1,'','1234','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(3,1,'etcd','127.0.0.1','2379',1,'','','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(4,1,'prometheus','127.0.0.1','3306',1,'','','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(5,1,'es','127.0.0.1','3306',1,'','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(6,1,'kibana','127.0.0.1','3306',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(7,1,'grafana','127.0.0.1','3306',1,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(8,1,'http','0.0.0.0','1111',1,'','','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(9,1,'mysql','8.142.177.235','3306',4,'root','mqzhifu','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(10,1,'redis','8.142.177.235','6370',4,'','1234','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(11,1,'etcd','8.142.177.235','2379',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(12,1,'prometheus','8.142.177.235','3306',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(13,1,'es','8.142.177.235','3306',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(14,1,'kibana','8.142.177.235','3306',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(15,1,'static','8.142.177.235','3306',4,'aaaa','bbbb','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(16,1,'http','127.0.0.1','6375',5,'','1234','',1,'小z',1650006845,1650006845,200,1650006845,0,NULL),(17,2,'mysql','rm-8vb10pi2gz9rma8p8.mysql.zhangbei.rds.aliyuncs.com','3306',5,'seed','willbeOK618','',1,'小z',0,0,0,0,0,NULL);
/*!40000 ALTER TABLE `instance` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-18 19:51:56
