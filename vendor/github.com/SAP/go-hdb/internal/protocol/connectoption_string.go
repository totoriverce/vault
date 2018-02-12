// Code generated by "stringer -type=connectOption"; DO NOT EDIT.

package protocol

import "strconv"

const (
	_connectOption_name_0 = "coConnectionIDcoCompleteArrayExecutioncoClientLocalecoSupportsLargeBulkOperations"
	_connectOption_name_1 = "coLargeNumberOfParameterSupportcoSystemID"
	_connectOption_name_2 = "coAbapVarcharModecoSelectForUpdateSupportedcoClientDistributionModecoEngineDataFormatVersioncoDistributionProtocolVersioncoSplitBatchCommandscoUseTransactionFlagsOnly"
	_connectOption_name_3 = "coIgnoreUnknownPartscoTableOutputParametercoDataFormatVersion2coItabParametercoDescribeTableOutputParametercoColumnarResultsetcoScrollablResultSetcoClientInfoNullValueSupportedcoAssociatedConnectionIdcoNoTransactionalPreparecoFDAEnabledcoOSUsercoRowslotImageResultcoEndianess"
	_connectOption_name_4 = "coImplicitLobStreaming"
)

var (
	_connectOption_index_0 = [...]uint8{0, 14, 38, 52, 81}
	_connectOption_index_1 = [...]uint8{0, 31, 41}
	_connectOption_index_2 = [...]uint8{0, 17, 43, 67, 92, 121, 141, 166}
	_connectOption_index_3 = [...]uint16{0, 20, 42, 62, 77, 107, 126, 146, 176, 200, 224, 236, 244, 264, 275}
	_connectOption_index_4 = [...]uint8{0, 22}
)

func (i connectOption) String() string {
	switch {
	case 1 <= i && i <= 4:
		i -= 1
		return _connectOption_name_0[_connectOption_index_0[i]:_connectOption_index_0[i+1]]
	case 10 <= i && i <= 11:
		i -= 10
		return _connectOption_name_1[_connectOption_index_1[i]:_connectOption_index_1[i+1]]
	case 13 <= i && i <= 19:
		i -= 13
		return _connectOption_name_2[_connectOption_index_2[i]:_connectOption_index_2[i+1]]
	case 21 <= i && i <= 34:
		i -= 21
		return _connectOption_name_3[_connectOption_index_3[i]:_connectOption_index_3[i+1]]
	case i == 37:
		return _connectOption_name_4
	default:
		return "connectOption(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
