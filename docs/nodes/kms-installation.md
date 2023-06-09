# KMS with YubiHSM2

## General concepts


## Hardware requirements


- SO: Ubuntu 18.04 (64-bit) or later  (preferably LTS version)
- Minimun amount of RAM: 4Gb
- Avaiable disk space: 100 Gb
- Connect to the your validator node in private network (LAN or VPN)


## Scenario

Schema di base/scenario
Per semplicità nella spiegazione possiamo supporre di avere una vpn che collega il Data Center con il Cloud con la classe di ip 10.1.1.0/24.
L’ip del nodo validatore sarà 10.1.1.254.
L’ip del kms sarà 10.1.1.1.
Inoltre supponiamo che i sentry node siano su delle altre classi di ip, 10.1.2.0/24 per Regione 1 e 10.1.3.0/24 per Regione 2
La protezione a livello di apertura porta sarà la seguente
KMS: nessuna porta aperta, solo accesso alla vpn verso il Validator node
Validator Node: porta in ascolto sulla vpn limitata al KMS 26658. Porta 26656 in ascolto sulle reti interne dei sentry node.
Sentry Node: Porta 26656 in ascolto su tutti gli IP. Porta 26657 in ascolto in locale o su tutte gli ip ma con accesso limitato a specifici client o mediata con un reverse proxy


## Installation

Upgrade your ubuntu system

```bash
sudo su -
apt update
apt upgrade -y
```

Create a new user to run the KMS server and reboot the system

```bash
mkdir /data_tmkms
useradd -m -d /data_tmkms/tmkms -G tmkms -s /bin/bash
echo 'SUBSYSTEMS=="usb", ATTRS{product}=="YubiHSM", GROUP=="tmkms"' >> /etc/udev/rules.d/10-yubihsm.rules
reboot
```

NB: it may be necessary to restart the server several times to apply the rules

Install the components

```bash
sudo su -
apt install gcc git libusb-1.0-0-dev pkg-config -y
```

Install rust: language for compiling TmKms

```bash
su - tmkms
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env

```

Install TmKms

```bash
cd $HOME
git clone https://github.com/tendermint/kms.git
cd $HOME/kms
cargo install tmkms --features=yubihsm --locked --force --version=0.10.0
```

Check the installation version

```bash
tmkms version
```

### YubiHSM2 Service Installation

The YubiHSM2 is a hardware device that must be connected to the server via USB. It is also necessary to install additional tools of the yubihsm in order to manage the connection as a service and not directly on usb.

