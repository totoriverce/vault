// Code generated by "stringer -type=topologyOption"; DO NOT EDIT.

package protocol

import "strconv"

const _topologyOption_name = "toHostNametoHostPortnumbertoTenantNametoLoadfactortoVolumeIDtoIsMastertoIsCurrentSessiontoServiceTypetoNetworkDomaintoIsStandbytoAllIPAddressestoAllHostNames"

var _topologyOption_index = [...]uint8{0, 10, 26, 38, 50, 60, 70, 88, 101, 116, 127, 143, 157}

func (i topologyOption) String() string {
	i -= 1
	if i < 0 || i >= topologyOption(len(_topologyOption_index)-1) {
		return "topologyOption(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _topologyOption_name[_topologyOption_index[i]:_topologyOption_index[i+1]]
}
