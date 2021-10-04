-- MariaDB dump 10.17  Distrib 10.4.7-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: store
-- ------------------------------------------------------
-- Server version	10.4.7-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `store`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `store` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `store`;

--
-- Table structure for table `book`
--

DROP TABLE IF EXISTS `book`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `book` (
  `book_id` int(11) NOT NULL AUTO_INCREMENT,
  `book_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `book_price` int(11) NOT NULL DEFAULT 0,
  `book_isbn` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_amount` int(11) NOT NULL DEFAULT 0,
  `book_detail` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `book`
--

LOCK TABLES `book` WRITE;
/*!40000 ALTER TABLE `book` DISABLE KEYS */;
INSERT INTO `book` VALUES (1,'The Man Who Changed China','Robert Lawrence Kuhn',4800,'7-5327-3654-7518',1926,'The Man Who Changed China: The Life and Legacy of Jiang Zemin a biography of former Chinese leader Jiang Zemin by Robert Lawrence Kuhn was published in 2005, in English and Chinese. It has generated both attention and controversy.'),(2,'Nineteen Eighty-Four','George Orwell',1950,'9-7875-4471-1647',64,'Nineteen Eighty-Four, often published as 1984, is a dystopian novel by English writer George Orwell published in June 1949. The novel is set in the year 1984 when most of the world population have become victims of perpetual war, omnipresent government surveillance and propaganda.'),(3,'No Longer Human','Osamu Dazai',2500,'9-7875-0638-0263',89,'No Longer Human is a Japanese novel by Osamu Dazai. Published after Run Melos and The Setting Sun, No Longer Human is considered Dazai\'s masterpiece and ranks as the second-best selling novel in Japan, behind Natsume Sōseki\'s Kokoro.');
/*!40000 ALTER TABLE `book` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `flag`
--

DROP TABLE IF EXISTS `flag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `flag` (
  `flag` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `flag`
--

LOCK TABLES `flag` WRITE;
/*!40000 ALTER TABLE `flag` DISABLE KEYS */;
INSERT INTO `flag` VALUES ('I move the flag to another place 23333...');
/*!40000 ALTER TABLE `flag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `flag_427`
--

DROP TABLE IF EXISTS `flag_427`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `flag_427` (
  `flag` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `flag_427`
--

LOCK TABLES `flag_427` WRITE;
/*!40000 ALTER TABLE `flag_427` DISABLE KEYS */;
INSERT INTO `flag_427` VALUES ('0ops{sample_flag}');
/*!40000 ALTER TABLE `flag_427` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'SHA256',
  `user_email` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_status` int(11) NOT NULL DEFAULT 0 COMMENT '0 for enabled',
  `user_permission` int(11) NOT NULL DEFAULT 0 COMMENT '0 for user',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`),
  UNIQUE KEY `user_email` (`user_email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin','384fde3636e6e01e0194d2976d8f26410af3e846e573379cb1a09e2f0752d8cc','me@imwxz.com',0,1),(2,'user','8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92','test@user.com',0,0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-10-07 21:30:54
