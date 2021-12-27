// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package zipkincore

import (
	"bytes"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var GoUnusedProtection__ int

type AnnotationType int64

const (
	AnnotationType_BOOL   AnnotationType = 0
	AnnotationType_BYTES  AnnotationType = 1
	AnnotationType_I16    AnnotationType = 2
	AnnotationType_I32    AnnotationType = 3
	AnnotationType_I64    AnnotationType = 4
	AnnotationType_DOUBLE AnnotationType = 5
	AnnotationType_STRING AnnotationType = 6
)

func (p AnnotationType) String() string {
	switch p {
	case AnnotationType_BOOL:
		return "BOOL"
	case AnnotationType_BYTES:
		return "BYTES"
	case AnnotationType_I16:
		return "I16"
	case AnnotationType_I32:
		return "I32"
	case AnnotationType_I64:
		return "I64"
	case AnnotationType_DOUBLE:
		return "DOUBLE"
	case AnnotationType_STRING:
		return "STRING"
	}
	return "<UNSET>"
}

func AnnotationTypeFromString(s string) (AnnotationType, error) {
	switch s {
	case "BOOL":
		return AnnotationType_BOOL, nil
	case "BYTES":
		return AnnotationType_BYTES, nil
	case "I16":
		return AnnotationType_I16, nil
	case "I32":
		return AnnotationType_I32, nil
	case "I64":
		return AnnotationType_I64, nil
	case "DOUBLE":
		return AnnotationType_DOUBLE, nil
	case "STRING":
		return AnnotationType_STRING, nil
	}
	return AnnotationType(0), fmt.Errorf("not a valid AnnotationType string")
}

func AnnotationTypePtr(v AnnotationType) *AnnotationType { return &v }

func (p AnnotationType) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *AnnotationType) UnmarshalText(text []byte) error {
	q, err := AnnotationTypeFromString(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}

// Indicates the network context of a service recording an annotation with two
// exceptions.
//
// When a BinaryAnnotation, and key is CLIENT_ADDR or SERVER_ADDR,
// the endpoint indicates the source or destination of an RPC. This exception
// allows zipkin to display network context of uninstrumented services, or
// clients such as web browsers.
//
// Attributes:
//  - Ipv4: IPv4 host address packed into 4 bytes.
//
// Ex for the ip 1.2.3.4, it would be (1 << 24) | (2 << 16) | (3 << 8) | 4
//  - Port: IPv4 port
//
// Note: this is to be treated as an unsigned integer, so watch for negatives.
//
// Conventionally, when the port isn't known, port = 0.
//  - ServiceName: Service name in lowercase, such as "memcache" or "zipkin-web"
//
// Conventionally, when the service name isn't known, service_name = "unknown".
//  - Ipv6: IPv6 host address packed into 16 bytes. Ex Inet6Address.getBytes()
type Endpoint struct {
	Ipv4        uint32 `thrift:"ipv4,1" json:"ipv4"`
	Port        int16  `thrift:"port,2" json:"port"`
	ServiceName string `thrift:"service_name,3" json:"service_name"`
	Ipv6        []byte `thrift:"ipv6,4" json:"ipv6,omitempty"`
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (p *Endpoint) GetIpv4() uint32 {
	return p.Ipv4
}

func (p *Endpoint) GetPort() int16 {
	return p.Port
}

func (p *Endpoint) GetServiceName() string {
	return p.ServiceName
}

var Endpoint_Ipv6_DEFAULT []byte

func (p *Endpoint) GetIpv6() []byte {
	return p.Ipv6
}
func (p *Endpoint) IsSetIpv6() bool {
	return p.Ipv6 != nil
}

func (p *Endpoint) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Endpoint) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Ipv4 = uint32(v)
	}
	return nil
}

func (p *Endpoint) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI16(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Port = v
	}
	return nil
}

func (p *Endpoint) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.ServiceName = v
	}
	return nil
}

func (p *Endpoint) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.Ipv6 = v
	}
	return nil
}

func (p *Endpoint) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Endpoint"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Endpoint) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("ipv4", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ipv4: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Ipv4)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.ipv4 (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ipv4: ", p), err)
	}
	return err
}

