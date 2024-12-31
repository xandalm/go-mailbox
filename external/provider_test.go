package external

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestNewProvider(t *testing.T) {
	pb := &stubProviderBridge{}
	p := NewProvider(pb)

	assert.NotNil(t, p)
	assert.Equal(t, p.(*provider).p, ProviderBridge(pb))
}

func TestProvider_Create(t *testing.T) {

	p := NewProvider(&stubProviderBridge{})

	got, err := p.Create("box_1")

	assert.Nil(t, err)
	assert.NotNil(t, got)
}
