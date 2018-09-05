-- phpMyAdmin SQL Dump
-- version 4.8.2
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: 2018-09-05 18:29:58
-- 服务器版本： 5.7.22
-- PHP Version: 7.2.7

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `comic`
--

-- --------------------------------------------------------

--
-- 表的结构 `tb_books`
--

CREATE TABLE `tb_books` (
  `id` int(11) NOT NULL COMMENT '编号',
  `name` varchar(52) NOT NULL COMMENT '漫画名',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(0正在更新,1暂停更新,2完结)',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `origin_web` varchar(48) NOT NULL COMMENT '源站名',
  `origin_book_id` int(11) NOT NULL COMMENT '源站漫画ID',
  `updated_at` int(11) NOT NULL COMMENT '最后更新时间',
  `created_at` int(11) NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漫画名';

-- --------------------------------------------------------

--
-- 表的结构 `tb_chapters`
--

CREATE TABLE `tb_chapters` (
  `id` int(11) NOT NULL COMMENT '编号',
  `bid` int(11) NOT NULL COMMENT '漫画编号',
  `chapter_id` int(11) NOT NULL COMMENT '章节ID(话)',
  `title` varchar(52) NOT NULL COMMENT '标题',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `created_at` int(11) NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='章节';

-- --------------------------------------------------------

--
-- 表的结构 `tb_images`
--

CREATE TABLE `tb_images` (
  `id` int(11) NOT NULL COMMENT '编号',
  `bid` int(11) NOT NULL COMMENT '漫画编号',
  `cid` int(11) NOT NULL COMMENT '章节编号',
  `chapter_id` int(11) NOT NULL COMMENT '章节ID',
  `image_url` varchar(255) NOT NULL COMMENT '图片地址',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `is_remote` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否远程图片',
  `created_at` int(11) NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漫画图库';

--
-- Indexes for dumped tables
--

--
-- Indexes for table `tb_books`
--
ALTER TABLE `tb_books`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `tb_chapters`
--
ALTER TABLE `tb_chapters`
  ADD PRIMARY KEY (`id`),
  ADD KEY `cid` (`chapter_id`);

--
-- Indexes for table `tb_images`
--
ALTER TABLE `tb_images`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `tb_books`
--
ALTER TABLE `tb_books`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号';

--
-- 使用表AUTO_INCREMENT `tb_chapters`
--
ALTER TABLE `tb_chapters`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号';

--
-- 使用表AUTO_INCREMENT `tb_images`
--
ALTER TABLE `tb_images`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号';
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
