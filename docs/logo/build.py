#!/usr/bin/env python3
"""Write the o family's logo set: assets/icon.svg, assets/logo-light.svg and
assets/logo-dark.svg.

The mark is the ekiben (concept fam-9 in ~/hub/o-plans/icons/build.py): a
real bento box whose compartments are not all the same shape or size. Rice
with an umeboshi takes the big compartment; a wide strip, a tall pickle slot,
a square of tamagoyaki, a round cup and a side dish take the rest. Each
compartment is one member's hue (ochazuke green, o-ai indigo, o-proj teal,
oden amber, matomeru violet, o-inventory vermilion), so no single tool owns
the family mark. The tile is a neutral ink, the walls and dividers cream.

The set follows the family's assets/: a 128x128 icon on a rounded tile (rx 28)
and a 112-high logo beside the name in the family's monoline lettering, drawn
as paths so no font is needed. Unlike the tools' logos, the logo keeps the
tile: the cream walls need the ink ground on a white page. Standard library
only, except --preview.

Regenerate from the repository root:

    python3 docs/logo/build.py                 # writes assets/*.svg
    python3 docs/logo/build.py --check         # the glyphs still redraw ochazuke's and matomeru's wordmarks
    python3 docs/logo/build.py --preview DIR   # also renders 16, 32 and 128 px PNGs, the
                                               # 1-colour icon, the dividers before and after,
                                               # and contact sheets (needs rsvg-convert, magick)
"""
import math
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[2]
ASSETS = ROOT / "assets"
NAME = "o family"
TITLE = "o-family"

TILE, CREAM, INK = "#18181d", "#f4efe6", "#18181d"

# The members' hues, as on their own tiles, and a light tint of each for the food hints.
HUE = {"ochazuke": "#3fa877", "o-ai": "#4a72f0", "o-proj": "#26b5c4", "oden": "#e8a33d",
       "matomeru": "#8a7cf5", "o-inventory": "#e5533d"}
TINT = {"ochazuke": "#a6e3c3", "o-ai": "#b6c8ff", "o-proj": "#a3e9f0", "oden": "#ffdc94",
        "matomeru": "#d3ccff", "o-inventory": "#ffb8a8"}

# Divider width, in icon units. The fam-9 concept used 5; 3.5 lets the hues
# take more of the box while the dividers still show as lines at 32 px.
DIVIDER = 3.5
R = 'stroke-linecap="round" stroke-linejoin="round"'


