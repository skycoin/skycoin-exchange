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

    //Constructor method for load HTTP object
    constructor(private http: Http) { }

    ngOnInit() {
        this.bidList = [];
        this.askList = [];
        this.loadBidList();
        this.loadAskList();
    }

    loadBidList() {
        var headers = new Headers();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        this.http.get('http://localhost:6060/api/v1/orders/bid?coin_pair=bitcoin/skycoin', { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
            console.log(data);
            this.bidList = data;
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
        this.http.get('http://localhost:6060/api/v1/orders/ask', { headers: headers })
            .map((res) => res.json())
            .subscribe(data => {
                console.log(data);
                this.askList = data;
            }, err => console.log("Error on load outputs: " + err), () => console.log('Connection load done'));

        this.askList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470193222,
            "coin_type":"SKY"
        });
        this.askList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470293222,
            "coin_type":"BTC"
        });
        this.askList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470393222,
            "coin_type":"SKY"
        });
        this.askList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470493222,
            "coin_type":"SKY"
        });
        this.askList.push({
            "id": 8,
            "type": "bid",
            "price": 25,
            "amount": 90000,
            "rest_amt": 90000,
            "created_at": 1470593222,
            "coin_type":"BTC"
        });
    }
}
