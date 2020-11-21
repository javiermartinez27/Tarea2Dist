package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tutorialedge/go-grpc-tutorial/chat"
	"google.golang.org/grpc"
)

func pedirLibros() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.154:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.PedirLibros(context.Background(), &chat.Message{Respuesta: "0"})
	if err != nil {
		log.Fatalf("Error when calling PedirLibros: %s", err)
	}
	log.Printf("Libros disponibles %s", response.Respuesta)
}

func pedirPartes(id string) string {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("10.10.28.154:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.PedirLibros(context.Background(), &chat.Message{Respuesta: id})
	if err != nil {
		log.Fatalf("Error when calling PedirLibros: %s", err)
	}
	return response.Respuesta
}

func pedirArchivo(archivo string, port string) []byte {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.EnviarArchivo(context.Background(), &chat.Message{Respuesta: archivo})
	if err != nil {
		log.Fatalf("Error when calling enviarArchivo: %s", err)
	}
	return response.Body
}

func fileWrite(chunk []byte, nombre string) {
	file, err := os.OpenFile(
		nombre,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	bytesWritten, err := file.Write(chunk)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}

func juntar(nombre string, ctadPartes int) {

	_, err := os.Create(nombre)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.OpenFile(nombre, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// en que parte del nuevo archivo estamos escribiendo
	var writePosition int64 = 0

	for j := uint64(0); j < uint64(ctadPartes); j++ {

		//leyendo un chunk
		currentChunkFileName := "aux/" + nombre + "_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Insertando en : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(nombre, chunkBufferBytes, os.ModeAppend)

		n, err := file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		fmt.Println("Escritos ", n, " bytes")

		fmt.Println("Recombinando parte [", j, "] en : ", nombre)

		err = os.Remove(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// now, we close the nombre
	file.Close()
}

func main() {
	pedirLibros()
	var id_libro string
	fmt.Println("Seleccione qué libro desea:")
	fmt.Scanf("%s", &id_libro)
	partes := pedirPartes(id_libro) //funcion que manda ubicacion de archivos de acuerdo al id
	listaPartes := strings.Split(partes, "#")
	ctadPartes := 0
	for _, parte := range listaPartes {
		if parte != "" {
			ctadPartes++
			part := strings.Split(parte, " ")
			file := pedirArchivo(part[0], part[1])
			fileWrite(file, "aux/"+part[0])
		}
	}

	juntar(strings.Split(partes, ".")[0]+".pdf", ctadPartes)

}