def mark(gap=DIVIDER, mono=False):
    """The ekiben in the 128x128 icon space. mono=True draws it in one colour
    (ink compartments and walls, the dividers and food hints as white gaps)."""
    fill = (lambda t: INK) if mono else (lambda t: HUE[t])
    tint = (lambda t: "#fff") if mono else (lambda t: TINT[t])
    ground = "#fff" if mono else CREAM
    wall = INK if mono else CREAM
    # Box: outer 18..110 x 26..102, wall 7 wide; inside 21.5..106.5 x 29.5..98.5.
    x0, x1, y0, y1 = 21.5 + gap, 106.5 - gap, 29.5 + gap, 98.5 - gap
    w, h = x1 - x0, y1 - y0
    rice_w = round(w * 0.5 - gap, 1)
    slot_w = 10.5
    right_x = x0 + rice_w + gap + slot_w + gap
    right_w = x1 - right_x
    strip_h = 12
    rice_h = h - strip_h - gap
    sq_h, side_h = 18, 20
    cup_y0, cup_y1 = y0 + sq_h + gap, y1 - side_h - gap
    cup_r = min(right_w, cup_y1 - cup_y0) / 2

    def cell(x, y, cw, ch, t, rx=6):
        return f'<rect x="{x:.1f}" y="{y:.1f}" width="{cw:.1f}" height="{ch:.1f}" rx="{rx}" fill="{fill(t)}"/>'

    rcx, rcy = x0 + rice_w / 2, y0 + rice_h / 2
    sesame = "".join(f"M{rcx+dx:.1f} {rcy+dy:.1f}h0.01" for dx, dy in ((-13, -14), (11, -11), (-10, 13), (11, 12), (1, -18)))
    slot_cx = x0 + rice_w + gap + slot_w / 2
    pick = "".join(f'<circle cx="{slot_cx:.1f}" cy="{y0 + h * k:.1f}" r="3" fill="{tint("o-proj")}"/>' for k in (0.16, 0.39, 0.62, 0.85))
    rcx2 = right_x + right_w / 2
    cup_cy = (cup_y0 + cup_y1) / 2
    side_y = y1 - side_h
    return (f'<rect x="18" y="26" width="92" height="76" rx="13" fill="{ground}"/>'
            f'<rect x="18" y="26" width="92" height="76" rx="13" fill="none" stroke="{wall}" stroke-width="7"/>'
            # rice and umeboshi
            + cell(x0, y0, rice_w, rice_h, "ochazuke", 7)
            + f'<circle cx="{rcx:.1f}" cy="{rcy:.1f}" r="6" fill="{"#fff" if mono else HUE["o-inventory"]}"/>'
            + f'<path d="{sesame}" stroke="{tint("ochazuke")}" stroke-width="3" {R}/>'
            # the wide shallow strip under the rice
            + cell(x0, y1 - strip_h, rice_w, strip_h, "o-ai", 5)
            + f'<path d="M{x0+7:.1f} {y1-strip_h/2:.1f}H{x0+rice_w*0.45:.1f}M{x0+rice_w*0.6:.1f} {y1-strip_h/2:.1f}H{x0+rice_w-7:.1f}" '
              f'stroke="{tint("o-ai")}" stroke-width="3" {R}/>'
            # the tall pickle slot
            + cell(x0 + rice_w + gap, y0, slot_w, h, "o-proj", 5) + pick
            # tamagoyaki, the round cup, the side
            + cell(right_x, y0, right_w, sq_h, "oden", 4)
            + f'<path d="M{rcx2-4.5:.1f} {y0+sq_h/2:.1f}H{rcx2+4.5:.1f}" stroke="{tint("oden")}" stroke-width="2.5" {R}/>'
            + f'<circle cx="{rcx2:.1f}" cy="{cup_cy:.1f}" r="{cup_r:.1f}" fill="{fill("matomeru")}"/>'
            + f'<circle cx="{rcx2:.1f}" cy="{cup_cy:.1f}" r="{cup_r*0.47:.1f}" fill="none" stroke="{tint("matomeru")}" stroke-width="2.5"/>'
            + cell(right_x, side_y, right_w, side_h, "o-inventory", 5)
            + f'<path d="M{rcx2-4.5:.1f} {side_y+5.5:.1f}H{rcx2+4.5:.1f}M{rcx2-4.5:.1f} {side_y+side_h-5.5:.1f}H{rcx2+4.5:.1f}" '
              f'stroke="{tint("o-inventory")}" stroke-width="3.5" {R}/>')


# --- The family lettering ---------------------------------------------------
# The glyph table is the one in oden's, o-proj's and o-inventory's
# docs/logo/build.py, whose shapes are copied from the ochazuke and matomeru
# logos. Added here, built the same way: f is t's stem turned over (the hook at
# the top, an ascender from y 18), l is h's ascender, and a space. Centre-line
# strokes 11 wide with round caps and joins, x-height y 44..84, ascenders from
# y 18, 15 between letters (centre line to centre line). Every family logo is
# 112 high (room for descenders); README width = 360 * viewBox width / 518,
# ochazuke's letter height.
STROKE, GAP, HEIGHT = 11, 15, 112


def f(v):
    s = f"{v:.1f}"
    return s[:-2] if s.endswith(".0") else s


