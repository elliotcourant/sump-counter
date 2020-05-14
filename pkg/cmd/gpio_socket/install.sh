#!/usr/bin/env bash

service gpio_socket stop
cp /home/ubuntu/bin/gpio_socket.service /lib/systemd/system/gpio_socket.service
systemctl enable gpio_socket.service
service gpio_socket start