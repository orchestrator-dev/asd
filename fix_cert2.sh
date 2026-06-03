sed -i -e '/switch k := pub.(type) {/,/}/c\
	switch k := pub.(type) {\
	case *rsa.PublicKey:\
		return "RSA", k.Size() * 8\
	case *ecdsa.PublicKey:\
		return "ECDSA", k.Params().BitSize\
	case ed25519.PublicKey:\
		return "Ed25519", 256\
	default:\
		return "Unknown", 0\
	}' handlers/cert.go
