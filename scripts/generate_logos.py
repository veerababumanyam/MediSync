#!/usr/bin/env python3
"""
MediSync Logo Generator
Generates various sizes of the logo for web, mobile, and other platforms.
"""

from PIL import Image
import os
import json

# Configuration
SOURCE_LOGO = "public/logo.png"
OUTPUT_DIR = "public/icons"
BACKUP_ORIGINAL_NAME = "logo-1024x1024.png"

# Define all required sizes with their use cases
SIZES = {
    # Favicon sizes (web)
    "favicon-16x16.png": 16,
    "favicon-32x32.png": 32,
    "favicon-48x48.png": 48,

    # PWA / Web app icons
    "icon-72x72.png": 72,
    "icon-96x96.png": 96,
    "icon-128x128.png": 128,
    "icon-144x144.png": 144,
    "icon-152x152.png": 152,
    "icon-192x192.png": 192,
    "icon-384x384.png": 384,
    "icon-512x512.png": 512,

    # iOS icons
    "ios-icon-60x60.png": 60,
    "ios-icon-60x60@2x.png": 120,
    "ios-icon-60x60@3x.png": 180,
    "ios-icon-76x76.png": 76,
    "ios-icon-76x76@2x.png": 152,
    "ios-icon-83.5x83.5@2x.png": 167,

    # Android icons
    "android-icon-36x36.png": 36,
    "android-icon-48x48.png": 48,
    "android-icon-72x72.png": 72,
    "android-icon-96x96.png": 96,
    "android-icon-144x144.png": 144,
    "android-icon-192x192.png": 192,
    "android-icon-512x512.png": 512,

    # App store / marketing
    "app-store-icon-1024x1024.png": 1024,

    # Misc sizes
    "logo-64x64.png": 64,
    "logo-128x128.png": 128,
    "logo-256x256.png": 256,
    "logo-512x512.png": 512,
}

# iOS icon set structure for Xcode
IOS_ICON_SET_CONTENT = '''{
  "images" : [
    {
      "filename" : "60.png",
      "idiom" : "iphone",
      "scale" : "2x",
      "size" : "60x60"
    },
    {
      "filename" : "60@2x.png",
      "idiom" : "iphone",
      "scale" : "3x",
      "size" : "60x60"
    },
    {
      "filename" : "76.png",
      "idiom" : "ipad",
      "scale" : "2x",
      "size" : "76x76"
    },
    {
      "filename" : "83.5@2x.png",
      "idiom" : "ipad",
      "scale" : "2x",
      "size" : "83.5x83.5"
    },
    {
      "filename" : "1024.png",
      "idiom" : "ios-marketing",
      "scale" : "1x",
      "size" : "1024x1024"
    }
  ],
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}'''

# Android adaptive icon content
ANDROID_ADAPTIVE_ICON = '''<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="@mipmap/ic_launcher_foreground"/>
</adaptive-icon>'''

# Web app manifest content
WEB_MANIFEST_CONTENT = {
    "name": "MediSync",
    "short_name": "MediSync",
    "description": "AI-Powered Conversational BI & Intelligent Accounting for Healthcare",
    "icons": [
        {
            "src": "/icons/icon-72x72.png",
            "sizes": "72x72",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-96x96.png",
            "sizes": "96x96",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-128x128.png",
            "sizes": "128x128",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-144x144.png",
            "sizes": "144x144",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-152x152.png",
            "sizes": "152x152",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-192x192.png",
            "sizes": "192x192",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-384x384.png",
            "sizes": "384x384",
            "type": "image/png",
            "purpose": "any"
        },
        {
            "src": "/icons/icon-512x512.png",
            "sizes": "512x512",
            "type": "image/png",
            "purpose": "any"
        }
    ],
    "start_url": "/",
    "display": "standalone",
    "theme_color": "#0056D2",
    "background_color": "#FFFFFF"
}


