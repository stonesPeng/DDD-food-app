/**
  @author: honor
  @since: 2021/7/8
  @desc: //TODO
**/
package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

/**
 * @Description:  hash256 for str
 * @param str source str
 * @return *string
 * @return error
 */
func Hash(str string) (*string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(str)); err != nil {
		return nil, err
	}
	result := hex.EncodeToString(h.Sum(nil))
	return &result, nil
}

func VerifyPassword(c1 string, c2 string) error {
	if hashStr, err := Hash(c1); err != nil {
		return err
	} else {
		if res := strings.Compare(*hashStr, c2); res != 0 {
			return errors.New("password not be equal")
		}
	}
	return nil
}
