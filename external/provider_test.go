package external

import (
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestNewProvider(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	assert.NotNil(t, p)
	assert.Equal(t, p.(*provider).p.(*spyProviderBridge), pb)
}

func TestProvider_CreateHasForwardedToOnCreate(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Create("box")

	assert.Equal(t, pb.OnCreateCalls, 1)
}

func TestProvider_ContainsHasForwardedToOnContains(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Contains("box")

	assert.Equal(t, pb.OnContainsCalls, 1)
}

func TestProvider_GetHasForwardedToOnGet(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Get("box")

	assert.Equal(t, pb.OnGetCalls, 1)
}

func TestProvider_DeleteHasForwardedToOnDelete(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Delete("box")

	assert.Equal(t, pb.OnDeleteCalls, 1)
}
