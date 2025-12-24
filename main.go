package main

import (
	"github.com/xugeaneeu/torrent-cli/shell"
)

func main() {
	shell := shell.New()
	shell.Run()

	// inPath := os.Args[1]
	// outPath := os.Args[2]

	// tf, err := torrentfile.Open(inPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = tf.DownloadToFile(outPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
