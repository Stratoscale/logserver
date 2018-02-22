#!/bin/bash

docker run -d --restart=always --name logserver2 --net host -v /opt/logserver2.0/logserver.json:/logserver.json -v /mnt/logs-netapp:/logs rackattack-nas.dc1:5000/logserver:latest -addr :8005 -debug -dynamic
