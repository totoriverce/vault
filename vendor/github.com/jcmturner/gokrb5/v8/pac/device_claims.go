package pac

import (
	"bytes"
	"fmt"

	"github.com/jcmturner/rpc/v2/mstypes"
	"github.com/jcmturner/rpc/v2/ndr"
)

// Claims reference: https://msdn.microsoft.com/en-us/library/hh553895.aspx

// DeviceClaimsInfo implements https://msdn.microsoft.com/en-us/library/hh554226.aspx
type DeviceClaimsInfo struct {
	ClaimsSetMetadata mstypes.ClaimsSetMetadata
	ClaimsSet         mstypes.ClaimsSet
}

// Unmarshal bytes into the ClientClaimsInfo struct
func (k *DeviceClaimsInfo) Unmarshal(b []byte) (err error) {
	dec := ndr.NewDecoder(bytes.NewReader(b))
	m := new(mstypes.ClaimsSetMetadata)
	err = dec.Decode(m)
	if err != nil {
		err = fmt.Errorf("error unmarshaling ClientClaimsInfo ClaimsSetMetadata: %v", err)
	}
	k.ClaimsSetMetadata = *m
	k.ClaimsSet, err = k.ClaimsSetMetadata.ClaimsSet()
	if err != nil {
		err = fmt.Errorf("error unmarshaling ClientClaimsInfo ClaimsSet: %v", err)
	}
	return
}
