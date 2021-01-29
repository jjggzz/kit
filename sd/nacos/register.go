// author: JGZ
// time:   2021-01-28 15:51
package nacos

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Param struct {
	Ip          string            `param:"ip"`          //required
	Port        uint64            `param:"port"`        //required
	ServiceName string            `param:"serviceName"` //required
	Weight      float64           `param:"weight"`      //required,it must be lager than 0
	Metadata    map[string]string `param:"metadata"`    //optional
	ClusterName string            `param:"clusterName"` //optional,default:DEFAULT
	GroupName   string            `param:"groupName"`   //optional,default:DEFAULT_GROUP
}

type Registrar struct {
	client naming_client.INamingClient
	param  Param
	logger log.Logger
}

func NewRegistrar(client naming_client.INamingClient, param Param, logger log.Logger) *Registrar {
	return &Registrar{
		client: client,
		param:  param,
		logger: log.With(logger, "service", param.ServiceName, "group", fmt.Sprint(param.GroupName), "address", fmt.Sprintf("%s:%d", param.Ip, param.Port)),
	}
}

// Register implements sd.Registrar interface.
func (p *Registrar) Register() {
	param := vo.RegisterInstanceParam{}
	// field copy
	param.Ip = p.param.Ip
	param.Port = p.param.Port
	param.ServiceName = p.param.ServiceName
	param.Weight = p.param.Weight
	param.Metadata = p.param.Metadata
	if p.param.ClusterName != "" {
		param.ClusterName = p.param.ClusterName
	}
	if p.param.GroupName != "" {
		param.GroupName = p.param.GroupName
	}
	param.Enable = true
	param.Healthy = true
	param.Ephemeral = true

	success, err := p.client.RegisterInstance(param)
	if err != nil {
		p.logger.Log("err", err)
	} else {
		p.logger.Log("register", success)
	}
}

// Deregister implements sd.Registrar interface.
func (p *Registrar) Deregister() {
	param := vo.DeregisterInstanceParam{}
	// field copy
	param.Ip = p.param.Ip
	param.Port = p.param.Port
	param.ServiceName = p.param.ServiceName
	if p.param.ClusterName != "" {
		param.Cluster = p.param.ClusterName
	}
	if p.param.GroupName != "" {
		param.GroupName = p.param.GroupName
	}
	param.Ephemeral = true

	success, err := p.client.DeregisterInstance(param)
	if err != nil {
		p.logger.Log("err", err)
	} else {
		p.logger.Log("deregister", success)
	}
}
