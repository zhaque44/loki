package v1

import (
	"github.com/grafana/loki/pkg/util/encoding"
	"testing"
)

func TestEncodeDecode(_ *testing.T) {
	numSeries := 100
	numKeysPerSeries := 10000
	data, _ := mkBasicSeriesWithBlooms(numSeries, numKeysPerSeries, 0, 0xffff, 0, 10000)
	enc := &encoding.Encbuf{}
	_ = data[0].Bloom.Encode(enc)
	dec := encoding.DecWith(enc.Get())
	var bloom Bloom
	_ = bloom.Decode(&dec)
}
