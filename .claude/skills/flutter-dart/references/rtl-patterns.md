# RTL (Right-to-Left) Layout Patterns for Flutter

## Understanding RTL in Flutter

Flutter provides built-in RTL support through `Directionality` widget and logical properties.

### Core Concepts

```dart
// Text direction
enum TextDirection { rtl, ltr }

// Directionality widget
Directionality(
  textDirection: TextDirection.rtl,
  child: Text('مرحبا'),  // Arabic text
)

// Check current direction
bool isRTL = Directionality.of(context) == TextDirection.rtl;
```

## Layout Patterns

### Directional Margins and Padding

```dart
// Use directional geometry
Padding(
  padding: EdgeInsetsDirectional.only(
    start: 16.0,  // Left in LTR, Right in RTL
    end: 8.0,     // Right in LTR, Left in RTL
    top: 8.0,
    bottom: 8.0,
  ),
  child: Text('Hello'),
)

// Alternative: Builder pattern
Container(
  margin: EdgeInsetsDirectional.fromSTEB(16, 8, 8, 8),  // start, top, end, bottom
  child: Text('Hello'),
)
```

### Directional Alignment

```dart
// Use AlignmentDirectional
Align(
  alignment: AlignmentDirectional.centerStart,  // Left in LTR, Right in RTL
  child: Text('Hello'),
)

// Other options:
// - AlignmentDirectional.centerEnd
// - AlignmentDirectional.topStart
// - AlignmentDirectional.topEnd
// - AlignmentDirectional.bottomStart
// - AlignmentDirectional.bottomEnd
```

### Row and Column with RTL

```dart
// Row respects direction automatically
Directionality(
  textDirection: TextDirection.rtl,
  child: Row(
    children: [
      Icon(Icons.arrow_back),  // Appears on right in RTL
      Text('Back'),            // Next to icon
    ],
  ),
)

// For reversed order in RTL, use textDirection on Row
Row(
  textDirection: isRTL ? TextDirection.ltr : TextDirection.rtl,
  children: [...],
)
```

## Icon Handling

### Icons That Should Mirror

```dart
// Navigation arrows, back/forward, etc.
class DirectionalIcon extends StatelessWidget {
  final IconData icon;
  final double? size;
  final Color? color;

  const DirectionalIcon({
    super.key,
    required this.icon,
    this.size,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    final isRTL = Directionality.of(context) == TextDirection.rtl;
    final shouldMirror = _shouldMirrorIcon(icon);

    return Transform(
      transform: Matrix4.identity()
        ..scale(shouldMirror && isRTL ? -1.0 : 1.0, 1.0),
      alignment: Alignment.center,
      child: Icon(icon, size: size, color: color),
    );
  }

  bool _shouldMirrorIcon(IconData icon) {
    // Icons that should flip in RTL
    const mirrorIcons = {
      Icons.arrow_back,
      Icons.arrow_forward,
      Icons.arrow_left,
      Icons.arrow_right,
      Icons.chevron_left,
      Icons.chevron_right,
      Icons.first_page,
      Icons.last_page,
      Icons.keyboard_arrow_left,
      Icons.keyboard_arrow_right,
      Icons.navigate_before,
      Icons.navigate_next,
      Icons.reply,
      Icons.forward,
      Icons.undo,
      Icons.redo,
    };
    return mirrorIcons.contains(icon);
  }
}
```

### Icons That Should NOT Mirror

```dart
// Icons that should stay the same regardless of direction
const noMirrorIcons = {
  Icons.language,       // Globe
  Icons.settings,       // Settings
  Icons.help,           // Help
  Icons.search,         // Search
  Icons.refresh,        // Refresh
  Icons.phone,          // Phone
  Icons.email,          // Email
  Icons.camera,         // Camera
  Icons.access_time,    // Clock
  Icons.schedule,       // Schedule
  Icons.timer,          // Timer
  Icons.volume_up,      // Volume
  Icons.volume_down,
  Icons.play_arrow,     // Media controls
  Icons.pause,
  Icons.stop,
};
```

## Text Patterns

### Bidi Text Handling

```dart
// Mixed LTR and RTL text
RichText(
  text: TextSpan(
    children: [
      TextSpan(
        text: 'Hello ',
        style: TextStyle(locale: Locale('en')),
      ),
      TextSpan(
        text: 'مرحبا',
        style: TextStyle(locale: Locale('ar')),
      ),
      TextSpan(
        text: ' World',
        style: TextStyle(locale: Locale('en')),
      ),
    ],
  ),
)
```

### Text Alignment

