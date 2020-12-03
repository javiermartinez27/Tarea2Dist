package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

var propuesta string

func sender() {
	dirname := "aux"

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		chunk, err := ioutil.ReadFile("aux/" + file.Name())
		if err != nil {
			fmt.Print(err)
		}

		sendFile(chunk, file.Name())

		err = os.Remove("aux/" + file.Name())

		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func sendFile(chunk []byte, nombre string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.154:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.RecibirArchivo(context.Background(), &chat.Message{Body: chunk, Respuesta: nombre})
	if err != nil {
		log.Fatalf("Error when calling RecibirArchivo: %s", err)
	}
	log.Printf("Response from server: %s", response.Respuesta)

}

func sendPropuesta(propuesta string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.154:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	resp, err := c.ProponerPropuesta(context.Background(), &chat.Message2{Mensaje: propuesta})
	if err != nil {
		log.Fatalf("Error when calling ProponerPropuesta: %s", err)
	}
	log.Printf("Response from server: %s", resp.Mensaje)
}

func readFiles() string {
	dirname := "libros"

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seleccione un libro: ")

	for i, file := range files {
		fmt.Println(strconv.Itoa(i+1) + ") " + file.Name())
	}

	var libro int
	fmt.Scanf("%d", &libro)

	for i, file := range files {
		if (i + 1) == libro {
			return file.Name()
		}
	}
	return ""
}

func separar(libro string) {
	fileToBeChunked := "libros/" + libro // NOMBRE

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()
	const fileChunk = 250000 // 250 kB

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Dividiendo en %d partes.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		fileName := "aux/" + libro + "_" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

		fmt.Println("Dividido a : ", fileName)
	}
}

func main() {
	libro := readFiles()
	separar(libro)
	sender()
	fmt.Println("Ahora seleccione el tipo de propuesta:\n1) Exclusion mutua centralizada\n2) Exclusion mutua distribuida")
	fmt.Scanf("%s", &propuesta)
	sendPropuesta(propuesta)

}
