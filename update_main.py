import os

html_file = 'playground.html'
main_file = r'cmd\api\main.go'

with open(html_file, 'r', encoding='utf-8') as f:
    html_content = f.read()

with open(main_file, 'r', encoding='utf-8') as f:
    main_content = f.read()

start_marker = "const htmlPlayground = `"
# We assume the end marker is simply the backtick at the very end of the string
start_idx = main_content.find(start_marker)

if start_idx != -1:
    new_content = main_content[:start_idx + len(start_marker)] + html_content + "`\n"
    with open(main_file, 'w', encoding='utf-8') as f:
        f.write(new_content)
    print("Updated main.go")
else:
    print("Could not find start marker")
