package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStratoFormat(t *testing.T) {
	line, err := GetParser("stratolog")(`{"process": 36154, "args": ["<disk: hostname=stratonode0.node.strato, ID=f3d510c7-1185-4942-b349-0de055165f78, path=/dev/sdc, type=mancala>", 0, 1], "module": "distributor", "funcName": "_updateDisksInTableStates", "exc_text": null, "name": "mancala.management.distributor.distributor", "thread": 140189066323712, "created": 1514211785.448693, "threadName": "DistributorThread", "filename": "distributor.py", "levelno": 20, "processName": "MainProcess", "pathname": "/usr/share/stratostorage/mancala_management_service.egg/mancala/management/distributor/distributor.py", "lineno": 162, "msg": "data disk %s was found in distributionID:%s table version:%s, setting inTable=True", "levelname": "INFO"}`)
	fmt.Println(line)
	require.Nil(t, err)
}
