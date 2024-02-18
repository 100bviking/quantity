# quantity
交易所量化交易,目标1万到1000万,每月收益需要实现100%,10个月完成目标。

## 目标:
   * 支持交易所包括:biance,okx,bitget,huobi,gate,mexc 
   * 支持交易类型包括: 现货交易，合约交易,交易所套利交易,异常(暴涨/暴跌)跟单交易
   * 支持策略包括:k线形态分析下单，手动指定下单，外部预言机热度下单
   * 支持同时监控交易对数量：1k-10k之间
   * 目前不支持不同交易对自动应用不同策略，所有交易对都从策略中心加载

## K线数据模块:
   * 封装k线api，底层数据来源可能是交易所或本地数据库
   * 交易所api出现异常，需要马上发送通知到告警中心，转为人工处理。同时暂停所有策略应用模块的执行。

## 策略相关模块:
   * 策略模块: 
     * 支持比如移动均线，MACD，双顶形态等k线形态分析，识别买入/卖出点位
     * 使用策略模式，每个策略单独开发
   * 策略管理模块:
     * 支持添加/删除策略到策略中心
     * 支持策略聚合，多个策略聚合为1个策略，最终应用只有1个策略。
     * 支持策略权重设置，可能高权重才会应用，低权重会进行忽略。
   * 策略应用模块：
     * 交易对从策略中心下载策略进行分析，识别出买/卖点，形成通知指令，到下单模块。
     * 所有交易对并发执行，提高并发度。
     
## 下单相关模块:
   * 指令中心模块:
     * 优先级队列。高优先级执行优先执行，低优先级随后执行。人工指令拥有最高优先级。
     * 指令定义:
       * 交易对 SPOT BUY  price number weight strategy 对应于 现货 交易对 买入 数量 (权重) 策略id
       * 交易对 SPOT SELL price number weight strategy 对应于 现货 交易对 卖出 数量 (权重) 策略id
     * 指令执行模块:
       *  指令解析器: 翻译指令，加载指令执行单元进行执行。
     * 交易记录模块：
       * 记录指令
       * 记录交易对开单记录，策略来源
       
## 收益分析模块：
   * 实时分析完成的策略收益，策略收益过低，自动降低其权重，策略收益增加，增加其权重
   * 策略告警，如果策略收益低于某个阈值，发送告警通知到通知中心
   * 收益报表，每日统计当日收益，总收益,年化等指标,形成报表，通知到通知中心。
   * 根据不同账号进行分析。

## 通知中心模块:
   * 支持多种通知渠道，包括telegram，slack，微信，邮件等。
## 账号中心模块:
   * 支持添加账号,目前只支持1个账号。

## 前台模块:
   * 支持查看收益报表
   * 支持策略排名
   * 支持策略中心管理

## 技术层面:
   * 使用 go-zero框架，对外api支持grpc，jsonrpc，以及http接口完成所有操作。
   * k线模块使用clickhouse集群，优点是压缩率高，api丰富。
   * 数据库使用mysqldb，存储策略，账号等信息。
   * 缓存使用redis集群
   * 消息通知使用nsq集群
   * 底层使用k8s进行部署
   * 日志使用es进行聚合
   * 代码使用golang-clint进行分析
   * 需要实现优先级队列，高并发，时间轮等。
   * 策略分析到指令执行，控制在10秒之内完成。
   * 支持promethues进行资源监控告警,和定时汇报信息.
   