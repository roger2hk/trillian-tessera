// Copyright 2017 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: configpb/config.proto

package configpb

import (
	keyspb "github.com/google/trillian/crypto/keyspb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// LogConfig describes the configuration options for a log instance.
//
// NEXT_ID: 15
type LogConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// origin identifies the log. It will be used in its checkpoint, and
	// is also its submission prefix, as per https://c2sp.org/static-ct-api
	Origin string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
	// Paths to the files containing root certificates that are acceptable to the
	// log. The certs are served through get-roots endpoint.
	RootsPemFile []string `protobuf:"bytes,2,rep,name=roots_pem_file,json=rootsPemFile,proto3" json:"roots_pem_file,omitempty"`
	// The private key used for signing Checkpoints or SCTs.
	PrivateKey *anypb.Any `protobuf:"bytes,3,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	// The public key matching the above private key (if both are present).
	// It can be specified for the convenience of test tools, but it not used
	// by the server.
	PublicKey *keyspb.PublicKey `protobuf:"bytes,4,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// If reject_expired is true then the certificate validity period will be
	// checked against the current time during the validation of submissions.
	// This will cause expired certificates to be rejected.
	RejectExpired bool `protobuf:"varint,5,opt,name=reject_expired,json=rejectExpired,proto3" json:"reject_expired,omitempty"`
	// If reject_unexpired is true then CTFE rejects certificates that are either
	// currently valid or not yet valid.
	RejectUnexpired bool `protobuf:"varint,6,opt,name=reject_unexpired,json=rejectUnexpired,proto3" json:"reject_unexpired,omitempty"`
	// If set, ext_key_usages will restrict the set of such usages that the
	// server will accept. By default all are accepted. The values specified
	// must be ones known to the x509 package.
	ExtKeyUsages []string `protobuf:"bytes,7,rep,name=ext_key_usages,json=extKeyUsages,proto3" json:"ext_key_usages,omitempty"`
	// not_after_start defines the start of the range of acceptable NotAfter
	// values, inclusive.
	// Leaving this unset implies no lower bound to the range.
	NotAfterStart *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=not_after_start,json=notAfterStart,proto3" json:"not_after_start,omitempty"`
	// not_after_limit defines the end of the range of acceptable NotAfter values,
	// exclusive.
	// Leaving this unset implies no upper bound to the range.
	NotAfterLimit *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=not_after_limit,json=notAfterLimit,proto3" json:"not_after_limit,omitempty"`
	// accept_only_ca controls whether or not *only* certificates with the CA bit
	// set will be accepted.
	AcceptOnlyCa bool `protobuf:"varint,10,opt,name=accept_only_ca,json=acceptOnlyCa,proto3" json:"accept_only_ca,omitempty"`
	// The Maximum Merge Delay (MMD) of this log in seconds. See RFC6962 section 3
	// for definition of MMD. If zero, the log does not provide an MMD guarantee
	// (for example, it is a frozen log).
	MaxMergeDelaySec int32 `protobuf:"varint,11,opt,name=max_merge_delay_sec,json=maxMergeDelaySec,proto3" json:"max_merge_delay_sec,omitempty"`
	// The merge delay that the underlying log implementation is able/targeting to
	// provide. This option is exposed in CTFE metrics, and can be particularly
	// useful to catch when the log is behind but has not yet violated the strict
	// MMD limit.
	// Log operator should decide what exactly EMD means for them. For example, it
	// can be a 99-th percentile of merge delays that they observe, and they can
	// alert on the actual merge delay going above a certain multiple of this EMD.
	ExpectedMergeDelaySec int32 `protobuf:"varint,12,opt,name=expected_merge_delay_sec,json=expectedMergeDelaySec,proto3" json:"expected_merge_delay_sec,omitempty"`
	// A list of X.509 extension OIDs, in dotted string form (e.g. "2.3.4.5")
	// which should cause submissions to be rejected.
	RejectExtensions []string `protobuf:"bytes,13,rep,name=reject_extensions,json=rejectExtensions,proto3" json:"reject_extensions,omitempty"`
	// TODO(phboneff): re-think how this could work with other storage systems
	// storage_config describes Trillian Tessera storage config
	//
	// Types that are assignable to StorageConfig:
	//
	//	*LogConfig_Gcp
	StorageConfig isLogConfig_StorageConfig `protobuf_oneof:"storage_config"`
}

