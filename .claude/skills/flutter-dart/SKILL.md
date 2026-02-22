---
name: flutter-dart
description: This skill should be used when the user asks to "create Flutter widgets", "build mobile features", "implement Dart patterns", "add Flutter screens", "Flutter state management", "offline sync with PowerSync", or mentions Flutter/Dart mobile development for MediSync.
---

# Flutter/Dart Mobile Development for MediSync

Flutter 3.42 powers the MediSync mobile app with iOS and Android support, PowerSync for offline-first architecture, and comprehensive i18n with ARB files.

★ Insight ─────────────────────────────────────
MediSync mobile architecture:
1. **Offline-first** - PowerSync local-first database
2. **RTL support** - Arabic layout mirroring
3. **State management** - Riverpod for reactive state
4. **Code generation** - freezed, json_serializable, riverpod_generator
─────────────────────────────────────────────────

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Flutter Version** | 3.42 |
| **Dart Version** | 3.x |
| **State Management** | Riverpod 2.x |
| **Offline DB** | PowerSync |
| **i18n** | ARB files + flutter_localizations |
| **HTTP Client** | Dio 5.x |

## Project Structure

```
mobile/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── core/
│   │   ├── constants/
│   │   ├── theme/
│   │   ├── router/
│   │   └── di/
│   ├── features/
│   │   ├── dashboard/
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   ├── chat/
│   │   ├── reports/
│   │   └── settings/
│   ├── shared/
│   │   ├── widgets/
│   │   ├── providers/
│   │   └── utils/
│   └── l10n/
│       ├── app_en.arb
│       └── app_ar.arb
├── pubspec.yaml
└── test/
```

## Widget Patterns

### Stateless Widget with i18n

```dart
class DashboardCard extends StatelessWidget {
  const DashboardCard({
    super.key,
    required this.title,
    required this.value,
    required this.trend,
  });

  final String title;
  final String value;
  final double trend;

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            Row(
              children: [
                Icon(
                  trend >= 0 ? Icons.trending_up : Icons.trending_down,
                  color: trend >= 0 ? Colors.green : Colors.red,
                  size: 16,
                ),
                Text(
                  '${trend.abs().toStringAsFixed(1)}%',
                  style: TextStyle(
                    color: trend >= 0 ? Colors.green : Colors.red,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
```

### Stateful Widget with Controllers

```dart
class ChatScreen extends StatefulWidget {
  const ChatScreen({super.key});

  @override
  State<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  final _controller = TextEditingController();
  final _scrollController = ScrollController();
  bool _isLoading = false;

  @override
  void dispose() {
    _controller.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(S.of(context).chatTitle)),
      body: Column(
        children: [
          Expanded(child: _buildMessageList()),
          _buildInputArea(),
        ],
      ),
    );
  }
}
```

## Riverpod State Management

### Provider Definition

```dart
// dashboard_provider.dart
@riverpod
class DashboardNotifier extends _$DashboardNotifier {
  @override
  DashboardState build() {
    return const DashboardState.loading();
  }

  Future<void> loadDashboard() async {
    state = const DashboardState.loading();

    try {
      final repository = ref.read(dashboardRepositoryProvider);
      final data = await repository.fetchDashboard();
      state = DashboardState.loaded(data);
    } catch (e, st) {
      state = DashboardState.error(e.toString());
    }
  }
}

@riverpod
DashboardData? dashboardData(DashboardDataRef ref) {
  final state = ref.watch(dashboardNotifierProvider);
  return state.maybeWhen(
    loaded: (data) => data,
    orElse: () => null,
  );
}
```

### Freezed State Classes

```dart
@freezed
class DashboardState with _$DashboardState {
  const factory DashboardState.loading() = _Loading;
  const factory DashboardState.loaded(DashboardData data) = _Loaded;
  const factory DashboardState.error(String message) = _Error;
}

@freezed
class DashboardData with _$DashboardData {
  const factory DashboardData({
    required double revenue,
    required int patientCount,
    required List<Appointment> appointments,
    required List<ChartPoint> trends,
  }) = _DashboardData;

  factory DashboardData.fromJson(Map<String, dynamic> json) =>
      _$DashboardDataFromJson(json);
}
```

### Consuming Providers

```dart
class DashboardView extends ConsumerWidget {
  const DashboardView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(dashboardNotifierProvider);
    final l10n = AppLocalizations.of(context)!;

    return state.when(
      loading: () => const Center(child: CircularProgressIndicator()),
      loaded: (data) => _buildDashboard(context, ref, data),
      error: (message) => Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(message),
            ElevatedButton(
              onPressed: () => ref.read(dashboardNotifierProvider.notifier).loadDashboard(),
              child: Text(l10n.retry),
            ),
          ],
        ),
      ),
    );
  }
}
```

## PowerSync Offline-First

### Database Schema

```dart
// lib/core/database/schema.dart
const kSchema = '''
  CREATE TABLE patients (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    name TEXT NOT NULL,
    patient_number TEXT,
    date_of_birth TEXT,
    gender TEXT,
    phone TEXT,
    email TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    _deleted_at TEXT
  );

  CREATE TABLE appointments (
    id TEXT PRIMARY KEY,
    company_id TEXT NOT NULL,
    patient_id TEXT NOT NULL,
    scheduled_at TEXT NOT NULL,
    status TEXT NOT NULL,
    notes TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );

  CREATE INDEX idx_patients_company ON patients(company_id);
  CREATE INDEX idx_appointments_patient ON appointments(patient_id);
''';
```