func (p *Endpoint) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("port", thrift.I16, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:port: ", p), err)
	}
	if err := oprot.WriteI16(int16(p.Port)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.port (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:port: ", p), err)
	}
	return err
}

func (p *Endpoint) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("service_name", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:service_name: ", p), err)
	}
	if err := oprot.WriteString(string(p.ServiceName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.service_name (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:service_name: ", p), err)
	}
	return err
}

func (p *Endpoint) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetIpv6() {
		if err := oprot.WriteFieldBegin("ipv6", thrift.STRING, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:ipv6: ", p), err)
		}
		if err := oprot.WriteBinary(p.Ipv6); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.ipv6 (4) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:ipv6: ", p), err)
		}
	}
	return err
}

func (p *Endpoint) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Endpoint(%+v)", *p)
}

// An annotation is similar to a log statement. It includes a host field which
// allows these events to be attributed properly, and also aggregatable.
//
// Attributes:
//  - Timestamp: Microseconds from epoch.
//
// This value should use the most precise value possible. For example,
// gettimeofday or syncing nanoTime against a tick of currentTimeMillis.
//  - Value
//  - Host: Always the host that recorded the event. By specifying the host you allow
// rollup of all events (such as client requests to a service) by IP address.
type Annotation struct {
	Timestamp int64     `thrift:"timestamp,1" json:"timestamp"`
	Value     string    `thrift:"value,2" json:"value"`
	Host      *Endpoint `thrift:"host,3" json:"host,omitempty"`
}

func NewAnnotation() *Annotation {
	return &Annotation{}
}

func (p *Annotation) GetTimestamp() int64 {
	return p.Timestamp
}

func (p *Annotation) GetValue() string {
	return p.Value
}

var Annotation_Host_DEFAULT *Endpoint

func (p *Annotation) GetHost() *Endpoint {
	if !p.IsSetHost() {
		return Annotation_Host_DEFAULT
	}
	return p.Host
}
func (p *Annotation) IsSetHost() bool {
	return p.Host != nil
}

func (p *Annotation) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Annotation) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Timestamp = v
	}
	return nil
}

func (p *Annotation) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Value = v
	}
	return nil
}

func (p *Annotation) readField3(iprot thrift.TProtocol) error {
	p.Host = &Endpoint{}
	if err := p.Host.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Host), err)
	}
	return nil
}

func (p *Annotation) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Annotation"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Annotation) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timestamp", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:timestamp: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.Timestamp)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.timestamp (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:timestamp: ", p), err)
	}
	return err
}

func (p *Annotation) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("value", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:value: ", p), err)
	}
	if err := oprot.WriteString(string(p.Value)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.value (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:value: ", p), err)
	}
	return err
}

func (p *Annotation) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetHost() {
		if err := oprot.WriteFieldBegin("host", thrift.STRUCT, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:host: ", p), err)
		}
		if err := p.Host.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Host), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:host: ", p), err)
		}
	}
	return err
}

func (p *Annotation) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Annotation(%+v)", *p)
}

// Binary annotations are tags applied to a Span to give it context. For
// example, a binary annotation of "http.uri" could the path to a resource in a
// RPC call.
//
// Binary annotations of type STRING are always queryable, though more a
// historical implementation detail than a structural concern.
//
// Binary annotations can repeat, and vary on the host. Similar to Annotation,
// the host indicates who logged the event. This allows you to tell the
// difference between the client and server side of the same key. For example,
// the key "http.uri" might be different on the client and server side due to
// rewriting, like "/api/v1/myresource" vs "/myresource. Via the host field,
// you can see the different points of view, which often help in debugging.
//
// Attributes:
//  - Key
//  - Value
//  - AnnotationType
//  - Host: The host that recorded tag, which allows you to differentiate between
// multiple tags with the same key. There are two exceptions to this.
//
// When the key is CLIENT_ADDR or SERVER_ADDR, host indicates the source or
// destination of an RPC. This exception allows zipkin to display network
// context of uninstrumented services, or clients such as web browsers.
type BinaryAnnotation struct {
	Key            string         `thrift:"key,1" json:"key"`
	Value          []byte         `thrift:"value,2" json:"value"`
	AnnotationType AnnotationType `thrift:"annotation_type,3" json:"annotation_type"`
	Host           *Endpoint      `thrift:"host,4" json:"host,omitempty"`
}

