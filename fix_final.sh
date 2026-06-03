sed -i 's|"github.com/pelletier/go-toml/v2"|toml "github.com/pelletier/go-toml/v2"|' handlers/toml.go
sed -i 's|"gopkg.in/yaml.v3"|yaml "gopkg.in/yaml.v3"|' handlers/yaml.go
sed -i 's/case \*ecdsa\.PrivateKey:/\/\//' handlers/cert.go
sed -i 's/case \*ecdsa\.PublicKey:/\/\//' handlers/cert.go
