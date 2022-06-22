create table t_base_data
(
    offline_time bigint       null comment '离线时间',
    is_online    tinyint(1)   null comment '是否在线',
    score        bigint       null comment '积分',
    avatar_url   varchar(100) null comment '第三方头像链接',
    nickname     varchar(32)  not null,
    user_id      bigint auto_increment
        primary key
);

create table t_account_data
(
    account      varchar(8)  not null,
    equipment_id varchar(36) not null comment '设备ID',
    passwd       varchar(36) not null,
    user_id      int auto_increment
        primary key
);


create table g_game_data
(
    game_data blob null,
    user_id   int  not null
        primary key
);





