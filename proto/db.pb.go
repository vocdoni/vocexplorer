// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.12.4
// source: db.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type BlockchainInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network           string               `protobuf:"bytes,1,opt,name=Network,proto3" json:"Network,omitempty"`
	Version           string               `protobuf:"bytes,2,opt,name=Version,proto3" json:"Version,omitempty"`
	LatestBlockHeight int64                `protobuf:"varint,3,opt,name=LatestBlockHeight,proto3" json:"LatestBlockHeight,omitempty"`
	GenesisTimeStamp  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=GenesisTimeStamp,proto3" json:"GenesisTimeStamp,omitempty"`
	ChainID           string               `protobuf:"bytes,5,opt,name=ChainID,proto3" json:"ChainID,omitempty"`
	BlockTime         []int32              `protobuf:"varint,6,rep,packed,name=BlockTime,proto3" json:"BlockTime,omitempty"`
	BlockTimeStamp    int32                `protobuf:"varint,7,opt,name=BlockTimeStamp,proto3" json:"BlockTimeStamp,omitempty"`
	Height            int64                `protobuf:"varint,8,opt,name=Height,proto3" json:"Height,omitempty"`
	MaxBytes          int64                `protobuf:"varint,9,opt,name=MaxBytes,proto3" json:"MaxBytes,omitempty"`
	Syncing           bool                 `protobuf:"varint,10,opt,name=syncing,proto3" json:"syncing,omitempty"`
}

func (x *BlockchainInfo) Reset() {
	*x = BlockchainInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockchainInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockchainInfo) ProtoMessage() {}

func (x *BlockchainInfo) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockchainInfo.ProtoReflect.Descriptor instead.
func (*BlockchainInfo) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{0}
}

func (x *BlockchainInfo) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *BlockchainInfo) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *BlockchainInfo) GetLatestBlockHeight() int64 {
	if x != nil {
		return x.LatestBlockHeight
	}
	return 0
}

func (x *BlockchainInfo) GetGenesisTimeStamp() *timestamp.Timestamp {
	if x != nil {
		return x.GenesisTimeStamp
	}
	return nil
}

func (x *BlockchainInfo) GetChainID() string {
	if x != nil {
		return x.ChainID
	}
	return ""
}

func (x *BlockchainInfo) GetBlockTime() []int32 {
	if x != nil {
		return x.BlockTime
	}
	return nil
}

func (x *BlockchainInfo) GetBlockTimeStamp() int32 {
	if x != nil {
		return x.BlockTimeStamp
	}
	return 0
}

func (x *BlockchainInfo) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BlockchainInfo) GetMaxBytes() int64 {
	if x != nil {
		return x.MaxBytes
	}
	return 0
}

func (x *BlockchainInfo) GetSyncing() bool {
	if x != nil {
		return x.Syncing
	}
	return false
}

type Height struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height int64 `protobuf:"varint,1,opt,name=Height,proto3" json:"Height,omitempty"`
}

func (x *Height) Reset() {
	*x = Height{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Height) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Height) ProtoMessage() {}

func (x *Height) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Height.ProtoReflect.Descriptor instead.
func (*Height) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{1}
}

