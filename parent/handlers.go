package parent

import (
	"amogus/common"
	"amogus/config"
	"amogus/next_value"

	"encoding/json"
	"fmt"
	"strconv"

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
		partNo, err := strconv.Atoi(m.Content)

		if err != nil {
			return nil, err
		}

		hashesInfo, err := config.GetHashesInfo(hashesPath, nil)

		if err != nil {
			return nil, err
		}

		if hashesInfo.Parts <= int64(partNo) || partNo < 0 {
			return nil, fmt.Errorf("part number %d is out of range (max %d)", partNo, hashesInfo.Parts-1)
		}

		part, err := config.GetHashesPart(hashesPath, int64(partNo))

		if err != nil {
			return nil, err
		}

		return m.CreateResponse(string(part)), nil
	}
}

// format: 'hash origin'
func hashCracked(oa config.OutputAppender) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		hash, origin := "", ""
		_, err := fmt.Sscanf(m.Content, "%s %s", &hash, &origin)

		if err != nil {
			return nil, err
		}

		fmt.Printf("CRACKED! hash '%s' origin '%s'\n", hash, origin)
		oa(m.Content)

		return m.CreateResponse(""), nil
	}
}

func getNextAssignment(cfg *config.AmogusConfig, s *parentState, hashesPath string) pvm_rpc.RpcHandler {
	return func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
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
