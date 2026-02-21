// Complete dashboard screen implementation for MediSync
// Demonstrates Riverpod, PowerSync, RTL support, and offline patterns

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:intl/intl.dart';

part 'dashboard_screen.freezed.dart';
part 'dashboard_screen.g.dart';

// =============================================================================
// STATE CLASSES
// =============================================================================

@freezed
class DashboardState with _$DashboardState {
  const factory DashboardState.loading() = _Loading;
  const factory DashboardState.loaded(DashboardData data) = _Loaded;
  const factory DashboardState.error(String message) = _Error;
  const factory DashboardState.offline(DashboardData data) = _Offline;
}

@freezed
class DashboardData with _$DashboardData {
  const factory DashboardData({
    required double totalRevenue,
    required int patientCount,
    required int appointmentCount,
    required double pendingAmount,
    required List<RevenuePoint> revenueTrend,
    required List<AppointmentSummary> todayAppointments,
    required DateTime lastUpdated,
  }) = _DashboardData;

  factory DashboardData.fromJson(Map<String, dynamic> json) =>
      _$DashboardDataFromJson(json);
}

@freezed
class RevenuePoint with _$RevenuePoint {
  const factory RevenuePoint({
    required DateTime date,
    required double amount,
  }) = _RevenuePoint;

  factory RevenuePoint.fromJson(Map<String, dynamic> json) =>
      _$RevenuePointFromJson(json);
}

@freezed
class AppointmentSummary with _$AppointmentSummary {
  const factory AppointmentSummary({
    required String id,
    required String patientName,
    required DateTime scheduledAt,
    required String status,
  }) = _AppointmentSummary;

  factory AppointmentSummary.fromJson(Map<String, dynamic> json) =>
      _$AppointmentSummaryFromJson(json);
}

// =============================================================================
// PROVIDERS
// =============================================================================

@riverpod
class DashboardNotifier extends _$DashboardNotifier {
  @override
  DashboardState build() {
    // Load data on init
    Future.microtask(() => loadDashboard());
    return const DashboardState.loading();
  }

  Future<void> loadDashboard() async {
    state = const DashboardState.loading();

    try {
      final repository = ref.read(dashboardRepositoryProvider);
      final data = await repository.fetchDashboard();
      state = DashboardState.loaded(data);
    } on OfflineException catch (_) {
      // Load from local cache
      final cachedData = await ref.read(localCacheProvider).getDashboard();
      if (cachedData != null) {
        state = DashboardState.offline(cachedData);
      } else {
        state = const DashboardState.error('No cached data available');
      }
    } catch (e) {
      state = DashboardState.error(e.toString());
    }
  }

  Future<void> refresh() async {
    await loadDashboard();
  }
}

// =============================================================================
// SCREEN
// =============================================================================

class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(dashboardNotifierProvider);
    final l10n = AppLocalizations.of(context)!;
    final isRTL = Directionality.of(context) == TextDirection.rtl;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.dashboardTitle),
        actions: [
          // Sync status indicator
          const SyncStatusIndicator(),
          const SizedBox(width: 8),
          // Refresh button
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () =>
                ref.read(dashboardNotifierProvider.notifier).refresh(),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () =>
            ref.read(dashboardNotifierProvider.notifier).refresh(),
        child: state.when(
          loading: () => const Center(child: CircularProgressIndicator()),
          loaded: (data) => _buildContent(context, ref, data, false),
          offline: (data) => _buildContent(context, ref, data, true),
          error: (message) => _buildError(context, ref, message),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/chat'),
        child: const Icon(Icons.chat),
      ),
    );
  }

  Widget _buildContent(
    BuildContext context,
    WidgetRef ref,
    DashboardData data,
    bool isOffline,
  ) {
    return SingleChildScrollView(
      physics: const AlwaysScrollableScrollPhysics(),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Offline banner
          if (isOffline) _buildOfflineBanner(context),
          const SizedBox(height: 8),

          // Key metrics
          _buildMetricsRow(context, data),
          const SizedBox(height: 24),

          // Revenue chart
          _buildRevenueChart(context, data),
          const SizedBox(height: 24),

          // Today's appointments
          _buildAppointmentsList(context, data),
        ],
      ),
    );
  }

  Widget _buildOfflineBanner(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.orange.shade100,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: Colors.orange.shade300),
      ),
      child: Row(
        children: [
          Icon(Icons.cloud_off, color: Colors.orange.shade700),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              l10n.offlineDataWarning,
              style: TextStyle(color: Colors.orange.shade900),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMetricsRow(BuildContext context, DashboardData data) {
    final l10n = AppLocalizations.of(context)!;
    final locale = Localizations.localeOf(context);

    return Row(
      children: [
        Expanded(
          child: _MetricCard(
            title: l10n.totalRevenue,
            value: _formatCurrency(data.totalRevenue, locale),
            trend: 12.5,
            icon: Icons.attach_money,
            color: Colors.green,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _MetricCard(
            title: l10n.patients,
            value: _formatNumber(data.patientCount, locale),
            trend: 8.2,
            icon: Icons.people,
            color: Colors.blue,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _MetricCard(
            title: l10n.appointments,
            value: _formatNumber(data.appointmentCount, locale),
            trend: -2.1,
            icon: Icons.calendar_today,
            color: Colors.purple,
          ),
        ),
      ],
    );
  }

  Widget _buildRevenueChart(BuildContext context, DashboardData data) {
    final l10n = AppLocalizations.of(context)!;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              l10n.revenueTrend,
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 16),
            SizedBox(
              height: 200,
              child: RevenueChart(points: data.revenueTrend),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAppointmentsList(BuildContext context, DashboardData data) {
    final l10n = AppLocalizations.of(context)!;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  l10n.todayAppointments,
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                TextButton(
                  onPressed: () => context.go('/appointments'),
                  child: Text(l10n.viewAll),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (data.todayAppointments.isEmpty)
              Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Text(l10n.noAppointmentsToday),
                ),
              )
            else
              ListView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: data.todayAppointments.length,
                itemBuilder: (context, index) {
                  final appointment = data.todayAppointments[index];
                  return _AppointmentTile(appointment: appointment);
                },
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildError(BuildContext context, WidgetRef ref, String message) {
    final l10n = AppLocalizations.of(context)!;

    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 48, color: Colors.red),
          const SizedBox(height: 16),
          Text(message, textAlign: TextAlign.center),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: () =>
                ref.read(dashboardNotifierProvider.notifier).refresh(),
            child: Text(l10n.retry),
          ),
        ],
      ),
    );
  }

  String _formatCurrency(double amount, Locale locale) {
    return NumberFormat.currency(
      locale: locale.toString(),
      symbol: '\$',
    ).format(amount);
  }

  String _formatNumber(int number, Locale locale) {
    return NumberFormat.compact(locale: locale.toString()).format(number);
  }
}

