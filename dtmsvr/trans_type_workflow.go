package dtmsvr

import (
	"github.com/sllt/dtm/client/dtmcli"
	"github.com/sllt/dtm/client/dtmcli/dtmimp"
	"github.com/sllt/dtm/client/dtmgrpc/dtmgimp"
	"github.com/sllt/dtm/client/workflow/wfpb"
)

type transWorkflowProcessor struct {
	*TransGlobal
}

func init() {
	registorProcessorCreator("workflow", func(trans *TransGlobal) transProcessor { return &transWorkflowProcessor{TransGlobal: trans} })
}

func (t *transWorkflowProcessor) GenBranches() []TransBranch {
	return []TransBranch{}
}

type cWorkflowCustom struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (t *transWorkflowProcessor) ProcessOnce(branches []TransBranch) error {
	if t.Status == dtmcli.StatusFailed || t.Status == dtmcli.StatusSucceed {
		return nil
	}

	cmc := cWorkflowCustom{}
	dtmimp.MustUnmarshalString(t.CustomData, &cmc)
	data := cmc.Data
	if t.Protocol == dtmimp.ProtocolGRPC {
		wd := wfpb.WorkflowData{Data: cmc.Data}
		data = dtmgimp.MustProtoMarshal(&wd)
	}
	return t.getURLResult(t.QueryPrepared, "00", cmc.Name, data)
}
