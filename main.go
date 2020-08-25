package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pgavlin/goldmark"
	mdtext "github.com/pgavlin/goldmark/text"
)

func writeMimetype(zw *zip.Writer) error {
	f, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(f, "application/vnd.oasis.opendocument.text")
	return err
}

func writeManifest(zw *zip.Writer) error {
	const manifest = `<?xml version="1.0" encoding="UTF-8"?>
<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0" manifest:version="1.3" xmlns:loext="urn:org:documentfoundation:names:experimental:office:xmlns:loext:1.0">
	<manifest:file-entry manifest:full-path="/" manifest:version="1.3" manifest:media-type="application/vnd.oasis.opendocument.text"/>
	<!--<manifest:file-entry manifest:full-path="styles.xml" manifest:media-type="text/xml"/>-->
	<manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml"/>
</manifest:manifest>
`

	f, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "META-INF/manifest.xml",
		Method: zip.Deflate,
	})
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(manifest))
	return err
}

func main() {
	doc, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to read markdown: %v", err)
	}

	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("failed to open output file: %v", err)
	}

	zw := zip.NewWriter(f)
	defer zw.Close()

	if err = writeMimetype(zw); err != nil {
		log.Fatalf("failed to write mimetype: %v", err)
	}
	if err = writeManifest(zw); err != nil {
		log.Fatalf("failed to write manifest: %v", err)
	}

	content, err := zw.Create("content.xml")
	if err != nil {
		log.Fatalf("failed to create content.xml: %v", err)
	}

	parser := goldmark.DefaultParser()
	renderer := NewRenderer("", "", nil, BlockStyle{})
	if err = renderer.Render(content, doc, parser.Parse(mdtext.NewReader(doc))); err != nil {
		log.Fatalf("failed to render output: %v", err)
	}
}