Install yubico utilities. Search in [https://developers.yubico.com/YubiHSM2/Releases/](https://developers.yubico.com/YubiHSM2/Releases/) the right version for your system

```bash
sudo su -
cd /tmp
VERSION="yubihsm2-sdk-2019-12-ubuntu1804-amd64.tar.gz"
wget "https://developers.yubico.com/YubiHSM2/Releases/$VERSION"
tar zxf $VERSION 
sudo apt install ./yubihsm2-sdk/*.deb
rm -rf $VERSION
```

Activate the service

```bash
sudo tee /etc/systemd/system/yubihsm-connector.service > /dev/null <<EOF 
[Unit]
Description=YubiHSM connector
Documentation=https://developers.yubico.com/YubiHSM2/Component_Reference/yubihsm-connector/
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service

[Service]
Restart=on-abnormal
User=tmkms
Group=tmkms
ExecStart=/usr/bin/yubihsm-connector -c /etc/yubihsm-connector.yaml
PrivateTmp=true
ProtectHome=true
ProtectSystem=full

[Install]
WantedBy=multi-user.target
EOF

systemctl enable yubihsm-connector.service
systemctl start yubihsm-connector.service
```
Check the service status and the listening port

```bash
systemctl status yubihsm-connector.service
ss -tlpn | grep 12345
```


### HSM Installation

Insert yubiHSM 2 into the usb slot. If the device has been used previously, reset it by holding down the metal ring for more than 3 seconds when the device is inserted.

You need to create a configuration file tmkms.toml containing the information of the validator node to which the Kms must connect.

You can find an example file below
Create the folder for the configuration

```bash
mkdir $HOME/kms/commercio
touch $HOME/kms/commercio/tmkms.toml

```
Using the tmkms user the `$HOME` variable corresponds to the `/data_tmkms/tmkms` directory. Eventually use the full path /data_tmkms/tmkms.
 

```toml
[[chain]]
id = "commercio-3"
key_format = { type = "bech32", account_key_prefix = "did:com:", consensus_key_prefix = "did:com:valconspub" }
state_file = "/data_tmkms/tmkms/kms/commercio/commercio_priv_validator_state.json"

[[validator]]
addr = "tcp://10.1.1.254:26658" #ip del Validator Node
chain_id = "commercio-3"
reconnect = true # true is the default
secret_key = "/data_tmkms/tmkms/kms/commercio/secret_connection.key"

[[providers.yubihsm]]
adapter = { type = "http", addr = "tcp://127.0.0.1:12345" }
auth = { key = 1, password_file = "/data_tmkms/tmkms/kms/password" } # it is possible to enter the password directly using the password parameter instead of password_file
keys = [{ chain_ids = ["commercio-3"], key = 1 }]
serial_number = "9876543210" # identify serial number of a specific YubiHSM to connect to
```

For the creation of the file you need to have the following data
- Chain-id: is the identifier of the chain for which the node is being configured. In the case of testnet it will be commercio-testnet6002, in the case of mainnet commercio-mainnet
- Prefix of public addresses of the chain: in the case of commercio it will be did:com:
- Prefix of the public addresses of the nodes: in the case of commercio it will be did:com:valconspub
- Address within the vpn of the validator node: for simplicity we have assumed to have the address
- The password of our HSM device: initially the password is "password"
- The id of the key to use of the HSM: for the single configuration is `1`
- The serial Number of our device: generally it is what is indicated on the label of the YubiHSM2. They must be 10 digits. For the missing digits add zeros at the beginning
- The ip of the validator node: in the case of the example
- 
- Chain-id: è l’identificativo della chain per cui si sta configurando il nodo. Nel caso di testnet sarà commercio-testnet6002, nel caso di mainnet commercio-mainnet
- Prefisso degli indirizzi pubblici della chain: nel caso di commercio sarà did:com:
- Prefisso degli indirizzi pubblici dei nodi: nel caso di commercio sarà did:com:valconspub
- Indirizzo all’interno della vpn del nodo validatore: per semplicità abbiamo supposto di avere l’indirizzo 10.1.1.254.
- La password del nostro dispositivo HSM: inizialmente la password è “password”
- L’id della chiave da utilizzare del HSM: per la configurazione singola è 1
- Il serial Number del nostro dispositivo: generalmente è quanto indicato nell’etichetta dello YubiHSM2. Devono essere 10 cifre. Per le cifre mancanti aggiungere degli zeri davanti al seriale.
- Path delle configurazioni: configurare il path /data_tmkms/tmkms/kms/commercio.

Create the file /data_tmkms/tmkms/kms/password and enter the password "password" inside it

```bash
printf "password" > /data_tmkms/tmkms/kms/password
```


#### HSM Reset
:warning:  Warning :warning: : this procedure will do a complete reset of the device. It must not be done for a possible second installation of another node that uses the same hsm

```bash
tmkms yubihsm setup -c /data_tmkms/tmkms/kms/commercio/tmkms.toml
```

An output of this type should appear
```
This process will *ERASE* the configured YubiHSM2 and reinitialize it:

- YubiHSM serial: 9876543210

Authentication keys with the following IDs and passwords will be created:

- key 0x0001: admin:

    double section release consider diet pilot flip shell mother alone what fantasy
    much answer lottery crew nut reopen stereo square popular addict just animal

- authkey 0x0002 [operator]:  kms-operator-password-1k02vtxh4ggxct5tngncc33rk9yy5yjhk
- authkey 0x0003 [auditor]:   kms-auditor-password-1s0ynq69ezavnqgq84p0rkhxvkqm54ks9
- authkey 0x0004 [validator]: kms-validator-password-1x4anf3n8vqkzm0klrwljhcx72sankcw0
- wrapkey 0x0001 [primary]:   21a6ca8cfd5dbe9c26320b5c4935ff1e63b9ab54e2dfe24f66677aba8852be13

*** Are you SURE you want erase and reinitialize this HSM? (y/N): y
21:08:09 [WARN] factory resetting HSM device! all data will be lost!

21:08:11 [WARN] deleting temporary setup authentication key from slot 65534
     Success reinitialized YubiHSM (serial: 9876543210)
```

Confirm when asked to reinitialize the device (in bold in the output). Take note of the output, especially the 24 words provided as a new password (also in bold in the output above).
Save the new password of the file

```bash
printf "double section release consider diet pilot flip shell mother alone what fantasy much answer lottery crew nut reopen stereo square popular addict just animal" >/data_tmkms/tmkms/kms/password
```
**NB**: The password is provided on two separate lines but must be set in the file on a single line


### Key generation
The keys are generated using the tmkms command. The command is structured as follows

```bash
tmkms yubihsm keys generate 1 \
  -b steakz4u-validator-key.enc \
  -c /data_tmkms/tmkms/kms/commercio/tmkms.toml
```
Launch the command to produce the new key. The `-b` option allows you to save the key

:::danger  
:warning: The steakz4u-validator-key.enc backup file must be immediately transferred to the multiple off-line supports and secured. If it is lost, the key generated at this time can no longer be replicated.
:::

In the case of a new installation of the node, the key must be imported into the HSM. To do this, use the following command

```bash
tmkms yubihsm keys import \
  -t json \
  -i 1 priv_validator.json \
  -c /data_tmkms/tmkms/kms/commercio/tmkms.toml
```
:::danger  
:warning: The priv_validator.json backup file must be immediately transferred to the multiple, off-line, encrypted and secured supports. If it is lost, the key generated at this time can no longer be replicated.
:::

### Key confirmation
To confirm that the keys are present and configured in the HSM use the command

```bash
tmkms yubihsm keys list  -c /data_tmkms/tmkms/kms/commercio/tmkms.toml
```
Should be presented an output like the following

```
Listing keys in YubiHSM #9876543210:
- 0x0001: did:com:valconspub1zcjduepq592mn6xucyqvfrvjegruhnx55cruffkrfq0rryu809fzkgwg684qmetxxs
```
`did:com:valconspub1zcjduepq592mn6xucyqvfrvjegruhnx55cruffkrfq0rryu809fzkgwg684qmetxxs` will be the public key of the validator node. This key will be needed to perform the validator node creation transaction


At this point the KMS can be connected to the Node

Test the KMS service
```bash
tmkms start -c /data_tmkms/tmkms/kms/commercio/tmkms.toml
```
An output like this should appear.
```
Mar 05 12:20:26.781  INFO tmkms::commands::start: tmkms 0.10.0 starting up…
Mar 05 12:20:27.280  INFO tmkms::keyring: [keyring:yubihsm] added consensus key did:com:valconspub1zcjduepq592mn6xucyqvfrvjegruhnx55cruffkrfq0rryu809fzkgwg684qmetxxs
Mar 05 12:20:27.280  INFO tmkms::connection::tcp: KMS node ID: 4248B5C7755600D694C47ECEA710A2DAB743AA38
Mar 05 12:20:58.682 ERROR tmkms::client: [commercio-3@tcp://10.1.1.254:26658] I/O error: Connection timed out (os error 110)
Mar 05 12:20:59.683  INFO tmkms::connection::tcp: KMS node ID: 4248B5C7755600D694C47ECEA710A2DAB743AA38
...
```

Se l’output riporta errori diversi dal semplice fallimento della connessione allora deve essere controllata l’installazione.
NB: Il tentativo di connessione fallisce perché non abbiamo ancora configurato il nodo a cui il kms dovrebbe connettersi.
crtl+c per interrompere il processo.
