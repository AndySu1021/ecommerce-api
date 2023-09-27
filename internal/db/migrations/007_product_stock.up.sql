CREATE TABLE IF NOT EXISTS `product_stock`
(
    id          bigint unsigned auto_increment primary key,
    merchant_id bigint unsigned not null comment '商戶ID',
    product_id  bigint unsigned not null comment '商品ID',
    spec_1_id   bigint unsigned not null comment '第一層規格ID',
    spec_2_id   bigint unsigned not null comment '第二層規格ID',
    quantity    int             not null comment '庫存數量',
    `code`      varchar(255)    not null comment 'SKU 貨號',
    created_at  datetime        not null comment '創建時間',
    updated_at  datetime        not null comment '更新時間',
    constraint uidx_product_spec
        unique (product_id, spec_1_id, spec_2_id)
) COMMENT ='商品庫存資料表' COLLATE utf8mb4_general_ci;
