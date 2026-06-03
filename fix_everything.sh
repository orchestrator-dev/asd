sed -i 's|case \*rsa\.PrivateKey:|case *ecdsa.PrivateKey:|2' handlers/cert.go
sed -i 's|case \*rsa\.PublicKey:|case *ecdsa.PublicKey:|2' handlers/cert.go
sed -i 's/var paraBuf strings\.Builder/paraBuf := strings.Builder{}/g' handlers/office.go
sed -i 's/var linesCount int/linesCount := 0/g' handlers/office.go
sed -i 's|var node yaml\.Node|var node interface{}|' handlers/yaml.go
