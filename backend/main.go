package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Estructura Partition
type Partition struct {
	Name   string `json:"name"`
	SizeKB int    `json:"sizeKB"`
	Type   string `json:"type"` // P o E
	Fit    string `json:"fit"`
}

// Estructura Disk
type Disk struct {
	Letter     string      `json:"letter"`
	SizeMB     int         `json:"sizeMB"`
	Partitions []Partition `json:"partitions"`
}

// Base de datos simulada en memoria
var disks = []Disk{
	{
		Letter: "A",
		SizeMB: 10,
		Partitions: []Partition{
			{"A1", 1000, "P", "FF"},
			{"A2", 1500, "P", "WF"},
			{"A3", 1200, "P", "BF"},
			{"A4", 800, "P", "FF"},
		},
	},
	{
		Letter: "B",
		SizeMB: 15,
		Partitions: []Partition{
			{"B1", 2000, "P", "BF"},
			{"B2", 1000, "P", "FF"},
			{"B3", 500, "P", "WF"},
			{"B4", 2500, "P", "FF"},
		},
	},
	{
		Letter: "C",
		SizeMB: 20,
		Partitions: []Partition{
			{"C1", 3000, "P", "FF"},
			{"C2", 1500, "P", "WF"},
			{"CEXT", 4000, "E", "BF"},
		},
	},
	{
		Letter: "D",
		SizeMB: 25,
		Partitions: []Partition{
			{"D1", 3500, "P", "BF"},
			{"D2", 2000, "P", "WF"},
			{"DEXT", 5000, "E", "FF"},
		},
	},
	{
		Letter: "E",
		SizeMB: 20,
		Partitions: []Partition{
			{"E1", 3000, "P", "FF"},
			{"EEXT", 1500, "E", "WF"},
		},
	},
	{
		Letter: "F",
		SizeMB: 25,
		Partitions: []Partition{
			{"F1", 3000, "P", "FF"},
			{"F2", 1500, "P", "WF"},
		},
	},
}

// Handler para obtener todos los discos
// Declara una función llamada getDisks que manejará solicitudes HTTP.
func getDisks(w http.ResponseWriter, r *http.Request) {
	//Establece el header Content-Type a application/json, para indicar que la respuesta será en formato JSON.
	w.Header().Set("Content-Type", "application/json")
	//Crea un codificador JSON y escribe el contenido de la variable disks directamente en la respuesta.
	json.NewEncoder(w).Encode(disks)
}

// Handler para obtener particiones de un disco por letra
// Define otro handler llamado getPartitionsByDisk, que toma la letra del disco desde la URL.
func getPartitionsByDisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	/*
		Extrae los parámetros de la URL usando mux.Vars.
		Por ejemplo, si se accede a /api/discos/C/particiones, extrae {letter: "C"}.
	*/
	params := mux.Vars(r)
	//Guarda en la variable letter el valor de la letra del disco enviada en la URL
	letter := params["letter"]

	//Recorre todos los discos disponibles (slice disks).
	for _, d := range disks {
		//Compara si la letra del disco coincide con la solicitada
		if d.Letter == letter {
			//Codifica y envía las particiones del disco como JSON.
			json.NewEncoder(w).Encode(d.Partitions)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Disco no encontrado"})
}

func main() {
	//Crea un nuevo router usando gorilla/mux
	r := mux.NewRouter()

	//Define que una petición GET a /api/discos será manejada por getDisks.
	r.HandleFunc("/api/discos", getDisks).Methods("GET")
	r.HandleFunc("/api/discos/{letter}/particiones", getPartitionsByDisk).Methods("GET")

	// Habilitar CORS para conexión con Angular (localhost:4200)
	//Crea una política de CORS (Cross-Origin Resource Sharing) para permitir que Angular (desde localhost:4200) se conecte.
	//Esto es necesario para evitar bloqueos del navegador al hacer peticiones desde otro dominio/puerto.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	//Muestra un mensaje en consola indicando que el servidor está listo.
	fmt.Println("Servidor corriendo en http://localhost:8080")
	//log.Fatal(...) detiene el programa si hay un error.
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
