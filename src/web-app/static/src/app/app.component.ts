import {Component, OnInit} from "@angular/core";
import {RouteConfig, RouterLink, ROUTER_DIRECTIVES} from "@angular/router-deprecated";
import {Http, HTTP_BINDINGS, Response} from '@angular/http';
import {HTTP_PROVIDERS, Headers} from '@angular/http';

import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Component({
    selector: "app",
    templateUrl: "./app/index.html",
    directives: [RouterLink, ROUTER_DIRECTIVES]
})

export class AppComponent implements OnInit {
    bidList: Array<any>;
    askList: Array<any>;
    accountId:any;
    key:any;
    balance:any;
    OrderInputIsVisible:boolean;
    orderType:number;
    orderPrice:number;
    orderAmount:number;

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    ngOnInit() {
        this.bidList = [];
        this.askList = [];
        this.accountId = null;
        this.key = null;
        this.orderPrice = 0;
        this.orderAmount = 0;
        this.OrderInputIsVisible = false;
        this.balance = {
          skycoin:0,
          bitcoin:0
        };
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var self = this;
        this.http.post('/api/v1/accounts', '')
            .map((res) => res.json())
            .subscribe(data => {
              console.log('request account', data);
              if (data.result.success) {
                self.accountId = data.account_id;
                self.key = data.key;
                self.loadBidList();
                self.loadAskList();
                self.getBalance();
                /*
                this.http.post('/api/v1/account/deposit_address?id=self.accountId&key=self.key&coin_type=skycoin', '')
                    .map((res) => res.json())
                    .subscribe(data => {
                    console.log("deposite", data);
                    self.loadBidList();
                    self.loadAskList();
                    }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
                    */
              } else {
                alert("Cannot get account from server. please check connection with server");
                return;
              }
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    loadBidList() {
        var self = this;
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/orders/bid?coin_pair=bitcoin/skycoin&id=' + this.accountId + "&key=" + this.key + '&start=1&end=10';
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
        var url = '/api/v1/orders/ask?coin_pair=bitcoin/skycoin&id=' + this.accountId + '&key=' + this.key + '&start=1&end=10';
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
        var url = '/api/v1/account/balance?id=' + this.accountId + '&key=' + this.key + '&coin_type=skycoin';
        var self = this;
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get skycoin balance", url, data);
              self.balance.skycoin = data.balance;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
        url = '/api/v1/account/balance?id=' + this.accountId + '&key=' + this.key + '&coin_type=bitcoin';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get bitcoin balance", url, data);
              self.balance.bitcoin = data.balance;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
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
      this.http.post('/api/v1/account/order?id=' + this.accountId + '&key=' + this.key, JSON.stringify(data))
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
}
