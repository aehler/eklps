-- MySQL dump 10.13  Distrib 5.5.31, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: eklps
-- ------------------------------------------------------
-- Server version	5.5.31-0+wheezy1-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `race_conditions`
--

DROP TABLE IF EXISTS `race_conditions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `race_conditions` (
  `id_races` bigint(20) NOT NULL,
  `class` varchar(2) DEFAULT NULL,
  `class_eq` varchar(2) DEFAULT NULL,
  `time` varchar(5) DEFAULT NULL,
  `time_eq` varchar(2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `race_conditions`
--

LOCK TABLES `race_conditions` WRITE;
/*!40000 ALTER TABLE `race_conditions` DISABLE KEYS */;
INSERT INTO `race_conditions` VALUES (2014111409,'','','',''),(2014111406,'','','',''),(2014111401,'E','=','',''),(2014111402,'D','=','',''),(2014111403,'C','=','',''),(2014111408,'E','<=','',''),(2014111407,'F','<=','',''),(2014111404,'E','<=','',''),(2014111405,'C','<=','',''),(2014111410,'','','',''),(2014111514,'','','',''),(2014111510,'','','',''),(2014111507,'','','',''),(2014111516,'','','',''),(2014111501,'D','=','',''),(2014111502,'D','=','',''),(2014111503,'C','=','',''),(2014111504,'B','=','',''),(2014111513,'F','<=','',''),(2014111511,'B','<=','',''),(2014111512,'E','<=','',''),(2014111509,'C','<=','',''),(2014111508,'C','<=','',''),(2014111515,'G','<=','',''),(2014111505,'C','<=','',''),(2014111506,'D','<=','',''),(2014111517,'','','',''),(2014111612,'','','',''),(2014111609,'','','',''),(2014111615,'','','',''),(2014111601,'D','>=','',''),(2014111602,'D','>=','',''),(2014111603,'D','>=','',''),(2014111604,'D','>=','',''),(2014111605,'D','>=','',''),(2014111608,'E','<=','',''),(2014111611,'C','<=','',''),(2014111610,'D','<=','',''),(2014111613,'D','<=','',''),(2014111607,'E','<=','',''),(2014111606,'C','<=','',''),(2014111614,'','','',''),(2014111711,'','','',''),(2014111712,'','','',''),(2014111710,'','','',''),(2014111706,'','','',''),(2014111701,'E','=','',''),(2014111702,'F','<=','',''),(2014111703,'B','=','',''),(2014111707,'E','<=','',''),(2014111708,'B','<=','',''),(2014111704,'C','<=','',''),(2014111705,'D','<=','',''),(2014111709,'','','',''),(2014111811,'','','',''),(2014111808,'','','',''),(2014111810,'','','',''),(2014111809,'','','',''),(2014111801,'A','>=','',''),(2014111802,'B','=','',''),(2014111803,'A','>=','',''),(2014111804,'A','>=','',''),(2014111812,'C','<=','',''),(2014111805,'F','<=','',''),(2014111806,'F','<=','',''),(2014111807,'','','',''),(2014111901,'C','=','',''),(2014111902,'A','>=','',''),(2014111903,'A','>=','',''),(2014111909,'D','<=','',''),(2014111906,'C','<=','',''),(2014111907,'B','<=','',''),(2014111908,'E','<=','',''),(2014111904,'C','<=','',''),(2014111905,'C','<=','',''),(2014112001,'C','=','',''),(2014112002,'B','=','',''),(2014112003,'B','=','',''),(2014112006,'C','<=','',''),(2014112004,'C','<=','',''),(2014112005,'E','<=','',''),(2014112101,'F','<=','',''),(2014112102,'C','=','',''),(2014112103,'A','>=','',''),(2014112105,'D','<=','',''),(2014112104,'G','<=','','');
/*!40000 ALTER TABLE `race_conditions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `races`
--

