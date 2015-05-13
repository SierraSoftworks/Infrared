// Code generated by protoc-gen-go.
// source: lib/infrared.proto
// DO NOT EDIT!

/*
Package lib is a generated protocol buffer package.

It is generated from these files:
	lib/infrared.proto

It has these top-level messages:
	Heartbeat
*/
package infrared

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Heartbeat struct {
	Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	NodeType         *string `protobuf:"bytes,2,req,name=node_type" json:"node_type,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Heartbeat) Reset()         { *m = Heartbeat{} }
func (m *Heartbeat) String() string { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()    {}

func (m *Heartbeat) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Heartbeat) GetNodeType() string {
	if m != nil && m.NodeType != nil {
		return *m.NodeType
	}
	return ""
}

func init() {
}
