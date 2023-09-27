CREATE TABLE IF NOT EXISTS `member`
(
    id              bigint unsigned auto_increment primary key,
    merchant_id     bigint unsigned         not null comment '商戶ID',
    email           varchar(255)            not null comment '信箱',
    password        char(32)                not null comment '密碼',
    real_name       varchar(30)  default '' not null comment '真實姓名',
    mobile          varchar(30)  default '' not null comment '手機號',
    sex             tinyint      default 1  not null comment '性別 1男 2女',
    birthday        date                    null comment '生日',
    city            varchar(255) default '' not null comment '城市',
    district        varchar(255) default '' not null comment '區域',
    address         varchar(255) default '' not null comment '剩餘地址',
    zip_code        varchar(50)  default '' not null comment '郵遞區號',
    last_login_time datetime                null comment '上次登入時間',
    is_enabled      tinyint      default 1  not null comment '狀態 1開啟 2關閉',
    created_at      datetime                not null comment '創建時間',
    updated_at      datetime                not null comment '更新時間',
    constraint idx_email
        unique (email)
) COMMENT ='會員資料表' COLLATE utf8mb4_general_ci;
