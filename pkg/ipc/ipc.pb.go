package ipc

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GstPipelineDebugDotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GstPipelineDebugDotRequest) Reset() {
	*x = GstPipelineDebugDotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GstPipelineDebugDotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GstPipelineDebugDotRequest) ProtoMessage() {}

func (x *GstPipelineDebugDotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GstPipelineDebugDotRequest.ProtoReflect.Descriptor instead.
func (*GstPipelineDebugDotRequest) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{0}
}

type GstPipelineDebugDotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DotFile string `protobuf:"bytes,1,opt,name=dot_file,json=dotFile,proto3" json:"dot_file,omitempty"`
}

func (x *GstPipelineDebugDotResponse) Reset() {
	*x = GstPipelineDebugDotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GstPipelineDebugDotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GstPipelineDebugDotResponse) ProtoMessage() {}

func (x *GstPipelineDebugDotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GstPipelineDebugDotResponse.ProtoReflect.Descriptor instead.
func (*GstPipelineDebugDotResponse) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{1}
}

func (x *GstPipelineDebugDotResponse) GetDotFile() string {
	if x != nil {
		return x.DotFile
	}
	return ""
}

type PProfRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfileName string `protobuf:"bytes,1,opt,name=profile_name,json=profileName,proto3" json:"profile_name,omitempty"`
	Timeout     int32  `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Debug       int32  `protobuf:"varint,3,opt,name=debug,proto3" json:"debug,omitempty"`
}

func (x *PProfRequest) Reset() {
	*x = PProfRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PProfRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PProfRequest) ProtoMessage() {}

func (x *PProfRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PProfRequest.ProtoReflect.Descriptor instead.
func (*PProfRequest) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{2}
}

func (x *PProfRequest) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

func (x *PProfRequest) GetTimeout() int32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *PProfRequest) GetDebug() int32 {
	if x != nil {
		return x.Debug
	}
	return 0
}

type PProfResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PprofFile []byte `protobuf:"bytes,1,opt,name=pprof_file,json=pprofFile,proto3" json:"pprof_file,omitempty"`
}

func (x *PProfResponse) Reset() {
	*x = PProfResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PProfResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PProfResponse) ProtoMessage() {}

func (x *PProfResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PProfResponse.ProtoReflect.Descriptor instead.
func (*PProfResponse) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{3}
}

func (x *PProfResponse) GetPprofFile() []byte {
	if x != nil {
		return x.PprofFile
	}
	return nil
}

type GatherMediaStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GatherMediaStatsRequest) Reset() {
	*x = GatherMediaStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatherMediaStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatherMediaStatsRequest) ProtoMessage() {}

func (x *GatherMediaStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatherMediaStatsRequest.ProtoReflect.Descriptor instead.
func (*GatherMediaStatsRequest) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{4}
}

type GatherMediaStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats *MediaStats `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats,omitempty"`
}

func (x *GatherMediaStatsResponse) Reset() {
	*x = GatherMediaStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatherMediaStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatherMediaStatsResponse) ProtoMessage() {}

func (x *GatherMediaStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatherMediaStatsResponse.ProtoReflect.Descriptor instead.
func (*GatherMediaStatsResponse) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{5}
}

func (x *GatherMediaStatsResponse) GetStats() *MediaStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

type UpdateMediaStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats *MediaStats `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats,omitempty"`
}

func (x *UpdateMediaStatsRequest) Reset() {
	*x = UpdateMediaStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMediaStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMediaStatsRequest) ProtoMessage() {}

func (x *UpdateMediaStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMediaStatsRequest.ProtoReflect.Descriptor instead.
func (*UpdateMediaStatsRequest) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateMediaStatsRequest) GetStats() *MediaStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

type MediaStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TrackStats map[string]*TrackStats `protobuf:"bytes,1,rep,name=track_stats,json=trackStats,proto3" json:"track_stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MediaStats) Reset() {
	*x = MediaStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MediaStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MediaStats) ProtoMessage() {}

func (x *MediaStats) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MediaStats.ProtoReflect.Descriptor instead.
func (*MediaStats) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{7}
}

func (x *MediaStats) GetTrackStats() map[string]*TrackStats {
	if x != nil {
		return x.TrackStats
	}
	return nil
}

type TrackStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AverageBitrate  uint32       `protobuf:"varint,1,opt,name=average_bitrate,json=averageBitrate,proto3" json:"average_bitrate,omitempty"`
	CurrentBitrate  uint32       `protobuf:"varint,2,opt,name=current_bitrate,json=currentBitrate,proto3" json:"current_bitrate,omitempty"`
	TotalPackets    uint64       `protobuf:"varint,4,opt,name=total_packets,json=totalPackets,proto3" json:"total_packets,omitempty"`
	CurrentPackets  uint64       `protobuf:"varint,5,opt,name=current_packets,json=currentPackets,proto3" json:"current_packets,omitempty"`
	TotalLossRate   float64      `protobuf:"fixed64,6,opt,name=total_loss_rate,json=totalLossRate,proto3" json:"total_loss_rate,omitempty"`
	CurrentLossRate float64      `protobuf:"fixed64,7,opt,name=current_loss_rate,json=currentLossRate,proto3" json:"current_loss_rate,omitempty"`
	TotalPli        uint64       `protobuf:"varint,8,opt,name=total_pli,json=totalPli,proto3" json:"total_pli,omitempty"`
	CurrentPli      uint64       `protobuf:"varint,9,opt,name=current_pli,json=currentPli,proto3" json:"current_pli,omitempty"`
	Jitter          *JitterStats `protobuf:"bytes,10,opt,name=jitter,proto3" json:"jitter,omitempty"`
}

func (x *TrackStats) Reset() {
	*x = TrackStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrackStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackStats) ProtoMessage() {}

func (x *TrackStats) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackStats.ProtoReflect.Descriptor instead.
func (*TrackStats) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{8}
}

func (x *TrackStats) GetAverageBitrate() uint32 {
	if x != nil {
		return x.AverageBitrate
	}
	return 0
}

func (x *TrackStats) GetCurrentBitrate() uint32 {
	if x != nil {
		return x.CurrentBitrate
	}
	return 0
}

func (x *TrackStats) GetTotalPackets() uint64 {
	if x != nil {
		return x.TotalPackets
	}
	return 0
}

func (x *TrackStats) GetCurrentPackets() uint64 {
	if x != nil {
		return x.CurrentPackets
	}
	return 0
}

func (x *TrackStats) GetTotalLossRate() float64 {
	if x != nil {
		return x.TotalLossRate
	}
	return 0
}

func (x *TrackStats) GetCurrentLossRate() float64 {
	if x != nil {
		return x.CurrentLossRate
	}
	return 0
}

func (x *TrackStats) GetTotalPli() uint64 {
	if x != nil {
		return x.TotalPli
	}
	return 0
}

func (x *TrackStats) GetCurrentPli() uint64 {
	if x != nil {
		return x.CurrentPli
	}
	return 0
}

func (x *TrackStats) GetJitter() *JitterStats {
	if x != nil {
		return x.Jitter
	}
	return nil
}

type JitterStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	P50 float64 `protobuf:"fixed64,1,opt,name=p50,proto3" json:"p50,omitempty"` // in ms
	P90 float64 `protobuf:"fixed64,2,opt,name=p90,proto3" json:"p90,omitempty"`
	P99 float64 `protobuf:"fixed64,3,opt,name=p99,proto3" json:"p99,omitempty"`
}

func (x *JitterStats) Reset() {
	*x = JitterStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ipc_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JitterStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JitterStats) ProtoMessage() {}

func (x *JitterStats) ProtoReflect() protoreflect.Message {
	mi := &file_ipc_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JitterStats.ProtoReflect.Descriptor instead.
func (*JitterStats) Descriptor() ([]byte, []int) {
	return file_ipc_proto_rawDescGZIP(), []int{9}
}

func (x *JitterStats) GetP50() float64 {
	if x != nil {
		return x.P50
	}
	return 0
}

func (x *JitterStats) GetP90() float64 {
	if x != nil {
		return x.P90
	}
	return 0
}

func (x *JitterStats) GetP99() float64 {
	if x != nil {
		return x.P99
	}
	return 0
}

var File_ipc_proto protoreflect.FileDescriptor

var file_ipc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x69, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x69, 0x70, 0x63,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1c, 0x0a,
	0x1a, 0x47, 0x73, 0x74, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x65, 0x62, 0x75,
	0x67, 0x44, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x38, 0x0a, 0x1b, 0x47,
	0x73, 0x74, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x65, 0x62, 0x75, 0x67, 0x44,
	0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f,
	0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x61, 0x0a, 0x0c, 0x50, 0x50, 0x72, 0x6f, 0x66, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x22, 0x2e, 0x0a, 0x0d, 0x50, 0x50, 0x72, 0x6f,
	0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x70, 0x72,
	0x6f, 0x66, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70,
	0x70, 0x72, 0x6f, 0x66, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x19, 0x0a, 0x17, 0x47, 0x61, 0x74, 0x68,
	0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x41, 0x0a, 0x18, 0x47, 0x61, 0x74, 0x68, 0x65, 0x72, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x25, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x69, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x22, 0x40, 0x0a, 0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x25, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x69, 0x70, 0x63, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x22, 0x9e, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x40, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x6b,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x69,
	0x70, 0x63, 0x2e, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x1a, 0x4e, 0x0a, 0x0f, 0x54, 0x72, 0x61,
	0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x25,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x69, 0x70, 0x63, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xe8, 0x02, 0x0a, 0x0a, 0x54, 0x72,
	0x61, 0x63, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x76, 0x65, 0x72,
	0x61, 0x67, 0x65, 0x5f, 0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0e, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x42, 0x69, 0x74, 0x72, 0x61, 0x74,
	0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x69, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x42, 0x69, 0x74, 0x72, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12,
	0x27, 0x0a, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x6c, 0x6f, 0x73, 0x73, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4c, 0x6f, 0x73, 0x73, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x2a, 0x0a, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x6f, 0x73, 0x73,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x73, 0x73, 0x52, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x69, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x6c, 0x69, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x6c, 0x69, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x6c, 0x69, 0x12, 0x28, 0x0a, 0x06, 0x6a, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x69, 0x70, 0x63,
	0x2e, 0x4a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x06, 0x6a, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x22, 0x43, 0x0a, 0x0b, 0x4a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x35, 0x30, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x03, 0x70, 0x35, 0x30, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x39, 0x30, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x70, 0x39, 0x30, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x39, 0x39, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x70, 0x39, 0x39, 0x32, 0xbb, 0x02, 0x0a, 0x0e, 0x49, 0x6e,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x6f, 0x74, 0x12, 0x1f,
	0x2e, 0x69, 0x70, 0x63, 0x2e, 0x47, 0x73, 0x74, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x44, 0x65, 0x62, 0x75, 0x67, 0x44, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x69, 0x70, 0x63, 0x2e, 0x47, 0x73, 0x74, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x44, 0x65, 0x62, 0x75, 0x67, 0x44, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x50, 0x72, 0x6f, 0x66, 0x12,
	0x11, 0x2e, 0x69, 0x70, 0x63, 0x2e, 0x50, 0x50, 0x72, 0x6f, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x12, 0x2e, 0x69, 0x70, 0x63, 0x2e, 0x50, 0x50, 0x72, 0x6f, 0x66, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x10, 0x47, 0x61, 0x74, 0x68,
	0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x2e, 0x69,
	0x70, 0x63, 0x2e, 0x47, 0x61, 0x74, 0x68, 0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x70, 0x63,
	0x2e, 0x47, 0x61, 0x74, 0x68, 0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x10, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12,
	0x1c, 0x2e, 0x69, 0x70, 0x63, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x76, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x69, 0x6e,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x69, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ipc_proto_rawDescOnce sync.Once
	file_ipc_proto_rawDescData = file_ipc_proto_rawDesc
)

func file_ipc_proto_rawDescGZIP() []byte {
	file_ipc_proto_rawDescOnce.Do(func() {
		file_ipc_proto_rawDescData = protoimpl.X.CompressGZIP(file_ipc_proto_rawDescData)
	})
	return file_ipc_proto_rawDescData
}

var file_ipc_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_ipc_proto_goTypes = []interface{}{
	(*GstPipelineDebugDotRequest)(nil),  // 0: ipc.GstPipelineDebugDotRequest
	(*GstPipelineDebugDotResponse)(nil), // 1: ipc.GstPipelineDebugDotResponse
	(*PProfRequest)(nil),                // 2: ipc.PProfRequest
	(*PProfResponse)(nil),               // 3: ipc.PProfResponse
	(*GatherMediaStatsRequest)(nil),     // 4: ipc.GatherMediaStatsRequest
	(*GatherMediaStatsResponse)(nil),    // 5: ipc.GatherMediaStatsResponse
	(*UpdateMediaStatsRequest)(nil),     // 6: ipc.UpdateMediaStatsRequest
	(*MediaStats)(nil),                  // 7: ipc.MediaStats
	(*TrackStats)(nil),                  // 8: ipc.TrackStats
	(*JitterStats)(nil),                 // 9: ipc.JitterStats
	nil,                                 // 10: ipc.MediaStats.TrackStatsEntry
	(*emptypb.Empty)(nil),               // 11: google.protobuf.Empty
}
var file_ipc_proto_depIdxs = []int32{
	7,  // 0: ipc.GatherMediaStatsResponse.stats:type_name -> ipc.MediaStats
	7,  // 1: ipc.UpdateMediaStatsRequest.stats:type_name -> ipc.MediaStats
	10, // 2: ipc.MediaStats.track_stats:type_name -> ipc.MediaStats.TrackStatsEntry
	9,  // 3: ipc.TrackStats.jitter:type_name -> ipc.JitterStats
	8,  // 4: ipc.MediaStats.TrackStatsEntry.value:type_name -> ipc.TrackStats
	0,  // 5: ipc.IngressHandler.GetPipelineDot:input_type -> ipc.GstPipelineDebugDotRequest
	2,  // 6: ipc.IngressHandler.GetPProf:input_type -> ipc.PProfRequest
	4,  // 7: ipc.IngressHandler.GatherMediaStats:input_type -> ipc.GatherMediaStatsRequest
	6,  // 8: ipc.IngressHandler.UpdateMediaStats:input_type -> ipc.UpdateMediaStatsRequest
	1,  // 9: ipc.IngressHandler.GetPipelineDot:output_type -> ipc.GstPipelineDebugDotResponse
	3,  // 10: ipc.IngressHandler.GetPProf:output_type -> ipc.PProfResponse
	5,  // 11: ipc.IngressHandler.GatherMediaStats:output_type -> ipc.GatherMediaStatsResponse
	11, // 12: ipc.IngressHandler.UpdateMediaStats:output_type -> google.protobuf.Empty
	9,  // [9:13] is the sub-list for method output_type
	5,  // [5:9] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_ipc_proto_init() }
func file_ipc_proto_init() {
	if File_ipc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ipc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GstPipelineDebugDotRequest); i {
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
		file_ipc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GstPipelineDebugDotResponse); i {
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
		file_ipc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PProfRequest); i {
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
		file_ipc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PProfResponse); i {
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
		file_ipc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatherMediaStatsRequest); i {
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
		file_ipc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatherMediaStatsResponse); i {
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
		file_ipc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMediaStatsRequest); i {
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
		file_ipc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MediaStats); i {
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
		file_ipc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrackStats); i {
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
		file_ipc_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JitterStats); i {
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
			RawDescriptor: file_ipc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ipc_proto_goTypes,
		DependencyIndexes: file_ipc_proto_depIdxs,
		MessageInfos:      file_ipc_proto_msgTypes,
	}.Build()
	File_ipc_proto = out.File
	file_ipc_proto_rawDesc = nil
	file_ipc_proto_goTypes = nil
	file_ipc_proto_depIdxs = nil
}
