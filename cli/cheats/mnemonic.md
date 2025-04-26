# Generate a default mnemonic (12 words, English)
cliu mnemonic generate

# Generate with a specific language
cliu mnemonic generate --lang=Spanish

# Generate with custom word count
cliu mnemonic generate --words=24

# Generate 3 mnemonics
cliu mnemonic generate --count 3

# Generate with a specific format
cliu mnemonic generate --count 3 --format=json

# Output as JSON
cliu mnemonic generate --format=json

# Output as CSV
cliu mnemonic generate --format=csv

# Custom template output
cliu mnemonic generate --format=template --template="MNEMONIC: {{.Sentence}}"
