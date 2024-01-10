package v1

import (
	"github.com/grafana/loki/pkg/util/encoding"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	numSeries := 100
	numKeysPerSeries := 10000
	data, _ := mkBasicSeriesWithBlooms(numSeries, numKeysPerSeries, 0, 0xffff, 0, 10000)
	enc := &encoding.Encbuf{}
	data[0].Bloom.Encode(enc)
	dec := encoding.DecWith(enc.Get())
	var bloom Bloom
	bloom.Decode(&dec)

}