// =============================================================================
// COMPONENTS
// =============================================================================

class _MetricCard extends StatelessWidget {
  const _MetricCard({
    required this.title,
    required this.value,
    required this.trend,
    required this.icon,
    required this.color,
  });

  final String title;
  final String value;
  final double trend;
  final IconData icon;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final isPositiveTrend = trend >= 0;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Icon(icon, color: color, size: 24),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: isPositiveTrend
                        ? Colors.green.withOpacity(0.1)
                        : Colors.red.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isPositiveTrend ? Icons.trending_up : Icons.trending_down,
                        size: 14,
                        color: isPositiveTrend ? Colors.green : Colors.red,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${trend.abs().toStringAsFixed(1)}%',
                        style: TextStyle(
                          fontSize: 12,
                          color: isPositiveTrend ? Colors.green : Colors.red,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            const SizedBox(height: 4),
            Text(
              title,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey.shade600,
                  ),
            ),
          ],
        ),
      ),
    );
  }
}

class _AppointmentTile extends StatelessWidget {
  const _AppointmentTile({required this.appointment});

  final AppointmentSummary appointment;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final timeFormat = DateFormat.jm(Localizations.localeOf(context).toString());

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: _getStatusColor(appointment.status).withOpacity(0.1),
        child: Icon(
          Icons.person,
          color: _getStatusColor(appointment.status),
        ),
      ),
      title: Text(appointment.patientName),
      subtitle: Text(timeFormat.format(appointment.scheduledAt)),
      trailing: _buildStatusChip(context, appointment.status),
      onTap: () => context.go('/appointments/${appointment.id}'),
    );
  }

  Widget _buildStatusChip(BuildContext context, String status) {
    final l10n = AppLocalizations.of(context)!;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: _getStatusColor(status).withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        _getStatusLabel(status, l10n),
        style: TextStyle(
          fontSize: 12,
          color: _getStatusColor(status),
        ),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'confirmed':
        return Colors.green;
      case 'pending':
        return Colors.orange;
      case 'cancelled':
        return Colors.red;
      default:
        return Colors.grey;
    }
  }

  String _getStatusLabel(String status, AppLocalizations l10n) {
    switch (status) {
      case 'confirmed':
        return l10n.statusConfirmed;
      case 'pending':
        return l10n.statusPending;
      case 'cancelled':
        return l10n.statusCancelled;
      default:
        return status;
    }
  }
}

// =============================================================================
// SYNC STATUS INDICATOR
// =============================================================================

class SyncStatusIndicator extends ConsumerWidget {
  const SyncStatusIndicator({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final syncStatus = ref.watch(syncStatusProvider);
    final l10n = AppLocalizations.of(context)!;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: _getStatusColor(syncStatus).withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _getStatusIcon(syncStatus),
            size: 14,
            color: _getStatusColor(syncStatus),
          ),
          const SizedBox(width: 4),
          Text(
            _getStatusText(syncStatus, l10n),
            style: TextStyle(
              fontSize: 12,
              color: _getStatusColor(syncStatus),
            ),
          ),
        ],
      ),
    );
  }

  Color _getStatusColor(SyncStatus status) {
    switch (status) {
      case SyncStatus.connected:
        return Colors.green;
      case SyncStatus.syncing:
        return Colors.blue;
      case SyncStatus.offline:
        return Colors.orange;
      case SyncStatus.error:
        return Colors.red;
    }
  }

  IconData _getStatusIcon(SyncStatus status) {
    switch (status) {
      case SyncStatus.connected:
        return Icons.cloud_done;
      case SyncStatus.syncing:
        return Icons.sync;
      case SyncStatus.offline:
        return Icons.cloud_off;
      case SyncStatus.error:
        return Icons.error_outline;
    }
  }

  String _getStatusText(SyncStatus status, AppLocalizations l10n) {
    switch (status) {
      case SyncStatus.connected:
        return l10n.syncConnected;
      case SyncStatus.syncing:
        return l10n.syncing;
      case SyncStatus.offline:
        return l10n.offline;
      case SyncStatus.error:
        return l10n.syncError;
    }
  }
}