func NewBinaryAnnotation() *BinaryAnnotation {
	return &BinaryAnnotation{}
}

func (p *BinaryAnnotation) GetKey() string {
	return p.Key
}

func (p *BinaryAnnotation) GetValue() []byte {
	return p.Value
}

func (p *BinaryAnnotation) GetAnnotationType() AnnotationType {
	return p.AnnotationType
}

var BinaryAnnotation_Host_DEFAULT *Endpoint

func (p *BinaryAnnotation) GetHost() *Endpoint {
	if !p.IsSetHost() {
		return BinaryAnnotation_Host_DEFAULT
	}
	return p.Host
}
func (p *BinaryAnnotation) IsSetHost() bool {
	return p.Host != nil
}

func (p *BinaryAnnotation) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *BinaryAnnotation) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Key = v
	}
	return nil
}

func (p *BinaryAnnotation) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.Value = v
	}
	return nil
}

func (p *BinaryAnnotation) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		temp := AnnotationType(v)
		p.AnnotationType = temp
	}
	return nil
}

func (p *BinaryAnnotation) readField4(iprot thrift.TProtocol) error {
	p.Host = &Endpoint{}
	if err := p.Host.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Host), err)
	}
	return nil
}

func (p *BinaryAnnotation) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("BinaryAnnotation"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BinaryAnnotation) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("key", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:key: ", p), err)
	}
	if err := oprot.WriteString(string(p.Key)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.key (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:key: ", p), err)
	}
	return err
}

func (p *BinaryAnnotation) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("value", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:value: ", p), err)
	}
	if err := oprot.WriteBinary(p.Value); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.value (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:value: ", p), err)
	}
	return err
}

func (p *BinaryAnnotation) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("annotation_type", thrift.I32, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:annotation_type: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.AnnotationType)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.annotation_type (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:annotation_type: ", p), err)
	}
	return err
}

func (p *BinaryAnnotation) writeField4(oprot thrift.TProtocol) (err error) {
	if p.IsSetHost() {
		if err := oprot.WriteFieldBegin("host", thrift.STRUCT, 4); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:host: ", p), err)
		}
		if err := p.Host.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Host), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 4:host: ", p), err)
		}
	}
	return err
}

func (p *BinaryAnnotation) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BinaryAnnotation(%+v)", *p)
}

// A trace is a series of spans (often RPC calls) which form a latency tree.
//
// The root span is where trace_id = id and parent_id = Nil. The root span is
// usually the longest interval in the trace, starting with a SERVER_RECV
// annotation and ending with a SERVER_SEND.
//
// Attributes:
//  - TraceID
//  - Name: Span name in lowercase, rpc method for example
//
// Conventionally, when the span name isn't known, name = "unknown".
//  - ID
//  - ParentID
//  - Annotations
//  - BinaryAnnotations
//  - Debug
//  - Timestamp: Microseconds from epoch of the creation of this span.
//
// This value should be set directly by instrumentation, using the most
// precise value possible. For example, gettimeofday or syncing nanoTime
// against a tick of currentTimeMillis.
//
// For compatibilty with instrumentation that precede this field, collectors
// or span stores can derive this via Annotation.timestamp.
// For example, SERVER_RECV.timestamp or CLIENT_SEND.timestamp.
//
// This field is optional for compatibility with old data: first-party span
// stores are expected to support this at time of introduction.
//  - Duration: Measurement of duration in microseconds, used to support queries.
//
// This value should be set directly, where possible. Doing so encourages
// precise measurement decoupled from problems of clocks, such as skew or NTP
// updates causing time to move backwards.
//
// For compatibilty with instrumentation that precede this field, collectors
// or span stores can derive this by subtracting Annotation.timestamp.
// For example, SERVER_SEND.timestamp - SERVER_RECV.timestamp.
//
// If this field is persisted as unset, zipkin will continue to work, except
// duration query support will be implementation-specific. Similarly, setting
// this field non-atomically is implementation-specific.
//
// This field is i64 vs i32 to support spans longer than 35 minutes.
//  - TraceIDHigh: Optional unique 8-byte additional identifier for a trace. If non zero, this
// means the trace uses 128 bit traceIds instead of 64 bit.
type Span struct {
	TraceID int64 `thrift:"trace_id,1" json:"trace_id"`
	// unused field # 2
	Name        string        `thrift:"name,3" json:"name"`
	ID          int64         `thrift:"id,4" json:"id"`
	ParentID    *int64        `thrift:"parent_id,5" json:"parent_id,omitempty"`
	Annotations []*Annotation `thrift:"annotations,6" json:"annotations"`
	// unused field # 7
	BinaryAnnotations []*BinaryAnnotation `thrift:"binary_annotations,8" json:"binary_annotations"`
	Debug             bool                `thrift:"debug,9" json:"debug,omitempty"`
	Timestamp         *int64              `thrift:"timestamp,10" json:"timestamp,omitempty"`
	Duration          *int64              `thrift:"duration,11" json:"duration,omitempty"`
	TraceIDHigh       *int64              `thrift:"trace_id_high,12" json:"trace_id_high,omitempty"`
}

