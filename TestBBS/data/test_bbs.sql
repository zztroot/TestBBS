-- MySQL dump 10.13  Distrib 8.0.19, for Win64 (x86_64)
--
-- Host: localhost    Database: test_bbs
-- ------------------------------------------------------
-- Server version	8.0.19

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
-- Table structure for table `tb_article`
--

DROP TABLE IF EXISTS `tb_article`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tb_article` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL DEFAULT '',
  `content` varchar(16000) NOT NULL DEFAULT '',
  `user_id` bigint NOT NULL,
  `type_id` bigint NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_article`
--

LOCK TABLES `tb_article` WRITE;
/*!40000 ALTER TABLE `tb_article` DISABLE KEYS */;
INSERT INTO `tb_article` VALUES (2,'python监控手机内存信息，当手机内存快要消耗完时，结束自动化测试脚本','我们在做手机camera 自动化压力测试时，当手机提示内存已满，而脚本还在继续执行。解决思路：我们要实现监控手机内存信息，当内存快满时将自动化测试脚本停止执行。\n\n内存监控部分：\n```python\ndef  get_memory():\n	mermory_data = os.popen(\"adb shell dumpsys meminfo \").read()\n	com = re.compile(\'Total RAM:(.*?)K \\\\(.*?\')\n	com_1 = re.compile(\' Used RAM:(.*?)K \\\\(.*?\')\n	total = re.findall(com, mermory_data)\n	used = re.findall(com_1, mermory_data)\n	for t in total:\n		t = filter(str.isdigit, t)\n		t = int(\"\".join(list(t)))\n	for u in used:\n		u = filter(str.isdigit, u)\n		u = int(\"\".join(list(u)))\n	return t, u\n\n```\n\n使用内存和总内存对比，得到当前手机内存情况。\n```python\ndef main():\n	while True:\n		mermory = get_memory()\n		if mermory[1] < (mermory[0] - 1000):\n			camera_test() #测试用例函数\n		else:\n			print(\"手机内存快满了，测试即将停止\")\n			break\n\n```',4,2,'2020-07-15 16:58:17'),(4,'python实现VTS和CTS-ON-GSI自动flash system.img脚本','直接上代码:\n\n```python\nimport os \nimport time\nimport sys\nimport re\nfrom colorama import init\n \ninit(autoreset=True)\na = sys.argv[1]\ntry:\n	boot = sys.argv[2]\nexcept:\n	boot = \" \"\nsystem = \"system.img\"\n\ndef flashGsi(a, boot, sn, s):\n	os.system(\"adb -s {} reboot bootloader\".format(sn))\n	time.sleep(1)\n	if a == \"vts\":\n		os.system(\"fastboot flash boot {}\".format(boot))\n		time.sleep(1)\n		os.system(\"fastboot -w\")\n	else:\n		os.system(\"fastboot -w\")\n	time.sleep(1)\n	os.system(\"fastboot reboot fastboot\")\n	time.sleep(2)\n	os.system(\"fastboot flash system {}\".format(system))\n	time.sleep(5)\n	os.system(\"fastboot reboot\")\n	print(\"\\n\")\n	#output = \'*\'*int((s/2-37)) + \"\\033[0;32;40m\\t{}：此设备刷GSI成功，正在重启中\\033[0m\".format(sn) + \'*\'*int((s/2-37))\n	print(\"#####\\033[0;32;40m{}:此设备刷GSI成功,正在重启中\\033[0m#####\".format(sn).center(s, \'*\'))\n\ndef getDevicesSn():\n    SN_list = []\n    device_info = os.popen(\'adb devices\').read()\n    for line in device_info.splitlines():\n        if line == \'List of devices attached\':\n            continue\n        else:\n            com = re.compile(\'(.*?)\\tde.*?\')\n            SN = re.findall(com, line)\n            for i in SN:\n                SN_list.append(i)\n    return SN_list\n\nif __name__ == \'__main__\':\n    width = os.get_terminal_size().columns\n    sn = getDevicesSn()\n    for i in sn:\n        print(\"\\n\"+ \"#####\\033[0;32;40m{}:此设备正在刷GSI,请稍等!\\033[0m#####\".format(i).center(width, \'*\'))\n        flashGsi(a, boot, i, width)\n        time.sleep(10)\n    print(\"\\n\"+ \"#####\\033[0;32;40m共{}台手机刷GSI完成!\\033[0m#####\".format(len(sn)).center(width, \'*\'))\n    time.sleep(2)\n    if a == \'gsi\':\n        print(\"\\n\"+ \"#####\\033[0;31;40m请手动点击Allow USB debugging弹框\\033[0m#####\".center(width))\n        time.sleep(75)\n        os.system(\"python3 setting.py\")\n        time.sleep(1)\n        os.system(\"python3 auto_media_push.py\")\n    else:\n        time.sleep(75)\n        os.system(\"python3 setting.py\")\n\n```',5,3,'2020-07-15 18:34:45'),(5,'python实现安装当前路径所有apk','直接上代码:\n```python\nimport os\nimport time\n\ndef install():\n	filepath = os.getcwd()\n	files = os.listdir(filepath)\n	a = 1\n	for file in files:\n		if file.endswith(\'.apk\'):\n			os.popen(\"adb install \\\"{}\\\\{}\\\"\".format(filepath, file))\n			time.sleep(5)\n			print(\"{}.........安装完成\".format(file))\n			a += 1\n		else:\n			continue\n	print(\'\\n总共安装了{}个APK\\n\'.format(a))\n	os.system(\"@echo ****安装全部完成，请点击任意键退出控制台****\")\n	os.popen(\"pause 5\")\n\ninstall()\n```',4,2,'2020-07-15 18:39:05'),(6,'golang 对mysql数据库的常用操作','导入包：\n```go\nimport (\n	\"database/sql\"\n	\"fmt\"\n	_\"github.com/go-sql-driver/mysql\"\n)\n\n```\n\n连接数据：\n```go\ndb, err := sql.Open(\"mysql\", \"root:123456789@/mydb?charset=utf8\")\n	if err != nil{\n		fmt.Println(\"connce mysql fialed\", err)\n	}\n```\n\n查询数据：\n```go\nvar id int\n	var username, password string\n	rows, err := db.Query(\"SELECT * FROM mydb.`user-login`;\")\n	if err != nil{\n		fmt.Println(err)\n	}\n	for rows.Next(){\n		rows.Scan(&id, &username, &password)\n		fmt.Println(id, username, password)\n	}\n\n```\n\n插入数据：\n```go\nret, _ := db.Exec(\"insert into mydb.`user-login` (username, password) values(\'dandan\', \'123456\')\")\n	insID, _ := ret.LastInsertId()\n	fmt.Println(insID)\n\n```\n\n修改数据：\n```go\nret2, _ := db.Exec(\"update mydb.`user-login` set username=\'zzt\' where userid=?\", 1)\n	//ret2, _ := db.Exec(\"update mydb.`user-login` set username=\'zzt\' where userid=1\")\n	affNums, _ := ret2.RowsAffected()\n	fmt.Println(affNums)\n\n	defer  db.Close()\n\n```',6,4,'2020-07-16 14:07:32');
/*!40000 ALTER TABLE `tb_article` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tb_comments`
--

DROP TABLE IF EXISTS `tb_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tb_comments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `content` varchar(16000) NOT NULL DEFAULT '',
  `user_id` bigint NOT NULL,
  `article_id` bigint NOT NULL,
  `to_user` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_comments`
--

LOCK TABLES `tb_comments` WRITE;
/*!40000 ALTER TABLE `tb_comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `tb_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tb_sentence`
--

DROP TABLE IF EXISTS `tb_sentence`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tb_sentence` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `content` varchar(999) NOT NULL DEFAULT '',
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_sentence`
--

LOCK TABLES `tb_sentence` WRITE;
/*!40000 ALTER TABLE `tb_sentence` DISABLE KEYS */;
INSERT INTO `tb_sentence` VALUES (1,'越是无能的人，越喜欢挑剔别人的错儿。',5),(3,'盛年不重来,一日难再晨。及时宜自勉,岁月不待人。',6),(23,'最灵繁的人也看不见自己的背脊。',5),(24,'在人生的道路上，当你的希望一个个落空的时候，你也要坚定，要沉着。',5),(25,'没有一件工作是旷日持久的，除了那件你不敢拌着手进行的工作。',5),(26,'一知半解的人，多不谦虚；见多识广有本领的人，一定谦虚。',5),(27,'没有人不爱惜他的生命，但很少人珍视他的时间。',5);
/*!40000 ALTER TABLE `tb_sentence` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tb_type`
--

DROP TABLE IF EXISTS `tb_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tb_type` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(500) NOT NULL DEFAULT '',
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_type`
--

LOCK TABLES `tb_type` WRITE;
/*!40000 ALTER TABLE `tb_type` DISABLE KEYS */;
INSERT INTO `tb_type` VALUES (2,'Python',4),(3,'Python',5),(4,'golang',6);
/*!40000 ALTER TABLE `tb_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tb_user`
--

DROP TABLE IF EXISTS `tb_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tb_user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '',
  `pwd` varchar(255) NOT NULL DEFAULT '',
  `phone` bigint NOT NULL DEFAULT '0',
  `sex` varchar(255) NOT NULL DEFAULT '',
  `img` varchar(255) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tb_user`
--

LOCK TABLES `tb_user` WRITE;
/*!40000 ALTER TABLE `tb_user` DISABLE KEYS */;
INSERT INTO `tb_user` VALUES (4,'zhongtian','123456',0,'男生','../static/img/man.jpg','2020-07-15 16:56:36'),(5,'dandan','123456',18483658580,'女生','../static/img/woman.jpg','2020-07-15 17:57:50'),(6,'fangai','123456',18483658580,'女生','../static/img/woman.jpg','2020-07-16 14:02:45'),(7,'yanglinhao','123456',0,'男生','../static/img/man.jpg','2020-07-16 15:04:36');
/*!40000 ALTER TABLE `tb_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-07-17 15:45:24
