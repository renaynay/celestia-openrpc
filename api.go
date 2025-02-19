package client

import (
	"context"

	"github.com/celestiaorg/go-fraud"
	libhead "github.com/celestiaorg/go-header"
	"github.com/celestiaorg/go-header/sync"
	"github.com/celestiaorg/rsmt2d"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"

	"github.com/rollkit/celestia-openrpc/types/blob"
	"github.com/rollkit/celestia-openrpc/types/das"
	"github.com/rollkit/celestia-openrpc/types/header"
	"github.com/rollkit/celestia-openrpc/types/node"
	"github.com/rollkit/celestia-openrpc/types/share"
	"github.com/rollkit/celestia-openrpc/types/state"
)

// Proof embeds the fraud.Proof interface type to provide a concrete type for JSON serialization.
type Proof struct {
	fraud.Proof[*header.ExtendedHeader]
}

type FraudAPI struct {
	Subscribe func(context.Context, fraud.ProofType) (<-chan *Proof, error) `perm:"read"`
	Get       func(context.Context, fraud.ProofType) ([]Proof, error)       `perm:"read"`
}

type DASAPI struct {
	SamplingStats func(ctx context.Context) (das.SamplingStats, error) `perm:"read"`
	WaitCatchUp   func(ctx context.Context) error                      `perm:"read"`
}

type BlobAPI struct {
	Submit   func(context.Context, []*blob.Blob, *SubmitOptions) (uint64, error)                        `perm:"write"`
	Get      func(context.Context, uint64, share.Namespace, blob.Commitment) (*blob.Blob, error)        `perm:"read"`
	GetAll   func(context.Context, uint64, []share.Namespace) ([]*blob.Blob, error)                     `perm:"read"`
	GetProof func(context.Context, uint64, share.Namespace, blob.Commitment) (*blob.Proof, error)       `perm:"read"`
	Included func(context.Context, uint64, share.Namespace, *blob.Proof, blob.Commitment) (bool, error) `perm:"read"`
}

