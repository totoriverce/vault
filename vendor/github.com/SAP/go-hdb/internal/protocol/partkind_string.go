// Code generated by "stringer -type=partKind"; DO NOT EDIT.

package protocol

import "strconv"

const _partKind_name = "pkNilpkCommandpkResultsetpkErrorpkStatementIDpkTransactionIDpkRowsAffectedpkResultsetIDpkTopologyInformationpkTableLocationpkReadLobRequestpkReadLobReplypkAbapIStreampkAbapOStreampkCommandInfopkWriteLobRequestpkWriteLobReplypkParameterspkAuthenticationpkSessionContextpkStatementContextpkPartitionInformationpkOutputParameterspkConnectOptionspkCommitOptionspkFetchOptionspkFetchSizepkParameterMetadatapkResultMetadatapkFindLobRequestpkFindLobReplypkItabSHMpkItabChunkMetadatapkItabMetadatapkItabResultChunkpkClientInfopkStreamDatapkOStreamResultpkFDARequestMetadatapkFDAReplyMetadatapkTransactionFlags"

var _partKind_map = map[partKind]string{
	0:  _partKind_name[0:5],
	3:  _partKind_name[5:14],
	5:  _partKind_name[14:25],
	6:  _partKind_name[25:32],
	10: _partKind_name[32:45],
	11: _partKind_name[45:60],
	12: _partKind_name[60:74],
	13: _partKind_name[74:87],
	15: _partKind_name[87:108],
	16: _partKind_name[108:123],
	17: _partKind_name[123:139],
	18: _partKind_name[139:153],
	25: _partKind_name[153:166],
	26: _partKind_name[166:179],
	27: _partKind_name[179:192],
	28: _partKind_name[192:209],
	30: _partKind_name[209:224],
	32: _partKind_name[224:236],
	33: _partKind_name[236:252],
	34: _partKind_name[252:268],
	39: _partKind_name[268:286],
	40: _partKind_name[286:308],
	41: _partKind_name[308:326],
	42: _partKind_name[326:342],
	43: _partKind_name[342:357],
	44: _partKind_name[357:371],
	45: _partKind_name[371:382],
	47: _partKind_name[382:401],
	48: _partKind_name[401:417],
	49: _partKind_name[417:433],
	50: _partKind_name[433:447],
	51: _partKind_name[447:456],
	53: _partKind_name[456:475],
	55: _partKind_name[475:489],
	56: _partKind_name[489:506],
	57: _partKind_name[506:518],
	58: _partKind_name[518:530],
	59: _partKind_name[530:545],
	60: _partKind_name[545:565],
	61: _partKind_name[565:583],
	64: _partKind_name[583:601],
}

func (i partKind) String() string {
	if str, ok := _partKind_map[i]; ok {
		return str
	}
	return "partKind(" + strconv.FormatInt(int64(i), 10) + ")"
}
