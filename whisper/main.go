package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	// Handler for file uploads
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Set max memory to limit the size of the uploaded file

		fileBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := os.Create("./audio.mp3")
		if err != nil {
			fmt.Println("Error creating temporary file")
			fmt.Fprintf(w, "Error creating temporary file")
			return
		}
		defer tempFile.Close()
		// Write this byte array to our temporary file
		tempFile.Write(fileBytes)

		fmt.Println("ABOUT TO TRANSCRIBE")
		out, err := exec.Command("whisper",
			"audio.mp3",
			"--device",
			"cuda",
			"--model",
			"small",
			"--language",
			"English",
			"--output_dir",
			"/app",
			"--output_format",
			"txt",
			"--beam_size", "5",
			// "--temperature", "0.5",
		).CombinedOutput()

		fmt.Println(err, string(out))
		final, err := os.ReadFile("audio.txt")
		if err != nil {
			w.WriteHeader(500)
			return
		}

		fmt.Println("TRANS: ", string(final))
		w.Write(final)
	})

	// Set up the server to listen on port 8080
	port := ":8080"
	log.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
