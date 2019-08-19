-- MySQL dump 10.13  Distrib 5.7.22, for osx10.13 (x86_64)
--
-- Host: localhost    Database: comic
-- ------------------------------------------------------
-- Server version	5.7.22

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
-- Table structure for table `tb_books`
--

DROP TABLE IF EXISTS `tb_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tb_books` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(52) NOT NULL COMMENT '漫画名',
  `image_url` varchar(255) DEFAULT NULL COMMENT '漫画图标',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(0正在更新,1暂停更新,2完结)',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `origin_web` varchar(48) NOT NULL COMMENT '源站名',
  `origin_web_type` varchar(24) NOT NULL DEFAULT 'pc' COMMENT '站点类型(pc, mobile, api)',
  `origin_flag` varchar(24) NOT NULL COMMENT '源站标识',
  `origin_image_url` varchar(255) DEFAULT NULL COMMENT '源站漫画图标地址',
  `origin_path_url` varchar(255) NOT NULL COMMENT '上次采集图片实际路径',
  `origin_book_id` int(11) NOT NULL COMMENT '源站漫画ID',
  `updated_at` int(11) NOT NULL COMMENT '最后更新时间',
  `created_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COMMENT='漫画名';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tb_chapters`
--

DROP TABLE IF EXISTS `tb_chapters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tb_chapters` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `book_id` int(11) NOT NULL COMMENT '漫画ID',
  `episode_id` int(11) DEFAULT '0' COMMENT '第几话',
  `title` varchar(52) NOT NULL COMMENT '标题',
  `order_id` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `origin_id` int(11) NOT NULL COMMENT '源章节ID',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 (0.采集成功, 1.采集失败, 2. 停止采集)',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `cid` (`episode_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2772 DEFAULT CHARSET=utf8mb4 COMMENT='章节';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tb_images`
--

DROP TABLE IF EXISTS `tb_images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tb_images` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `book_id` int(11) NOT NULL COMMENT '漫画ID',
  `chapter_id` int(11) NOT NULL COMMENT '章节ID',
  `episode_id` int(11) NOT NULL DEFAULT '0' COMMENT '第几话',
  `image_url` varchar(255) NOT NULL COMMENT '图片地址',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `order_id` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `is_remote` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否远程图片',
  `created_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5916 DEFAULT CHARSET=utf8mb4 COMMENT='漫画图库';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-08-19 21:23:26
