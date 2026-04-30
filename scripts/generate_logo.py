#!/usr/bin/env python3
"""Generate EZPay logo"""

from PIL import Image, ImageDraw, ImageFont

def create_logo(size=512):
    """Create EZPay logo with gradient background and text"""
    # Create image with gradient background
    img = Image.new('RGB', (size, size), '#000000')
    draw = ImageDraw.Draw(img)

    # Gradient effect (simplified as circles)
    for i in range(size, 0, -2):
        alpha = int(255 * (1 - i / size))
        color = (0, max(0, 100 - alpha // 3), max(0, 150 - alpha // 2))
        draw.ellipse([(size//2 - i//2, size//2 - i//2),
                      (size//2 + i//2, size//2 + i//2)],
                     fill=color)

    # Try to use a nice font, fall back to default
    try:
        font_large = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", size // 4)
        font_small = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", size // 8)
    except:
        font_large = ImageFont.load_default()
        font_small = ImageFont.load_default()

    # Draw "EZ" text
    text = "EZ"
    bbox = draw.textbbox((0, 0), text, font=font_large)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    x = (size - text_width) // 2
    y = (size - text_height) // 2 - size // 16
    draw.text((x, y), text, fill='white', font=font_large)

    # Draw "Pay" text
    text2 = "Pay"
    bbox2 = draw.textbbox((0, 0), text2, font=font_small)
    text_width2 = bbox2[2] - bbox2[0]
    x2 = (size - text_width2) // 2
    y2 = y + text_height + size // 32
    draw.text((x2, y2), text2, fill='#4ade80', font=font_small)

    return img

# Generate multiple sizes
sizes = [192, 512]
for size in sizes:
    logo = create_logo(size)
    logo.save(f'/home/dev/epusdt/src/www/pwa-{size}x{size}.png')
    if size == 512:
        logo.save('/home/dev/epusdt/src/www/pwa-maskable-512x512.png')
        logo.save('/home/dev/epusdt/src/www/apple-touch-icon.png')
    else:
        logo.save('/home/dev/epusdt/src/www/pwa-maskable-192x192.png')

# Generate main logo
logo = create_logo(512)
logo.save('/home/dev/epusdt/src/www/images/logo.png')

# Generate favicon sizes
logo_16 = create_logo(16)
logo_16.save('/home/dev/epusdt/src/www/favicon-16x16.png')

logo_32 = create_logo(32)
logo_32.save('/home/dev/epusdt/src/www/favicon-32x32.png')

# Generate ICO
logo_48 = create_logo(48)
logo_48.save('/home/dev/epusdt/src/www/favicon.ico', format='ICO')

print("Logo generated successfully!")
