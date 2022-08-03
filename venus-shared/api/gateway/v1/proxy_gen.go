// Code generated by github.com/filecoin-project/venus/venus-devtool/api-gen. DO NOT EDIT.
package gateway

import (
	"context"

	address "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/specs-storage/storage"
	cid "github.com/ipfs/go-cid"

	"github.com/filecoin-project/venus/venus-shared/actors/builtin"
	"github.com/filecoin-project/venus/venus-shared/types"
	gtypes "github.com/filecoin-project/venus/venus-shared/types/gateway"
)

type IProofClientStruct struct {
	Internal struct {
		ComputeProof        func(ctx context.Context, miner address.Address, sectorInfos []builtin.ExtendedSectorInfo, rand abi.PoStRandomness, height abi.ChainEpoch, nwVersion network.Version) ([]builtin.PoStProof, error) `perm:"admin"`
		ListConnectedMiners func(ctx context.Context) ([]address.Address, error)                                                                                                                                               `perm:"admin"`
		ListMinerConnection func(ctx context.Context, addr address.Address) (*gtypes.MinerState, error)                                                                                                                        `perm:"admin"`
	}
}

func (s *IProofClientStruct) ComputeProof(p0 context.Context, p1 address.Address, p2 []builtin.ExtendedSectorInfo, p3 abi.PoStRandomness, p4 abi.ChainEpoch, p5 network.Version) ([]builtin.PoStProof, error) {
	return s.Internal.ComputeProof(p0, p1, p2, p3, p4, p5)
}
func (s *IProofClientStruct) ListConnectedMiners(p0 context.Context) ([]address.Address, error) {
	return s.Internal.ListConnectedMiners(p0)
}
func (s *IProofClientStruct) ListMinerConnection(p0 context.Context, p1 address.Address) (*gtypes.MinerState, error) {
	return s.Internal.ListMinerConnection(p0, p1)
}

type IProofServiceProviderStruct struct {
	Internal struct {
		ListenProofEvent   func(ctx context.Context, policy *gtypes.ProofRegisterPolicy) (<-chan *gtypes.RequestEvent, error) `perm:"read"`
		ResponseProofEvent func(ctx context.Context, resp *gtypes.ResponseEvent) error                                        `perm:"read"`
	}
}

func (s *IProofServiceProviderStruct) ListenProofEvent(p0 context.Context, p1 *gtypes.ProofRegisterPolicy) (<-chan *gtypes.RequestEvent, error) {
	return s.Internal.ListenProofEvent(p0, p1)
}
func (s *IProofServiceProviderStruct) ResponseProofEvent(p0 context.Context, p1 *gtypes.ResponseEvent) error {
	return s.Internal.ResponseProofEvent(p0, p1)
}

type IProofEventStruct struct {
	IProofClientStruct
	IProofServiceProviderStruct
}

type IWalletClientStruct struct {
	Internal struct {
		ListWalletInfo         func(ctx context.Context) ([]*gtypes.WalletDetail, error)                                                                     `perm:"admin"`
		ListWalletInfoByWallet func(ctx context.Context, wallet string) (*gtypes.WalletDetail, error)                                                        `perm:"admin"`
		WalletHas              func(ctx context.Context, supportAccount string, addr address.Address) (bool, error)                                          `perm:"admin"`
		WalletSign             func(ctx context.Context, account string, addr address.Address, toSign []byte, meta types.MsgMeta) (*crypto.Signature, error) `perm:"admin"`
	}
}

func (s *IWalletClientStruct) ListWalletInfo(p0 context.Context) ([]*gtypes.WalletDetail, error) {
	return s.Internal.ListWalletInfo(p0)
}
func (s *IWalletClientStruct) ListWalletInfoByWallet(p0 context.Context, p1 string) (*gtypes.WalletDetail, error) {
	return s.Internal.ListWalletInfoByWallet(p0, p1)
}
func (s *IWalletClientStruct) WalletHas(p0 context.Context, p1 string, p2 address.Address) (bool, error) {
	return s.Internal.WalletHas(p0, p1, p2)
}
func (s *IWalletClientStruct) WalletSign(p0 context.Context, p1 string, p2 address.Address, p3 []byte, p4 types.MsgMeta) (*crypto.Signature, error) {
	return s.Internal.WalletSign(p0, p1, p2, p3, p4)
}