def generate_logo_sizes(source_path, output_dir):
    """Generate all required logo sizes from the source image."""
    print(f"Loading source logo from: {source_path}")

    # Open the source image
    with Image.open(source_path) as img:
        # Convert to RGBA if needed
        if img.mode != 'RGBA':
            img = img.convert('RGBA')

        # Get original dimensions
        original_width, original_height = img.size
        print(f"Original size: {original_width}x{original_height}")

        # Create output directory if it doesn't exist
        os.makedirs(output_dir, exist_ok=True)

        # Generate each size
        for filename, size in SIZES.items():
            output_path = os.path.join(output_dir, filename)

            # Use high-quality resampling
            resized = img.resize((size, size), Image.Resampling.LANCZOS)
            resized.save(output_path, 'PNG')
            print(f"Generated: {filename} ({size}x{size})")

        # Create iOS icon set structure
        ios_dir = os.path.join(output_dir, "ios-iconset")
        os.makedirs(ios_dir, exist_ok=True)

        ios_mappings = {
            "60.png": 60,
            "60@2x.png": 120,
            "76.png": 76,
            "83.5@2x.png": 167,
            "1024.png": 1024
        }

        for filename, size in ios_mappings.items():
            output_path = os.path.join(ios_dir, filename)
            resized = img.resize((size, size), Image.Resampling.LANCZOS)
            resized.save(output_path, 'PNG')
            print(f"Generated iOS: {filename} ({size}x{size})")

        # Create Contents.json for iOS
        contents_path = os.path.join(ios_dir, "Contents.json")
        with open(contents_path, 'w') as f:
            f.write(IOS_ICON_SET_CONTENT)
        print(f"Generated: ios-iconset/Contents.json")

        # Create Android mipmap directories
        android_sizes = {
            "mipmap-mdpi": 36,
            "mipmap-hdpi": 48,
            "mipmap-xhdpi": 72,
            "mipmap-xxhdpi": 96,
            "mipmap-xxxhdpi": 144,
            "mipmap-xxxhdpi-512": 512
        }

        for dpi, size in android_sizes.items():
            android_dir = os.path.join(output_dir, "android", dpi)
            os.makedirs(android_dir, exist_ok=True)

            output_path = os.path.join(android_dir, "ic_launcher.png")
            resized = img.resize((size, size), Image.Resampling.LANCZOS)
            resized.save(output_path, 'PNG')

            # Also create foreground for adaptive icon
            fg_path = os.path.join(android_dir, "ic_launcher_foreground.png")
            resized.save(fg_path, 'PNG')

        # Create adaptive icon XML
        mipmap_anydpi = os.path.join(output_dir, "android", "mipmap-anydpi-v26")
        os.makedirs(mipmap_anydpi, exist_ok=True)

        adaptive_path = os.path.join(mipmap_anydpi, "ic_launcher.xml")
        with open(adaptive_path, 'w') as f:
            f.write(ANDROID_ADAPTIVE_ICON)
        print(f"Generated: Android adaptive icon XML")

        # Create colors.xml for Android
        colors_path = os.path.join(output_dir, "android", "values", "colors.xml")
        os.makedirs(os.path.dirname(colors_path), exist_ok=True)
        with open(colors_path, 'w') as f:
            f.write('<?xml version="1.0" encoding="utf-8"?>\n')
            f.write('<resources>\n')
            f.write('    <color name="ic_launcher_background">#0056D2</color>\n')
            f.write('</resources>\n')
        print(f"Generated: Android colors.xml")

        # Generate favicon.ico (multi-size ICO file)
        favicon_sizes = [16, 32, 48]
        favicon_images = []
        for size in favicon_sizes:
            resized = img.resize((size, size), Image.Resampling.LANCZOS)
            favicon_images.append(resized)

        favicon_path = os.path.join(output_dir, "favicon.ico")
        favicon_images[0].save(
            favicon_path,
            format='ICO',
            sizes=[(size, size) for size in favicon_sizes]
        )
        print(f"Generated: favicon.ico (contains {favicon_sizes})")

        # Create web manifest
        manifest_path = os.path.join(output_dir, "..", "manifest.json")
        with open(manifest_path, 'w') as f:
            json.dump(WEB_MANIFEST_CONTENT, f, indent=2)
        print(f"Generated: public/manifest.json")

        # Create favicon HTML snippet
        html_path = os.path.join(output_dir, "favicon-snippet.html")
        with open(html_path, 'w') as f:
            f.write('<!-- Favicon and app icons -->\n')
            f.write('<link rel="icon" type="image/png" sizes="16x16" href="/icons/favicon-16x16.png">\n')
            f.write('<link rel="icon" type="image/png" sizes="32x32" href="/icons/favicon-32x32.png">\n')
            f.write('<link rel="icon" type="image/png" sizes="48x48" href="/icons/favicon-48x48.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="60x60" href="/icons/ios-icon-60x60.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="60x60" href="/icons/ios-icon-60x60@2x.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="76x76" href="/icons/ios-icon-76x76.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="76x76" href="/icons/ios-icon-76x76@2x.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="152x152" href="/icons/ios-icon-152x152.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="167x167" href="/icons/ios-icon-83.5x83.5@2x.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="180x180" href="/icons/ios-icon-60x60@3x.png">\n')
            f.write('<link rel="apple-touch-icon" sizes="192x192" href="/icons/icon-192x192.png">\n')
            f.write('<link rel="icon" type="image/png" sizes="512x512" href="/icons/icon-512x512.png">\n')
            f.write('<link rel="manifest" href="/manifest.json">\n')
            f.write('<meta name="theme-color" content="#0056D2">\n')
            f.write('<meta name="apple-mobile-web-app-capable" content="yes">\n')
            f.write('<meta name="apple-mobile-web-app-status-bar-style" content="default">\n')
            f.write('<meta name="apple-mobile-web-app-title" content="MediSync">\n')
        print(f"Generated: favicon-snippet.html")

        print(f"\nSummary: Generated {len(SIZES) + len(ios_mappings) + len(android_sizes) + 3} files")
        print(f"Output directory: {output_dir}")


def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    source_path = os.path.join(project_root, SOURCE_LOGO)
    output_dir = os.path.join(project_root, OUTPUT_DIR)

    if not os.path.exists(source_path):
        print(f"Error: Source logo not found at {source_path}")
        return 1

    try:
        generate_logo_sizes(source_path, output_dir)
        print("\nAll logo sizes generated successfully!")
        return 0
    except Exception as e:
        print(f"Error generating logos: {e}")
        return 1


if __name__ == "__main__":
    exit(main())
