What it provides?\
1.) Quick on-demand containers and provisioning\
2.) Stronger Authentication mechanism\
3.) Storage mechanism via docker volumes\
4.) Suitable for large scale deployments (Currently single node)\
5.) More features coming soon...

      Following one shows some features:
      
Setup: \
       Setup Virtual machine with Host-bridged network\
       Run this as daemon inside virtual machine\
       Make sure to give virtual machine static ip (not required) to be consistent regarding deployment\
       
Limitations (Currently)

Currently all containers share same port pool from the host, if one container already using port x, other service running in another container cannot use it

Assigning containers separate ip addresses from dhcp is still in progress to make them on the same network as host, that feature coming soon...
      

https://user-images.githubusercontent.com/86964576/218269144-a1405ff5-0fad-4c00-843d-1dec5c323137.mp4

