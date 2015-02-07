package packages

import "regexp"

func replaceAllSubmatchFunc(re *regexp.Regexp, data []byte, fn func(s []byte) []byte) []byte {
	idxs := re.FindAllSubmatchIndex(data, -1)
	if len(idxs) == 0 {
		return data
	}
	n := len(idxs)
	ret := append([]byte{}, data[:idxs[0][0]]...)
	for i, pair := range idxs {
		// replace internal submatch with result of user supplied function
		ret = append(ret, fn(data[pair[2]:pair[3]])...)
		if i+1 < n {
			ret = append(ret, data[pair[1]:idxs[i+1][0]]...)
		}
	}
	ret = append(ret, data[idxs[len(idxs)-1][1]:]...)
	return ret
}
