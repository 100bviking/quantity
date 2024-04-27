create database klines;
use klines;
create table `cursor`
(
    symbol     varchar(15)                        not null
        primary key,
    timestamp  datetime                           not null comment 'k线当前同步到的时间',
    balance    double                             not null comment '当前余额,账户买卖以后自动更新',
    created_at datetime default CURRENT_TIMESTAMP not null,
    updated_at datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
    engine = InnoDB;

create database account;
use account;
create table user
(
    id         int auto_increment
        primary key,
    user_name  varchar(255) not null,
    api_key    varchar(255) not null,
    api_secret varchar(255) not null
)
    engine = InnoDB;


create database `order`;
use `order`;
create table `order`
(
    id           int auto_increment
        primary key,
    symbol       varchar(15)                        not null,
    order_price  double                             null comment '订单成交时候的价格',
    submit_price double                             not null comment '订单提交时候的价格',
    amount       varchar(255)                       not null comment '成交数量',
    money        double                             not null comment '成交usdt数量',
    action       int                                not null comment '0: hold;1:buy;-1:sell',
    order_time   datetime                           not null,
    status       int                                not null comment '0:未知;1:成功;2失败',
    created_at   datetime default CURRENT_TIMESTAMP not null,
    updated_at   datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
)
    engine = InnoDB;




