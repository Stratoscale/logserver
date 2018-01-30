package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Parallel()
	time1, err := time.Parse(time.RFC3339, "2017-12-25T16:23:05+02:00")
	require.Nil(t, err)
	time2, err := time.Parse("2006-01-02T15:04:05.000", "2017-12-25T16:23:05.123")
	require.Nil(t, err)
	time3, err := time.Parse("2006-01-02T15:04:05.000", "0000-12-25T16:23:05.123")
	require.Nil(t, err)

	parsers, err := New([]Config{
		{
			Glob: "*.jsonlog",
			JsonMapping: map[string]string{
				"msg":   "msg",
				"level": "levelname",
				"time":  "created",
				"args":  "args",
			},
			TimeFormats: []string{"unix_float"},
		},
		{
			Regexp:      `(?P<time>\d{4}-\d{2}-\d{2}\W\d{2}:\d{2}:\d{2},\d{3}).\d{3}\W\d+\W(?P<level>[[:alpha:]]+)\W(?P<module>[^\.]+)\.(?P<function>[^\W]+)\W\[[^]+]\]\W(?P<msg>.*)`,
			TimeFormats: []string{"2006-01-02 15:04:05.000"},
		},
		{
			Glob: "*.jsonappend",
			JsonMapping: map[string]string{
				"msg":   "msg",
				"level": "level",
				"time":  "ts",
			},
			AppendArgs:  true,
			TimeFormats: []string{"01/02-15:04:05.000"},
		},
	})
	require.Nil(t, err)

	tests := []struct {
		name    string
		logName string
		line    string
		want    *Log
	}{
		{
			name:    "jsonlog/list args",
			logName: "file.jsonlog",
			line:    `{"process": 36154, "args": ["val1", 0, true], "module": "distributor", "funcName": "_updateDisksInTableStates", "exc_text": null, "name": "my-name", "thread": 140189066323712, "created": 1514211785.448693, "threadName": "DistributorThread", "filename": "distributor.py", "levelno": 20, "processName": "MainProcess", "pathname": "/var/lib/file.py", "lineno": 162, "msg": "hello %s you got:%s yes? %s, sure", "levelname": "INFO"}`,
			want: &Log{
				Msg:   `hello val1 you got:0 yes? true, sure`,
				Time:  &time1,
				Level: "INFO",
			},
		},
		{
			name:    "jsonlog/map args",
			logName: "file.jsonlog",
			line:    `{"process": 30319, "args": {"hostname": "rabbitmq-server.service.strato", "userid": "guest", "password": "guest", "virtual_host": "/", "port": 5672, "insist": false, "ssl": false, "transport": "amqp", "connect_timeout": 5, "transport_options": {"on_blocked": "<function _on_connection_blocked at 0x343cd70>", "on_unblocked": "<function _on_connection_unblocked at 0x343cde8>", "confirm_publish": true}, "login_method": "AMQPLAIN", "uri_prefix": null, "heartbeat": 60.0, "failover_strategy": "shuffle", "alternates": []}, "module": "impl_rabbit", "funcName": "__init__", "exc_text": null, "extra_keys": ["project", "version"], "project": "unknown", "name": "oslo.messaging._drivers.impl_rabbit", "thread": 140161415309056, "created": 1514211785.448693, "threadName": "Thread-46", "filename": "impl_rabbit.py", "levelno": 20, "processName": "MainProcess", "version": "unknown", "pathname": "/usr/lib/python2.7/site-packages/oslo_messaging/_drivers/impl_rabbit.py", "lineno": 483, "msg": "Connecting to AMQP server on %(hostname)s:%(port)s", "levelname": "INFO"}`,
			want: &Log{
				Msg:   "Connecting to AMQP server on rabbitmq-server.service.strato:5672",
				Time:  &time1,
				Level: "INFO",
			},
		},
		{
			name:    "jsonlog/map into %s",
			logName: "file.jsonlog",
			line:    `{"args": {"key": "value", "key1": "value1"}, "levelname": "INFO", "msg": "Message: %s", "created": 1514211785.448693}`,
			want: &Log{
				Msg:   `Message: {"key":"value","key1":"value1"}`,
				Time:  &time1,
				Level: "INFO",
			},
		},
		{
			name:    "jsonlog/list into %s",
			logName: "file.jsonlog",
			line:    `{"args": ["arg1", "arg2"], "levelname": "INFO", "msg": "Message: %s", "created": 1514211785.448693}`,
			want: &Log{
				Msg:   `Message: ["arg1","arg2"]`,
				Time:  &time1,
				Level: "INFO",
			},
		},
		{
			name:    "openstack",
			logName: "optnstack.log",
			line:    "2017-12-25 16:23:05,123.123 33983 WARN oslo_service.periodic_task [-] Skipping periodic task _periodic_update_dns because its interval is negative",
			want: &Log{
				Msg:   "Skipping periodic task _periodic_update_dns because its interval is negative",
				Time:  &time2,
				Level: "WARN",
			},
		},
		{
			name:    "append",
			logName: "log.jsonappend",
			line:    `{"msg":"hi", "ts":"12/25-16:23:05.123","level":"info","arg1":"value1"}`,
			want: &Log{
				Msg:   "hi arg1=value1",
				Time:  &time3,
				Level: "info",
			},
		},
		{
			name:    "default",
			logName: "optnstack.jsonlog",
			line:    "some log",
			want:    &Log{Msg: "some log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parsers.Parse(tt.logName, []byte(tt.line), &Memory{}))
		})
	}
}
