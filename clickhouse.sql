create table default.kline
(
    symbol             String comment '符号',
    k_start_time       DateTime comment 'k线开始时间',
    k_end_time         DateTime comment 'k线结束时间',
    start_price        String comment '开盘价',
    end_price          String comment '收盘价',
    high_price         String comment '最高价',
    low_price          String comment '最低价',
    volume_total_usd   String comment '成交总金额',
    volume_total_count Int64 comment '成交总笔数'
)
    engine = MergeTree PARTITION BY symbol
        PRIMARY KEY (symbol, k_start_time)
        ORDER BY (symbol, k_start_time)
        SETTINGS index_granularity = 8192
        comment 'k线数据,按照1小时存储';

create view day7_avg_kline as
select symbol,
       k_start_time as timestamp,
       avg(toFloat64OrZero(end_price))
           over (PARTITION BY symbol ORDER BY k_start_time range between 167 preceding and current row ) as price
from kline;


create view day25_avg_kline as
select symbol,
       k_start_time as timestamp,
       avg(toFloat64OrZero(end_price))
           over (PARTITION BY symbol ORDER BY k_start_time range between 599 preceding and current row ) as price
from kline;

create view day99_avg_kline as
select symbol,
       k_start_time as timestamp,
       avg(toFloat64OrZero(end_price))
           over (PARTITION BY symbol ORDER BY k_start_time range between 2375 preceding and current row ) as price
from kline

