# 游戏用户表
create table t_player
(
    userId bigint auto_increment,
    nickname varchar(32) not null,
    avatarURL varchar(100) null comment '第三方头像链接',
    score bigint null comment '积分',
    isOnline bool null comment '是否在线',
    offlineTime bigint null comment '离线时间',
    constraint player_pk
        primary key (userId)
);

#  账户表
create table t_account
(
    userId int auto_increment,
    passwd varchar(36) null,
    equipmentId int null comment '设备ID',
    account varchar(8) null,
    constraint t_account_pk
        primary key (userId)
);


create table g_game_data
(
    id int auto_increment,
    userId int null,
    game_data blob null,
    constraint g_game_data_pk
        primary key (id)
);