GLYPHS = {
    "o": (40, lambda x: f"M{f(x+20)} 44a20 20 0 1 0 0.01 0"),
    "c": (34.1, lambda x: f"M{f(x+34.1)} 49.9A20 20 0 1 0 {f(x+34.1)} 78.1"),
    "h": (36, lambda x: f"M{f(x)} 18V84M{f(x)} 62a18 18 0 0 1 36 0V84"),
    "a": (40, lambda x: f"M{f(x+20)} 44a20 20 0 1 0 0.01 0M{f(x+40)} 44V84"),
    "d": (40, lambda x: f"M{f(x+20)} 44a20 20 0 1 0 0.01 0M{f(x+40)} 18V84"),
    "z": (34, lambda x: f"M{f(x)} 44h34L{f(x)} 84h34"),
    "u": (36, lambda x: f"M{f(x)} 44V66a18 18 0 0 0 36 0M{f(x+36)} 44V84"),
    "k": (32, lambda x: f"M{f(x)} 18V84M{f(x+31)} 44L{f(x)} 70M{f(x+11)} 61L{f(x+32)} 84"),
    "e": (40, lambda x: f"M{f(x)} 64H{f(x+40)}A20 20 0 1 0 {f(x+34.1)} 78.1"),
    "m": (52, lambda x: f"M{f(x)} 44V84M{f(x)} 57a13 13 0 0 1 26 0V84M{f(x+26)} 57a13 13 0 0 1 26 0V84"),
    "t": (23, lambda x: f"M{f(x+7)} 26V73a11 11 0 0 0 11 11h4M{f(x)} 44H{f(x+23)}"),
    "r": (24, lambda x: f"M{f(x)} 44V84M{f(x)} 64a20 20 0 0 1 20 -20h4"),
    "i": (0, lambda x: f"M{f(x)} 44V84M{f(x)} 24h0.01"),
    "n": (36, lambda x: f"M{f(x)} 44V84M{f(x)} 62a18 18 0 0 1 36 0V84"),
    "v": (34, lambda x: f"M{f(x)} 44L{f(x+17)} 84L{f(x+34)} 44"),
    "-": (16, lambda x: f"M{f(x)} 64h16"),
    "p": (40, lambda x: f"M{f(x)} 44V106M{f(x+20)} 44a20 20 0 1 0 0.01 0"),
    "j": (0, lambda x: f"M{f(x)} 44V95a11 11 0 0 1 -11 11h-4M{f(x)} 24h0.01"),
    "y": (36, lambda x: f"M{f(x)} 44V66a18 18 0 0 0 36 0M{f(x+36)} 44V95a11 11 0 0 1 -11 11h-4"),
    "f": (23, lambda x: f"M{f(x+7)} 84V29a11 11 0 0 1 11 -11h4M{f(x)} 44H{f(x+23)}"),
    "l": (0, lambda x: f"M{f(x)} 18V84"),
    " ": (12, lambda x: ""),
}


def wordmark(name, start=112):
    """Return (path data, x of each glyph, right edge of the last centre line)."""
    x, parts, xs = start, [], []
    for ch in name:
        w, d = GLYPHS[ch]
        parts.append(d(x))
        xs.append(x)
        x += w + GAP
    return "".join(parts), xs, x - GAP


def width_for(right):
    return int(right + STROKE / 2 + 3.5)  # ochazuke 509.1 -> 518, matomeru 524 -> 533


def word(d, ink):
    return (f'<path d="{d}" fill="none" stroke="{ink}" stroke-width="{STROKE}" '
            f'stroke-linecap="round" stroke-linejoin="round"/>\n')


def head(w, h, title):
    return (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {w} {h}" width="{w}" height="{h}" '
            f'role="img" aria-label="{title}">\n<title>{title}</title>\n')


def tile(body):
    return f'<rect width="128" height="128" rx="28" fill="{TILE}"/>\n' + body


def icon(gap=DIVIDER, mono=False):
    if mono:
        return head(128, 128, TITLE) + mark(gap, mono=True) + "\n</svg>\n"
    return head(128, 128, TITLE) + tile(mark(gap)) + "\n</svg>\n"


# The tiled icon is scaled to 90 high and set 3 from the top, in the 96-unit
# band the family's marks share; the name starts the family's 24 after it.
LOGO_MARK = "translate(3 3) scale(0.703125)"
MARK_RIGHT = 3 + 90


