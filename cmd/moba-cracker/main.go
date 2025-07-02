package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Merith-TK/utils/pkg/debug"
)

func init() {
	fmt.Println("MobaXterm License Keygen")
	fmt.Println("Author: Merith-TK")
	fmt.Println("Version: 1.0")

	flag.Parse()
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("    MobaXterm-Keygen.py <UserName> <Version>")
	fmt.Println()
	fmt.Println("    <UserName>:      The Name licensed to")
	fmt.Println("    <Version>:       The Version of MobaXterm")
	fmt.Println("                     Example:    10.9")
	fmt.Println()
}

func main() {
	if len(flag.Args()) != 2 {
		help()
		os.Exit(0)
	} else {
		Version := strings.Split(flag.Args()[1], ".")
		MajorVersionInt, _ := strconv.Atoi(Version[0])
		MinorVersionInt, _ := strconv.Atoi(Version[1])
		GenerateLicense(1,
			1,
			flag.Args()[0],
			MajorVersionInt,
			MinorVersionInt)
		fmt.Println("[*] Success!")
		fmt.Println("[*] File generated:", filepath.Join("Pro.key"))
		fmt.Println("[*] Please manually compress Pro.key to a zipfile and rename it to Custom.mxtpro")
		fmt.Println()
	}
}

func GenerateLicense(Type int, Count int, UserName string, MajorVersion int, MinorVersion int) {
	debug.Print("Generating license")
	if Count < 0 {
		log.Fatal("Count must be greater than or equal to 0")
	}
	debug.Print("Count is valid")
	LicenseString := fmt.Sprintf("%d#%s|%d%d#%d#%d3%d6%d#%d#%d#%d#", Type, UserName, MajorVersion, MinorVersion, Count, MajorVersion, MinorVersion, MinorVersion, 0, 0, 0)
	EncryptedLicenseString := EncryptBytes(0x787, []byte(LicenseString))
	EncodedLicenseString := VariantBase64Encode(EncryptedLicenseString)

	// write Pro.key to zip file without compression
	ZipFile, err := os.Create("Golang.mxtpro")
	if err != nil {
		log.Fatal(err)
	}

	ZipWriter := zip.NewWriter(ZipFile)
	defer ZipWriter.Close()

	header := &zip.FileHeader{
		Name:          "Pro.key",
		Method:        zip.Store, // Use Store (no compression)
		Modified:      time.Time{},
		ExternalAttrs: 0x1800000,
		// TODO: Strip characteristics
	}

	zipFileWriter, err := ZipWriter.CreateHeader(header)
	if err != nil {
		log.Fatal(err)
	}
	// Write the encoded license string to the file
	_, err = io.Copy(zipFileWriter, strings.NewReader(string(EncodedLicenseString)))
	if err != nil {
		log.Fatal(err)
	}

}
