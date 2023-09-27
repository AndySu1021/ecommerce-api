CREATE TABLE IF NOT EXISTS `product_spec`
(
    id          bigint unsigned auto_increment primary key,
    merchant_id bigint unsigned not null comment '商戶ID',
    product_id  bigint unsigned not null comment '商品ID',
    `level`     tinyint         not null comment '規格層級',
    type        tinyint         not null comment '類型 1規格標題 2規格選項',
    `name`      varchar(30)     not null comment '名稱',
    created_at  datetime        not null comment '創建時間',
    updated_at  datetime        not null comment '更新時間',
    constraint uidx_product_level_name
        unique (product_id, `level`, `name`)
) COMMENT ='商品規格資料表' COLLATE utf8mb4_general_ci;