func (x *LogConfig) Reset() {
	*x = LogConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configpb_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogConfig) ProtoMessage() {}

func (x *LogConfig) ProtoReflect() protoreflect.Message {
	mi := &file_configpb_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogConfig.ProtoReflect.Descriptor instead.
func (*LogConfig) Descriptor() ([]byte, []int) {
	return file_configpb_config_proto_rawDescGZIP(), []int{0}
}

func (x *LogConfig) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *LogConfig) GetRootsPemFile() []string {
	if x != nil {
		return x.RootsPemFile
	}
	return nil
}

func (x *LogConfig) GetPrivateKey() *anypb.Any {
	if x != nil {
		return x.PrivateKey
	}
	return nil
}

func (x *LogConfig) GetPublicKey() *keyspb.PublicKey {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *LogConfig) GetRejectExpired() bool {
	if x != nil {
		return x.RejectExpired
	}
	return false
}

func (x *LogConfig) GetRejectUnexpired() bool {
	if x != nil {
		return x.RejectUnexpired
	}
	return false
}

func (x *LogConfig) GetExtKeyUsages() []string {
	if x != nil {
		return x.ExtKeyUsages
	}
	return nil
}

func (x *LogConfig) GetNotAfterStart() *timestamppb.Timestamp {
	if x != nil {
		return x.NotAfterStart
	}
	return nil
}

func (x *LogConfig) GetNotAfterLimit() *timestamppb.Timestamp {
	if x != nil {
		return x.NotAfterLimit
	}
	return nil
}

func (x *LogConfig) GetAcceptOnlyCa() bool {
	if x != nil {
		return x.AcceptOnlyCa
	}
	return false
}

func (x *LogConfig) GetMaxMergeDelaySec() int32 {
	if x != nil {
		return x.MaxMergeDelaySec
	}
	return 0
}

func (x *LogConfig) GetExpectedMergeDelaySec() int32 {
	if x != nil {
		return x.ExpectedMergeDelaySec
	}
	return 0
}

func (x *LogConfig) GetRejectExtensions() []string {
	if x != nil {
		return x.RejectExtensions
	}
	return nil
}

func (m *LogConfig) GetStorageConfig() isLogConfig_StorageConfig {
	if m != nil {
		return m.StorageConfig
	}
	return nil
}

func (x *LogConfig) GetGcp() *GCPConfig {
	if x, ok := x.GetStorageConfig().(*LogConfig_Gcp); ok {
		return x.Gcp
	}
	return nil
}

type isLogConfig_StorageConfig interface {
	isLogConfig_StorageConfig()
}

type LogConfig_Gcp struct {
	Gcp *GCPConfig `protobuf:"bytes,14,opt,name=gcp,proto3,oneof"`
}

func (*LogConfig_Gcp) isLogConfig_StorageConfig() {}

