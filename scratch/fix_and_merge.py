import re

print("Starting clean fix and merge of htmlPlayground...")

# 1. Read and clean playground_updated.html
with open("playground_updated.html", "r", encoding="utf-8") as f:
    html_content = f.read()

# Strip leading/trailing backticks and whitespace
html_content = html_content.strip()
if html_content.startswith("`"):
    html_content = html_content[1:]
if html_content.endswith("`"):
    html_content = html_content[:-1]
html_content = html_content.strip()

# Scan for any remaining backticks in the HTML which would break Go raw strings
backtick_count = html_content.count("`")
if backtick_count > 0:
    print(f"Warning: Found {backtick_count} backticks inside the HTML content. Replacing them with single quotes to avoid Go compilation errors...")
    html_content = html_content.replace("`", "'")

# 2. Read main.go
with open("cmd/api/main.go", "r", encoding="utf-8") as f:
    main_content = f.read()

# Find the start of the htmlPlayground definition in main.go
match = re.search(r'const\s+htmlPlayground\s*=\s*`*', main_content)
if match:
    start_idx = match.start()
    
    # Construct clean main.go content
    # We slice main.go up to const htmlPlayground = and append the new raw string literal
    new_main_content = main_content[:start_idx] + "const htmlPlayground = `" + html_content + "`\n"
    
    with open("cmd/api/main.go", "w", encoding="utf-8") as f:
        f.write(new_main_content)
    print("SUCCESS: Clean htmlPlayground successfully merged back into cmd/api/main.go!")
else:
    print("ERROR: Could not find const htmlPlayground in main.go")
