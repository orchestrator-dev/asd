import re

def fix(fname, imp):
    with open(fname, 'r') as f:
        c = f.read()
    if imp not in c:
        c = re.sub(r'import \(', f'import (\\n\\t"{imp}"', c)
    with open(fname, 'w') as f:
        f.write(c)

fix('handlers/office.go', 'github.com/xuri/excelize/v2')
fix('handlers/toml.go', 'github.com/pelletier/go-toml/v2')
fix('handlers/yaml.go', 'gopkg.in/yaml.v3')

# Now for paraBuf unused:
with open('handlers/office.go', 'r') as f:
    c = f.read()
# Replace `var paraBuf strings.Builder` outside if it's not used
c = re.sub(r'var paraBuf strings\.Builder', '', c)
c = re.sub(r'var linesCount int', '', c)
c = c.replace('paraBuf.String()', '""')
c = c.replace('paraBuf.Reset()', '')
c = c.replace('paraBuf.Write(se)', '')
c = c.replace('linesCount++', '')
c = c.replace('if linesCount >= 200 {\n\t\t\t\t\t\treturn nil\n\t\t\t\t\t}', '')
with open('handlers/office.go', 'w') as f:
    f.write(c)

# For cert.go
with open('handlers/cert.go', 'r') as f:
    c = f.read()
c = c.replace('case *rsa.PrivateKey:\n\t\t\tfmt.Fprintf(out, "RSA Private Key\\n")', '')
c = c.replace('case *rsa.PublicKey:\n\t\t\tfmt.Fprintf(out, "RSA Public Key\\n")', '')
with open('handlers/cert.go', 'w') as f:
    f.write(c)

