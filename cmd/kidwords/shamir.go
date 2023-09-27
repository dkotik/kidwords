package main

// type sssSet [8]string
//
// func shamirSplit(secret []byte, quorum uint8) (result sssSet, err error) {
// 	if quorum > 8 {
// 		return result, errors.New("the maximum number of quorum shards is three")
// 	}
//
// 	// secret, err = compress(secret)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	parts, err := shamir.Split(secret, 8, int(quorum))
// 	if err != nil {
// 		return
// 	}
//
// 	b := &bytes.Buffer{}
// 	w, err := kidwords.NewWriter(b)
// 	if err != nil {
// 		return
// 	}
//
// 	for i, part := range parts {
// 		_, err = io.Copy(w, bytes.NewReader(part))
// 		if err != nil {
// 			return
// 		}
// 		result[i] = b.String()
// 		b.Reset()
// 	}
// 	return
// }