func (x *Height) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EncryptionKeyIndexes []int32 `protobuf:"varint,1,rep,packed,name=EncryptionKeyIndexes,proto3" json:"EncryptionKeyIndexes,omitempty"`
	Nullifier            string  `protobuf:"bytes,2,opt,name=Nullifier,proto3" json:"Nullifier,omitempty"`
	ProcessID            string  `protobuf:"bytes,3,opt,name=ProcessID,proto3" json:"ProcessID,omitempty"`
	Package              string  `protobuf:"bytes,4,opt,name=Package,proto3" json:"Package,omitempty"`
	ProcessHeight        int64   `protobuf:"varint,5,opt,name=ProcessHeight,proto3" json:"ProcessHeight,omitempty"`
	GlobalHeight         int64   `protobuf:"varint,6,opt,name=GlobalHeight,proto3" json:"GlobalHeight,omitempty"`
	TxHeight             int64   `protobuf:"varint,7,opt,name=TxHeight,proto3" json:"TxHeight,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{2}
}

func (x *Envelope) GetEncryptionKeyIndexes() []int32 {
	if x != nil {
		return x.EncryptionKeyIndexes
	}
	return nil
}

func (x *Envelope) GetNullifier() string {
	if x != nil {
		return x.Nullifier
	}
	return ""
}

func (x *Envelope) GetProcessID() string {
	if x != nil {
		return x.ProcessID
	}
	return ""
}

func (x *Envelope) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

func (x *Envelope) GetProcessHeight() int64 {
	if x != nil {
		return x.ProcessHeight
	}
	return 0
}

func (x *Envelope) GetGlobalHeight() int64 {
	if x != nil {
		return x.GlobalHeight
	}
	return 0
}

func (x *Envelope) GetTxHeight() int64 {
	if x != nil {
		return x.TxHeight
	}
	return 0
}

type StoreBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash     []byte               `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height   int64                `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	NumTxs   int64                `protobuf:"varint,3,opt,name=NumTxs,proto3" json:"NumTxs,omitempty"`
	Time     *timestamp.Timestamp `protobuf:"bytes,4,opt,name=Time,proto3" json:"Time,omitempty"`
	Proposer []byte               `protobuf:"bytes,5,opt,name=Proposer,proto3" json:"Proposer,omitempty"`
}

func (x *StoreBlock) Reset() {
	*x = StoreBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreBlock) ProtoMessage() {}

func (x *StoreBlock) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreBlock.ProtoReflect.Descriptor instead.
func (*StoreBlock) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{3}
}

func (x *StoreBlock) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *StoreBlock) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *StoreBlock) GetNumTxs() int64 {
	if x != nil {
		return x.NumTxs
	}
	return 0
}

func (x *StoreBlock) GetTime() *timestamp.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *StoreBlock) GetProposer() []byte {
	if x != nil {
		return x.Proposer
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height    int64  `protobuf:"varint,1,opt,name=Height,proto3" json:"Height,omitempty"`
	Index     uint32 `protobuf:"varint,2,opt,name=Index,proto3" json:"Index,omitempty"`
	Tx        []byte `protobuf:"bytes,3,opt,name=Tx,proto3" json:"Tx,omitempty"`
	TxHeight  int64  `protobuf:"varint,4,opt,name=TxHeight,proto3" json:"TxHeight,omitempty"`
	Nullifier string `protobuf:"bytes,5,opt,name=Nullifier,proto3" json:"Nullifier,omitempty"`
	Hash      []byte `protobuf:"bytes,6,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{4}
}

func (x *Transaction) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Transaction) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Transaction) GetTx() []byte {
	if x != nil {
		return x.Tx
	}
	return nil
}

func (x *Transaction) GetTxHeight() int64 {
	if x != nil {
		return x.TxHeight
	}
	return 0
}

func (x *Transaction) GetNullifier() string {
	if x != nil {
		return x.Nullifier
	}
	return ""
}

func (x *Transaction) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type ItemList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items [][]byte `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ItemList) Reset() {
	*x = ItemList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ItemList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemList) ProtoMessage() {}

func (x *ItemList) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemList.ProtoReflect.Descriptor instead.
func (*ItemList) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{5}
}

func (x *ItemList) GetItems() [][]byte {
	if x != nil {
		return x.Items
	}
	return nil
}

type StringList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []string `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *StringList) Reset() {
	*x = StringList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringList) ProtoMessage() {}

func (x *StringList) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringList.ProtoReflect.Descriptor instead.
func (*StringList) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{6}
}

func (x *StringList) GetItems() []string {
	if x != nil {
		return x.Items
	}
	return nil
}

type Validator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address []byte `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	// Bytes-representation of PubKey. Could be nested pubkey struct
	PubKey           []byte  `protobuf:"bytes,2,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	VotingPower      int64   `protobuf:"varint,3,opt,name=VotingPower,proto3" json:"VotingPower,omitempty"`
	ProposerPriority int64   `protobuf:"varint,4,opt,name=ProposerPriority,proto3" json:"ProposerPriority,omitempty"`
	Height           *Height `protobuf:"bytes,5,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *Validator) Reset() {
	*x = Validator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Validator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Validator) ProtoMessage() {}

func (x *Validator) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Validator.ProtoReflect.Descriptor instead.
func (*Validator) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{7}
}

func (x *Validator) GetAddress() []byte {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Validator) GetPubKey() []byte {
	if x != nil {
		return x.PubKey
	}
	return nil
}

func (x *Validator) GetVotingPower() int64 {
	if x != nil {
		return x.VotingPower
	}
	return 0
}

func (x *Validator) GetProposerPriority() int64 {
	if x != nil {
		return x.ProposerPriority
	}
	return 0
}

func (x *Validator) GetHeight() *Height {
	if x != nil {
		return x.Height
	}
	return nil
}

type Process struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	EntityID    string  `protobuf:"bytes,2,opt,name=EntityID,proto3" json:"EntityID,omitempty"`
	LocalHeight *Height `protobuf:"bytes,3,opt,name=LocalHeight,proto3" json:"LocalHeight,omitempty"`
}

func (x *Process) Reset() {
	*x = Process{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Process) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Process) ProtoMessage() {}

func (x *Process) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Process.ProtoReflect.Descriptor instead.
func (*Process) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{8}
}

func (x *Process) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Process) GetEntityID() string {
	if x != nil {
		return x.EntityID
	}
	return ""
}

func (x *Process) GetLocalHeight() *Height {
	if x != nil {
		return x.LocalHeight
	}
	return nil
}

type HeightMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Heights map[string]int64 `protobuf:"bytes,1,rep,name=heights,proto3" json:"heights,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *HeightMap) Reset() {
	*x = HeightMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeightMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeightMap) ProtoMessage() {}