func NewSpan() *Span {
	return &Span{}
}

func (p *Span) GetTraceID() int64 {
	return p.TraceID
}

func (p *Span) GetName() string {
	return p.Name
}

func (p *Span) GetID() int64 {
	return p.ID
}

var Span_ParentID_DEFAULT int64

func (p *Span) GetParentID() int64 {
	if !p.IsSetParentID() {
		return Span_ParentID_DEFAULT
	}
	return *p.ParentID
}

func (p *Span) GetAnnotations() []*Annotation {
	return p.Annotations
}

func (p *Span) GetBinaryAnnotations() []*BinaryAnnotation {
	return p.BinaryAnnotations
}

var Span_Debug_DEFAULT bool = false

func (p *Span) GetDebug() bool {
	return p.Debug
}

var Span_Timestamp_DEFAULT int64

func (p *Span) GetTimestamp() int64 {
	if !p.IsSetTimestamp() {
		return Span_Timestamp_DEFAULT
	}
	return *p.Timestamp
}

var Span_Duration_DEFAULT int64

func (p *Span) GetDuration() int64 {
	if !p.IsSetDuration() {
		return Span_Duration_DEFAULT
	}
	return *p.Duration
}

var Span_TraceIDHigh_DEFAULT int64

func (p *Span) GetTraceIDHigh() int64 {
	if !p.IsSetTraceIDHigh() {
		return Span_TraceIDHigh_DEFAULT
	}
	return *p.TraceIDHigh
}
func (p *Span) IsSetParentID() bool {
	return p.ParentID != nil
}

func (p *Span) IsSetDebug() bool {
	return p.Debug != Span_Debug_DEFAULT
}

func (p *Span) IsSetTimestamp() bool {
	return p.Timestamp != nil
}

func (p *Span) IsSetDuration() bool {
	return p.Duration != nil
}

func (p *Span) IsSetTraceIDHigh() bool {
	return p.TraceIDHigh != nil
}

func (p *Span) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.readField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.readField6(iprot); err != nil {
				return err
			}
		case 8:
			if err := p.readField8(iprot); err != nil {
				return err
			}
		case 9:
			if err := p.readField9(iprot); err != nil {
				return err
			}
		case 10:
			if err := p.readField10(iprot); err != nil {
				return err
			}
		case 11:
			if err := p.readField11(iprot); err != nil {
				return err
			}
		case 12:
			if err := p.readField12(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *Span) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.TraceID = v
	}
	return nil
}

func (p *Span) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Name = v
	}
	return nil
}

func (p *Span) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 4: ", err)
	} else {
		p.ID = v
	}
	return nil
}

func (p *Span) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.ParentID = &v
	}
	return nil
}

func (p *Span) readField6(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*Annotation, 0, size)
	p.Annotations = tSlice
	for i := 0; i < size; i++ {
		_elem0 := &Annotation{}
		if err := _elem0.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem0), err)
		}
		p.Annotations = append(p.Annotations, _elem0)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Span) readField8(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*BinaryAnnotation, 0, size)
	p.BinaryAnnotations = tSlice
	for i := 0; i < size; i++ {
		_elem1 := &BinaryAnnotation{}
		if err := _elem1.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem1), err)
		}
		p.BinaryAnnotations = append(p.BinaryAnnotations, _elem1)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *Span) readField9(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 9: ", err)
	} else {
		p.Debug = v
	}
	return nil
}

