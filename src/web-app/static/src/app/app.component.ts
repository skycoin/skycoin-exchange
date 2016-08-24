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

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    ngOnInit() {
        this.bidList = [];
        this.askList = [];
        this.accountId = null;
        this.key = null;
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
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/orders/bid?coin_pair=bitcoin/skycoin&id=' + this.accountId + "&key=" + this.key + '&start=1&end=10';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get bid list", url, data);
              if (data.result.success) {
                this.bidList = data.orders;
              } else {
                return;
              }

            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));

        //input data format
        /*
         {
         "id": 8,
         "type": "bid",
         "price": 25,
         "amount": 90000,
         "rest_amt": 90000,
         "created_at": 1470193222
         },
         */
        this.bidList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470193222,
            "coin_type":"SKY"
        });
        this.bidList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470293222,
            "coin_type":"BTC"
        });
        this.bidList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470393222,
            "coin_type":"SKY"
        });
        this.bidList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470493222,
            "coin_type":"SKY"
        });
        this.bidList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470593222,
            "coin_type":"BTC"
        });
    }

    loadAskList() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/orders/ask?coin_pair=bitcoin/skycoin&id=' + this.accountId + '&key=' + this.key + '&start=1&end=10';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get ask list", url, data);
              if (data.result.success) {
                this.askList = data.orders;
              } else {
                return;
              }
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }

    getBalance() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        var url = '/api/v1/account/balance?id=' + this.accountId + '&key=' + this.key + '&coin_type=skycoin';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get skycoin balance", url, data);
              this.balance.skycoin = data.balance;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
        url = '/api/v1/account/balance?id=' + this.accountId + '&key=' + this.key + '&coin_type=bitcoin';
        this.http.get(url, { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
              console.log("get bitcoin balance", url, data);
              this.balance.bitcoin = data.balance;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));
    }
}
