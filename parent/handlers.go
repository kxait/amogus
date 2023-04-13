package parent

import (
	"amogus/common"
	"amogus/config"
	"amogus/next_value"

	"encoding/json"
	"fmt"

	pvm_rpc "github.com/kxait/pvm-rpc"
)

func registerParentHandlers(
	rs *pvm_rpc.RpcServer,
	config *config.AmogusConfig,
	hashesPath string,
	oa config.OutputAppender,
	s *parentState) {
	rs.Handlers[common.GetConfig] = getConfig(config)
	rs.Handlers[common.GetHashesInfo] = getHashesInfo(hashesPath, s)
	rs.Handlers[common.GetHashesPart] = getHashesPart(hashesPath)
	rs.Handlers[common.HashCracked] = hashCracked(oa)
	rs.Handlers[common.GetNextAssignment] = getNextAssignment(config, s, hashesPath)
}

func getConfig(config *config.AmogusConfig) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		serialized, err := json.Marshal(*config)
		if err != nil {
			return nil, err
		}

		return m.CreateResponse(string(serialized)), nil
	}
}

func getHashesInfo(hashesPath string, state *parentState) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		hashesInfo, err := config.GetHashesInfo(hashesPath, &state.shadowMode)

		if err != nil {
			return nil, err
		}

		serialized, err := json.Marshal(hashesInfo)

		return m.CreateResponse(string(serialized)), nil
	}
}

func getHashesPart(hashesPath string) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		partArgs := common.GetHashesPartArgs{}
		err := json.Unmarshal([]byte(m.Content), &partArgs)

		if err != nil {
			return nil, err
		}

		hashesInfo, err := config.GetHashesInfo(hashesPath, nil)

		if err != nil {
			return nil, err
		}

		if hashesInfo.Parts <= int64(partArgs.Part) || partArgs.Part < 0 {
			return nil, fmt.Errorf("part number %d is out of range (max %d)", partArgs.Part, hashesInfo.Parts-1)
		}

		part, err := config.GetHashesPart(hashesPath, int64(partArgs.Part))

		if err != nil {
			return nil, err
		}

		return m.CreateResponse(string(part)), nil
	}
}

// format: 'hash origin'
func hashCracked(oa config.OutputAppender) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		hashCrackedArgs := common.HashCrackedArgs{}
		err := json.Unmarshal([]byte(m.Content), &hashCrackedArgs)

		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		fmt.Printf("CRACKED! hash '%s' origin '%s'\n", hashCrackedArgs.Hash, hashCrackedArgs.Origin)
		oa(m.Content)

		return m.CreateResponse(""), nil
	}
}

func getNextAssignment(cfg *config.AmogusConfig, s *parentState, hashesPath string) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		gnaArgs := common.GetNextAssignmentArgs{}
		err := json.Unmarshal([]byte(m.Content), &gnaArgs)

		if err != nil {
			return nil, err
		}

		seconds := float64(gnaArgs.ChunkTimeMillis) / 1000.0
		hashRatePerSecond := int64(float64(cfg.ChunkSize) / seconds)

		s.hashrate.pushHashRate(m.CallerTaskId, hashRatePerSecond)
		//fmt.Printf("cracked %d hashes in %g seconds\n", cfg.ChunkSize, seconds)

		var next string
		if s.lastOrigin == "" && !s.ranOut {
			next = next_value.GetNextValue(cfg, s.lastOrigin)
		} else if s.lastOrigin != "" && !s.ranOut {
			next = next_value.GetNextValueOffset(cfg, s.lastOrigin, int64(cfg.ChunkSize))
		}

		if next == "" {
			s.ranOut = true
			return nil, fmt.Errorf("finished!")
		}

		res := m.CreateResponse(next)

		s.lastOrigin = next

		return res, nil
	}
}
