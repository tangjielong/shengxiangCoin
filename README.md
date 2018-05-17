# shengxiangCoin
使用区块链技术实现的类似bitcoin的分布式数字货币

# 命令
## createblockchain -address ADDRESS
创建区块链<br/>
创建一个coinbase交易，里面记录了奖励给ADDRESS账户的10个coin<br/>
把这个coinbase存储到创世块（genesis block）中<br/>
把区块链信息（所有block的数据和最后一个Block的hash值）存储到hello.db中，并且创建UTXO集提高查询效率

## printchain
打印区块链中所有block和transaction的信息，json格式

## reindexutxo
重新创建UTXO集（blockchain的索引）

## createwallet
创建一个新的钱包地址并且保存到wallet.bin中

## listaddresses
打印wallet.bin中所有钱包的地址

## getbalance -address ADDRESS
查询ADDRESS账户的余额

## send -from FROM -to TO -amount AMOUNT
从FROM账户转帐AMOUNT个coin到TO账户<br/>
这个命令同时会挖一个新的block来存储这个交易以及coinbase交易（作为奖励）

# TODO
实现分布式，把各个角色分离开<br/>
矿工 —— 负责接收用户的交易信息，计算新的block<br/>
全结点 —— 验证交易的有效性<br/>
SPV —— 客户端，用来简单支付验证

