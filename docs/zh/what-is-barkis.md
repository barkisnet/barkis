# Barkis是什么

`barkis`是作为BarkisNet应用程序的barkis的名称。它有两个主要的入口：

+ `barkisd` : Barkis的服务进程，运行着`barkis`程序的全节点。
+ `barkiscli` : Barkis的命令行界面，用于同一个Barkis的全节点交互。

`barkis`基于barkis构建，使用了如下模块:

+ `x/auth` : 账户和签名
+ `x/bank` : token转账
+ `x/staking` : 抵押逻辑
+ `x/mint` : 增发通胀逻辑
+ `x/distribution` : 费用分配逻辑
+ `x/slashing` : 处罚逻辑
+ `x/gov` : 治理逻辑
+ `x/ibc` : 跨链交易
+ `x/params` : 处理应用级别的参数

接着，学习如何[安装Barkis](installation.md)