```dart
// Use textAlign with direction awareness
Text(
  'مرحبا بالعالم',
  textAlign: TextAlign.start,  // Start respects direction
  textDirection: TextDirection.rtl,
)

// For numbers in RTL text
Directionality(
  textDirection: TextDirection.rtl,
  child: Text(
    'المريض رقم 123',
    // Numbers will display correctly within RTL context
  ),
)
```

## List and Grid Patterns

### RTL-Aware Lists

```dart
ListView.builder(
  itemCount: items.length,
  itemBuilder: (context, index) {
    return ListTile(
      leading: DirectionalIcon(icon: Icons.person),
      title: Text(items[index].name),
      trailing: DirectionalIcon(icon: Icons.chevron_right),
    );
  },
)
```

### Bidirectional Scroll

```dart
// GridView with RTL support
GridView.builder(
  scrollDirection: Axis.horizontal,
  reverse: isRTL,  // Start from right in RTL
  gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
    crossAxisCount: 2,
    childAspectRatio: 1.0,
  ),
  itemBuilder: (context, index) => ItemWidget(items[index]),
)
```

## Form Patterns

### RTL Form Fields

```dart
Column(
  children: [
    // Text field with RTL hint
    TextFormField(
      textAlign: TextAlign.start,
      textDirection: TextDirection.rtl,
      decoration: InputDecoration(
        hintText: 'أدخل اسم المريض',
        hintStyle: TextStyle(
          textDirection: TextDirection.rtl,
        ),
        prefixIcon: Icon(Icons.person),  // Right side in RTL
        suffixIcon: Icon(Icons.clear),   // Left side in RTL
      ),
    ),

    // Dropdown
    DropdownButtonFormField<String>(
      decoration: InputDecoration(
        labelText: 'نوع الموعد',
      ),
      items: ['عادي', 'عاجل', 'متابعة'].map((type) {
        return DropdownMenuItem(
          value: type,
          child: Text(type, textDirection: TextDirection.rtl),
        );
      }).toList(),
      onChanged: (value) {},
    ),
  ],
)
```

### Date/Time Pickers

```dart
// DatePicker respects locale automatically
showDatePicker(
  context: context,
  locale: Locale('ar'),  // Arabic locale
  initialDate: DateTime.now(),
  firstDate: DateTime(2020),
  lastDate: DateTime(2030),
);

// TimePicker
showTimePicker(
  context: context,
  locale: Locale('ar'),
  initialTime: TimeOfDay.now(),
);
```

## Navigation Patterns

### RTL-Aware Navigation Bar

```dart
BottomNavigationBar(
  items: [
    BottomNavigationBarItem(
      icon: DirectionalIcon(icon: Icons.home),
      label: AppLocalizations.of(context)!.navHome,
    ),
    BottomNavigationBarItem(
      icon: DirectionalIcon(icon: Icons.chat),
      label: AppLocalizations.of(context)!.navChat,
    ),
    BottomNavigationBarItem(
      icon: DirectionalIcon(icon: Icons.reports),
      label: AppLocalizations.of(context)!.navReports,
    ),
  ],
)
```

### App Bar with RTL

```dart
AppBar(
  // Title automatically centers in RTL
  title: Text('ميديسينك'),
  // leading is on right in RTL
  leading: IconButton(
    icon: DirectionalIcon(icon: Icons.menu),
    onPressed: () => scaffoldKey.currentState?.openDrawer(),
  ),
  // actions are on left in RTL
  actions: [
    IconButton(
      icon: Icon(Icons.search),
      onPressed: () {},
    ),
    IconButton(
      icon: Icon(Icons.notifications),
      onPressed: () {},
    ),
  ],
)
```

## Custom Painter with RTL

```dart
class DirectionalDividerPainter extends CustomPainter {
  final bool isRTL;

  DirectionalDividerPainter({required this.isRTL});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.grey
      ..strokeWidth = 1;

    final startX = isRTL ? size.width : 0;
    final endX = isRTL ? 0 : size.width;

    canvas.drawLine(
      Offset(startX, size.height / 2),
      Offset(endX, size.height / 2),
      paint,
    );
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
```

## Testing RTL

```dart
void main() {
  testWidgets('Widget renders correctly in RTL', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        localizationsDelegates: [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
        ],
        supportedLocales: [
          Locale('ar'),
        ],
        locale: Locale('ar'),
        home: Directionality(
          textDirection: TextDirection.rtl,
          child: Scaffold(
            body: PatientCard(patient: testPatient),
          ),
        ),
      ),
    );

    // Verify icon is mirrored
    final icon = tester.widget<Icon>(find.byType(Icon).first);
    // Check transform or visual position

    // Verify text alignment
    final text = tester.widget<Text>(find.byType(Text).first);
    expect(text.textAlign, equals(TextAlign.start));
  });
}
```
