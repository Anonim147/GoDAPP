package temp

import (
	"bytes"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func download() {
	log.SetFlags(log.Lshortfile)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		modtime := time.Now()
		content := randomContent(modtime.UnixNano(), 1024)

		// ServeContent uses the name for mime detection
		const name = "random.txt"

		// tell the browser the returned content should be downloaded
		w.Header().Add("Content-Disposition", "Attachment")

		http.ServeContent(w, req, name, modtime, content)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func randomContent(seed int64, length int) io.ReadSeeker {
	r := rand.New(rand.NewSource(seed))

	content := make([]byte, length, length)
	for i := range content {
		b := byte(r.Intn(math.MaxUint8))

		b = b%('~'-' ') + ' ' // make it a visible character

		content[i] = b
	}

	return bytes.NewReader(content)
}
