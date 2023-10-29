#!/bin/bash


#Remove previous images

docker rmi varadbelwalkar/golang_server
docker rmi varadbelwalkar/odc_swarm_server


#Build the image

docker image build -t varadbelwalkar/golang_server .
docker image tag varadbelwalkar/golang_server varadbelwalkar/odc_swarm_server

#Push the images to repository

docker image push varadbelwalkar/golang_server
docker image push varadbelwalkar/odc_swarm_server