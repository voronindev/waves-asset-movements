WAVES asset movements
-
Allows scan blockchain Waves for asset X movements from height Y (for issue height for example).
Results will be saved in the file, located in the same directory.
Every movement will be printed as new line with the following pattern:
```
{height:10} {tx_hash:42} {tx_type:4} {tx_amount:16} {tx_sender:42} {tx_recipient:42}
```

### Supported transactions:
- TransferV1
- TransferV2
- MassTransferV1
- ExchangeV1
- ExchangeV2

### How to use:
If you want to track every WBTC (`8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS`) movement from `257457` (issue height), for example:
```
./go_build_waves_movements -a 8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS -h 257457
```

If you want to use custom waves node (your own, local):
```
./go_build_waves_movements -a 8LQW8f7P5d5PZM7GtZEBgaqRPGSzS3DfPuiXrURJ4AJS -h 257457 -n http://127.0.0.1:6869
```