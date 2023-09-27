CREATE TABLE IF NOT EXISTS `merchant`
(
    id             bigint unsigned auto_increment primary key,
    `name`         varchar(255)      not null comment '名稱',
    `code`         varchar(255)      not null comment '編號',
    `host`         varchar(255)      not null comment '域名',
    encrypt_salt   varchar(255)      not null comment '加密鹽',
    is_enabled     tinyint default 1 not null comment '是否啟用 0否 1是',
    created_at     datetime          not null comment '創建時間',
    updated_at     datetime          not null comment '更新時間',
    constraint idx_host
        unique (`host`)
) COMMENT ='商戶資料表' COLLATE utf8mb4_general_ci;
