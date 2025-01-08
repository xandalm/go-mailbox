package external

import (
	"context"
	"testing"

	"github.com/xandalm/go-testing/assert"
)

func TestNewProvider(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	assert.NotNil(t, p)
	assert.Equal(t, p.(*provider).p.(*spyProviderBridge), pb)
}

func TestProvider_CreateHasForwardedToHandleCreate(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Create(context.TODO(), "box")

	assert.Equal(t, pb.HandleCreateCalls, 1)
}

func TestProvider_ContainsHasForwardedToHandleContains(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Contains(context.TODO(), "box")

	assert.Equal(t, pb.HandleContainsCalls, 1)
}

func TestProvider_GetHasForwardedToHandleGet(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Get(context.TODO(), "box")

	assert.Equal(t, pb.HandleGetCalls, 1)
}

func TestProvider_DeleteHasForwardedToHandleDelete(t *testing.T) {
	pb := &spyProviderBridge{}
	p := NewProvider(pb)

	p.Delete(context.TODO(), "box")

	assert.Equal(t, pb.HandleDeleteCalls, 1)
}
