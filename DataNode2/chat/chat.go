package chat

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

func buscarArchivo(archivo string) []byte {
	file, err := ioutil.ReadFile("book_parts/" + archivo)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func fileWrite(chunk []byte, nombre string) { //aun no lo uso al 100
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

func leerLista() string {
	file, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var retorno string

	scanner := bufio.NewScanner(file)
	n := 0
	id_libro := 1
	tope := 0 //cantidad de partes
	for scanner.Scan() {
		if n == 0 || n == tope+1 {
			n = 0
			linea := strings.Split(scanner.Text(), " ")
			tope, err = strconv.Atoi(linea[1])
			retorno = retorno + "\n" + strconv.Itoa(id_libro) + ") " + linea[0]
			id_libro++
			n++
		} else {
			n++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retorno
}

func enviarPartes(id string) string {
	file, err := os.Open("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var retorno string

	scanner := bufio.NewScanner(file)
	n := 0
	id_libro := 1
	tope := 0 //cantidad de partes
	enLibro := false
	for scanner.Scan() {
		if enLibro == true && n != tope+1 {
			retorno = retorno + scanner.Text() + "#" //junta el nombre de las partes junto a su ip
			n++
		} else if n == 0 || n == tope+1 {

			if enLibro == true {
				enLibro = false
			}
			n = 0
			linea := strings.Split(scanner.Text(), " ")
			tope, err = strconv.Atoi(linea[1])
			if strconv.Itoa(id_libro) == id {
				enLibro = true
			}
			id_libro++
			n++
		} else {
			n++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retorno
}

type Server struct {
}

func (s *Server) RecibirArchivo(ctx context.Context, in *Message) (*Message, error) { //cuando el uploader envia un archivo
	log.Printf("Recibido archivo desde el cliente: %s", in.Respuesta)
	fileWrite(in.Body, in.Respuesta)
	return &Message{Respuesta: "Recibido"}, nil
}

func (s *Server) EnviarArchivo(ctx context.Context, in *Message) (*Message, error) { //cuando el downloader pide un archivo
	log.Printf("Archivo solicitado por el cliente: %s", in.Respuesta)
	file := buscarArchivo(in.Respuesta)
	return &Message{Body: file}, nil
}

func (s *Server) PedirLibros(ctx context.Context, in *Message) (*Message, error) { //cuando el downloader solicita la lista de libros
	id := in.Respuesta
	if id == "0" {
		log.Printf("Solicitud de lista de libros")
		texto := leerLista()
		return &Message{Respuesta: texto}, nil //enviar lista
	} else {
		log.Printf(id)
		texto := enviarPartes(id)
		return &Message{Respuesta: texto}, nil
	}
}

func (s *Server) PedirArchivo(ctx context.Context, in *Message) (*Message, error) { //cuando el downloader solicita las partes de un libro
	log.Printf("Cliente solicita archivo: %s", in.Respuesta)
	fileWrite(in.Body, in.Respuesta)
	return &Message{Respuesta: "Recbido"}, nil
}

func (s *Server) ProponerPropuesta(ctx context.Context, in *Message2) (*Message2, error) {
	prop := in.Mensaje
	if prop != "1" {
		return &Message2{Mensaje: "Propuesta Recibidaaaaa"}, nil

		log.Printf("Propuesta Aceptada por NameNode")
	}

	return &Message2{Mensaje: "Propuesta Recibida"}, nil
}
