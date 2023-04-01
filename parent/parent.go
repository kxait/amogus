package parent

import (
	"amogus/config"
	"crypto/sha512"
	"fmt"
)

func RunParent(hashesPath string, configPath string, output string) error {

	cfg, err := config.GetConfig(configPath)

	if err != nil {
		return err
	}

	_, err = config.ReadHashesFile(hashesPath, cfg.Mode)

	if err != nil {
		return err
	}

	_ = config.CreateOutputAppender(output)

	fmt.Printf("%+v\n", cfg)

	return nil
}

type hashPair struct {
	hash   string
	origin string
}

func generateHashes(cfg *config.AmogusConfig, last string, amount int) []hashPair {
	var result []hashPair

	workingLast := last

	for i := 0; i < amount; i++ {
		next := GetNextValue(cfg, workingLast)
		var hash *hashPair
		if cfg.Mode == config.Sha512 {
			hash = hashSha512(next)
		} else {
			panic("at the disco")
		}

		result = append(result, *hash)
		workingLast = next
	}

	return result
}

func hashSha512(origin string) *hashPair {
	bytes := []byte(origin)

	hash := sha512.Sum512(bytes)
	result := &hashPair{
		hash:   string(hash[:]),
		origin: origin,
	}

	return result
}
