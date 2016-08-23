/* Avoid: 'error TS2304: Cannot find name <type>' during compilation */
///<reference path="../../typings/index.d.ts"/>

import {AppComponent} from "./app.component";
import {bootstrap} from "@angular/platform-browser-dynamic";
import {provide} from "@angular/core";
import {LocationStrategy, HashLocationStrategy} from "@angular/common";
import {ROUTER_PROVIDERS} from "@angular/router-deprecated";
import {HTTP_BINDINGS} from '@angular/http';

bootstrap(AppComponent, [
    ROUTER_PROVIDERS,
    HTTP_BINDINGS,
    provide(LocationStrategy, {useClass: HashLocationStrategy})
]);
