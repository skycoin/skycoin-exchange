//
//  ViewController.m
//  wallet-api-test
//
//  Created by iketheadore on 05/03/2018.
//  Copyright © 2018 iketheadore. All rights reserved.
//

#import "ViewController.h"
#import "mobile/Mobile.h"

@interface ViewController ()

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view, typically from a nib.
    NSString *docPath = [NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES) lastObject];
    
    MobileConfig *cfg = MobileNewConfig();
    cfg.walletDirPath = docPath;
    cfg.serverAddr = @"127.0.0.1:8080";
    cfg.serverPubkey = @"";

    NSLog(@"%@", cfg.walletDirPath);
    NSLog(@"%@", cfg.serverAddr);
    
    // Mobile environment initialize
    MobileInit(cfg);
    
    // Skycoin
    [self testCoin:docPath
          coinType:@"skycoin"
              seed:@""
  balanceCheckAddr:@"ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"
     transactionID:@"f73994cf833676f25766fdc88834e8c5019c34db66f78d8a5f3ca63bc9250d44"
          outputID:@"0b263bb6c7eb21d920846853a7b0c0bafa48e1b55dcc9c527fd4a2e1637057ab"];

    [self testCoin:docPath
          coinType:@"mzcoin"
              seed:@""
  balanceCheckAddr:@"HLEwsx3ndWwzBXLDcBom9N94bMFhx4ULup"
     transactionID:@"de76b20b71c9b29eacd3f0d64bbbf68f5a00804dafcbe08aa5f35aadb9207441"
          outputID:@"33dc904e0e697509b216d13cafbae49cc4b3da7073bab19e85a6d428645f429e"];

    // Aynrandcoin
    [self testCoin:docPath
          coinType:@"aynrandcoin"
              seed:@""
  balanceCheckAddr:@"Dd2ACzYi43bV9DQPNxKobwJ8FSLu6T2NNx"
     transactionID:@"ba484370bc8ab466cd94518b71772c7315989c1b44c0607cd973b6862952d242"
          outputID:@"5523a528f1162175f9a565da4b774eb5ce9ead9321bdd2bba9b172c7e9e69ec7"];

    // Shellcoin2
    [self testCoin:docPath
          coinType:@"shellcoin"
              seed:@""
  balanceCheckAddr:@"2Q5FxueDfDpCbQrVGxL1D9aU5Ruk9osg1RS"
     transactionID:@"bdf5cbdf4a84cbb8882ea79383450fae7cf16e5b3c9be98a339c69ab5193f3ed"
          outputID:@"36593c2d86839ecb9fc3383584cde5dc6c0fc36d87be044a8bd5fc1bf20b67cf"];

    // Suncoin
    [self testCoin:docPath
          coinType:@"suncoin"
              seed:@""
  balanceCheckAddr:@"aaRVsumKypkXtvYR5NP1NNwZJ9LFY7vEpn"
     transactionID:@"8f3bbb92f4f7137aa0b0546575edb1248fed47a2ea29e97a9493b1090f94a08a"
          outputID:@"4f8140af3f5dfb583c4eb60a9fb0e26e73a3f35dff544a26bacf9314bf7d9e6b"];

    // Metalicoin
    [self testCoin:docPath
          coinType:@"metalicoin"
              seed:@""
  balanceCheckAddr:@"22xpqmFsUiUv7EQSoY6ANVkuGZ5n7zvuVc6"
     transactionID:@"ec917bad743b8f950ca81a0ea9dc834b152ff8fcf8695ac453ccb1663812ea5d"
          outputID:@"4a65db65aa5362724392b206dcc3958f4ee558c96b6d97eed9d024c8a441b611"];

    // Lifecoin
    [self testCoin:docPath
          coinType:@"lifecoin"
              seed:@""
  balanceCheckAddr:@"R4kaEiHYoD8sKsgnnA4HQAuEheT9NTzL2s"
     transactionID:@"70597e175d4dc41b4dcd7ca9ca074b4cd121be231875771f876c0ab199681b5f"
          outputID:@"95f8bc6cc1e92f6a0c17fb85fddf4a0bc42d8d37dfbef4a24b8a0860f3f0064d"];
    
    // Fishercoin
    [self testCoin:docPath
          coinType:@"fishercoin"
              seed:@""
  balanceCheckAddr:@"2QDwdNMLXA2PEztTZBdS7233HY5Dy3L6pgM"
     transactionID:@"d585a387a594b489e8ad9ad7b5b5dbb3f1224fb0cd4e1152617e57f4df5add6e"
          outputID:@"8b2083cdf10b037846b24a791660dff46cfd00ae7bb4bece6c3830f24ff9b096"];
}

- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (NSString*)decodeAddressFromJSON: (NSString*) json index:(NSInteger) at {
    NSError *error = nil;
    NSData* data = [json dataUsingEncoding:NSUTF8StringEncoding];
    NSDictionary *dic = [NSJSONSerialization JSONObjectWithData:data options:kNilOptions error:&error];
    NSArray *addrs = [dic objectForKey:@"addresses"];
    
    NSDictionary *addrEntry = [addrs objectAtIndex:at];
    NSString *addr = [addrEntry objectForKey:@"address"];
    return addr;
}

- (void)testCoin: (NSString*) docPath
        coinType:(NSString*) ct
            seed:(NSString*) sd
balanceCheckAddr:(NSString*) baddr
   transactionID:(NSString*) txid
        outputID:(NSString*) outid {
    NSError *error = nil;
    NSString *wltID = MobileNewWallet(ct, sd, &error);
    if (error) {
        NSLog(@"Failed to create new wallet:%@", error);
        return;
    }
    NSLog(@"wallet id: %@", wltID);
    
    // Create new addresses
    NSString *addressEntries = MobileNewAddress(wltID, 1, &error);
    NSLog(@"new addresses: %@", addressEntries);

    // Get addresses
    NSString *addressEntries2 = MobileGetAddresses(wltID, &error);
    NSLog(@"get addresses: %@", addressEntries2);
    
    // Deocde json to get address
    NSString *addr = [self decodeAddressFromJSON:addressEntries index:0];
    
    // GetKeyPairOfAddr
    NSString *keyPairs = MobileGetKeyPairOfAddr(wltID, addr, &error);
    NSLog(@"key pairs: %@", keyPairs);
    
    // NewSeed
    NSString *seed = MobileNewSeed();
    NSLog(@"new seed: %@", seed);
    
    // GetBalance
    NSString* bal = MobileGetBalance(ct, baddr, &error);
    NSLog(@"balance: %@", bal);

    // GetWalletBalance
    NSString *wbal = MobileGetWalletBalance(ct, wltID, &error);
    NSLog(@"wallet balance: %@", wbal);
    
    // GetTransactionByID
    NSString *tx = MobileGetTransactionByID(ct, txid , &error);
    NSLog(@"tx of txid: %@ %@", txid, tx);
    
    // GetOutputByID
    NSString *output = MobileGetOutputByID(ct, outid, &error);
    NSLog(@"output: %@", output);
    
    // Send
    MobileSendOption *opt = MobileNewSendOption();

    NSString *sendTxID = MobileSend(ct, wltID, addr, @"100000", opt, &error);
    if (error) {
        NSLog(@"Send %@ failed: %@", ct, error);
    }
    NSLog(@"txid: %@", sendTxID);
    
    
    
    // Release resources
    // Delete the wallet file
    NSFileManager *manager = [NSFileManager defaultManager];
    NSString *wltFile = [wltID stringByAppendingString:@".wlt"];
    [manager removeItemAtPath:[docPath stringByAppendingPathComponent:wltFile] error:&error];
    
    // Delete wallet bak file
    NSString *bakWltFile = [wltFile stringByAppendingString:@".bak"];
    [manager removeItemAtPath:[docPath stringByAppendingPathComponent:bakWltFile] error:&error];
}

@end
