-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2019-08-02 17:42:02
-- 服务器版本： 5.7.26-log
-- PHP 版本： 7.3.6

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

--
-- 数据库： `comic`
--

-- --------------------------------------------------------

--
-- 表的结构 `tb_books`
--

CREATE TABLE `tb_books` (
  `id` int(11) NOT NULL COMMENT '编号',
  `name` varchar(52) NOT NULL COMMENT '漫画名',
  `image_url` varchar(255) DEFAULT NULL COMMENT '漫画图标',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(0正在更新,1暂停更新,2完结)',
  `origin_url` varchar(255) NOT NULL COMMENT '采集地址',
  `origin_web` varchar(48) NOT NULL COMMENT '源站名',
  `origin_web_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '站点类型(0.pc, 1.mobile, 3.api)',
  `origin_flag` varchar(24) NOT NULL COMMENT '源站标识',
  `origin_image_url` varchar(255) DEFAULT NULL COMMENT '源站漫画图标地址',
  `origin_path_url` varchar(255) NOT NULL COMMENT '上次采集图片实际路径',
  `origin_book_id` int(11) NOT NULL COMMENT '源站漫画ID',
  `updated_at` int(11) NOT NULL COMMENT '最后更新时间',
  `created_at` int(11) NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漫画名';

--
-- 转存表中的数据 `tb_books`
--

INSERT INTO `tb_books` (`id`, `name`, `image_url`, `status`, `origin_url`, `origin_web`, `origin_web_type`, `origin_flag`, `origin_image_url`, `origin_path_url`, `origin_book_id`, `updated_at`, `created_at`) VALUES
(1, '白莲妖姬 ', '', 0, 'https://www.manhuaniu.com/manhua/11684/', '漫画牛', 0, 'manhuaniu', 'https://res.nbhbzl.com/images/cover/201907/1564358077LnafsebVBintTlPE.jpg!cover-400', '', 11684, 1552496203, 1552494553);

-- --------------------------------------------------------

--
-- 表的结构 `tb_chapters`
--

CREATE TABLE `tb_chapters` (
  `id` int(11) NOT NULL COMMENT '编号',
  `bid` int(11) NOT NULL COMMENT '漫画编号',
  `chapter_id` int(11) NOT NULL COMMENT '章节ID(话)',
  `title` varchar(52) NOT NULL COMMENT '标题',
  `order_id` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `origin_id` int(11) NOT NULL COMMENT '源章节ID',
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
  `size` int(11) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `order_id` int(11) NOT NULL COMMENT '排序',
  `is_remote` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否远程图片',
  `created_at` int(11) NOT NULL COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='漫画图库';

--
-- 转储表的索引
--

--
-- 表的索引 `tb_books`
--
ALTER TABLE `tb_books`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `tb_chapters`
--
ALTER TABLE `tb_chapters`
  ADD PRIMARY KEY (`id`),
  ADD KEY `cid` (`chapter_id`);

--
-- 表的索引 `tb_images`
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
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号', AUTO_INCREMENT=7;

--
-- 使用表AUTO_INCREMENT `tb_chapters`
--
ALTER TABLE `tb_chapters`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号', AUTO_INCREMENT=2177;

--
-- 使用表AUTO_INCREMENT `tb_images`
--
ALTER TABLE `tb_images`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号', AUTO_INCREMENT=13290;
COMMIT;

