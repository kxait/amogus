package common

import (
	"amogus/config"
	"fmt"
	"strings"
)

func GetNextValue(cfg *config.AmogusConfig, value string) string {
	if value == "" {
		return repeatString(string(cfg.Characters[0]), int(cfg.LengthStart))
	}

	return getNextValueInternal(cfg, value, 1)
}

func GetNextValueOffset(cfg *config.AmogusConfig, value string, offset int64) string {
	if value == "" {
		return getNextValueInternal(cfg, repeatString(string(cfg.Characters[0]), int(cfg.LengthStart)), offset)
	}

	return getNextValueInternal(cfg, value, offset)
}

func getNextValueInternal(cfg *config.AmogusConfig, value string, offset int64) string {
	chars := strings.Split(value, "")

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	indices := make([]int, cfg.Base())
	for i := 0; i < len(value); i++ {
		indices[i] = strings.Index(cfg.Characters, chars[i]) + 1
	}

	maxValue := cfg.Base()

	var count int64
	for count = 0; count < offset; count++ {
		for i := 0; i < cfg.Base(); i++ {
			if indices[i] == maxValue {
				indices[i] = 1
				if i == len(indices)-2 {
					return ""
				}
			} else {
				indices[i]++
				break
			}
		}
	}

	result := ""
	for i := 0; indices[i] != 0; i++ {
		result = fmt.Sprintf("%c%s", cfg.Characters[indices[i]-1], result)
	}

	if len(result) > int(cfg.LengthEnd) {
		return ""
	}

	return result
}

func repeatString(subject string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result = fmt.Sprintf("%s%s", result, subject)
	}
	return result
}
