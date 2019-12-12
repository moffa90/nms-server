#!/bin/bash

service="$(connmanctl services | grep ethernet | awk '{print $3}')"

connmanctl config $service --ipv4 dhcp
