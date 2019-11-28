package ecs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMapper(t *testing.T) {
	a := assert.New(t)

	_, err := NewMapper("./ecs_translating_mapping.csv")
	a.NoError(err)

}

func TestMapper_Size(t *testing.T) {

	a := assert.New(t)

	mapper, err := NewMapper("./ecs_translating_mapping.csv")
	a.NoError(err)

	a.Equal(27, mapper.Size())
}

func TestMapper_Get(t *testing.T) {

	a := assert.New(t)

	mapper, err := NewMapper("./ecs_translating_mapping.csv")
	a.NoError(err)

	ecsKey := mapper.EcsField("connection.src_name")
	a.Equal("source.host.name", ecsKey)

	ecsKey2 := mapper.EcsField("connection.dst_ip")
	a.Equal("destination.ip", ecsKey2)

	nonEcsKey := mapper.EcsField("feature_name")
	a.Equal("feature_name", nonEcsKey)
}
