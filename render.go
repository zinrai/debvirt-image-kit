package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"text/template"
)

func cmdRender(args []string) error {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	var (
		debianVersion  = fs.String("version", "12.7.0", "Debian version")
		debianArch     = fs.String("arch", "amd64", "Debian architecture")
		diskSize       = fs.String("disk-size", "20000M", "Disk size (e.g. 5000M, 10G)")
		memorySize     = fs.String("memory", "1024", "Memory size in MB")
		sshUsername    = fs.String("ssh-username", "debian", "SSH username")
		sshPassword    = fs.String("ssh-password", "", "SSH password (random if empty)")
		isoBaseURL     = fs.String("iso-base-url", "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/", "Base URL for ISO download")
		checksumFile   = fs.String("checksum-file", "SHA256SUMS", "Checksum file name under the ISO base URL")
		isoFileName    = fs.String("iso-file", "", "ISO file name (derived if empty)")
		preseedFile    = fs.String("preseed-file", "preseed.cfg.tpl", "Path to the preseed template")
		packerTemplate = fs.String("packer-template", "debian.pkr.hcl.tpl", "Path to the Packer template")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}

	password := *sshPassword
	if password == "" {
		p, err := randomPassword(16)
		if err != nil {
			return err
		}
		password = p
		fmt.Printf("Generated random SSH password: %s\n", password)
	}

	dir := buildDirName(*debianVersion, *debianArch)
	httpDir := filepath.Join(dir, "http")
	if err := os.MkdirAll(httpDir, 0o755); err != nil {
		return fmt.Errorf("create http directory: %w", err)
	}

	preseed := preseedData{SSHUsername: *sshUsername, SSHPassword: password}
	if err := renderTemplate(*preseedFile, filepath.Join(httpDir, "preseed.cfg"), preseed); err != nil {
		return fmt.Errorf("render preseed file: %w", err)
	}

	iso := *isoFileName
	if iso == "" {
		iso = fmt.Sprintf("debian-%s-%s-netinst.iso", *debianVersion, *debianArch)
	}
	packer := packerData{
		DebianVersion: *debianVersion,
		DebianArch:    *debianArch,
		DiskSize:      *diskSize,
		MemorySize:    *memorySize,
		SSHUsername:   *sshUsername,
		SSHPassword:   password,
		ISOURL:        fmt.Sprintf("%s%s", *isoBaseURL, iso),
		ISOChecksum:   fmt.Sprintf("file:%s%s", *isoBaseURL, *checksumFile),
	}
	if err := renderTemplate(*packerTemplate, filepath.Join(dir, renderedPackerTemplate), packer); err != nil {
		return fmt.Errorf("render Packer template: %w", err)
	}

	fmt.Printf("Rendered build inputs into %s\n", dir)
	fmt.Printf("Build it with: debvirt-image-kit build --version %s --arch %s\n", *debianVersion, *debianArch)
	return nil
}

func buildDirName(version, arch string) string {
	return fmt.Sprintf("build-debian-%s-%s", version, arch)
}

type preseedData struct {
	SSHUsername string
	SSHPassword string
}

type packerData struct {
	DebianVersion string
	DebianArch    string
	DiskSize      string
	MemorySize    string
	SSHUsername   string
	SSHPassword   string
	ISOURL        string
	ISOChecksum   string
}

func renderTemplate(srcPath, dstPath string, data any) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(filepath.Base(srcPath)).Parse(string(content))
	if err != nil {
		return err
	}
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(out, data); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func randomPassword(n int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("generate password: %w", err)
		}
		b[i] = chars[idx.Int64()]
	}
	return string(b), nil
}
