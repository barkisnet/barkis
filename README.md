# Barkis

Barkisnet establishes a next-generation distributed application based on blockchain technology.

**Note**: Requires [Go 1.12+](https://golang.org/dl/)

## Deploy Mainnet Node

1. Get code

    ```shell script
    git clone https://github.com/barkisnet/barkis.git
    ```

2. Build binarie

    ```shell script
    make build
    ```
    
    The generated binaries locate in **build** directory.

3. Init mainnet node PathToNodeHomeDirectory
    
    ```shell script
    ./build/barkisd init {node name} --home {PathToNodeHomeDirectory}
    ```
4. Get mainnet genesis file

    ```shell script
    wget https://raw.githubusercontent.com/barkisnet/barkisnet-binary/master/barkisnet-mainnet/genesis.json -O {PathToNodeHomeDirectory}/config/genesis.json
    ```
   
5. Get mainent configuration file: `networkConfig.json`

    ```shell script
    wget https://raw.githubusercontent.com/barkisnet/barkisnet-binary/master/barkisnet-mainnet/networkConfig.json -O networkConfig.json
    ```

6. Edit `{PathToNodeHomeDirectory}/config.toml` to config `seeds` and `persistent_peers` according to `networkConfig.json`

   ```toml
    # Comma separated list of seed nodes to connect to
    seeds = ""
    
    # Comma separated list of nodes to keep persistent connections to
    persistent_peers = ""
    ```
 
7. Edit `{PathToNodeHomeDirectory}/app.toml` to change the upgrade heights according to `networkConfig.json`
    
   ```toml
    [upgrade]
    # Upgrade to change reward rules
    RewardUpgrade = 9223372036854775807
    
    # Upgrade to change reward rules
    TokenIssueHeight = 9223372036854775807
    
    # Upgrade to update voting period
    UpdateVotingPeriodHeight = 9223372036854775807
    ```

8. Edit `{PathToNodeHomeDirectory}/app.toml` to config minimum gas prices
    
    ```toml
    minimum-gas-prices = ""
    ```
    Recommended value: `0.01ubarkis`

9. Start mainnet node

    ```shell script
    nohup ./build/barkisd start --home {PathToNodeHomeDirectory} &
    ```
 
10. Check running log

    ```shell script
    tail -f nohup.out
    ```
