/**
  @author: honor
  @since: 2021/7/8
  @desc: //TODO
**/
package security

import "crypto/sha256"

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
	result := string(h.Sum(nil))
	return &result, nil
}