func (x *HeightMap) ProtoReflect() protoreflect.Message {
	mi := &file_db_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeightMap.ProtoReflect.Descriptor instead.
func (*HeightMap) Descriptor() ([]byte, []int) {
	return file_db_proto_rawDescGZIP(), []int{9}
}

func (x *HeightMap) GetHeights() map[string]int64 {
	if x != nil {
		return x.Heights
	}
	return nil
}

var File_db_proto protoreflect.FileDescriptor

var file_db_proto_rawDesc = []byte{
	0x0a, 0x08, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xe8, 0x02, 0x0a, 0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12,
	0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x11, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x46, 0x0a, 0x10, 0x47, 0x65, 0x6e, 0x65, 0x73,
	0x69, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x10, 0x47,
	0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x18, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28, 0x05, 0x52, 0x09, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x4d, 0x61, 0x78, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4d, 0x61, 0x78, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x69, 0x6e, 0x67, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x69, 0x6e, 0x67, 0x22, 0x20, 0x0a,
	0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22,
	0xfa, 0x01, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x14,
	0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x14, 0x45, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1c,
	0x0a, 0x09, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07,
	0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x22, 0x0a, 0x0c,
	0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x54, 0x78, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x54, 0x78, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x9c, 0x01, 0x0a,
	0x0a, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x54, 0x78,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x54, 0x78, 0x73, 0x12,
	0x2e, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x08, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x22, 0x99, 0x01, 0x0a, 0x0b,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x48, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x78, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x54, 0x78, 0x12, 0x1a, 0x0a, 0x08, 0x54, 0x78, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x54, 0x78, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x22, 0x20, 0x0a, 0x08, 0x49, 0x74, 0x65, 0x6d, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0c, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x22, 0x0a, 0x0a, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xb2, 0x01,
	0x0a, 0x09, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a,
	0x0b, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0b, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12,
	0x2a, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x50, 0x72, 0x69, 0x6f, 0x72,
	0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x65, 0x72, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x25, 0x0a, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x22, 0x66, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a,
	0x08, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x44, 0x12, 0x2f, 0x0a, 0x0b, 0x4c, 0x6f, 0x63,
	0x61, 0x6c, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x52, 0x0b, 0x4c,
	0x6f, 0x63, 0x61, 0x6c, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x80, 0x01, 0x0a, 0x09, 0x48,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x4d, 0x61, 0x70, 0x12, 0x37, 0x0a, 0x07, 0x68, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x4d, 0x61, 0x70, 0x2e, 0x48, 0x65, 0x69, 0x67,
	0x68, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x08, 0x5a,
	0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_proto_rawDescOnce sync.Once
	file_db_proto_rawDescData = file_db_proto_rawDesc
)

func file_db_proto_rawDescGZIP() []byte {
	file_db_proto_rawDescOnce.Do(func() {
		file_db_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_proto_rawDescData)
	})
	return file_db_proto_rawDescData
}

var file_db_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_db_proto_goTypes = []interface{}{
	(*BlockchainInfo)(nil),      // 0: proto.BlockchainInfo
	(*Height)(nil),              // 1: proto.Height
	(*Envelope)(nil),            // 2: proto.Envelope
	(*StoreBlock)(nil),          // 3: proto.StoreBlock
	(*Transaction)(nil),         // 4: proto.Transaction
	(*ItemList)(nil),            // 5: proto.ItemList
	(*StringList)(nil),          // 6: proto.StringList
	(*Validator)(nil),           // 7: proto.Validator
	(*Process)(nil),             // 8: proto.Process
	(*HeightMap)(nil),           // 9: proto.HeightMap
	nil,                         // 10: proto.HeightMap.HeightsEntry
	(*timestamp.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_db_proto_depIdxs = []int32{
	11, // 0: proto.BlockchainInfo.GenesisTimeStamp:type_name -> google.protobuf.Timestamp
	11, // 1: proto.StoreBlock.Time:type_name -> google.protobuf.Timestamp
	1,  // 2: proto.Validator.height:type_name -> proto.Height
	1,  // 3: proto.Process.LocalHeight:type_name -> proto.Height
	10, // 4: proto.HeightMap.heights:type_name -> proto.HeightMap.HeightsEntry
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_db_proto_init() }
func file_db_proto_init() {
	if File_db_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockchainInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Height); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ItemList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Validator); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Process); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeightMap); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_proto_goTypes,
		DependencyIndexes: file_db_proto_depIdxs,
		MessageInfos:      file_db_proto_msgTypes,
	}.Build()
	File_db_proto = out.File
	file_db_proto_rawDesc = nil
	file_db_proto_goTypes = nil
	file_db_proto_depIdxs = nil
}