func (p *Span) readField10(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 10: ", err)
	} else {
		p.Timestamp = &v
	}
	return nil
}

func (p *Span) readField11(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 11: ", err)
	} else {
		p.Duration = &v
	}
	return nil
}

func (p *Span) readField12(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 12: ", err)
	} else {
		p.TraceIDHigh = &v
	}
	return nil
}

func (p *Span) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Span"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := p.writeField6(oprot); err != nil {
		return err
	}
	if err := p.writeField8(oprot); err != nil {
		return err
	}
	if err := p.writeField9(oprot); err != nil {
		return err
	}
	if err := p.writeField10(oprot); err != nil {
		return err
	}
	if err := p.writeField11(oprot); err != nil {
		return err
	}
	if err := p.writeField12(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Span) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("trace_id", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:trace_id: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.TraceID)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.trace_id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:trace_id: ", p), err)
	}
	return err
}

func (p *Span) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:name: ", p), err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.name (3) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:name: ", p), err)
	}
	return err
}

func (p *Span) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.I64, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:id: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.ID)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (4) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:id: ", p), err)
	}
	return err
}

func (p *Span) writeField5(oprot thrift.TProtocol) (err error) {
	if p.IsSetParentID() {
		if err := oprot.WriteFieldBegin("parent_id", thrift.I64, 5); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:parent_id: ", p), err)
		}
		if err := oprot.WriteI64(int64(*p.ParentID)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.parent_id (5) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 5:parent_id: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("annotations", thrift.LIST, 6); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:annotations: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Annotations)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Annotations {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 6:annotations: ", p), err)
	}
	return err
}

func (p *Span) writeField8(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("binary_annotations", thrift.LIST, 8); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 8:binary_annotations: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.BinaryAnnotations)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.BinaryAnnotations {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 8:binary_annotations: ", p), err)
	}
	return err
}

func (p *Span) writeField9(oprot thrift.TProtocol) (err error) {
	if p.IsSetDebug() {
		if err := oprot.WriteFieldBegin("debug", thrift.BOOL, 9); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 9:debug: ", p), err)
		}
		if err := oprot.WriteBool(bool(p.Debug)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.debug (9) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 9:debug: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField10(oprot thrift.TProtocol) (err error) {
	if p.IsSetTimestamp() {
		if err := oprot.WriteFieldBegin("timestamp", thrift.I64, 10); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 10:timestamp: ", p), err)
		}
		if err := oprot.WriteI64(int64(*p.Timestamp)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.timestamp (10) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 10:timestamp: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField11(oprot thrift.TProtocol) (err error) {
	if p.IsSetDuration() {
		if err := oprot.WriteFieldBegin("duration", thrift.I64, 11); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 11:duration: ", p), err)
		}
		if err := oprot.WriteI64(int64(*p.Duration)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.duration (11) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 11:duration: ", p), err)
		}
	}
	return err
}

func (p *Span) writeField12(oprot thrift.TProtocol) (err error) {
	if p.IsSetTraceIDHigh() {
		if err := oprot.WriteFieldBegin("trace_id_high", thrift.I64, 12); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 12:trace_id_high: ", p), err)
		}
		if err := oprot.WriteI64(int64(*p.TraceIDHigh)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.trace_id_high (12) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 12:trace_id_high: ", p), err)
		}
	}
	return err
}

func (p *Span) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Span(%+v)", *p)
}

// Attributes:
//  - Ok
type Response struct {
	Ok bool `thrift:"ok,1,required" json:"ok"`
}

func NewResponse() *Response {
	return &Response{}
}

func (p *Response) GetOk() bool {
	return p.Ok
}
func (p *Response) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetOk bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
			issetOk = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetOk {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Ok is not set"))
	}
	return nil
}

func (p *Response) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBool(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Ok = v
	}
	return nil
}

func (p *Response) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Response"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *Response) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("ok", thrift.BOOL, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:ok: ", p), err)
	}
	if err := oprot.WriteBool(bool(p.Ok)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.ok (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:ok: ", p), err)
	}
	return err
}

func (p *Response) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Response(%+v)", *p)
}
