#!/bin/bash

PNAME=go-ws

ps -ef | grep $PNAME | grep -v grep | awk '{print $2}' | xargs kill
