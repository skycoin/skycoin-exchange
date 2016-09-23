import {Component, OnInit} from "@angular/core";
import {RouteConfig, RouterLink, ROUTER_DIRECTIVES} from "@angular/router-deprecated";
import {Http, HTTP_BINDINGS, Response} from '@angular/http';
import {HTTP_PROVIDERS, Headers} from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

declare var moment: any;
declare var _: any;

@Component({
    selector: "app",
    templateUrl: "./app/index.html",
    directives: [RouterLink, ROUTER_DIRECTIVES]
})

export class AppComponent implements OnInit {
    bidList: Array<any>;
    askList: Array<any>;
    depositList: Array<any>;
    accountList: Array<any>;
    eventList: Array<any>;
    pubkey:any;
    balance:any;
    OrderInputIsVisible:boolean;
    orderType:number;
    orderPrice:number;
    orderAmount:number;
    createWalletVisible:boolean;
    walletSeed:any;
    wallets:any;
    updateWalletVisible:boolean;
    walletAmount:number;
    tempWallet:any;

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    ngOnInit() {
        this.bidList = [];
        this.askList = [];
        this.depositList = [];
        this.accountList = [];
        this.eventList = [];
        this.orderPrice = 0;
        this.orderAmount = 0;
        this.OrderInputIsVisible = false;
        this.createWalletVisible = false;
        this.updateWalletVisible = false;
        this.walletSeed = '';
        this.walletAmount = 0;
        this.wallets = [];
        this.balance = {
          skycoin:0,
          bitcoin:0
        };
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var self = this;
        this.getAccountList();
        this.http.get('/api/v1/account?active=1')
          .map((res) => res.json())
          .subscribe(data => {
            if (data.result.success) {
              //got active user
              console.log("Found active user", data.accounts[0]);
              self.pubkey = data.accounts[0].pubkey;
              if (data.accounts[0].wallet_ids) {
                this.wallets = [];
                if (data.accounts[0].wallet_ids.bitcoin) {
                  this.wallets.push({
                    type:'bitcoin',
                    id:data.accounts[0].wallet_ids.bitcoin,
                    balance:{amount:0}
                  });
                }
                if (data.accounts[0].wallet_ids.skycoin) {
                  this.wallets.push({
                    type:'skycoin',
                    id:data.accounts[0].wallet_ids.skycoin,
                    balance:{amount:0}
                  });
                }
              }
              this.getWalletBalance();
              this.loadBidList();
              this.loadAskList();
              this.getBalance();
              this.getDepositList();
              this.getEventList();
            } else {
              //create new account
              this.http.post('/api/v1/accounts', '')
                .map((res) => res.json())
                .subscribe(data => {
                  console.log('request account', data);
                  if (data.result.success) {
                    this.pubkey = data.pubkey;
                    this.wallets = [];
                    this.loadBidList();
                    this.loadAskList();
                    this.getBalance();
                    this.getDepositList();
                    this.getEventList();
                  } else {
                    alert("Cannot get account from server. please check connection with server");
                    return;
                  }
                }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
            }
          }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    loadBidList() {
        var self = this;
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/orders/bid?coin_pair=bitcoin/skycoin&pubkey=' + this.pubkey + '&start=1&end=10';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get bid list", url, data);
              if (data.result.success) {
                self.bidList = data.orders;
              } else {
                return;
              }

            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    loadAskList() {
        var self = this;
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/orders/ask?coin_pair=bitcoin/skycoin&pubkey=' + this.pubkey + '&start=1&end=10';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get ask list", url, data);
              if (data.result.success) {
                self.askList = data.orders;
              } else {
                return;
              }
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    getBalance() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/account/balance?coin_type=skycoin';
        var self = this;
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get skycoin balance", url, data);
              self.balance.skycoin = data.balance.amount;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
        url = '/api/v1/account/balance?coin_type=bitcoin';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get bitcoin balance", url, data);
              self.balance.bitcoin = data.balance.amount;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    getWalletBalance() {
        var self = this;
        this.wallets.map(function(o){
          var url = '/api/v1/wallet/balance?id=' + o.id;
          self.http.get(url, {})
              .map((res) => res.json())
              .subscribe(data => {
                console.log("get wallet balance", url, data);
                o.balance = data.balance;
              }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
        });
    }

    getDepositList() {
      var self = this;
      var headers = new Headers();
      headers.append('Content-Type', 'application/x-www-form-urlencoded');
      this.http.post('/api/v1/account/deposit_address?pubkey=' + this.pubkey + '&coin_type=skycoin', '')
          .map((res) => res.json())
          .subscribe(data => {
            if (data.result.success) {
                self.depositList['SKY'] = [];
                self.depositList['SKY'].push({
                  "address": data.address
                });
            }
          }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
      this.http.post('/api/v1/account/deposit_address?pubkey=' + this.pubkey + '&coin_type=bitcoin', '')
          .map((res) => res.json())
          .subscribe(data => {
          console.log("deposite-bitcoin", data);
            if (data.result.success) {
                self.depositList['BTC'] = [];
                self.depositList['BTC'].push({
                  "address": data.address
                });
            }
          }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    getAccountList() {
      var self = this;
      this.http.get('/api/v1/account')
        .map((res) => res.json())
        .subscribe(data => {
          if (data.result.success) {
            console.log("getAccountList", data);
            data.accounts.map(function(o){
              self.accountList.push({
              "pubkey": o.pubkey,
              "wallet_ids": o.wallet_ids
              });
            });
          }
        }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    getEventList() {
      var self = this;
      this.eventList.push({
      "event_type": "deposit",
      "timestamp": Date.now() - Math.round(Math.random() * 1000000),
      "coin_type": "BTC",
      "amount": 50000
      });
      this.eventList.push({
      "event_type": "deposit",
      "timestamp": Date.now() - Math.round(Math.random() * 1000000),
      "coin_type": "SKY",
      "amount": 85000
      });
      this.eventList.push({
      "event_type": "withdraw",
      "timestamp": Date.now() - Math.round(Math.random() * 1000000),
      "coin_type": "BTC",
      "amount": 20000
      });
    }

    getDateTimeString(ts) {
        return moment.unix(ts / 1000).format("YYYY-MM-DD HH:mm");
    }

    createOrder(type) {
      this.orderType = type;
      this.orderAmount = 0;
      this.orderPrice = 0;
      this.OrderInputIsVisible = true;
    }

    createOrderDo(type, amount, price) {
      var data = {
                     "type": (type === 1 ? 'bid' : 'ask'),
                     "coin_pair":'bitcoin/skycoin',
                     "amount":Number(amount),
                     "price":Number(price)
                  };
      var self = this;
      this.http.post('/api/v1/account/order?pubkey=' + this.pubkey, JSON.stringify(data))
          .map((res) => res.json())
          .subscribe(data => {
          console.log("create order", data);
            if (data.result.success) {
              if (type === 1) {
                self.loadBidList();
              } else {
                self.loadAskList();
              }
            } else {
              alert(data.result.reason);
            }
          }, err => console.log("Error on load outputs: " + err), () => {
            console.log('Connection load done');
            self.hideOrderInputDialog();
          }
          );
    }

    hideOrderInputDialog() {
      this.OrderInputIsVisible = false;
    }

    createWallet() {
      this.walletSeed = this.randomString(16, 36);
      this.createWalletVisible = true;
    }

    hideCreateWallet() {
      this.createWalletVisible = false;
    }

    createWalletDo(seed, type) {
      //create wallet
      this.http.post('/api/v1/wallet?type=' + type + '&seed=' + seed, '')
          .map((res) => res.json())
          .subscribe(data => {
          console.log("create wallet", data);
            if (data.result.success) {
              var oldWallet = _.find(this.wallets, function(o){
                return o.type === type;
              });
              if (oldWallet) {
                oldWallet.id = data.id;
              } else {
                this.wallets.push({
                  type:type,
                  id:data.id,
                  balance:{amount:0}
                });
              }

              //create addrsss
              this.http.post('/api/v1/wallet/address?id=' + data.id, '')
                  .map((res) => res.json())
                  .subscribe(data => {
                    if (data.result.success) {
                      console.log("create address", data.address);
                      this.getWalletBalance();
                    }
                    this.hideCreateWallet();
                  }, err => console.log("Error on create wallet: " + err), () => {

                  });
            } else {
              alert(data.result.reason);
            }
          }, err => console.log("Error on create wallet: " + err), () => {

          }
          );
    }

    updateWallet(wallet) {
      this.walletAmount = 0;
      this.updateWalletVisible = true;
      this.tempWallet = wallet;
    }

    hideUpdateWallet() {
      this.updateWalletVisible = false;
    }

    updateWalletDo(amount) {
      //create wallet
      var url = '/api/v1/admin/account/balance?coin_type=' + this.tempWallet.type + '&dst=' + this.pubkey + '&amt=' + amount;
      console.log(url);
      this.http.put(url, '')
          .map((res) => res.json())
          .subscribe(data => {
          console.log("update wallet", data);
            if (data.result.success) {

            } else {
              alert(data.result.reason);
            }
          }, err => console.log("Error on update wallet: " + err), () => {

          }
          );
    }

    randomString(len, bits) {
        bits = bits || 36;
        var outStr = "", newStr;
        while (outStr.length < len) {
            newStr = Math.random().toString(bits).slice(2);
            outStr += newStr.slice(0, Math.min(newStr.length, (len - outStr.length)));
        }
        return outStr.toUpperCase();
    }

    convertCoinType(type) {
      if (type === 'bitcoin') {
        return 'BTC';
      } else {
        return 'SKY';
      }
    }
}
