package node

import (
	"ants/http"
	"time"
)

/*
what a Distributer do
*	status,running|parse|stop
*	distribute a request,by some strategy
*	*
*/
const (
	DISTRIBUTE_RUNING = iota
	DISTRIBUTE_PAUSE
	DISTRIBUTE_STOP
)

type Distributer struct {
	Status    int
	Cluster   *Cluster
	Node      *Node
	LastIndex int
}

func NewDistributer(cluster *Cluster, node *Node) *Distributer {
	return &Distributer{DISTRIBUTE_STOP, cluster, node, 0}
}

func (this *Distributer) IsStop() bool {
	return this.Status == DISTRIBUTE_STOP
}

func (this *Distributer) IsPause() bool {
	return this.Status == DISTRIBUTE_PAUSE
}
func (this *Distributer) Pause() {
	this.Status = DISTRIBUTE_PAUSE
}
func (this *Distributer) Stop() {
	this.Status = DISTRIBUTE_STOP
}

func (this *Distributer) Start() {
	if this.Status == DISTRIBUTE_RUNING {
		return
	}
	this.Status = DISTRIBUTE_RUNING
	this.Run()
}

// dead loop cluster pop request
func (this *Distributer) Run() {
	for {
		if this.IsStop() {
			break
		}
		if this.IsPause() {
			time.Sleep(1 * time.Second)
			continue
		}
		request := this.Cluster.PopRequest()
		if request == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		this.Distribute(request)
		this.Node.DistributeRequest(request)
	}
}

// if cookiejar > 0 means it require cookie context ,so we should send it to where it come from
// else distribute it by order
func (this *Distributer) Distribute(request *http.Request) {
	if request.CookieJar > 0 {
		return
	} else {
		if this.LastIndex >= len(this.Cluster.ClusterInfo.NodeList) {
			this.LastIndex = 0
		}
		nodeName := this.Cluster.ClusterInfo.NodeList[this.LastIndex].Name
		request.NodeName = nodeName
		this.LastIndex += 1
	}
}