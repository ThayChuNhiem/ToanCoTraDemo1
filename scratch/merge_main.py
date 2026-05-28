import re

print("Merging playground_updated.html back into cmd/api/main.go...")

with open("cmd/api/main.go", "r", encoding="utf-8") as f:
    main_content = f.read()

with open("playground_updated.html", "r", encoding="utf-8") as f:
    html_content = f.read()

# Let's find the start of const htmlPlayground = `
match = re.search(r'const\s+htmlPlayground\s*=\s*`', main_content)
if match:
    start_idx = match.start()
    # Find the last backtick of the file
    end_idx = main_content.rfind('`')
    
    new_main_content = main_content[:start_idx] + "const htmlPlayground = `" + html_content + "`\n"
    
    with open("cmd/api/main.go", "w", encoding="utf-8") as f:
        f.write(new_main_content)
    print("SUCCESS: htmlPlayground successfully merged back into cmd/api/main.go!")
else:
    print("ERROR: Could not find const htmlPlayground in main.go")
