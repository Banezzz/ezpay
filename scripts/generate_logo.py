#!/usr/bin/env python3
"""Generate EZPay logo assets."""

from pathlib import Path

from PIL import Image, ImageDraw


ROOT = Path(__file__).resolve().parents[1]
WWW = ROOT / "src" / "www"

CANVAS = 512
BG_TOP = (5, 17, 19, 255)
BG_BOTTOM = (14, 32, 38, 255)
BORDER = (54, 211, 153, 72)
WHITE = (248, 250, 252, 255)
MINT = (52, 241, 145, 255)
TEAL = (31, 200, 184, 255)
SHADOW = (0, 0, 0, 76)


SVG_LOGO = """<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" role="img" aria-labelledby="title desc">
  <title id="title">EZPay logo</title>
  <desc id="desc">A dark rounded square with an abstract EZ monogram in white and mint green.</desc>
  <defs>
    <linearGradient id="bg" x1="80" y1="36" x2="432" y2="476" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="#051113"/>
      <stop offset="1" stop-color="#0e2026"/>
    </linearGradient>
    <linearGradient id="mint" x1="206" y1="140" x2="398" y2="366" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="#34f191"/>
      <stop offset="1" stop-color="#1fc8b8"/>
    </linearGradient>
    <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
      <feDropShadow dx="0" dy="10" stdDeviation="10" flood-color="#000000" flood-opacity=".3"/>
    </filter>
  </defs>
  <rect x="20" y="20" width="472" height="472" rx="96" fill="url(#bg)"/>
  <rect x="23" y="23" width="466" height="466" rx="93" fill="none" stroke="#36d399" stroke-opacity=".28" stroke-width="6"/>
  <g filter="url(#shadow)">
    <rect x="116" y="148" width="44" height="216" rx="14" fill="#f8fafc"/>
    <rect x="116" y="148" width="154" height="40" rx="14" fill="#f8fafc"/>
    <rect x="116" y="236" width="128" height="40" rx="14" fill="#f8fafc"/>
    <rect x="116" y="324" width="144" height="40" rx="14" fill="#f8fafc"/>
    <path d="M251 148h141v40H349L210 364h-48l139-176h-50z" fill="url(#mint)"/>
    <rect x="204" y="324" width="188" height="40" rx="14" fill="url(#mint)"/>
  </g>
</svg>
"""


def lerp(a: int, b: int, t: float) -> int:
    return round(a + (b - a) * t)


def scaled_box(box, scale):
    return tuple(round(value * scale) for value in box)


def scaled_points(points, scale):
    return [(round(x * scale), round(y * scale)) for x, y in points]


def rounded_rectangle(draw, box, radius, fill):
    draw.rounded_rectangle(box, radius=radius, fill=fill)


def create_logo(size=512):
    scale = size / CANVAS
    supersample = 4
    work_size = size * supersample
    draw_scale = scale * supersample

    bg = Image.new("RGBA", (work_size, work_size), (0, 0, 0, 0))
    gradient = Image.new("RGBA", (work_size, work_size), BG_TOP)
    grad_pixels = gradient.load()
    for y in range(work_size):
        t = y / max(1, work_size - 1)
        row = (
            lerp(BG_TOP[0], BG_BOTTOM[0], t),
            lerp(BG_TOP[1], BG_BOTTOM[1], t),
            lerp(BG_TOP[2], BG_BOTTOM[2], t),
            255,
        )
        for x in range(work_size):
            grad_pixels[x, y] = row

    mask = Image.new("L", (work_size, work_size), 0)
    mask_draw = ImageDraw.Draw(mask)
    rounded_rectangle(
        mask_draw,
        scaled_box((20, 20, 492, 492), draw_scale),
        round(96 * draw_scale),
        255,
    )
    bg.paste(gradient, (0, 0), mask)
    draw = ImageDraw.Draw(bg)

    draw.rounded_rectangle(
        scaled_box((23, 23, 489, 489), draw_scale),
        radius=round(93 * draw_scale),
        outline=BORDER,
        width=max(1, round(6 * draw_scale)),
    )

    shadow_offset = round(10 * draw_scale)
    shadow_shapes(draw, draw_scale, shadow_offset)
    monogram(draw, draw_scale)

    return bg.resize((size, size), Image.Resampling.LANCZOS)


def shadow_shapes(draw, scale, offset):
    shadow_rects = [
        (116, 148 + offset / scale, 160, 364 + offset / scale, 14),
        (116, 148 + offset / scale, 270, 188 + offset / scale, 14),
        (116, 236 + offset / scale, 244, 276 + offset / scale, 14),
        (116, 324 + offset / scale, 260, 364 + offset / scale, 14),
        (204, 324 + offset / scale, 392, 364 + offset / scale, 14),
    ]
    for x1, y1, x2, y2, radius in shadow_rects:
        draw.rounded_rectangle(
            scaled_box((x1, y1, x2, y2), scale),
            radius=round(radius * scale),
            fill=SHADOW,
        )
    draw.polygon(
        scaled_points(
            [
                (251, 148 + offset / scale),
                (392, 148 + offset / scale),
                (392, 188 + offset / scale),
                (349, 188 + offset / scale),
                (210, 364 + offset / scale),
                (162, 364 + offset / scale),
                (301, 188 + offset / scale),
                (251, 188 + offset / scale),
            ],
            scale,
        ),
        fill=SHADOW,
    )


def monogram(draw, scale):
    for box in [
        (116, 148, 160, 364),
        (116, 148, 270, 188),
        (116, 236, 244, 276),
        (116, 324, 260, 364),
    ]:
        draw.rounded_rectangle(
            scaled_box(box, scale),
            radius=round(14 * scale),
            fill=WHITE,
        )

    draw.polygon(
        scaled_points(
            [
                (251, 148),
                (392, 148),
                (392, 188),
                (349, 188),
                (210, 364),
                (162, 364),
                (301, 188),
                (251, 188),
            ],
            scale,
        ),
        fill=MINT,
    )
    draw.rounded_rectangle(
        scaled_box((204, 324, 392, 364), scale),
        radius=round(14 * scale),
        fill=TEAL,
    )


def save_logo(path, size):
    path.parent.mkdir(parents=True, exist_ok=True)
    create_logo(size).save(path)


def main():
    (WWW / "images" / "logo.svg").write_text(SVG_LOGO, encoding="utf-8")

    for size in (192, 512):
        save_logo(WWW / f"pwa-{size}x{size}.png", size)
        save_logo(WWW / f"pwa-maskable-{size}x{size}.png", size)

    save_logo(WWW / "images" / "logo.png", 512)
    save_logo(WWW / "apple-touch-icon.png", 180)
    save_logo(WWW / "favicon-16x16.png", 16)
    save_logo(WWW / "favicon-32x32.png", 32)

    icons = [create_logo(size).convert("RGBA") for size in (16, 32, 48)]
    icons[-1].save(WWW / "favicon.ico", format="ICO", sizes=[(16, 16), (32, 32), (48, 48)])
    print("Logo assets generated.")


if __name__ == "__main__":
    main()
