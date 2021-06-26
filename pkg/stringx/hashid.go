package stringx

import "github.com/speps/go-hashids/v2"

type simpleHashId struct {
	hid *hashids.HashID
}

func NewSimpleHashId(salt string, minLength int, alphabet ...string) *simpleHashId {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	if len(alphabet) > 0 {
		hd.Alphabet = alphabet[0]
	}
	hid, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	return &simpleHashId{
		hid: hid,
	}
}

func (t *simpleHashId) EncodeInt64(id int64) (string, error) {
	return t.hid.EncodeInt64([]int64{id})
}

func (t *simpleHashId) DecodeInt64(hash string) (int64, error) {
	r, err := t.hid.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}
	return r[0], nil
}