type HeaderAPI struct {
	LocalHead func(context.Context) (*header.ExtendedHeader, error) `perm:"read"`
	GetByHash func(
		ctx context.Context,
		hash libhead.Hash,
	) (*header.ExtendedHeader, error) `perm:"read"`
	GetVerifiedRangeByHeight func(
		context.Context,
		*header.ExtendedHeader,
		uint64,
	) ([]*header.ExtendedHeader, error) `perm:"read"`
	GetByHeight func(context.Context, uint64) (*header.ExtendedHeader, error)    `perm:"read"`
	SyncState   func(ctx context.Context) (sync.State, error)                    `perm:"read"`
	SyncWait    func(ctx context.Context) error                                  `perm:"read"`
	NetworkHead func(ctx context.Context) (*header.ExtendedHeader, error)        `perm:"read"`
	Subscribe   func(ctx context.Context) (<-chan *header.ExtendedHeader, error) `perm:"read"`
}
type StateAPI struct {
	AccountAddress    func(ctx context.Context) (state.Address, error)                      `perm:"read"`
	IsStopped         func(ctx context.Context) bool                                        `perm:"read"`
	Balance           func(ctx context.Context) (*state.Balance, error)                     `perm:"read"`
	BalanceForAddress func(ctx context.Context, addr state.Address) (*state.Balance, error) `perm:"read"`
	Transfer          func(
		ctx context.Context,
		to state.AccAddress,
		amount,
		fee state.Int,
		gasLimit uint64,
	) (*state.TxResponse, error) `perm:"write"`
	SubmitTx         func(ctx context.Context, tx state.Tx) (*state.TxResponse, error) `perm:"write"`
	SubmitPayForBlob func(
		ctx context.Context,
		fee state.Int,
		gasLim uint64,
		blobs []*blob.Blob,
	) (*state.TxResponse, error) `perm:"write"`
	CancelUnbondingDelegation func(
		ctx context.Context,
		valAddr state.ValAddress,
		amount,
		height,
		fee state.Int,
		gasLim uint64,
	) (*state.TxResponse, error) `perm:"write"`
	BeginRedelegate func(
		ctx context.Context,
		srcValAddr,
		dstValAddr state.ValAddress,
		amount,
		fee state.Int,
		gasLim uint64,
	) (*state.TxResponse, error) `perm:"write"`
	Undelegate func(
		ctx context.Context,
		delAddr state.ValAddress,
		amount,
		fee state.Int,
		gasLim uint64,
	) (*state.TxResponse, error) `perm:"write"`
	Delegate func(
		ctx context.Context,
		delAddr state.ValAddress,
		amount,
		fee state.Int,
		gasLim uint64,
	) (*state.TxResponse, error) `perm:"write"`
	QueryDelegation func(
		ctx context.Context,
		valAddr state.ValAddress,
	) (*state.QueryDelegationResponse, error) `perm:"read"`
	QueryUnbonding func(
		ctx context.Context,
		valAddr state.ValAddress,
	) (*state.QueryUnbondingDelegationResponse, error) `perm:"read"`
	QueryRedelegations func(
		ctx context.Context,
		srcValAddr,
		dstValAddr state.ValAddress,
	) (*state.QueryRedelegationsResponse, error) `perm:"read"`
}
type ShareAPI struct {
	SharesAvailable func(context.Context, *header.ExtendedHeader) error `perm:"read"`
	GetShare        func(
		ctx context.Context,
		eh *header.ExtendedHeader,
		row, col int,
	) (share.Share, error) `perm:"read"`
	GetEDS func(
		ctx context.Context,
		eh *header.ExtendedHeader,
	) (*rsmt2d.ExtendedDataSquare, error) `perm:"read"`
	GetSharesByNamespace func(
		ctx context.Context,
		eh *header.ExtendedHeader,
		namespace share.Namespace,
	) (share.NamespacedShares, error) `perm:"read"`
}
type P2PAPI struct {
	Peers                func(context.Context) ([]peer.ID, error)                             `perm:"admin"`
	PeerInfo             func(ctx context.Context, id peer.ID) (peer.AddrInfo, error)         `perm:"admin"`
	Connect              func(ctx context.Context, pi peer.AddrInfo) error                    `perm:"admin"`
	ClosePeer            func(ctx context.Context, id peer.ID) error                          `perm:"admin"`
	Connectedness        func(ctx context.Context, id peer.ID) (network.Connectedness, error) `perm:"admin"`
	NATStatus            func(context.Context) (network.Reachability, error)                  `perm:"admin"`
	BlockPeer            func(ctx context.Context, p peer.ID) error                           `perm:"admin"`
	UnblockPeer          func(ctx context.Context, p peer.ID) error                           `perm:"admin"`
	ListBlockedPeers     func(context.Context) ([]peer.ID, error)                             `perm:"admin"`
	Protect              func(ctx context.Context, id peer.ID, tag string) error              `perm:"admin"`
	Unprotect            func(ctx context.Context, id peer.ID, tag string) (bool, error)      `perm:"admin"`
	IsProtected          func(ctx context.Context, id peer.ID, tag string) (bool, error)      `perm:"admin"`
	BandwidthStats       func(context.Context) (metrics.Stats, error)                         `perm:"admin"`
	BandwidthForPeer     func(ctx context.Context, id peer.ID) (metrics.Stats, error)         `perm:"admin"`
	BandwidthForProtocol func(ctx context.Context, proto protocol.ID) (metrics.Stats, error)  `perm:"admin"`
	ResourceState        func(context.Context) (rcmgr.ResourceManagerStat, error)             `perm:"admin"`
	PubSubPeers          func(ctx context.Context, topic string) ([]peer.ID, error)           `perm:"admin"`
}
type NodeAPI struct {
	Info        func(context.Context) (node.Info, error)                           `perm:"admin"`
	LogLevelSet func(ctx context.Context, name, level string) error                `perm:"admin"`
	AuthVerify  func(ctx context.Context, token string) ([]auth.Permission, error) `perm:"admin"`
	AuthNew     func(ctx context.Context, perms []auth.Permission) ([]byte, error) `perm:"admin"`
}

// SubmitOptions contains the information about fee and gasLimit price in order to configure the Submit request.
type SubmitOptions struct {
	Fee      int64
	GasLimit uint64
}

// DefaultSubmitOptions creates a default fee and gas price values.
func DefaultSubmitOptions() *SubmitOptions {
	return &SubmitOptions{
		Fee:      -1,
		GasLimit: 0,
	}
}
