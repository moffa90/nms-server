package worker

import (
	"github.com/moffa90/nms-server/db/models"
	"sync"
	"time"
)

type CommonWorkRequest struct {
	Action    uint
	Wg        *sync.WaitGroup
	CreatedAt time.Time
}

type SerialWorkRequest struct {
	*CommonWorkRequest
	Device     models.Hardware
	RawCommand string
	Response   map[string]string
}

type SNMPWorkRequest struct {
	*CommonWorkRequest
	Remote    models.Remote
	Community string
	Response   string
}

type WorkRequest interface {
	GetAction() uint
	GetWg() *sync.WaitGroup
	GetCreatedAt() time.Time
	GetRequestInfo() *CommonWorkRequest
}

func (cwq *CommonWorkRequest) GetAction() uint {
	return cwq.Action
}

func (cwq *CommonWorkRequest) GetWg() *sync.WaitGroup {
	return cwq.Wg
}

func (cwq *CommonWorkRequest) GetCreatedAt() time.Time {
	return cwq.CreatedAt
}

func (cwq *CommonWorkRequest) GetRequestInfo() *CommonWorkRequest {
	return cwq
}

func NewSerialWorkRequest(device models.Hardware, action uint, rawCommand string, wg *sync.WaitGroup) *SerialWorkRequest {
	return &SerialWorkRequest{
		CommonWorkRequest: &CommonWorkRequest{
			Action:    action,
			Wg:        wg,
			CreatedAt: time.Now(),
		},
		Device:     device,
		RawCommand: rawCommand,
	}
}

func NewSNMPWorkRequest(remote models.Remote, community string, action uint, wg *sync.WaitGroup) *SNMPWorkRequest {
	return &SNMPWorkRequest{
		CommonWorkRequest: &CommonWorkRequest{
			Action:    action,
			Wg:        wg,
			CreatedAt: time.Now(),
		},
		Remote: remote,
		Community: community,
	}
}

type WorkQueue struct {
	Pending chan WorkRequest
}

func (r *WorkQueue)AddWork(request WorkRequest){
	r.Pending <- request
}