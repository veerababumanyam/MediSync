# MediSync Logo Assets

This directory contains all generated logo sizes for web, mobile (iOS/Android), and app store use.

## Directory Structure

```
public/icons/
├── favicon.ico                      # Multi-size favicon (16, 32, 48)
├── favicon-16x16.png               # Browser tab icon
├── favicon-32x32.png               # Browser tab icon (Retina)
├── favicon-48x48.png               # Desktop shortcut
├── icon-*.png                      # PWA icons (72, 96, 128, 144, 152, 192, 384, 512)
├── ios-icon-*.png                  # iOS touch icons
├── android-icon-*.png              # Android launcher icons
├── app-store-icon-1024x1024.png    # App Store / Play Store icon
├── logo-*.png                      # General use logos (64, 128, 256, 512)
├── ios-iconset/                    # Xcode icon asset catalog
│   ├── Contents.json
│   ├── 60.png, 60@2x.png           # iPhone icons
│   ├── 76.png                      # iPad icons
│   ├── 83.5@2x.png                 # iPad Pro icons
│   └── 1024.png                    # App Store icon
├── android/                        # Android resource structure
│   ├── mipmap-mdpi/                # 36x36
│   ├── mipmap-hdpi/                # 48x48
│   ├── mipmap-xhdpi/               # 72x72
│   ├── mipmap-xxhdpi/              # 96x96
│   ├── mipmap-xxxhdpi/             # 144x144
│   ├── mipmap-anydpi-v26/          # Adaptive icon XML
│   └── values/colors.xml           # Icon background color
├── favicon-snippet.html            # HTML snippet to include
└── README.md                       # This file
```

## Usage

### Web / PWA

Include the following in your HTML `<head>`:

```html
<link rel="icon" type="image/png" sizes="16x16" href="/icons/favicon-16x16.png">
<link rel="icon" type="image/png" sizes="32x32" href="/icons/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="48x48" href="/icons/favicon-48x48.png">
<link rel="apple-touch-icon" sizes="180x180" href="/icons/ios-icon-60x60@3x.png">
<link rel="icon" type="image/png" sizes="192x192" href="/icons/icon-192x192.png">
<link rel="icon" type="image/png" sizes="512x512" href="/icons/icon-512x512.png">
<link rel="manifest" href="/manifest.json">
<meta name="theme-color" content="#0056D2">
```

Or copy the snippet from `favicon-snippet.html`.

### iOS (React Native / Flutter)

For iOS apps, copy the `ios-iconset/` directory to:
- **React Native**: `ios/Assets.xcassets/AppIcon.appiconset/`
- **Flutter**: `ios/Runner/Assets.xcassets/AppIcon.appiconset/`

### Android (Flutter)

For Android apps, copy the `android/mipmap-*` directories to:
- **Flutter**: `android/app/src/main/res/`

Copy the adaptive icon XML to:
- **Flutter**: `android/app/src/main/res/mipmap-anydpi-v26/ic_launcher.xml`

Copy the colors to:
- **Flutter**: `android/app/src/main/res/values/colors.xml`

### App Store / Play Store

Use `app-store-icon-1024x1024.png` for:
- Apple App Store (1024x1024 required)
- Google Play Store (512x512 required - will be downscaled)

## Icon Sizes Reference

| Use Case | Size | File |
|----------|------|------|
| Favicon (standard) | 16x16 | `favicon-16x16.png` |
| Favicon (Retina) | 32x32 | `favicon-32x32.png` |
| PWA Icon (small) | 192x192 | `icon-192x192.png` |
| PWA Icon (large) | 512x512 | `icon-512x512.png` |
| iOS Touch Icon | 180x180 | `ios-icon-60x60@3x.png` |
| Android Launcher (hdpi) | 48x48 | `android-icon-48x48.png` |
| Android Launcher (xxhdpi) | 96x96 | `android-icon-96x96.png` |
| App Store | 1024x1024 | `app-store-icon-1024x1024.png` |
| General UI (small) | 64x64 | `logo-64x64.png` |
| General UI (medium) | 256x256 | `logo-256x256.png` |

## Regenerating Icons

To regenerate all icons from the source logo:

```bash
python3 scripts/generate_logos.py
```

This requires:
- Python 3
- Pillow (PIL): `pip install Pillow`

## Brand Colors

The icons use the following brand colors:

- **Primary Blue**: `#0056D2` (Trust Blue)
- **Theme Color**: `#0056D2` (used for browser chrome)
- **Background**: `#FFFFFF` (white)

For dark mode or glassmorphism effects, use:
- **Midnight Navy**: `#0F172A`

## License

These logo assets are part of the MediSync project and follow the same license terms.
