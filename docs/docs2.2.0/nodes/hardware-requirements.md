# Validator hardware (**WIP**)
This page contains three main hardware configurations that can be used in order to create a new Commercio.network 
validator machine. Please note that higher the effort you will put into creating a stable and robust machine and lower 
the chances of getting slashed due to downtime. 

## Requirements
- [Testnet](#testnet)
- [Mainnet](#mainnet)

## Testnet
While running a testnet node, you can use any hardware you want. 

During this guide we will assume that you will use a simple cloud server, to avoid having to deal with 
complicated stuff that is not really necessary.  
Due to this reason, the following lines are written for Digital Ocean, but everything you will read is applicable to 
any other cloud provider like Amazon AWS, Google Cloud, Microsoft Azure, Alibaba Cloud or Scaleway.

Here's a friendly Digital Ocean $50 credit Coupon link: https://m.do.co/c/132ef6958ef7

For the sake of simplicity, we will assume you have selected the following DigitalOcean configuration. 
Please not that this is just an example, but any configuration similar to this one will work perfectly fine.      

| Characteristic | Specification |
| :------------: | :-----------: |
| Operative System | Ubuntu 18.04 |
| Number of CPUs | 2 |
| RAM | 4GB |
| SSD | 80GB | 

Also, we need to make sure the following requirements are met: 
* Allow incoming connections on port `26656`
* Have a static IP address
* Have access to the root user

## Mainnet
### Low level (Not Recommended)
This configuration is a basic level with low security level and low high availability configuration.    
It is prone to intrusion and out of service in any moment due to hardware failure or network line problems.    
     
* 1 Server with
  * **CPU**: Bare minimum to support the last version of most common operative systems
  * **Ram**: 32GB
  * **Storage**: 1 x 1TB SSD
  * **Power supply**: 1 power supply
  * **Network**: 1 ethernet port with 100 Mbit speed capability
* 1 internet connection with 30/30Mbit bidirectional capability
* 1 [Yubi HSM2](https://www.yubico.com/product/yubihsm-2/)
* 50,000 Commercio Tokens


### Mid level  (Minimum Necessary)
This configuration guarantees a basic security level due to the fact that the validator node is not directly 
attached to the internet. It also guarantees a fair availability due to the double power supply and ethernet connection. 

* 1 Server with
  * **CPU**: Bare minimum to support the last version of most common operative systems
  * **Ram**: 32GB
  * **Storage**: 1 x 1TB SSD
  * **Power supply**: 2 x power supplies (to prevent power down if one breaks)
  * **Net**: 2 x Ethernet port with 100 Mbit speed capabilities
* 1 internet connection with 30/30Mbit bidirectional capability
* 1 firewall with an on-site VPN capability
* 1 sentry node configured on a major service provider (AWS, Azure, Google)
* 1 [Yubi HSM2](https://www.yubico.com/product/yubihsm-2/)
* 1 UPS with 1000VA capacity for minimum resistance to power surges and out of power services
* 50,000 Commercio Tokens
 

#### Nodes configuration:    
##### Validator node
| Config        | Setting       |
| ------------- |:-------------:|
| pex      | false |
| persistent_peers     | private sentry node      |
| private_peer_ids | --     |
| addr_book_strict | false     |

##### Private sentry node
| Config        | Setting       |
| ------------- |:-------------:|
| pex      | true |
| persistent_peers     | validator node, own public node     |
| private_peer_ids | validator node id    |
| addr_book_strict | false     |

##### Public sentry node
| Config        | Setting       |
| ------------- |:-------------:|
| pex      | true |
| persistent_peers     | other public sentry nodes    |
| private_peer_ids | validator node id, private node id    |
| addr_book_strict | false     |


### Top Level (Recommended)
This configuration ensures high availability and keeps the chances of getting slashed at the minimum possible. It works
with two different server to ensure that even if one fails completely the other one can kick in and replace it until 
the first is properly fixed.

* 2 Servers
  * **CPU**: Bare minimum to support the last version of most common operative systems
  * **Ram**: 32GB
  * **Storage**: 2/3 x 1TB SSD with RAID 1 or RAID 5 configuration 
  * **Power supply**: 2 x power supply (to prevent power down if one breaks)
  * **Net**: 2 x ethernet port with 1000 Mbit speed capabilities
* 2 internet connections with 100/100Mbit bidirectional capability
* 2 dedicated switches
* 2 firewalls with an on-site VPN capability
* 1 AWS Load Balancer
* 2 sentry node configured on major service providers (AWS, Azure, Google)
* 2 [Yubi HSM2](https://www.yubico.com/product/yubihsm-2/) (one per server)
* 50,000 Commercio Tokens
