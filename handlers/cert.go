package handlers

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/crypto/ssh"
)

type CertHandler struct{}

func (h *CertHandler) CanHandle(mime, ext string) bool {
	switch ext {
	case ".pem", ".crt", ".cer", ".der", ".key", ".p12", ".pfx", ".pub":
		return true
	}
	return false
}

func (h *CertHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	if opts.Flat {
		_, err := io.Copy(w, r)
		return err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	ext := strings.ToLower(meta.Name)

	if strings.HasSuffix(ext, ".pub") {
		return h.renderSSH(w, data)
	}
	if strings.HasSuffix(ext, ".p12") || strings.HasSuffix(ext, ".pfx") {
		return h.renderPKCS12(w, data)
	}
	if strings.HasSuffix(ext, ".key") {
		return h.renderPrivateKey(w, data)
	}

	return h.renderCert(w, data, ext)
}

func (h *CertHandler) renderCert(w io.Writer, data []byte, ext string) error {
	var certs []*x509.Certificate

	if ext == ".der" {
		cert, err := x509.ParseCertificate(data)
		if err != nil {
			return err
		}
		certs = append(certs, cert)
	} else {
		// PEM parsing
		for {
			block, rest := pem.Decode(data)
			if block == nil {
				break
			}
			if block.Type == "CERTIFICATE" {
				cert, err := x509.ParseCertificate(block.Bytes)
				if err == nil {
					certs = append(certs, cert)
				}
			}
			data = rest
		}
	}

	if len(certs) == 0 {
		return fmt.Errorf("no certificates found")
	}

	for i, cert := range certs {
		if i > 0 {
			fmt.Fprintln(w, "---")
		}
		fmt.Fprintf(w, "Subject: %s\n", cert.Subject.String())
		fmt.Fprintf(w, "Issuer: %s\n", cert.Issuer.String())
		fmt.Fprintf(w, "Serial Number: %x\n", cert.SerialNumber)
		fmt.Fprintf(w, "Validity:\n")
		fmt.Fprintf(w, "  Not Before: %s\n", cert.NotBefore.String())
		fmt.Fprintf(w, "  Not After:  %s\n", cert.NotAfter.String())
		if len(cert.DNSNames) > 0 || len(cert.IPAddresses) > 0 {
			fmt.Fprintf(w, "SANs: ")
			var sans []string
			sans = append(sans, cert.DNSNames...)
			for _, ip := range cert.IPAddresses {
				sans = append(sans, ip.String())
			}
			fmt.Fprintln(w, strings.Join(sans, ", "))
		}

		keyType, keySize := publicKeyInfo(cert.PublicKey)
		fmt.Fprintf(w, "Public Key: %s (%d bits)\n", keyType, keySize)

		fingerprint := sha256.Sum256(cert.Raw)
		fmt.Fprintf(w, "SHA256 Fingerprint: %s\n", hex.EncodeToString(fingerprint[:]))
	}

	return nil
}

func (h *CertHandler) renderPrivateKey(w io.Writer, data []byte) error {
	block, _ := pem.Decode(data)
	if block == nil {
		return fmt.Errorf("no PEM data found")
	}

	var key interface{}
	var err error
	if strings.Contains(block.Type, "RSA PRIVATE KEY") {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if strings.Contains(block.Type, "EC PRIVATE KEY") {
		key, err = x509.ParseECPrivateKey(block.Bytes)
	} else {
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	if err != nil {
		return err
	}

	var keyType string
	var keySize int

	switch k := key.(type) {
	case *rsa.PrivateKey:
		keyType = "RSA Private Key"
		keySize = k.Size() * 8
	case *ecdsa.PrivateKey:
		keyType = "ECDSA Private Key"
		keySize = k.Params().BitSize
		// case ed25519.PrivateKey:
		keyType = "Ed25519 Private Key"
		keySize = 256
	default:
		keyType = "Unknown Private Key"
	}

	fmt.Fprintf(w, "Type: %s\n", keyType)
	fmt.Fprintf(w, "Size: %d bits\n", keySize)
	return nil
}

func (h *CertHandler) renderPKCS12(w io.Writer, data []byte) error {
	blocks, err := pkcs12.ToPEM(data, "")
	if err != nil {
		return fmt.Errorf("PKCS12 error (might require password): %w", err)
	}

	certCount := 0
	for _, block := range blocks {
		if block.Type == "CERTIFICATE" {
			certCount++
		}
	}
	fmt.Fprintf(w, "PKCS12 Archive\n")
	fmt.Fprintf(w, "Certificates in chain: %d\n", certCount)
	return nil
}

func (h *CertHandler) renderSSH(w io.Writer, data []byte) error {
	pub, _, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return err
	}

	keyType := pub.Type()
	parsed, err := ssh.ParsePublicKey(pub.Marshal())
	if err != nil {
		return err
	}

	var keySize int
	if cp, ok := parsed.(ssh.CryptoPublicKey); ok {
		_, keySize = publicKeyInfo(cp.CryptoPublicKey())
	}

	fp := ssh.FingerprintSHA256(pub)
	fmt.Fprintf(w, "Type: %s\n", keyType)
	if keySize > 0 {
		fmt.Fprintf(w, "Size: %d bits\n", keySize)
	}
	fmt.Fprintf(w, "Fingerprint: %s\n", fp)

	return nil
}

func publicKeyInfo(pub interface{}) (string, int) {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		return "RSA", k.Size() * 8
	case *ecdsa.PublicKey:
		return "ECDSA", k.Params().BitSize
		// case ed25519.PublicKey:
		return "Ed25519", 256
	default:
		return "Unknown", 0
	}
}
