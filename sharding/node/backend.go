package node

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/sharding"
	"github.com/ethereum/go-ethereum/sharding/mainchain"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/sharding/params"
	"github.com/ethereum/go-ethereum/sharding/txpool"
	cli "gopkg.in/urfave/cli.v1"
)

// ShardEthereum is a service that is registered and started when geth is launched.
// it contains APIs and fields that handle the different components of the sharded
// Ethereum network.
type ShardEthereum struct {
	shardConfig  *params.ShardConfig  // Holds necessary information to configure shards.
	txPool       *txpool.ShardTxPool  // Defines the sharding-specific txpool. To be designed.
	actor        sharding.Actor       // Either notary, proposer, or observer.
	shardChainDb ethdb.Database       // Access to the persistent db to store shard data.
	eventFeed    *event.Feed          // Used to enable P2P related interactions via different sharding actors.
	smcClient    *mainchain.SMCClient // Provides bindings to the SMC deployed on the Ethereum mainchain.

	// Lifecycle and service stores.
	serviceFuncs []sharding.ServiceConstructor     // Service constructors (in dependency order).
	services     map[reflect.Type]sharding.Service // Currently running services.
	lock         sync.RWMutex
}

// New creates a new sharding-enabled Ethereum instance. This is called in the main
// geth sharding entrypoint.
func New(ctx *cli.Context) (*ShardEthereum, error) {

	shardEthereum := &ShardEthereum{}

	path := node.DefaultDataDir()
	if ctx.GlobalIsSet(utils.DataDirFlag.Name) {
		path = ctx.GlobalString(utils.DataDirFlag.Name)
	}

	endpoint := ctx.Args().First()
	if endpoint == "" {
		endpoint = fmt.Sprintf("%s/%s.ipc", path, mainchain.ClientIdentifier)
	}
	if ctx.GlobalIsSet(utils.IPCPathFlag.Name) {
		endpoint = ctx.GlobalString(utils.IPCPathFlag.Name)
	}

	passwordFile := ctx.GlobalString(utils.PasswordFileFlag.Name)
	depositFlag := ctx.GlobalBool(utils.DepositFlag.Name)

	smcClient, err := mainchain.NewSMCClient(endpoint, path, depositFlag, passwordFile)
	if err != nil {
		return nil, err
	}

	// Adds the initialized SMCClient to the ShardEthereum instance.
	shardEthereum.smcClient = smcClient

	if err := shardEthereum.registerShardingServices(); err != nil {
		return nil, err
	}

	return shardEthereum, nil
}

// Start the ShardEthereum service and kicks off the p2p and actor's main loop.
func (s *ShardEthereum) Start() error {
	log.Println("Starting sharding service")
	if err := s.actor.Start(); err != nil {
		return err
	}
	defer s.actor.Stop()

	// TODO: start p2p and other relevant services.
	return nil
}

// Close handles graceful shutdown of the system.
func (s *ShardEthereum) Close() error {
	return nil
}

// SMCClient returns an instance of a client that communicates to a mainchain node via
// RPC and provides helpful bindings to the Sharding Manager Contract.
func (s *ShardEthereum) SMCClient() *mainchain.SMCClient {
	return s.smcClient
}

// Register appends a service constructor function to the service registry of the
// sharding node.
func (s *ShardEthereum) Register(constructor sharding.ServiceConstructor) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serviceFuncs = append(s.serviceFuncs, constructor)
	return nil
}

func (s *ShardEthereum) registerShardingServices() error {
	// TODO: Register the different services to the node.
	s.Register(func(ctx *sharding.ServiceContext) (sharding.Service, error) {
		return nil, nil
	})
	return nil
}
