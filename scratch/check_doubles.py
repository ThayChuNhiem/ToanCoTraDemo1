with open('playground.html', 'r', encoding='utf-8') as f:
    js = f.read()

import re
matches = [m.start() for m in re.finditer(r"''", js)]
print(f'Found {len(matches)} occurrences of double single quotes:')
for idx, start in enumerate(matches):
    print(f'Match {idx}:')
    print(js[max(0, start-100):start+100])
