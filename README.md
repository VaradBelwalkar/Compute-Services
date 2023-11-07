## Compute Services

### Build labs faster

Building a lab environment can be time-consuming.  Our platform provides,


- Quick on-demand containers and provisioning.
- Comprehensive list of various development containers.
- Simplified service deployment with services accessible throughout network
- Automatic DNS manipulation, thereby accessing services with choosen domain-names
- access the labs even from remote locations, no need to worry about behind-NAT communication anymore!
- Stronger Authentication mechanism.
- Storage mechanism via docker volumes.
- Suitable for large scale deployments (Currently single node).

### Wargames
- Host container-based wargames easily on this platform.

### Unified Networking
- Gone those days where you need to setup the whole network configuration yourself to make nodes work together, our **platform** lets you focus on your work, taking care of the networking.
- Best for learning networking/Cloud tools that normally require manual configuration and are time consuming.

### Run a honeypot

## Initial Setup:
- Setup host bridging with separate host virtual interface, and install docker plugin as explained here: [Initial Setup](https://github.com/VaradBelwalkar/Compute-Services/blob/master/configure/setup.md)

## How to Deploy?
**Follow [this](https://hub.docker.com/r/varadbelwalkar/golang_server) docker image for production deployment (Recommended)**

Simply create the files as mentioned on the docker hub, edit the files accordingly, and then simply run,
```
$docker-compose up
```