DROP TABLE IF EXISTS `races`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `races` (
  `id` bigint(20) NOT NULL,
  `date` int(11) NOT NULL,
  `conditions` varchar(255) DEFAULT NULL,
  `distance` int(11) DEFAULT NULL,
  `sc` tinyint(4) DEFAULT NULL,
  `age_conditions` varchar(25) DEFAULT NULL,
  `sex` varchar(3) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `fullname` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `races`
--

LOCK TABLES `races` WRITE;
/*!40000 ALTER TABLE `races` DISABLE KEYS */;
INSERT INTO `races` VALUES (2014111401,20141114,'  E',1300,0,'2+ Ж','f','Гр.III Приз Красноярска','Гр.III Приз Красноярска, 1300 м, 2+yo Ж'),(2014111402,20141114,'  D',4400,0,'4+ Ж','m','Гр.II Подвиг Геракла','Гр.II Подвиг Геракла, 4400 м, 4+yo Ж'),(2014111403,20141114,'  C',2600,1,'3+ К','f','Гр.II Sessill Cup','Гр.II Sessill Cup, ст-з, 2600 м, 3+yo К'),(2014111404,20141114,' максимальный  E и ниже',1000,1,'2+','all','Золотой класс','Золотой класс, торф, 1000 м, 2+yo'),(2014111405,20141114,' максимальный  C и ниже',3000,1,'3 К','f','Золотой класс','Золотой класс, ст-з, 3000 м, 3yo К'),(2014111406,20141114,'',4000,0,'4+','all','Тестовый класс','Тестовый класс, 4000 м, 4+yo'),(2014111407,20141114,' максимальный  F и ниже',2400,1,'3+','all','Бронзовый класс','Бронзовый класс, ст-з, 2400 м, 3+yo'),(2014111408,20141114,' максимальный  E и ниже',1300,0,'2+','all','Бронзовый класс','Бронзовый класс, 1300 м, 2+yo'),(2014111409,20141114,'',1100,0,'2+','all','Тестовый класс','Тестовый класс, 1100 м, 2+yo'),(2014111410,20141114,'',1000,0,'2+ К','f','Платиновый класс','Платиновый класс, 1000 м, 2+yo К'),(2014111501,20141115,'  D',2600,0,'3+ Ж','m','Гр.II Machinehead Stakes','Гр.II Machinehead Stakes, 2600 м, 3+yo Ж'),(2014111502,20141115,'  D',1400,1,'2+ Ж','m','Гр.II Парижские Огни','Гр.II Парижские Огни, ст-з, 1400 м, 2+yo Ж'),(2014111503,20141115,'  C',3400,1,'4+ К','f','Гр.II Neon Night Run','Гр.II Neon Night Run, ст-з, 3400 м, 4+yo К'),(2014111504,20141115,'  B',1200,0,'2+ К','f','Гр.I The Golden Flames Sprint','Гр.I The Golden Flames Sprint, 1200 м, 2+yo К'),(2014111505,20141115,' максимальный  C и ниже',1900,0,'2 К','f','Золотой класс','Золотой класс, 1900 м, 2yo К'),(2014111506,20141115,' максимальный  D и ниже',2000,0,'3 К','f','Золотой класс','Золотой класс, 2000 м, 3yo К'),(2014111507,20141115,'',4400,0,'4+','all','Тестовый класс','Тестовый класс, 4400 м, 4+yo'),(2014111508,20141115,' максимальный  C и ниже',2000,1,'3+','all','Медный класс','Медный класс, ст-з, 2000 м, 3+yo'),(2014111509,20141115,' максимальный  C и ниже',1800,1,'4+','all','Медный класс','Медный класс, ст-з, 1800 м, 4+yo'),(2014111510,20141115,'',1400,0,'2+','all','Тестовый класс','Тестовый класс, 1400 м, 2+yo'),(2014111511,20141115,' максимальный  B и ниже',1400,0,'4+','all','Медный класс','Медный класс, 1400 м, 4+yo'),(2014111512,20141115,' максимальный  E и ниже',1600,0,'2+','all','Медный класс','Медный класс, 1600 м, 2+yo'),(2014111513,20141115,' максимальный  F и ниже',1400,0,'2+','all','Медный класс','Медный класс, 1400 м, 2+yo'),(2014111514,20141115,'',1000,0,'2+','all','Тестовый класс','Тестовый класс, 1000 м, 2+yo'),(2014111515,20141115,' максимальный  G и ниже',2800,0,'3+','all','Серебряный класс','Серебряный класс, 2800 м, 3+yo'),(2014111516,20141115,'',1200,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1200 м, 2+yo'),(2014111517,20141115,'',4800,0,'4+ Ж','m','Платиновый класс','Платиновый класс, 4800 м, 4+yo Ж'),(2014111601,20141116,' минимальный  D и выше',2600,0,'3+ К','f','AA Karoline Trials','AA Karoline Trials, 2600 м, 3+yo К'),(2014111602,20141116,' минимальный  D и выше',1200,1,'3+ К','f','AA Oxbow Cup','AA Oxbow Cup, ст-з, 1200 м, 3+yo К'),(2014111603,20141116,' минимальный  D и выше',1600,1,'2 Ж','m','AA Royal Mile Stakes','AA Royal Mile Stakes, торф, 1600 м, 2yo Ж'),(2014111604,20141116,' минимальный  D и выше',1700,1,'3+ К','f','AA Queens Stakes','AA Queens Stakes, ст-з, 1700 м, 3+yo К'),(2014111605,20141116,' минимальный  D и выше',3200,1,'4+ Ж','m','AA Grandness of Waterfall Cup','AA Grandness of Waterfall Cup, ст-з, 3200 м, 4+yo Ж'),(2014111606,20141116,' максимальный  C и ниже',1900,0,'2','all','Золотой класс','Золотой класс, 1900 м, 2yo'),(2014111607,20141116,' максимальный  E и ниже',1400,0,'2+','all','Золотой класс','Золотой класс, 1400 м, 2+yo'),(2014111608,20141116,' максимальный  E и ниже',1100,0,'2+','all','Медный класс','Медный класс, 1100 м, 2+yo'),(2014111609,20141116,'',2200,0,'3+','all','Тестовый класс','Тестовый класс, 2200 м, 3+yo'),(2014111610,20141116,' максимальный  D и ниже',1800,0,'4+','all','Медный класс','Медный класс, 1800 м, 4+yo'),(2014111611,20141116,' максимальный  C и ниже',1600,0,'4+','all','Медный класс','Медный класс, 1600 м, 4+yo'),(2014111612,20141116,'',1600,0,'2+','all','Тестовый класс','Тестовый класс, 1600 м, 2+yo'),(2014111613,20141116,' максимальный  D и ниже',1100,1,'4+','all','Медный класс','Медный класс, ст-з, 1100 м, 4+yo'),(2014111614,20141116,'',2000,0,'3+ К','f','Платиновый класс','Платиновый класс, 2000 м, 3+yo К'),(2014111615,20141116,'',2400,1,'3+','all','Тестовый класс','Тестовый класс, ст-з, 2400 м, 3+yo'),(2014111701,20141117,'  E',3000,0,'3+ К','f','Гр.III Saroque Stakes','Гр.III Saroque Stakes, 3000 м, 3+yo К'),(2014111702,20141117,' максимальный  F и ниже',1900,1,'2+ К','f','Гр.III Young Ladies Cup','Гр.III Young Ladies Cup, ст-з, 1900 м, 2+yo К'),(2014111703,20141117,'  B',1200,1,'2+ Ж','m','Гр.I Вечерняя Москва','Гр.I Вечерняя Москва, ст-з, 1200 м, 2+yo Ж'),(2014111704,20141117,' максимальный  C и ниже',1900,0,'2,3 Ж','m','Золотой класс','Золотой класс, 1900 м, 2,3yo Ж'),(2014111705,20141117,' максимальный  D и ниже',1900,1,'3+','all','Золотой класс','Золотой класс, ст-з, 1900 м, 3+yo'),(2014111706,20141117,'',1600,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1600 м, 2+yo'),(2014111707,20141117,' максимальный  E и ниже',1400,0,'3+','all','Медный класс','Медный класс, 1400 м, 3+yo'),(2014111708,20141117,' максимальный  B и ниже',2400,1,'4+','all','Медный класс','Медный класс, ст-з, 2400 м, 4+yo'),(2014111709,20141117,'',2000,1,'3+ К','f','Платиновый класс','Платиновый класс, ст-з, 2000 м, 3+yo К'),(2014111710,20141117,'',1100,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1100 м, 2+yo'),(2014111711,20141117,'',2000,0,'3+','all','Тестовый класс','Тестовый класс, 2000 м, 3+yo'),(2014111712,20141117,'',2400,0,'3+','all','Тестовый класс','Тестовый класс, 2400 м, 3+yo'),(2014111801,20141118,' минимальный  A и выше',1600,0,'2+ К','f','Гр.I Ramada Stakes','Гр.I Ramada Stakes, 1600 м, 2+yo К'),(2014111802,20141118,'  B',4000,0,'4+ Ж','m','Гр.I Space Stakes','Гр.I Space Stakes, 4000 м, 4+yo Ж'),(2014111803,20141118,' минимальный  A и выше',2200,1,'3+ К','f','Гр.I Приз Императрицы','Гр.I Приз Императрицы, ст-з, 2200 м, 3+yo К'),(2014111804,20141118,' минимальный  A и выше',4800,1,'4+ Ж','m','Гр.I Scandinavia Pride','Гр.I Scandinavia Pride, ст-з, 4800 м, 4+yo Ж'),(2014111805,20141118,' максимальный  F и ниже',1200,0,'2,3','all','Золотой класс','Золотой класс, 1200 м, 2,3yo'),(2014111806,20141118,' максимальный  F и ниже',1800,0,'2+','all','Золотой класс','Золотой класс, 1800 м, 2+yo'),(2014111807,20141118,'',1200,0,'2+ К','f','Платиновый класс','Платиновый класс, 1200 м, 2+yo К'),(2014111808,20141118,'',2600,0,'3+','all','Тестовый класс','Тестовый класс, 2600 м, 3+yo'),(2014111809,20141118,'',2800,1,'3+','all','Тестовый класс','Тестовый класс, ст-з, 2800 м, 3+yo'),(2014111810,20141118,'',2800,0,'3+','all','Тестовый класс','Тестовый класс, 2800 м, 3+yo'),(2014111811,20141118,'',1800,0,'2+','all','Тестовый класс','Тестовый класс, 1800 м, 2+yo'),(2014111812,20141118,' максимальный  C и ниже',1100,0,'4+','all','Медный класс','Медный класс, 1100 м, 4+yo'),(2014111901,20141119,'  C',1600,1,'2+ К','f','Гр.II Turbulent Chase','Гр.II Turbulent Chase, ст-з, 1600 м, 2+yo К'),(2014111902,20141119,' минимальный  A и выше',1800,0,'2+ Ж','m','Гр.I Yung Yong Plate','Гр.I Yung Yong Plate, 1800 м, 2+yo Ж'),(2014111903,20141119,' минимальный  A и выше',3200,0,'4+ Ж','m','Гр.I Luna Plate','Гр.I Luna Plate, 3200 м, 4+yo Ж'),(2014111904,20141119,' максимальный  C и ниже',1800,0,'3+','all','Золотой класс','Золотой класс, 1800 м, 3+yo'),(2014111905,20141119,' максимальный  C и ниже',2200,0,'3+ К','f','Золотой класс','Золотой класс, 2200 м, 3+yo К'),(2014111906,20141119,' максимальный  C и ниже',1700,0,'2+','all','Медный класс','Медный класс, 1700 м, 2+yo'),(2014111907,20141119,' максимальный  B и ниже',1900,0,'4+','all','Медный класс','Медный класс, 1900 м, 4+yo'),(2014111908,20141119,' максимальный  E и ниже',1800,1,'3+','all','Медный класс','Медный класс, ст-з, 1800 м, 3+yo'),(2014111909,20141119,' максимальный  D и ниже',1600,0,'3+','all','Медный класс','Медный класс, 1600 м, 3+yo'),(2014112001,20141120,'  C',1000,0,'2+ К','f','Гр.II Земляничная Поляна','Гр.II Земляничная Поляна, 1000 м, 2+yo К'),(2014112002,20141120,'  B',2600,0,'3+ Ж','m','Гр.I Золотая Лихорадка','Гр.I Золотая Лихорадка, 2600 м, 3+yo Ж'),(2014112003,20141120,'  B',2400,1,'3+ К','f','Гр.I Ladies First Stakes','Гр.I Ladies First Stakes, ст-з, 2400 м, 3+yo К'),(2014112004,20141120,' максимальный  C и ниже',1400,0,'4+ К','f','Золотой класс','Золотой класс, 1400 м, 4+yo К'),(2014112005,20141120,' максимальный  E и ниже',2400,0,'3+','all','Золотой класс','Золотой класс, 2400 м, 3+yo'),(2014112006,20141120,' максимальный  C и ниже',1400,0,'2,3','all','Серебряный класс','Серебряный класс, 1400 м, 2,3yo'),(2014112101,20141121,' максимальный  F и ниже',2000,1,'3+ Ж','m','Гр.III Loose Me Not Run','Гр.III Loose Me Not Run, ст-з, 2000 м, 3+yo Ж'),(2014112102,20141121,'  C',1100,1,'2+ К','f','Гр.II Приз Двойная Звезда','Гр.II Приз Двойная Звезда, ст-з, 1100 м, 2+yo К'),(2014112103,20141121,' минимальный  A и выше',2400,0,'3+ Ж','m','Гр.I In The Next Life','Гр.I In The Next Life, 2400 м, 3+yo Ж'),(2014112104,20141121,' максимальный  G и ниже',1900,1,'2+','all','Золотой класс','Золотой класс, торф, 1900 м, 2+yo'),(2014112105,20141121,' максимальный  D и ниже',1300,0,'2+','all','Золотой класс','Золотой класс, 1300 м, 2+yo');
/*!40000 ALTER TABLE `races` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2014-11-12 17:41:24
-- MySQL dump 10.13  Distrib 5.5.31, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: eklps
-- ------------------------------------------------------
-- Server version	5.5.31-0+wheezy1-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `race_conditions`
--

DROP TABLE IF EXISTS `race_conditions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `race_conditions` (
  `id_races` bigint(20) NOT NULL,
  `class` varchar(2) DEFAULT NULL,
  `class_eq` varchar(2) DEFAULT NULL,
  `time` varchar(5) DEFAULT NULL,
  `time_eq` varchar(2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `race_conditions`
--

LOCK TABLES `race_conditions` WRITE;
/*!40000 ALTER TABLE `race_conditions` DISABLE KEYS */;
INSERT INTO `race_conditions` VALUES (2014122122,'','','',''),(2014122121,'','','',''),(2014122120,'','','',''),(2014122123,'','','',''),(2014122124,'','','',''),(2014122101,'D','>=','',''),(2014122102,'D','>=','',''),(2014122103,'D','>=','',''),(2014122104,'D','>=','',''),(2014122105,'D','>=','',''),(2014122106,'D','>=','',''),(2014122107,'D','>=','',''),(2014122108,'D','>=','',''),(2014122109,'D','>=','',''),(2014122110,'D','>=','',''),(2014122111,'D','>=','',''),(2014122112,'D','>=','',''),(2014122113,'D','>=','',''),(2014122114,'D','>=','',''),(2014122115,'D','>=','',''),(2014122116,'D','>=','',''),(2014122119,'G','<=','',''),(2014122117,'G','<=','',''),(2014122118,'F','<=','',''),(2014122208,'','','',''),(2014122209,'','','',''),(2014122201,'F','<=','',''),(2014122202,'F','<=','',''),(2014122203,'D','=','',''),(2014122210,'C','<=','',''),(2014122206,'C','<=','',''),(2014122207,'E','<=','',''),(2014122205,'','','',''),(2014122204,'F','<=','',''),(2014122211,'','','',''),(2014122310,'','','',''),(2014122301,'F','<=','',''),(2014122302,'C','=','',''),(2014122303,'A','>=','',''),(2014122308,'D','<=','',''),(2014122307,'E','<=','',''),(2014122304,'B','<=','',''),(2014122305,'D','<=','',''),(2014122309,'','','',''),(2014122306,'','','',''),(2014122408,'','','',''),(2014122401,'D','=','',''),(2014122402,'B','=','',''),(2014122403,'B','=','',''),(2014122407,'G','<=','',''),(2014122409,'E','<=','',''),(2014122410,'C','<=','',''),(2014122411,'D','<=','',''),(2014122406,'B','<=','',''),(2014122404,'D','<=','',''),(2014122405,'F','<=','',''),(2014122501,'E','=','',''),(2014122502,'B','=','',''),(2014122503,'B','=','',''),(2014122504,'A','>=','',''),(2014122507,'D','<=','',''),(2014122508,'F','<=','',''),(2014122505,'G','<=','',''),(2014122506,'E','<=','',''),(2014122601,'E','=','',''),(2014122602,'D','=','',''),(2014122603,'B','=','',''),(2014122605,'B','<=','',''),(2014122604,'F','<=','',''),(2014122701,'E','=','',''),(2014122702,'F','<=','',''),(2014122703,'F','<=','',''),(2014122705,'C','<=','',''),(2014122704,'C','<=','','');
/*!40000 ALTER TABLE `race_conditions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `races`
--

DROP TABLE IF EXISTS `races`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `races` (
  `id` bigint(20) NOT NULL,
  `date` int(11) NOT NULL,
  `conditions` varchar(255) DEFAULT NULL,
  `distance` int(11) DEFAULT NULL,
  `sc` tinyint(4) DEFAULT NULL,
  `age_conditions` varchar(25) DEFAULT NULL,
  `sex` varchar(3) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `fullname` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `races`
--

LOCK TABLES `races` WRITE;
/*!40000 ALTER TABLE `races` DISABLE KEYS */;
INSERT INTO `races` VALUES (2014122101,20141221,' минимальный  D и выше',1100,0,'4+ (SprTC-2)','all','TC Oakleigh Plate','TC Oakleigh Plate, 1100 м, 4+yo (SprTC-2)'),(2014122102,20141221,' минимальный  D и выше',1300,0,'2 К (2FTC-2)','f','TC Matron Stakes','TC Matron Stakes, 1300 м, 2yo К (2yoFTC-2)'),(2014122103,20141221,' минимальный  D и выше',1600,0,'2 Ж (2TC-2)','m','TC Futurity Stakes','TC Futurity Stakes, 1600 м, 2yo Ж (2yoTC-2)'),(2014122104,20141221,' минимальный  D и выше',1600,0,'3 К (FTC-2)','f','TC Pimlico Oaks','TC Pimlico Oaks, 1600 м, 3yo К (FTC-2)'),(2014122105,20141221,' минимальный  D и выше',1800,0,'4+ (ClassTC-2)','all','TC Azur Plate','TC Azur Plate, 1800 м, 4+yo (ClassTC-2)'),(2014122106,20141221,' минимальный  D и выше',1900,0,'3 Ж (TC-2)','m','TC Pimlico Derby','TC Pimlico Derby, 1900 м, 3yo Ж (TC-2)'),(2014122107,20141221,' минимальный  D и выше',2000,0,'3 (TrTC-2)','all','TC Nakheel Cup','TC Nakheel Cup, 2000 м, 3yo (TrTC-2)'),(2014122108,20141221,' минимальный  D и выше',4000,0,'4+ (EndTC-2)','all','TC Albatros Plate','TC Albatros Plate, 4000 м, 4+yo (EndTC-2)'),(2014122109,20141221,' минимальный  D и выше',1000,1,'4+ (MelbSCTC-2)','all','TC Manikato Chase','TC Manikato Chase, ст-з, 1000 м, 4+yo (MelbSCTC-2)'),(2014122110,20141221,' минимальный  D и выше',1300,1,'2 К (2FTTC-2)','f','TC Pyramisa Stakes','TC Pyramisa Stakes, торф, 1300 м, 2yo К (2yoFTTC-2)'),(2014122111,20141221,' минимальный  D и выше',1600,1,'2 Ж (2TTC-2)','m','TC Ritz Stakes','TC Ritz Stakes, торф, 1600 м, 2yo Ж (2yoTTC-2)'),(2014122112,20141221,' минимальный  D и выше',1800,1,'3 К (FSCTC-2)','f','TC Filly Hurdle Challenge','TC Filly Hurdle Challenge, ст-з, 1800 м, 3yo К (FSCTC-2)'),(2014122113,20141221,' минимальный  D и выше',1900,1,'3 Ж (SCTC-2)','m','TC Hurdle Challenge','TC Hurdle Challenge, ст-з, 1900 м, 3yo Ж (SCTC-2)'),(2014122114,20141221,' минимальный  D и выше',2000,1,'3 (TrSCTC-2)','all','TC Grand Chase Cup','TC Grand Chase Cup, ст-з, 2000 м, 3yo (TrSCTC-2)'),(2014122115,20141221,' минимальный  D и выше',2400,1,'4+ (SCHTC-2)','all','TC Brooklyn SC Handicap','TC Brooklyn SC Handicap, ст-з, 2400 м, 4+yo (SCHTC-2)'),(2014122116,20141221,' минимальный  D и выше',4000,1,'4+ (SCEndTC-2)','all','TC Ambassador Chase','TC Ambassador Chase, ст-з, 4000 м, 4+yo (SCEndTC-2)'),(2014122117,20141221,' максимальный  G и ниже',1400,1,'2+','all','Золотой класс','Золотой класс, торф, 1400 м, 2+yo'),(2014122118,20141221,' максимальный  F и ниже',1800,1,'2+','all','Золотой класс','Золотой класс, торф, 1800 м, 2+yo'),(2014122119,20141221,' максимальный  G и ниже',1600,0,'3+','all','Медный класс','Медный класс, 1600 м, 3+yo'),(2014122120,20141221,'',1600,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1600 м, 2+yo'),(2014122121,20141221,'',1800,0,'2+','all','Тестовый класс','Тестовый класс, 1800 м, 2+yo'),(2014122122,20141221,'',1400,0,'2+','all','Тестовый класс','Тестовый класс, 1400 м, 2+yo'),(2014122123,20141221,'',1700,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1700 м, 2+yo'),(2014122124,20141221,'',1900,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1900 м, 2+yo'),(2014122201,20141222,' максимальный  F и ниже',1200,1,'2+ Ж','m','Гр.III Duke of Edinburgh Stakes','Гр.III Duke of Edinburgh Stakes, ст-з, 1200 м, 2+yo Ж'),(2014122202,20141222,' максимальный  F и ниже',4800,1,'4+ К','f','Гр.III Легенды Старого Леса','Гр.III Легенды Старого Леса, ст-з, 4800 м, 4+yo К'),(2014122203,20141222,'  D',2400,0,'3+ К','f','Гр.II Северная Звезда','Гр.II Северная Звезда, 2400 м, 3+yo К'),(2014122204,20141222,' максимальный  F и ниже',4800,1,'4+','all','Золотой класс','Золотой класс, ст-з, 4800 м, 4+yo'),(2014122205,20141222,'',1300,0,'4+','all','Золотой класс','Золотой класс, 1300 м, 4+yo'),(2014122206,20141222,' максимальный  C и ниже',1700,0,'3','all','Серебряный класс','Серебряный класс, 1700 м, 3yo'),(2014122207,20141222,' максимальный  E и ниже',2000,0,'3','all','Серебряный класс','Серебряный класс, 2000 м, 3yo'),(2014122208,20141222,'',1200,0,'2+','all','Тестовый класс','Тестовый класс, 1200 м, 2+yo'),(2014122209,20141222,'',2000,1,'3+','all','Тестовый класс','Тестовый класс, ст-з, 2000 м, 3+yo'),(2014122210,20141222,' максимальный  C и ниже',1400,0,'3+','all','Медный класс','Медный класс, 1400 м, 3+yo'),(2014122211,20141222,'',2200,0,'4+ К','f','Платиновый класс','Платиновый класс, 2200 м, 4+yo К'),(2014122301,20141223,' максимальный  F и ниже',1700,0,'2+ Ж','m','Гр.III Самый Стойкий','Гр.III Самый Стойкий, 1700 м, 2+yo Ж'),(2014122302,20141223,'  C',3400,0,'4+ Ж','m','Гр.II Mirovik Stakes','Гр.II Mirovik Stakes, 3400 м, 4+yo Ж'),(2014122303,20141223,' минимальный  A и выше',2400,1,'3+ Ж','m','Гр.I Grand Prix Du Mans','Гр.I Grand Prix Du Mans, ст-з, 2400 м, 3+yo Ж'),(2014122304,20141223,' максимальный  B и ниже',1100,0,'2+','all','Золотой класс','Золотой класс, 1100 м, 2+yo'),(2014122305,20141223,' максимальный  D и ниже',1200,0,'2+ К','f','Золотой класс','Золотой класс, 1200 м, 2+yo К'),(2014122306,20141223,'',1900,1,'4+','all','Платиновый класс','Платиновый класс, ст-з, 1900 м, 4+yo'),(2014122307,20141223,' максимальный  E и ниже',3000,0,'4','all','Серебряный класс','Серебряный класс, 3000 м, 4yo'),(2014122308,20141223,' максимальный  D и ниже',1000,0,'4+','all','Медный класс','Медный класс, 1000 м, 4+yo'),(2014122309,20141223,'',2400,0,'4+','all','Платиновый класс','Платиновый класс, 2400 м, 4+yo'),(2014122310,20141223,'',1400,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1400 м, 2+yo'),(2014122401,20141224,'  D',1200,0,'2+ Ж','m','Гр.II Приз Орбиты','Гр.II Приз Орбиты, 1200 м, 2+yo Ж'),(2014122402,20141224,'  B',1900,0,'2+ К','f','Гр.I Givenchy Cup','Гр.I Givenchy Cup, 1900 м, 2+yo К'),(2014122403,20141224,'  B',3400,1,'4+ Ж','m','Гр.I Sands Still Cup','Гр.I Sands Still Cup, ст-з, 3400 м, 4+yo Ж'),(2014122404,20141224,' максимальный  D и ниже',1800,0,'2+','all','Золотой класс','Золотой класс, 1800 м, 2+yo'),(2014122405,20141224,' максимальный  F и ниже',2800,1,'3+','all','Золотой класс','Золотой класс, ст-з, 2800 м, 3+yo'),(2014122406,20141224,' максимальный  B и ниже',1900,1,'4+','all','Медный класс','Медный класс, ст-з, 1900 м, 4+yo'),(2014122407,20141224,' максимальный  G и ниже',1000,0,'3+','all','Медный класс','Медный класс, 1000 м, 3+yo'),(2014122408,20141224,'',1800,1,'2+','all','Тестовый класс','Тестовый класс, ст-з, 1800 м, 2+yo'),(2014122409,20141224,' максимальный  E и ниже',1400,0,'2+','all','Медный класс','Медный класс, 1400 м, 2+yo'),(2014122410,20141224,' максимальный  C и ниже',1700,0,'4+','all','Медный класс','Медный класс, 1700 м, 4+yo'),(2014122411,20141224,' максимальный  D и ниже',1900,1,'3+','all','Медный класс','Медный класс, ст-з, 1900 м, 3+yo'),(2014122501,20141225,'  E',2600,1,'3+ К','f','Гр.III Sunday Stakes','Гр.III Sunday Stakes, ст-з, 2600 м, 3+yo К'),(2014122502,20141225,'  B',3200,0,'4+ К','f','Гр.I Pine Lane Stakes','Гр.I Pine Lane Stakes, 3200 м, 4+yo К'),(2014122503,20141225,'  B',1600,1,'2+ Ж','m','Гр.I Black Onyx Stakes','Гр.I Black Onyx Stakes, ст-з, 1600, 2+yo Ж'),(2014122504,20141225,' минимальный  A и выше',4400,1,'4+ Ж','m','Гр.I Great Taxis Cup','Гр.I Great Taxis Cup, ст-з, 4400 м, 4+yo Ж'),(2014122505,20141225,' максимальный  G и ниже',1000,0,'2+','all','Золотой класс','Золотой класс, 1000 м, 2+yo'),(2014122506,20141225,' максимальный  E и ниже',3000,1,'3+','all','Золотой класс','Золотой класс, ст-з, 3000 м, 3+yo'),(2014122507,20141225,' максимальный  D и ниже',1600,0,'2+','all','Медный класс','Медный класс, 1600 м, 2+yo'),(2014122508,20141225,' максимальный  F и ниже',1000,1,'2+','all','Медный класс','Медный класс, торф, 1000 м, 2+yo'),(2014122601,20141226,'  E',3200,1,'4+ Ж','m','Гр.III Lamodar Stakes','Гр.III Lamodar Stakes, ст-з, 3200 м, 4+yo Ж'),(2014122602,20141226,'  D',2800,1,'3+ К','f','Гр.II Краса Страны','Гр.II Краса Страны, ст-з, 2800 м, 3+yo К'),(2014122603,20141226,'  B',1000,1,'2+ К','f','Гр.I Lis-de-Fleur Stakes','Гр.I Lis-de-Fleur Stakes, ст-з, 1000 м, 2+yo К'),(2014122604,20141226,' максимальный  F и ниже',1700,1,'2+','all','Золотой класс','Золотой класс, торф, 1700 м, 2+yo'),(2014122605,20141226,' максимальный  B и ниже',2200,0,'4+ К','f','Золотой класс','Золотой класс, 2200 м, 4+yo К'),(2014122701,20141227,'  E',1800,0,'2+ К','f','Гр.III Приз Констанции','Гр.III Приз Констанции, 1800 м, 2+yo К'),(2014122702,20141227,' максимальный  F и ниже',1800,1,'2+ Ж','m','Гр.III Ледяное Сердце','Гр.III Ледяное Сердце, ст-з, 1800 м, 2+yo Ж'),(2014122703,20141227,' максимальный  F и ниже',4800,1,'4+ К','f','Гр.III Каникулы в Париже','Гр.III Каникулы в Париже, ст-з, 4800 м, 4+yo К'),(2014122704,20141227,' максимальный  C и ниже',2000,0,'4+ К','f','Золотой класс','Золотой класс, 2000 м, 4+yo К'),(2014122705,20141227,' максимальный  C и ниже',1900,0,'2,3 К','f','Золотой класс','Золотой класс, 1900 м, 2,3yo К');
/*!40000 ALTER TABLE `races` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `season`
--

DROP TABLE IF EXISTS `season`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `season` (
  `season` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `season`
--

LOCK TABLES `season` WRITE;
/*!40000 ALTER TABLE `season` DISABLE KEYS */;
INSERT INTO `season` VALUES (2015);
/*!40000 ALTER TABLE `season` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2014-12-19 10:54:52