### PowerSync Setup

```dart
// lib/core/database/powersync.dart
@riverpod
PowerSyncDatabase powerSyncDatabase(PowerSyncDatabaseRef ref) {
  final db = PowerSyncDatabase(
    schema: const Schema(kSchema),
    path: ref.watch(databasePathProvider),
  );

  ref.onDispose(db.close);
  return db;
}

@riverpod
Future<void> initializeDatabase(InitializeDatabaseRef ref) async {
  final db = ref.watch(powerSyncDatabaseProvider);
  await db.initialize();

  // Start sync
  final syncClient = ref.watch(syncClientProvider);
  await syncClient.connect();
}
```

### Repository with Offline Support

```dart
@riverpod
PatientRepository patientRepository(PatientRepositoryRef ref) {
  return PatientRepository(
    db: ref.watch(powerSyncDatabaseProvider),
    api: ref.watch(apiClientProvider),
  );
}

class PatientRepository {
  PatientRepository({required this.db, required this.api});

  final PowerSyncDatabase db;
  final ApiClient api;

  Stream<List<Patient>> watchPatients(String companyId) {
    return db.watch(
      'SELECT * FROM patients WHERE company_id = ? AND _deleted_at IS NULL ORDER BY name',
      parameters: [companyId],
      mapper: Patient.fromRow,
    );
  }

  Future<void> createPatient(Patient patient) async {
    await db.execute(
      '''INSERT INTO patients (id, company_id, name, patient_number, date_of_birth, gender, phone, email, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)''',
      [patient.id, patient.companyId, patient.name, patient.patientNumber,
       patient.dateOfBirth?.toIso8601String(), patient.gender, patient.phone,
       patient.email, patient.createdAt.toIso8601String(),
       patient.updatedAt.toIso8601String()],
    );
  }
}
```

## RTL Support

### Directional Layout

```dart
class DirectionalLayout extends StatelessWidget {
  const DirectionalLayout({
    super.key,
    required this.child,
    this.textAlign = TextAlign.start,
  });

  final Widget child;
  final TextAlign textAlign;

  @override
  Widget build(BuildContext context) {
    final isRTL = Directionality.of(context) == TextDirection.rtl;

    return Directionality(
      textDirection: isRTL ? TextDirection.rtl : TextDirection.ltr,
      child: child,
    );
  }
}

// Use in widgets
class PatientCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return DirectionalLayout(
      child: Row(
        children: [
          // Automatically flips for RTL
          const Icon(Icons.person),
          const SizedBox(width: 8),
          Expanded(child: Text(patient.name)),
        ],
      ),
    );
  }
}
```

### Mirrored Icons and Assets

```dart
Widget buildMirroredIcon(IconData icon, {bool mirror = true}) {
  return Transform(
    transform: Matrix4.identity()
      ..scale(mirror && _isRTL ? -1.0 : 1.0, 1.0),
    alignment: Alignment.center,
    child: Icon(icon),
  );
}
```

## HTTP with Dio

### API Client

```dart
@riverpod
Dio dio(DioRef ref) {
  final dio = Dio(BaseOptions(
    baseUrl: EnvironmentConfig.apiUrl,
    connectTimeout: const Duration(seconds: 30),
    receiveTimeout: const Duration(seconds: 60),
  ));

  dio.interceptors.addAll([
    AuthInterceptor(ref),
    LoggingInterceptor(),
    RetryInterceptor(dio: dio),
  ]);

  return dio;
}

class AuthInterceptor extends Interceptor {
  AuthInterceptor(this.ref);

  final Ref ref;

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final token = ref.read(authTokenProvider);
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 401) {
      ref.read(authNotifierProvider.notifier).logout();
    }
    handler.next(err);
  }
}
```

## Code Generation Setup

### pubspec.yaml Dependencies

```yaml
dependencies:
  flutter:
    sdk: flutter
  flutter_localizations:
    sdk: flutter
  riverpod: ^2.4.0
  flutter_riverpod: ^2.4.0
  riverpod_annotation: ^2.3.0
  freezed_annotation: ^2.4.0
  json_annotation: ^4.8.0
  dio: ^5.4.0
  powersync: ^1.0.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  build_runner: ^2.4.0
  riverpod_generator: ^2.3.0
  freezed: ^2.4.0
  json_serializable: ^6.7.0
  flutter_lints: ^3.0.0
```

### Build Commands

```bash
# Generate code
flutter pub run build_runner build --delete-conflicting-outputs

# Watch for changes
flutter pub run build_runner watch --delete-conflicting-outputs
```

## Additional Resources

### Reference Files
- **`references/offline-patterns.md`** - PowerSync sync strategies
- **`references/rtl-patterns.md`** - Comprehensive RTL layout patterns

### Example Files
- **`examples/dashboard_screen.dart`** - Complete dashboard implementation
- **`examples/chat_screen.dart`** - Chat interface with streaming

### Scripts
- **`scripts/generate.sh`** - Run code generation
- **`scripts/l10n-validate.sh`** - Validate ARB files