// GCPConfig describes Trillian Tessera GCP config
type GCPConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Bucket        string `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	SpannerDbPath string `protobuf:"bytes,3,opt,name=spanner_db_path,json=spannerDbPath,proto3" json:"spanner_db_path,omitempty"`
}

func (x *GCPConfig) Reset() {
	*x = GCPConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configpb_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GCPConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GCPConfig) ProtoMessage() {}

func (x *GCPConfig) ProtoReflect() protoreflect.Message {
	mi := &file_configpb_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GCPConfig.ProtoReflect.Descriptor instead.
func (*GCPConfig) Descriptor() ([]byte, []int) {
	return file_configpb_config_proto_rawDescGZIP(), []int{1}
}

func (x *GCPConfig) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *GCPConfig) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *GCPConfig) GetSpannerDbPath() string {
	if x != nil {
		return x.SpannerDbPath
	}
	return ""
}

var File_configpb_config_proto protoreflect.FileDescriptor

var file_configpb_config_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70,
	0x62, 0x1a, 0x1a, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x62,
	0x2f, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa8, 0x05, 0x0a, 0x09, 0x4c, 0x6f,
	0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12,
	0x24, 0x0a, 0x0e, 0x72, 0x6f, 0x6f, 0x74, 0x73, 0x5f, 0x70, 0x65, 0x6d, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x6f, 0x6f, 0x74, 0x73, 0x50, 0x65,
	0x6d, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x35, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x0a,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x25,
	0x0a, 0x0e, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x45, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x64, 0x12, 0x29, 0x0a, 0x10, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x75, 0x6e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x55, 0x6e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
	0x12, 0x24, 0x0a, 0x0e, 0x65, 0x78, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x42, 0x0a, 0x0f, 0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x66,
	0x74, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0d, 0x6e, 0x6f, 0x74,
	0x41, 0x66, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x42, 0x0a, 0x0f, 0x6e, 0x6f,
	0x74, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0d, 0x6e, 0x6f, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x24,
	0x0a, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x63, 0x61,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x4f, 0x6e,
	0x6c, 0x79, 0x43, 0x61, 0x12, 0x2d, 0x0a, 0x13, 0x6d, 0x61, 0x78, 0x5f, 0x6d, 0x65, 0x72, 0x67,
	0x65, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x73, 0x65, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x10, 0x6d, 0x61, 0x78, 0x4d, 0x65, 0x72, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x61, 0x79,
	0x53, 0x65, 0x63, 0x12, 0x37, 0x0a, 0x18, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x73, 0x65, 0x63, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x15, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4d,
	0x65, 0x72, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x53, 0x65, 0x63, 0x12, 0x2b, 0x0a, 0x11,
	0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x72, 0x65, 0x6a, 0x65, 0x63, 0x74, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x27, 0x0a, 0x03, 0x67, 0x63, 0x70,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70,
	0x62, 0x2e, 0x47, 0x43, 0x50, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x48, 0x00, 0x52, 0x03, 0x67,
	0x63, 0x70, 0x42, 0x10, 0x0a, 0x0e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x22, 0x6a, 0x0a, 0x09, 0x47, 0x43, 0x50, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x73, 0x70, 0x61, 0x6e,
	0x6e, 0x65, 0x72, 0x5f, 0x64, 0x62, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x73, 0x70, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x44, 0x62, 0x50, 0x61, 0x74, 0x68,
	0x42, 0x4b, 0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x2d, 0x64, 0x65, 0x76, 0x2f,
	0x74, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x61, 0x6e, 0x2d, 0x74, 0x65, 0x73, 0x73, 0x65, 0x72, 0x61,
	0x2f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x73,
	0x63, 0x74, 0x66, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_configpb_config_proto_rawDescOnce sync.Once
	file_configpb_config_proto_rawDescData = file_configpb_config_proto_rawDesc
)

func file_configpb_config_proto_rawDescGZIP() []byte {
	file_configpb_config_proto_rawDescOnce.Do(func() {
		file_configpb_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_configpb_config_proto_rawDescData)
	})
	return file_configpb_config_proto_rawDescData
}

var file_configpb_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_configpb_config_proto_goTypes = []any{
	(*LogConfig)(nil),             // 0: configpb.LogConfig
	(*GCPConfig)(nil),             // 1: configpb.GCPConfig
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
	(*keyspb.PublicKey)(nil),      // 3: keyspb.PublicKey
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_configpb_config_proto_depIdxs = []int32{
	2, // 0: configpb.LogConfig.private_key:type_name -> google.protobuf.Any
	3, // 1: configpb.LogConfig.public_key:type_name -> keyspb.PublicKey
	4, // 2: configpb.LogConfig.not_after_start:type_name -> google.protobuf.Timestamp
	4, // 3: configpb.LogConfig.not_after_limit:type_name -> google.protobuf.Timestamp
	1, // 4: configpb.LogConfig.gcp:type_name -> configpb.GCPConfig
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_configpb_config_proto_init() }
func file_configpb_config_proto_init() {
	if File_configpb_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_configpb_config_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LogConfig); i {
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
		file_configpb_config_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GCPConfig); i {
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
	file_configpb_config_proto_msgTypes[0].OneofWrappers = []any{
		(*LogConfig_Gcp)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_configpb_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_configpb_config_proto_goTypes,
		DependencyIndexes: file_configpb_config_proto_depIdxs,
		MessageInfos:      file_configpb_config_proto_msgTypes,
	}.Build()
	File_configpb_config_proto = out.File
	file_configpb_config_proto_rawDesc = nil
	file_configpb_config_proto_goTypes = nil
	file_configpb_config_proto_depIdxs = nil
}
