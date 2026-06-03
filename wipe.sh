sed -i 's|render.NewWriter(w, opts)|render.NewWriter(w, opts.Flat)|g' handlers/*.go
sed -i 's|"github.com/pelletier/go-toml/v2"|toml "github.com/pelletier/go-toml/v2"|g' handlers/toml.go
sed -i 's|"gopkg.in/yaml.v3"|yaml "gopkg.in/yaml.v3"|g' handlers/yaml.go
sed -i '/case \*ecdsa\.PrivateKey:/d' handlers/cert.go
sed -i '/case \*ecdsa\.PublicKey:/d' handlers/cert.go
sed -i '/var paraBuf strings\.Builder/d' handlers/office.go
sed -i '/var linesCount int/d' handlers/office.go
