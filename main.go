package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*")) //ruta para carpeta templates

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Printf("Servidor corriendo")         // Corrección del mensaje de log
	err := http.ListenAndServe(":8081", nil) // Separación adecuada de parámetros
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}
	defer registros.Close()

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	if err = registros.Err(); err != nil {
		panic(err.Error())
	}

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado) // Ejecuta el template inicio y pasa los datos
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil) // Ejecuta el template crear
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES (?, ?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insertarRegistros.Exec(nombre, correo)
		if err != nil {
			panic(err.Error())
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	_, err = borrarRegistros.Exec(idEmpleado)
	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	registro, err := conexionEstablecida.Query("SELECT id, nombre, correo FROM empleados WHERE id=?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}
	defer registro.Close()

	empleado := Empleado{}
	if registro.Next() {
		err = registro.Scan(&empleado.Id, &empleado.Nombre, &empleado.Correo)
		if err != nil {
			panic(err.Error())
		}
	}
	//fmt.Println(empleado) // Esto es solo para depuración; puedes eliminarlo más adelante.

	// Ejecutar el template con los datos del empleado para edición
	plantillas.ExecuteTemplate(w, "editar", empleado)
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		_, err = actualizarRegistro.Exec(nombre, correo, id)
		if err != nil {
			panic(err.Error())
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
