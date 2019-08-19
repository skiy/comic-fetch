### 站点标识变更
ALTER TABLE `tb_books` CHANGE `origin_web_type` `origin_web_type` VARCHAR(24) NOT NULL DEFAULT 'pc' COMMENT '站点类型(pc, mobile, api)';
