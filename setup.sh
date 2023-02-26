#!/bin/bash

################### ************IMPORTANT************ ###################
#   DO NOT CHANGE THE **PASSWORD_HASH_SECRET** ONCE RUNNING, IF CHANGED,AFTER RESTART,
#   YOU WON'T BE ABLE TO AUTHENTICATE USERS BASED ON THEIR PASSWORDS
#   RECOMMENDED: BACKUP **PASSWORD_HASH_SECRET** 

# Set MongoDB URL and password
export MONGODB_URI="mongodb://user:pwd@localhost:27017" #Must provide Username and Password in the URL

# Set Redis URL and password
export REDIS_URL="localhost:6379"
export REDIS_PASSWORD="DLKJ340jfdK0934jkljdfkldsdf09klsdfj"

# Set JWT secret key
export JWT_SECRET="DFSD,asdlfkjFSDFJRWFKDvkdlsLKjsde5tToDogge"
echo $JWT_SECRET
# Set password hashing secret key
export PASSWORD_HASH_SECRET="dsfJ34klasdjf0wjoi334FDLKJ23309j4dlj"


#Setup the Redis password by first locating /etc/redis/redis.conf file, then,

#Uncomment this line "# requirepass foobared"

#And make it to "requirepass <your_super_secure_password>"