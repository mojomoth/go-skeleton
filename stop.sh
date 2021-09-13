#!/bin/bash

PNAME=sos-detection-to-protection

ps -ef | grep $PNAME | grep -v grep | awk '{print $2}' | xargs kill
