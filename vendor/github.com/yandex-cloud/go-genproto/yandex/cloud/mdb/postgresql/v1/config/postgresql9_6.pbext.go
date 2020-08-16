// Code generated by protoc-gen-goext. DO NOT EDIT.

package postgresql

import (
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

func (m *PostgresqlConfig9_6) SetMaxConnections(v *wrappers.Int64Value) {
	m.MaxConnections = v
}

func (m *PostgresqlConfig9_6) SetSharedBuffers(v *wrappers.Int64Value) {
	m.SharedBuffers = v
}

func (m *PostgresqlConfig9_6) SetTempBuffers(v *wrappers.Int64Value) {
	m.TempBuffers = v
}

func (m *PostgresqlConfig9_6) SetMaxPreparedTransactions(v *wrappers.Int64Value) {
	m.MaxPreparedTransactions = v
}

func (m *PostgresqlConfig9_6) SetWorkMem(v *wrappers.Int64Value) {
	m.WorkMem = v
}

func (m *PostgresqlConfig9_6) SetMaintenanceWorkMem(v *wrappers.Int64Value) {
	m.MaintenanceWorkMem = v
}

func (m *PostgresqlConfig9_6) SetReplacementSortTuples(v *wrappers.Int64Value) {
	m.ReplacementSortTuples = v
}

func (m *PostgresqlConfig9_6) SetAutovacuumWorkMem(v *wrappers.Int64Value) {
	m.AutovacuumWorkMem = v
}

func (m *PostgresqlConfig9_6) SetTempFileLimit(v *wrappers.Int64Value) {
	m.TempFileLimit = v
}

func (m *PostgresqlConfig9_6) SetVacuumCostDelay(v *wrappers.Int64Value) {
	m.VacuumCostDelay = v
}

func (m *PostgresqlConfig9_6) SetVacuumCostPageHit(v *wrappers.Int64Value) {
	m.VacuumCostPageHit = v
}

func (m *PostgresqlConfig9_6) SetVacuumCostPageMiss(v *wrappers.Int64Value) {
	m.VacuumCostPageMiss = v
}

func (m *PostgresqlConfig9_6) SetVacuumCostPageDirty(v *wrappers.Int64Value) {
	m.VacuumCostPageDirty = v
}

func (m *PostgresqlConfig9_6) SetVacuumCostLimit(v *wrappers.Int64Value) {
	m.VacuumCostLimit = v
}

func (m *PostgresqlConfig9_6) SetBgwriterDelay(v *wrappers.Int64Value) {
	m.BgwriterDelay = v
}

func (m *PostgresqlConfig9_6) SetBgwriterLruMaxpages(v *wrappers.Int64Value) {
	m.BgwriterLruMaxpages = v
}

func (m *PostgresqlConfig9_6) SetBgwriterLruMultiplier(v *wrappers.DoubleValue) {
	m.BgwriterLruMultiplier = v
}

func (m *PostgresqlConfig9_6) SetBgwriterFlushAfter(v *wrappers.Int64Value) {
	m.BgwriterFlushAfter = v
}

func (m *PostgresqlConfig9_6) SetBackendFlushAfter(v *wrappers.Int64Value) {
	m.BackendFlushAfter = v
}

func (m *PostgresqlConfig9_6) SetOldSnapshotThreshold(v *wrappers.Int64Value) {
	m.OldSnapshotThreshold = v
}

func (m *PostgresqlConfig9_6) SetWalLevel(v PostgresqlConfig9_6_WalLevel) {
	m.WalLevel = v
}

func (m *PostgresqlConfig9_6) SetSynchronousCommit(v PostgresqlConfig9_6_SynchronousCommit) {
	m.SynchronousCommit = v
}

func (m *PostgresqlConfig9_6) SetCheckpointTimeout(v *wrappers.Int64Value) {
	m.CheckpointTimeout = v
}

func (m *PostgresqlConfig9_6) SetCheckpointCompletionTarget(v *wrappers.DoubleValue) {
	m.CheckpointCompletionTarget = v
}

func (m *PostgresqlConfig9_6) SetCheckpointFlushAfter(v *wrappers.Int64Value) {
	m.CheckpointFlushAfter = v
}

func (m *PostgresqlConfig9_6) SetMaxWalSize(v *wrappers.Int64Value) {
	m.MaxWalSize = v
}

func (m *PostgresqlConfig9_6) SetMinWalSize(v *wrappers.Int64Value) {
	m.MinWalSize = v
}

func (m *PostgresqlConfig9_6) SetMaxStandbyStreamingDelay(v *wrappers.Int64Value) {
	m.MaxStandbyStreamingDelay = v
}

func (m *PostgresqlConfig9_6) SetDefaultStatisticsTarget(v *wrappers.Int64Value) {
	m.DefaultStatisticsTarget = v
}

func (m *PostgresqlConfig9_6) SetConstraintExclusion(v PostgresqlConfig9_6_ConstraintExclusion) {
	m.ConstraintExclusion = v
}

func (m *PostgresqlConfig9_6) SetCursorTupleFraction(v *wrappers.DoubleValue) {
	m.CursorTupleFraction = v
}

func (m *PostgresqlConfig9_6) SetFromCollapseLimit(v *wrappers.Int64Value) {
	m.FromCollapseLimit = v
}

func (m *PostgresqlConfig9_6) SetJoinCollapseLimit(v *wrappers.Int64Value) {
	m.JoinCollapseLimit = v
}

func (m *PostgresqlConfig9_6) SetForceParallelMode(v PostgresqlConfig9_6_ForceParallelMode) {
	m.ForceParallelMode = v
}

func (m *PostgresqlConfig9_6) SetClientMinMessages(v PostgresqlConfig9_6_LogLevel) {
	m.ClientMinMessages = v
}

func (m *PostgresqlConfig9_6) SetLogMinMessages(v PostgresqlConfig9_6_LogLevel) {
	m.LogMinMessages = v
}

func (m *PostgresqlConfig9_6) SetLogMinErrorStatement(v PostgresqlConfig9_6_LogLevel) {
	m.LogMinErrorStatement = v
}

func (m *PostgresqlConfig9_6) SetLogMinDurationStatement(v *wrappers.Int64Value) {
	m.LogMinDurationStatement = v
}

func (m *PostgresqlConfig9_6) SetLogCheckpoints(v *wrappers.BoolValue) {
	m.LogCheckpoints = v
}

func (m *PostgresqlConfig9_6) SetLogConnections(v *wrappers.BoolValue) {
	m.LogConnections = v
}

func (m *PostgresqlConfig9_6) SetLogDisconnections(v *wrappers.BoolValue) {
	m.LogDisconnections = v
}

func (m *PostgresqlConfig9_6) SetLogDuration(v *wrappers.BoolValue) {
	m.LogDuration = v
}

func (m *PostgresqlConfig9_6) SetLogErrorVerbosity(v PostgresqlConfig9_6_LogErrorVerbosity) {
	m.LogErrorVerbosity = v
}

func (m *PostgresqlConfig9_6) SetLogLockWaits(v *wrappers.BoolValue) {
	m.LogLockWaits = v
}

func (m *PostgresqlConfig9_6) SetLogStatement(v PostgresqlConfig9_6_LogStatement) {
	m.LogStatement = v
}

func (m *PostgresqlConfig9_6) SetLogTempFiles(v *wrappers.Int64Value) {
	m.LogTempFiles = v
}

func (m *PostgresqlConfig9_6) SetSearchPath(v string) {
	m.SearchPath = v
}

func (m *PostgresqlConfig9_6) SetRowSecurity(v *wrappers.BoolValue) {
	m.RowSecurity = v
}

func (m *PostgresqlConfig9_6) SetDefaultTransactionIsolation(v PostgresqlConfig9_6_TransactionIsolation) {
	m.DefaultTransactionIsolation = v
}

func (m *PostgresqlConfig9_6) SetStatementTimeout(v *wrappers.Int64Value) {
	m.StatementTimeout = v
}

func (m *PostgresqlConfig9_6) SetLockTimeout(v *wrappers.Int64Value) {
	m.LockTimeout = v
}

func (m *PostgresqlConfig9_6) SetIdleInTransactionSessionTimeout(v *wrappers.Int64Value) {
	m.IdleInTransactionSessionTimeout = v
}

func (m *PostgresqlConfig9_6) SetByteaOutput(v PostgresqlConfig9_6_ByteaOutput) {
	m.ByteaOutput = v
}

func (m *PostgresqlConfig9_6) SetXmlbinary(v PostgresqlConfig9_6_XmlBinary) {
	m.Xmlbinary = v
}

func (m *PostgresqlConfig9_6) SetXmloption(v PostgresqlConfig9_6_XmlOption) {
	m.Xmloption = v
}

func (m *PostgresqlConfig9_6) SetGinPendingListLimit(v *wrappers.Int64Value) {
	m.GinPendingListLimit = v
}

func (m *PostgresqlConfig9_6) SetDeadlockTimeout(v *wrappers.Int64Value) {
	m.DeadlockTimeout = v
}

func (m *PostgresqlConfig9_6) SetMaxLocksPerTransaction(v *wrappers.Int64Value) {
	m.MaxLocksPerTransaction = v
}

func (m *PostgresqlConfig9_6) SetMaxPredLocksPerTransaction(v *wrappers.Int64Value) {
	m.MaxPredLocksPerTransaction = v
}

func (m *PostgresqlConfig9_6) SetArrayNulls(v *wrappers.BoolValue) {
	m.ArrayNulls = v
}

func (m *PostgresqlConfig9_6) SetBackslashQuote(v PostgresqlConfig9_6_BackslashQuote) {
	m.BackslashQuote = v
}

func (m *PostgresqlConfig9_6) SetDefaultWithOids(v *wrappers.BoolValue) {
	m.DefaultWithOids = v
}

func (m *PostgresqlConfig9_6) SetEscapeStringWarning(v *wrappers.BoolValue) {
	m.EscapeStringWarning = v
}

func (m *PostgresqlConfig9_6) SetLoCompatPrivileges(v *wrappers.BoolValue) {
	m.LoCompatPrivileges = v
}

func (m *PostgresqlConfig9_6) SetOperatorPrecedenceWarning(v *wrappers.BoolValue) {
	m.OperatorPrecedenceWarning = v
}

func (m *PostgresqlConfig9_6) SetQuoteAllIdentifiers(v *wrappers.BoolValue) {
	m.QuoteAllIdentifiers = v
}

func (m *PostgresqlConfig9_6) SetStandardConformingStrings(v *wrappers.BoolValue) {
	m.StandardConformingStrings = v
}

func (m *PostgresqlConfig9_6) SetSynchronizeSeqscans(v *wrappers.BoolValue) {
	m.SynchronizeSeqscans = v
}

func (m *PostgresqlConfig9_6) SetTransformNullEquals(v *wrappers.BoolValue) {
	m.TransformNullEquals = v
}

func (m *PostgresqlConfig9_6) SetExitOnError(v *wrappers.BoolValue) {
	m.ExitOnError = v
}

func (m *PostgresqlConfig9_6) SetSeqPageCost(v *wrappers.DoubleValue) {
	m.SeqPageCost = v
}

func (m *PostgresqlConfig9_6) SetRandomPageCost(v *wrappers.DoubleValue) {
	m.RandomPageCost = v
}

func (m *PostgresqlConfig9_6) SetSqlInheritance(v *wrappers.BoolValue) {
	m.SqlInheritance = v
}

func (m *PostgresqlConfig9_6) SetAutovacuumMaxWorkers(v *wrappers.Int64Value) {
	m.AutovacuumMaxWorkers = v
}

func (m *PostgresqlConfig9_6) SetAutovacuumVacuumCostDelay(v *wrappers.Int64Value) {
	m.AutovacuumVacuumCostDelay = v
}

func (m *PostgresqlConfig9_6) SetAutovacuumVacuumCostLimit(v *wrappers.Int64Value) {
	m.AutovacuumVacuumCostLimit = v
}

func (m *PostgresqlConfig9_6) SetAutovacuumNaptime(v *wrappers.Int64Value) {
	m.AutovacuumNaptime = v
}

func (m *PostgresqlConfig9_6) SetArchiveTimeout(v *wrappers.Int64Value) {
	m.ArchiveTimeout = v
}

func (m *PostgresqlConfig9_6) SetTrackActivityQuerySize(v *wrappers.Int64Value) {
	m.TrackActivityQuerySize = v
}

func (m *PostgresqlConfig9_6) SetEffectiveIoConcurrency(v *wrappers.Int64Value) {
	m.EffectiveIoConcurrency = v
}

func (m *PostgresqlConfig9_6) SetEffectiveCacheSize(v *wrappers.Int64Value) {
	m.EffectiveCacheSize = v
}

func (m *PostgresqlConfigSet9_6) SetEffectiveConfig(v *PostgresqlConfig9_6) {
	m.EffectiveConfig = v
}

func (m *PostgresqlConfigSet9_6) SetUserConfig(v *PostgresqlConfig9_6) {
	m.UserConfig = v
}

func (m *PostgresqlConfigSet9_6) SetDefaultConfig(v *PostgresqlConfig9_6) {
	m.DefaultConfig = v
}
