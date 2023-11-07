NOTE: Wireless interfaces do not work with this setup, so ensure you have physical connection (ethernet, tethering) to ensure the deployment works properly

### Install docker and docker-compose, configure user etc.
```
#apt install docker.io docker-compose
#usermod -aG docker $(whoami)
#newgrp docker
```

### (Assuming your default interface is eth0)
```
#ip link add dhcp-bridge type bridge
#ip link set dhcp-bridge up
#ip link set eth0 master dhcp-bridge
#dhcpcd dhcp-bridge
```
### Install network driver and configure network to work with.
```
$docker plugin install ghcr.io/devplayer0/docker-net-dhcp:release-linux-amd64
$docker network create -d ghcr.io/devplayer0/docker-net-dhcp:release-linux-amd64 --ipam-driver null -o bridge=docker-dhcp -o ipv6=true my-dhcp-net
```