def logo(theme):
    ink = INK if theme == "light" else CREAM
    d, _, right = wordmark(NAME, MARK_RIGHT + 24 + STROKE / 2)
    return (head(width_for(right), HEIGHT, TITLE)
            + f'<g transform="{LOGO_MARK}">{tile(mark())}</g>\n' + word(d, ink) + "</svg>\n")


def check():
    """The glyph table must redraw ochazuke's and matomeru's wordmarks exactly."""
    import re
    ok = True
    for fam in ("ochazuke", "matomeru"):
        p = pathlib.Path.home() / f"proj/{fam}/assets/logo-light.svg"
        if not p.exists():
            print(f"skip {fam}: {p} not found")
            continue
        t = p.read_text()
        want = re.findall(r'<path d="(M[^"]*)" fill="none" stroke="#[0-9a-f]+" stroke-width="11"', t)[-1]
        want_w = int(re.search(r'viewBox="0 0 (\d+) 96"', t).group(1))
        d, _, right = wordmark(fam)
        split = lambda s: re.sub(r"([A-Za-z])", r" \1 ", s).split()
        a, b = split(d), split(want)
        near = lambda u, v: u == v or (u[0] not in "MmAaHhVvLlZz" and abs(float(u) - float(v)) < 0.15)
        same = len(a) == len(b) and all(map(near, a, b)) and width_for(right) == want_w
        ok &= same
        print(f"{fam}: {'same' if same else 'DIFFERENT'} wordmark ({width_for(right)} wide)")
    return ok


def preview(out):
    """Render the icon at 16, 32 and 128 px (and blown up 8x and 4x, to see
    the pixels), the 1-colour icon, the dividers before (5, the concept) and
    after (DIVIDER), and both logos, into out (needs rsvg-convert and magick)."""
    import subprocess
    out.mkdir(parents=True, exist_ok=True)
    bg = "#f6f3ee"

    def render(svg_text, name, size):
        src = out / f"{name}.svg"
        src.write_text(svg_text)
        p = out / f"{name}-{size}.png"
        subprocess.run(["rsvg-convert", "-w", str(size), "-h", str(size), "-o", p, src], check=True)
        return p

    def strip(name, svg_text):
        pngs = [render(svg_text, name, s) for s in (16, 32, 128)]
        big = []
        for size, k in ((16, 8), (32, 4)):
            b = out / f"{name}-{size}-x{k}.png"
            subprocess.run(["magick", pngs[[16, 32].index(size)], "-filter", "point", "-resize", f"{k*100}%", b], check=True)
            big.append(b)
        row = out / f"{name}-row.png"
        subprocess.run(["magick", "-background", bg, *map(str, pngs), *map(str, big), "-bordercolor", bg,
                        "-border", "12", "-gravity", "center", "+append", row], check=True)
        return row

    rows = [strip("icon", icon()), strip("mono", icon(mono=True)),
            strip("dividers-before", icon(gap=5)), strip("dividers-after", icon())]
    subprocess.run(["magick", *map(str, rows[:2]), "-append", out / "icons.png"], check=True)
    subprocess.run(["magick", *map(str, rows[2:]), "-append", out / "dividers.png"], check=True)
    for theme, b in (("light", "white"), ("dark", "#0d1117")):
        (out / f"logo-{theme}.svg").write_text(logo(theme))
        subprocess.run(["rsvg-convert", "-w", "720", "-b", b, "-o", out / f"logo-{theme}.png",
                        out / f"logo-{theme}.svg"], check=True)
    subprocess.run(["magick", out / "logo-light.png", out / "logo-dark.png", "-append", out / "logos.png"], check=True)
    for n in ("icons.png", "dividers.png", "logos.png"):
        print(out / n)


def main():
    if sys.argv[1:] == ["--check"]:
        sys.exit(0 if check() else 1)
    ASSETS.mkdir(exist_ok=True)
    files = {"icon.svg": icon(), "logo-light.svg": logo("light"), "logo-dark.svg": logo("dark")}
    for name, text in files.items():
        (ASSETS / name).write_text(text)
        print(ASSETS / name)
    if len(sys.argv) == 3 and sys.argv[1] == "--preview":
        preview(pathlib.Path(sys.argv[2]))


if __name__ == "__main__":
    main()
