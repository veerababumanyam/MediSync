/**
 * CopilotKit Components Index
 *
 * Exports all CopilotKit-related components for MediSync.
 *
 * @module components/copilot
 */
export { MediSyncCopilot, CopilotFloatingButton } from './MediSyncCopilot'
export type { MediSyncCopilotProps } from './MediSyncCopilot'

export {
  medisyncTools,
  QueryResultComponent,
  SyncStatusComponent,
  NavigationComponent,
  AlertCreatedComponent,
  ReportCreatedComponent,
  ExportStatusComponent,
} from './MediSyncTools'

export type {
  QueryBIParams,
  SyncTallyParams,
  PinChartParams,
  NavigateParams,
  CreateAlertParams,
  CreateReportParams,
  ExportParams,
  ToolResult,
  QueryResult,
} from './MediSyncTools'
