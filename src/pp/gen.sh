protoc --go_out=./ pp.common.proto \
  pp.encrypt.proto \
  pp.account.proto \
  pp.deposit.proto \
  pp.withdrawal.proto \
  pp.balance.proto \
  pp.order.proto \
  pp.coin.proto \
  pp.request.proto \
  pp.utxo.proto \
  pp.transaction.proto
