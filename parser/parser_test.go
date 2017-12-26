package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStratoFormat(t *testing.T) {
	line, err := GetParser(".stratolog")([]byte(`{"process": 36154, "args": ["<disk: hostname=stratonode0.node.strato, ID=f3d510c7-1185-4942-b349-0de055165f78, path=/dev/sdc, type=mancala>", 0, 1], "module": "distributor", "funcName": "_updateDisksInTableStates", "exc_text": null, "name": "mancala.management.distributor.distributor", "thread": 140189066323712, "created": 1514211785.448693, "threadName": "DistributorThread", "filename": "distributor.py", "levelno": 20, "processName": "MainProcess", "pathname": "/usr/share/stratostorage/mancala_management_service.egg/mancala/management/distributor/distributor.py", "lineno": 162, "msg": "data disk %s was found in distributionID:%s table version:%s, setting inTable=True", "levelname": "INFO"}`))
	fmt.Println(line)
	require.Nil(t, err)

	// TODO: (filanov) add tests
}

func TestStratoFormat2(t *testing.T) {
	line, err := GetParser(".stratolog")([]byte(`{"process": 30319, "args": {"hostname": "rabbitmq-server.service.strato", "userid": "guest", "password": "guest", "virtual_host": "/", "port": 5672, "insist": false, "ssl": false, "transport": "amqp", "connect_timeout": 5, "transport_options": {"on_blocked": "<function _on_connection_blocked at 0x343cd70>", "on_unblocked": "<function _on_connection_unblocked at 0x343cde8>", "confirm_publish": true}, "login_method": "AMQPLAIN", "uri_prefix": null, "heartbeat": 60.0, "failover_strategy": "shuffle", "alternates": []}, "module": "impl_rabbit", "funcName": "__init__", "exc_text": null, "extra_keys": ["project", "version"], "project": "unknown", "name": "oslo.messaging._drivers.impl_rabbit", "thread": 140161415309056, "created": 1514286682.533927, "threadName": "Thread-46", "filename": "impl_rabbit.py", "levelno": 20, "processName": "MainProcess", "version": "unknown", "pathname": "/usr/lib/python2.7/site-packages/oslo_messaging/_drivers/impl_rabbit.py", "lineno": 483, "msg": "Connecting to AMQP server on %(hostname)s:%(port)s", "levelname": "INFO"}`))
	require.Nil(t, err)
	fmt.Println(line.Msg)
	// TODO: (filanov) add tests
}
