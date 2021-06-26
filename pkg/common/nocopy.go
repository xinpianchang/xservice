package common

// NoCopy may be embedded into structs which must not be copied
// after the first use.
//
// refer https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
type NoCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*NoCopy) Lock()   {}
func (*NoCopy) Unlock() {}
