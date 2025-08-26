package propagation

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

type GRPCMetadataCarrier metadata.MD

func (c GRPCMetadataCarrier) Get(key string) string {
	vals := metadata.MD(c).Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (c GRPCMetadataCarrier) Set(key string, value string) {
	metadata.MD(c).Set(strings.ToLower(key), value)
}

func (c GRPCMetadataCarrier) Keys() []string {
	out := make([]string, 0, len(c))
	for k := range c {
		out = append(out, k)
	}
	return out
}
