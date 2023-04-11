package dtmgrpc

import (
	"reflect"
	"testing"

	"github.com/sllt/dtm/client/dtmcli"
)

// TestNewMsgGrpc ut for NewMsgGrpc
func TestNewMsgGrpc(t *testing.T) {
	var (
		server            = "dmt_server_address"
		gidNoOptions      = "msg_no_setup_options"
		gidTraceIDXXX     = "msg_setup_options_trace_id_xxx"
		msgWithTraceIDXXX = &MsgGrpc{Msg: *dtmcli.NewMsg(server, gidTraceIDXXX)}
		traceIDHeaders    = map[string]string{
			"x-trace-id": "xxx",
		}
	)
	msgWithTraceIDXXX.BranchHeaders = traceIDHeaders
	type args struct {
		gid  string
		opts []TransBaseOption
	}
	tests := []struct {
		name string
		args args
		want *MsgGrpc
	}{
		{
			name: "no setup options",
			args: args{gid: gidNoOptions},
			want: &MsgGrpc{Msg: *dtmcli.NewMsg(server, gidNoOptions)},
		},
		{
			name: "msg with trace_id",
			args: args{
				gid: gidTraceIDXXX,
				opts: []TransBaseOption{
					WithBranchHeaders(traceIDHeaders),
				},
			},
			want: msgWithTraceIDXXX,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMsgGrpc(server, tt.args.gid, tt.args.opts...)
			t.Logf("TestNewMsgGrpc %s got %+v\n", tt.name, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMsgGrpc() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewSagaGrpc ut for NewSagaGrpc
func TestNewSagaGrpc(t *testing.T) {
	var (
		server             = "dmt_server_address"
		gidNoOptions       = "msg_no_setup_options"
		gidTraceIDXXX      = "msg_setup_options_trace_id_xxx"
		sagaWithTraceIDXXX = &SagaGrpc{Saga: *dtmcli.NewSaga(server, gidTraceIDXXX)}
		traceIDHeaders     = map[string]string{
			"x-trace-id": "xxx",
		}
	)
	sagaWithTraceIDXXX.BranchHeaders = traceIDHeaders
	type args struct {
		gid  string
		opts []TransBaseOption
	}
	tests := []struct {
		name string
		args args
		want *SagaGrpc
	}{
		{
			name: "no setup options",
			args: args{gid: gidNoOptions},
			want: &SagaGrpc{Saga: *dtmcli.NewSaga(server, gidNoOptions)},
		},
		{
			name: "msg with trace_id",
			args: args{
				gid: gidTraceIDXXX,
				opts: []TransBaseOption{
					WithBranchHeaders(traceIDHeaders),
				},
			},
			want: sagaWithTraceIDXXX,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSagaGrpc(server, tt.args.gid, tt.args.opts...)
			t.Logf("TestNewSagaGrpc %s got %+v\n", tt.name, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSagaGrpc() = %v, want %v", got, tt.want)
			}
		})
	}
}