type IWalletServiceProviderStruct struct {
	Internal struct {
		AddNewAddress       func(ctx context.Context, channelID types.UUID, newAddrs []address.Address) error                   `perm:"read"`
		ListenWalletEvent   func(ctx context.Context, policy *gtypes.WalletRegisterPolicy) (<-chan *gtypes.RequestEvent, error) `perm:"read"`
		RemoveAddress       func(ctx context.Context, channelID types.UUID, newAddrs []address.Address) error                   `perm:"read"`
		ResponseWalletEvent func(ctx context.Context, resp *gtypes.ResponseEvent) error                                         `perm:"read"`
		SupportNewAccount   func(ctx context.Context, channelID types.UUID, account string) error                               `perm:"read"`
	}
}

func (s *IWalletServiceProviderStruct) AddNewAddress(p0 context.Context, p1 types.UUID, p2 []address.Address) error {
	return s.Internal.AddNewAddress(p0, p1, p2)
}
func (s *IWalletServiceProviderStruct) ListenWalletEvent(p0 context.Context, p1 *gtypes.WalletRegisterPolicy) (<-chan *gtypes.RequestEvent, error) {
	return s.Internal.ListenWalletEvent(p0, p1)
}
func (s *IWalletServiceProviderStruct) RemoveAddress(p0 context.Context, p1 types.UUID, p2 []address.Address) error {
	return s.Internal.RemoveAddress(p0, p1, p2)
}
func (s *IWalletServiceProviderStruct) ResponseWalletEvent(p0 context.Context, p1 *gtypes.ResponseEvent) error {
	return s.Internal.ResponseWalletEvent(p0, p1)
}
func (s *IWalletServiceProviderStruct) SupportNewAccount(p0 context.Context, p1 types.UUID, p2 string) error {
	return s.Internal.SupportNewAccount(p0, p1, p2)
}

type IWalletEventStruct struct {
	IWalletClientStruct
	IWalletServiceProviderStruct
}

type IMarketClientStruct struct {
	Internal struct {
		IsUnsealed                 func(ctx context.Context, miner address.Address, pieceCid cid.Cid, sector storage.SectorRef, offset types.PaddedByteIndex, size abi.PaddedPieceSize) (bool, error)      `perm:"admin"`
		ListMarketConnectionsState func(ctx context.Context) ([]gtypes.MarketConnectionState, error)                                                                                                       `perm:"admin"`
		SectorsUnsealPiece         func(ctx context.Context, miner address.Address, pieceCid cid.Cid, sector storage.SectorRef, offset types.PaddedByteIndex, size abi.PaddedPieceSize, dest string) error `perm:"admin"`
	}
}

func (s *IMarketClientStruct) IsUnsealed(p0 context.Context, p1 address.Address, p2 cid.Cid, p3 storage.SectorRef, p4 types.PaddedByteIndex, p5 abi.PaddedPieceSize) (bool, error) {
	return s.Internal.IsUnsealed(p0, p1, p2, p3, p4, p5)
}
func (s *IMarketClientStruct) ListMarketConnectionsState(p0 context.Context) ([]gtypes.MarketConnectionState, error) {
	return s.Internal.ListMarketConnectionsState(p0)
}
func (s *IMarketClientStruct) SectorsUnsealPiece(p0 context.Context, p1 address.Address, p2 cid.Cid, p3 storage.SectorRef, p4 types.PaddedByteIndex, p5 abi.PaddedPieceSize, p6 string) error {
	return s.Internal.SectorsUnsealPiece(p0, p1, p2, p3, p4, p5, p6)
}

type IMarketServiceProviderStruct struct {
	Internal struct {
		ListenMarketEvent   func(ctx context.Context, policy *gtypes.MarketRegisterPolicy) (<-chan *gtypes.RequestEvent, error) `perm:"read"`
		ResponseMarketEvent func(ctx context.Context, resp *gtypes.ResponseEvent) error                                         `perm:"read"`
	}
}

func (s *IMarketServiceProviderStruct) ListenMarketEvent(p0 context.Context, p1 *gtypes.MarketRegisterPolicy) (<-chan *gtypes.RequestEvent, error) {
	return s.Internal.ListenMarketEvent(p0, p1)
}
func (s *IMarketServiceProviderStruct) ResponseMarketEvent(p0 context.Context, p1 *gtypes.ResponseEvent) error {
	return s.Internal.ResponseMarketEvent(p0, p1)
}

type IMarketEventStruct struct {
	IMarketClientStruct
	IMarketServiceProviderStruct
}

type IGatewayStruct struct {
	IProofEventStruct
	IWalletEventStruct
	IMarketEventStruct

	Internal struct {
		Version func(ctx context.Context) (types.Version, error) `perm:"read"`
	}
}

func (s *IGatewayStruct) Version(p0 context.Context) (types.Version, error) {
	return s.Internal.Version(p0)
}
