#!/bin/bash

PNAME=go-skeleton

ps -ef | grep $PNAME | grep -v grep | awk '{print $2}' | xargs kill
