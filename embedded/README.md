## Get CSRF Token

GET /api/v1/auth/csrf_token

## Login

POST /api/v1/auth/login  
headers: {'Accept':'*/*','X-CSRF-Token':<CSRF_TOKEN>}  
data: {'username':'','password':''}  

Returns session cookie in response header set-cookie  

## User Balances

GET /api/v1/user/balances  
headers: {'Accept':'*/*','Cookie':<set-cookie>}  

## POST Order

POST /api/v1/exchange/orders   
data: {"chain-id":"xar-chain-zafx","market_id":"1","direction":"BID|ASK","price":"100000000","quantity":"100000000","type":"LIMIT","time_in_force":100},  
headers:  {'Accept':'*/*','Cookie':<set-cookie>}  
