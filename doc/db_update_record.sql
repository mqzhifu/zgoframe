ALTER TABLE `user_reg` ADD `province` INT NOT NULL DEFAULT '0' COMMENT '省' ;
ALTER TABLE `user_reg` ADD `City` INT NOT NULL DEFAULT '0' COMMENT '市' ;
ALTER TABLE `user_reg` ADD `County` INT NOT NULL DEFAULT '0' COMMENT '县' ;
ALTER TABLE `user_reg` ADD `Town` INT NOT NULL DEFAULT '0' COMMENT '镇' ;
ALTER TABLE `user_reg` ADD `area_detail` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '详细地址' ;

ALTER TABLE `user_login` ADD `province` INT NOT NULL DEFAULT '0' COMMENT '省' ;
ALTER TABLE `user_login` ADD `City` INT NOT NULL DEFAULT '0' COMMENT '市' ;
ALTER TABLE `user_login` ADD `County` INT NOT NULL DEFAULT '0' COMMENT '县' ;
ALTER TABLE `user_login` ADD `Town` INT NOT NULL DEFAULT '0' COMMENT '镇' ;
ALTER TABLE `user_login` ADD `area_detail` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '详细地址' ;




user_reg os int =>String