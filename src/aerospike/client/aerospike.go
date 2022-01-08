package client

import (
	"aerospike-cli/src/logger"
	"aerospike-cli/src/property"
	"aerospike-cli/src/util"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	"os"
	"regexp"
	"time"
)

const (
	HostEnv = "aerospike.client.host.%s"
	PortEnv = "aerospike.client.port.%s"
	TtlSec  = "aerospike.connection.timeout.sec"
)

type Aerospike struct {
	client     *as.Client
	policy     *as.BasePolicy
	connection *as.Connection
}

func New(env string) *Aerospike {
	return &Aerospike{
		client:     getClient(env),
		connection: getConnection(env),
		//policy:           getPolicy(),
	}
}

func getConnection(env string) *as.Connection {
	hostProp := property.Props.MustGetString(fmt.Sprintf(HostEnv, env))
	portProp := property.Props.MustGetInt(fmt.Sprintf(PortEnv, env))
	cp := as.NewClientPolicy()
	cp.Timeout = 30 * time.Second
	conn, err := as.NewConnection(cp, as.NewHost(hostProp, portProp))
	if err != nil {
		util.PrintErrorToConsole("error while trying to connect aerospike, env': %s, host: %s, port: %d, error: %+v", env, hostProp, portProp, err.Error())
		os.Exit(-1)
	}
	conn.SetTimeout(time.Time{}, 30*time.Second)
	return conn
}

func getPolicy() *as.BasePolicy {
	connectionTtlProp := property.Props.MustGetInt(fmt.Sprintf(TtlSec))
	policy := as.NewPolicy()
	policy.SocketTimeout = time.Second * time.Duration(connectionTtlProp)
	return policy
}

func getClient(env string) *as.Client {
	hostProp := property.Props.MustGetString(fmt.Sprintf(HostEnv, env))
	portProp := property.Props.MustGetInt(fmt.Sprintf(PortEnv, env))
	connectionTtlProp := property.Props.MustGetInt(fmt.Sprintf(TtlSec))

	clientPolicy := as.NewClientPolicy()
	clientPolicy.ConnectionQueueSize = 64
	clientPolicy.LimitConnectionsToQueueSize = true
	clientPolicy.Timeout = time.Duration(connectionTtlProp) * time.Second
	client, err := as.NewClientWithPolicy(clientPolicy, hostProp, portProp)

	if err != nil {
		util.PrintErrorToConsole("error while trying to connect aerospike, env': %s, host: %s, port: %d, error: %+v", env, hostProp, portProp, err.Error())
		os.Exit(-1)
	}
	logger.Info.Printf("Success to create aerospike client, env': %s, host: %s, port: %d", env, hostProp, portProp)
	return client
}

func getClientPolicy() *as.ClientPolicy {
	policy := as.NewClientPolicy()
	policy.LoginTimeout = 30 * time.Second
	return policy
}

func (aerospike *Aerospike) ScanAll(namespace string, setName string) (*as.Recordset, error) {
	recordset, err := aerospike.client.ScanAll(nil, namespace, setName)
	if err != nil {
		logger.Error.Printf("Error while trying to scanAll, setName: %s, %s", setName, err)
		err := recordset.Close()
		if err != nil {
			logger.Error.Printf("Error while trying to close recordset, setName: %s, %s", setName, err)
			return nil, err
		}
		return nil, err
	}
	return recordset, nil
}

func (aerospike *Aerospike) FindObj(namespace string, setName string, keyName string, obj interface{}) error {
	key, _ := as.NewKey(namespace, setName, keyName)
	err := aerospike.client.GetObject(aerospike.policy, key, obj)
	if err != nil {
		logger.Error.Printf("Error while trying to get data, setName: %s, keyName: %s, %s", setName, keyName, err)
	}
	return err
}

func (aerospike *Aerospike) Find(namespace string, setName string, keyName string) (*as.Record, error) {
	key, _ := as.NewKey(namespace, setName, keyName)
	record, err := aerospike.client.Get(aerospike.policy, key)
	if err != nil {
		logger.Error.Printf("Error while trying to get data, setName: %s, keyName: %s, %s", setName, keyName, err)
		return nil, err
	}
	return record, nil
}

func (aerospike *Aerospike) Update(namespace string, setName string, keyName string, bins ...*as.Bin) error {
	key, _ := as.NewKey(namespace, setName, keyName)
	err := aerospike.client.PutBins(nil, key, bins...)
	if err != nil {
		logger.Error.Printf("Error while trying to updating bin, setName: %s, keyName: %s, %s", setName, keyName, err)
	}
	return err
}

func (aerospike *Aerospike) GetSets(namespace string) ([]string, error) {
	infoMap, err := as.RequestInfo(aerospike.connection, fmt.Sprintf("sets/%s", namespace))
	if err != nil {
		logger.Error.Printf("GetSets, failed, %s", err.Error())
		return nil, err
	}

	var sets []string
	for _, v := range infoMap {
		pat := regexp.MustCompile(`(set=)(.*?)(:)`)
		submatch := pat.FindAllStringSubmatch(v, -1)
		for _, set := range submatch {
			sets = append(sets, set[2])
		}
	}
	return sets, nil
}

func (aerospike *Aerospike) Close() {
	util.PrintInfoToConsoleNewLine("closing aerospike client...")
	aerospike.client.Close()
}

func (aerospike *Aerospike) GetAllKeysBySet(namespace string, setName string) (*as.Recordset, error) {
	policy := as.NewScanPolicy()
	policy.IncludeBinData = false
	recordset, err := aerospike.client.ScanAll(policy, namespace, setName)
	if err != nil {
		logger.Error.Printf("Error while trying to get all keys by set, setName: %s, %s", setName, err)
		err := recordset.Close()
		if err != nil {
			logger.Error.Printf("Error while trying to close recordset, setName: %s, %s", setName, err)
			return nil, err
		}
		return nil, err
	}
	return recordset, nil
}
