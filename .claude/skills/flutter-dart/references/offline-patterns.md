# PowerSync Offline-First Patterns

## Architecture Overview

PowerSync provides a local-first database that syncs with the backend when connectivity is available. This enables offline operation and instant UI updates.

### Sync Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    FLUTTER APP                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    UI Layer                              │   │
│  │         Riverpod Providers · Widgets                     │   │
│  └────────────────────────┬────────────────────────────────┘   │
│                           │                                     │
│  ┌────────────────────────▼────────────────────────────────┐   │
│  │               Repository Layer                           │   │
│  │         Read/Write to Local SQLite                       │   │
│  └────────────────────────┬────────────────────────────────┘   │
│                           │                                     │
│  ┌────────────────────────▼────────────────────────────────┐   │
│  │              PowerSync Database                          │   │
│  │         Local SQLite + Sync Engine                       │   │
│  └────────────────────────┬────────────────────────────────┘   │
└───────────────────────────┼─────────────────────────────────────┘
                            │ (when online)
                            ▼
              ┌─────────────────────────────┐
              │      Backend API            │
              │   PostgreSQL + Sync Service │
              └─────────────────────────────┘
```

## Sync Configuration

### Upload Queue

```dart
class PowerSyncBackendConnector extends PowerSyncBackendConnector {
  final PowerSyncDatabase db;
  final ApiClient api;

  @override
  Future<void> uploadData(PowerSyncDatabase database) async {
    // Get pending transactions
    final tx = await database.getNextCrudTransaction();

    if (tx == null) return;

    try {
      // Convert to backend format
      final operations = tx.crud.map((op) => {
        switch (op.op) {
          case CrudEntry.create:
            return {'type': 'INSERT', 'table': op.table, 'data': op.opData};
          case CrudEntry.update:
            return {'type': 'UPDATE', 'table': op.table, 'id': op.id, 'data': op.opData};
          case CrudEntry.delete:
            return {'type': 'DELETE', 'table': op.table, 'id': op.id};
        }
      }).toList();

      // Send to backend
      await api.bulkWrite('/sync/upload', {'operations': operations});

      // Mark as complete
      await database.completeCrudTransaction(tx);
    } catch (e) {
      // Will retry on next sync
      rethrow;
    }
  }
}
```

### Download Stream

```dart
@riverpod
Stream<void> syncStream(SyncStreamRef ref) async* {
  final db = ref.watch(powerSyncDatabaseProvider);
  final api = ref.watch(apiClientProvider);

  // Watch for remote changes
  await for (final update in api.watchSyncUpdates()) {
    // Apply to local database
    await db.writeTransaction((tx) async {
      for (final op in update.operations) {
        if (op.type == 'INSERT') {
          await tx.execute(
            'INSERT OR REPLACE INTO ${op.table} ...',
            op.data,
          );
        } else if (op.type == 'DELETE') {
          await tx.execute(
            'DELETE FROM ${op.table} WHERE id = ?',
            [op.id],
          );
        }
      }
    });

    yield;
  }
}
```

## Conflict Resolution

### Last-Write-Wins

```dart
// In migration/schema
CREATE TABLE patients (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  _ops TEXT,  -- PowerSync internal
  _deleted_at TEXT
);

// Conflict resolver on backend
func resolveConflict(local, remote map[string]any) map[string]any {
    localTime, _ := time.Parse(time.RFC3339, local["updated_at"].(string))
    remoteTime, _ := time.Parse(time.RFC3339, remote["updated_at"].(string))

    if localTime.After(remoteTime) {
        return local
    }
    return remote
}
```

### Custom Merge Strategy

```dart
// For complex merges (e.g., array fields)
Map<String, dynamic> mergePatientData(
  Map<String, dynamic> local,
  Map<String, dynamic> remote,
) {
  // Merge medical history arrays
  final localHistory = List.from(local['medical_history'] ?? []);
  final remoteHistory = List.from(remote['medical_history'] ?? []);

  // Combine and deduplicate by ID
  final mergedHistory = <String, dynamic>{};
  for (final item in [...remoteHistory, ...localHistory]) {
    mergedHistory[item['id']] = item;
  }

  return {
    ...remote,
    ...local,
    'medical_history': mergedHistory.values.toList(),
    'updated_at': DateTime.now().toIso8601String(),
  };
}
```

## Optimistic Updates

```dart
@riverpod
class AppointmentNotifier extends _$AppointmentNotifier {
  @override
  List<Appointment> build() {
    // Watch local database for appointments
    final db = ref.watch(powerSyncDatabaseProvider);
    return db.watch(
      'SELECT * FROM appointments WHERE status != ? ORDER BY scheduled_at',
      parameters: ['cancelled'],
      mapper: Appointment.fromRow,
    );
  }

