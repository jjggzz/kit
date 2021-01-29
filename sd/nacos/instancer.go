// author: JGZ
// time:   2021-01-28 10:26
package nacos

import (
	"errors"
	"fmt"
	"github.com/jjggzz/kit/log"
	"github.com/jjggzz/kit/sd"
	"github.com/jjggzz/kit/sd/internal/instance"
	"github.com/jjggzz/kit/util/conn"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"time"
)

// errStopped notifies the loop to quit. aka stopped via quitc
var errStopped = errors.New("quit and closed nacos instancer")

type Instancer struct {
	cache     *instance.Cache
	client    naming_client.INamingClient
	logger    log.Logger
	service   string
	groupName string   // 默认值DEFAULT_GROUP
	clusters  []string // 默认值DEFAULT
	quitc     chan struct{}
}

func (s *Instancer) Register(ch chan<- sd.Event) {
	s.cache.Register(ch)
}

func (s *Instancer) Deregister(ch chan<- sd.Event) {
	s.cache.Deregister(ch)
}

func (s *Instancer) Stop() {
	close(s.quitc)
}

func NewInstancer(client naming_client.INamingClient, logger log.Logger, service string, groupName string, clusters []string) *Instancer {
	s := &Instancer{
		cache:     instance.NewCache(),
		client:    client,
		logger:    logger,
		groupName: groupName,
		clusters:  clusters,
		service:   service,
		quitc:     make(chan struct{}),
	}
	instances, err := s.getInstances(nil)
	if err == nil {
		s.logger.Log("instances", len(instances))
	} else {
		s.logger.Log("err", err)
	}
	s.cache.Update(sd.Event{Instances: instances, Err: err})
	go s.loop()
	return s
}

func (s *Instancer) loop() {
	var (
		instances []string
		err       error
		d         = 10 * time.Millisecond
	)
	for {
		instances, err = s.getInstances(s.quitc)
		switch {
		case err == errStopped:
			return
		case err != nil:
			s.logger.Log("err", err)
			time.Sleep(d)
			d = conn.Exponential(d)
			s.cache.Update(sd.Event{Err: err})
		default:
			s.cache.Update(sd.Event{Instances: instances})
			d = 10 * time.Millisecond
		}
	}

}

func (s *Instancer) getInstances(interruptc chan struct{}) ([]string, error) {
	var (
		errc = make(chan error, 1)
		resc = make(chan []string, 1)
	)

	go func() {
		instanceList, err := s.client.SelectInstances(vo.SelectInstancesParam{
			ServiceName: s.service,
			GroupName:   s.groupName, // 默认值DEFAULT_GROUP
			Clusters:    s.clusters,  // 默认值DEFAULT
			HealthyOnly: true,
		})
		if err != nil {
			errc <- err
			return
		}
		var entities []string
		for _, e := range instanceList {
			entities = append(entities, fmt.Sprintf("%s:%d", e.Ip, e.Port))
		}
		resc <- entities
	}()

	select {
	case err := <-errc:
		return nil, err
	case res := <-resc:
		return res, nil
	case <-interruptc:
		return nil, errStopped
	}
}
