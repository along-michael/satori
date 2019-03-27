#!/usr/bin/python
# -*- coding: utf-8 -*-

# -- stdlib --
import json
import os
import socket
import subprocess
import sys
import time

# -- third party --
# -- own --

# -- code --
ts = int(time.time())

p = subprocess.Popen("which megacli", shell=True, stdout=subprocess.PIPE)
cli = p.stdout.read().strip()

if cli:
    CLI = cli
else:
    CLI = '/opt/MegaRAID/MegaCli/MegaCli64'

if not os.path.exists(CLI):
    print json.dumps([])
    sys.exit(0)

p = subprocess.Popen(
    CLI + " -LdPdInfo -a0 | grep -Ec 'Firmware state: (Failed|Offline)'",
    shell=True, stdout=subprocess.PIPE,
)

total_err = int(p.stdout.read())

print json.dumps([{
    "metric": "megaraid.offline",
    "timestamp": ts,
    "step": 600,
    "value": total_err,
}])
