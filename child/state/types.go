package state

import (
	"amogus/common"
	"amogus/config"
	"time"

	"github.com/nathanaelle/password/v2"
)

type ChildState struct {
	CurrentAssignment string
	CurrentState      common.ChildState
	Config            config.AmogusConfig
	HashesInfo        config.HashesInfo
	HashPartReceived  int64
	LastChunkStart    time.Time

	ShadowCrypter *password.Crypter
}
