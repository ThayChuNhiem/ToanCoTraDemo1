import os

html_file = 'playground.html'
main_file = r'cmd\api\main.go'

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

# Remove backticks from html_content if they exist at the start/end
html_content = html_content.strip()
if html_content.startswith('`'):
    html_content = html_content[1:]
if html_content.endswith('`'):
    html_content = html_content[:-1]
# Strip again just in case there are newlines after backtick
html_content = html_content.strip()

with open(main_file, 'r', encoding='utf-8') as f:
    main_content = f.read()

start_marker = "const htmlPlayground = `"
start_idx = main_content.find(start_marker)

if start_idx != -1:
    new_content = main_content[:start_idx + len(start_marker)] + html_content + "`\n"
    with open(main_file, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print("Updated main.go successfully without double backticks!")
else:
    print("Could not find start marker")
