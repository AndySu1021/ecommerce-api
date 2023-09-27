CREATE TABLE IF NOT EXISTS `product_category`
(
    id          bigint unsigned auto_increment primary key,
    merchant_id bigint unsigned not null comment '商戶ID',
    name        varchar(30)     not null comment '名稱',
    top_id      bigint unsigned not null comment '頂層ID',
    parent_id   bigint unsigned not null comment '父級ID',
    tree_left   bigint unsigned not null comment '代理數節點左編號',
    tree_right  bigint unsigned not null comment '代理數節點右編號',
    created_at  datetime        not null comment '創建時間',
    updated_at  datetime        not null comment '更新時間'
) COMMENT ='商品分類資料表' COLLATE utf8mb4_general_ci;
