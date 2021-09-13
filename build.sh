#!/bin/bash

PNAME=sos-detection-to-protection

go build cmd/main.go && mv main $PNAME