  Future<void> createAppointment(Appointment appointment) async {
    final db = ref.read(powerSyncDatabaseProvider);

    // Optimistic update - immediately visible in UI
    await db.execute(
      '''INSERT INTO appointments (id, company_id, patient_id, scheduled_at, status, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)''',
      [appointment.id, appointment.companyId, appointment.patientId,
       appointment.scheduledAt.toIso8601String(), appointment.status,
       DateTime.now().toIso8601String(), DateTime.now().toIso8601String()],
    );

    // Sync will happen automatically when online
  }

  Future<void> cancelAppointment(String id) async {
    final db = ref.read(powerSyncDatabaseProvider);

    // Optimistic update
    await db.execute(
      'UPDATE appointments SET status = ?, updated_at = ? WHERE id = ?',
      ['cancelled', DateTime.now().toIso8601String(), id],
    );
  }
}
```

## Sync Status UI

```dart
class SyncStatusIndicator extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final syncStatus = ref.watch(syncStatusProvider);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: _getStatusColor(syncStatus),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _getStatusIcon(syncStatus),
            size: 14,
            color: Colors.white,
          ),
          const SizedBox(width: 4),
          Text(
            _getStatusText(syncStatus),
            style: const TextStyle(
              color: Colors.white,
              fontSize: 12,
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
        return Icons.error;
    }
  }

  String _getStatusText(SyncStatus status) {
    final l10n = AppLocalizations.of(context)!;
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
```

## Pending Changes Indicator

```dart
@riverpod
int pendingChangesCount(PendingChangesCountRef ref) {
  final db = ref.watch(powerSyncDatabaseProvider);

  // Count unacknowledged CRUD entries
  return db.get(
    'SELECT COUNT(*) as count FROM ps_crud WHERE ack = 0',
  ).then((row) => row?['count'] ?? 0);
}

class PendingChangesBadge extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final pendingCount = ref.watch(pendingChangesCountProvider);

    if (pendingCount == 0) return const SizedBox.shrink();

    return Badge(
      label: Text('$pendingCount'),
      child: const Icon(Icons.sync),
    );
  }
}
```

## Error Recovery

```dart
@riverpod
class SyncRecoveryNotifier extends _$SyncRecoveryNotifier {
  @override
  AsyncValue<void> build() => const AsyncValue.data(null);

  Future<void> retryFailedSync() async {
    state = const AsyncValue.loading();

    state = await AsyncValue.guard(() async {
      final db = ref.read(powerSyncDatabaseProvider);

      // Clear failed transactions and retry
      await db.execute('DELETE FROM ps_crud WHERE error IS NOT NULL');

      // Trigger immediate sync
      final connector = ref.read(backendConnectorProvider);
      await connector.uploadData(db);
    });
  }

  Future<void> forceFullResync() async {
    state = const AsyncValue.loading();

    state = await AsyncValue.guard(() async {
      final db = ref.read(powerSyncDatabaseProvider);

      // Disconnect and clear local data
      await db.disconnect();
      await db.clear();

      // Reconnect and resync
      final connector = ref.read(backendConnectorProvider);
      await db.connect(connector);
    });
  }
}
```

## Testing Offline Behavior

```dart
void main() {
  testWidgets('App works offline', (tester) async {
    // Set up mock offline database
    final container = ProviderContainer(overrides: [
      powerSyncDatabaseProvider.overrideWithValue(mockOfflineDb),
    ]);

    await tester.pumpWidget(
      UncontrolledProviderScope(
        container: container,
        child: const MediSyncApp(),
      ),
    );

    // Verify sync status shows offline
    expect(find.text('Offline'), findsOneWidget);

    // Verify can still create data
    await tester.tap(find.byIcon(Icons.add));
    await tester.pumpAndSettle();

    // Fill form
    await tester.enterText(find.byKey(Key('patient_name')), 'John Doe');
    await tester.tap(find.text('Save'));
    await tester.pumpAndSettle();

    // Verify data appears locally
    expect(find.text('John Doe'), findsOneWidget);

    // Verify pending indicator
    expect(find.byType(PendingChangesBadge), findsOneWidget);
  });
}
```
