package config

import (
	"amogus/common"
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
)

const DefaultPartLength int64 = 32000

type HashesInfo struct {
	Parts      int64
	PartLength int64
	ShadowMode common.ShadowMode
}

func GetHashesInfo(filename string, shadowMode *common.ShadowMode) (*HashesInfo, error) {
	st, err := os.Stat(filename)

	if err != nil {
		return nil, err
	}

	var shMode common.ShadowMode
	if shadowMode != nil {
		shMode = *shadowMode
	}

	return &HashesInfo{
		Parts:      int64(math.Ceil(float64(st.Size()) / float64(DefaultPartLength))),
		PartLength: DefaultPartLength,
		ShadowMode: shMode,
	}, nil
}

func GetHashesPart(filename string, partNo int64) ([]byte, error) {
	st, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	offset := DefaultPartLength * partNo
	lengthToEnd := st.Size() - offset
	var bufsize int
	if lengthToEnd < DefaultPartLength {
		bufsize = int(lengthToEnd)
	} else {
		bufsize = int(DefaultPartLength)
	}

	file.Seek(DefaultPartLength*partNo, 0)

	buf := make([]byte, bufsize)

	file.Read(buf)

	return buf, nil
}

func ValidateHashes(filename string, mode Mode) (error, *common.ShadowMode) {
	if mode == Sha512 {
		return validateSha512(filename), nil
	} else if mode == Sha256 {
		return validateSha256(filename), nil
	} else if mode == Shadow {
		err, mode := validateShadow(filename)
		if err != nil {
			return err, nil
		}
		return nil, &mode
	}

	return fmt.Errorf("unsupported mode %s", mode), nil
}

func validateSha512(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if !validSha512(line) {
			return fmt.Errorf("line %s was invalid sha512", line)
		}
	}

	return nil
}

func validSha512(line string) bool {
	re, _ := regexp.Compile("^[a-fA-F0-9]{128}$")
	return re.MatchString(line)
}

func validateSha256(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if !validSha256(line) {
			return fmt.Errorf("line %s was invalid sha256", line)
		}
	}

	return nil
}

func validSha256(line string) bool {
	re, _ := regexp.Compile("^[a-fA-F0-9]{64}$")
	return re.MatchString(line)
}

func validateShadow(filename string) (error, common.ShadowMode) {
	file, err := os.Open(filename)
	if err != nil {
		return err, 0
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	modeRe, _ := regexp.Compile("^\\$[0-9a-z]+\\$")
	mode := modeRe.FindString(line)

	var modeType common.ShadowMode
	// sha512crypt
	if mode == "$6$" {
		if !validShadowSha512Line(line) {
			return fmt.Errorf("line %s was invalid sha512crypt", line), 0
		}
		modeType = common.ShadowSha512
	} else {
		return fmt.Errorf("mode %s is not supported", mode), 0
	}

	if scanner.Scan() {
		return fmt.Errorf("this mode only supports single (one line) cracking"), 0
	}

	return nil, modeType
}

func validShadowSha512Line(line string) bool {
	re, _ := regexp.Compile("^\\$6\\$[a-zA-Z0-9./]+\\$[a-zA-Z0-9./]{86}$")
	return re.MatchString(line)
}
