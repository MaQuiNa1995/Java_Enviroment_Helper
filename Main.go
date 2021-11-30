package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "github.com/dimiro1/banner/autoload"
)

func main() {

	destinationFolder := createDestinationFolder()

	programs := map[string]string{
		"Eclipse":   "https://rhlx01.hs-esslingen.de/pub/Mirrors/eclipse/technology/epp/downloads/release/2021-09/R/eclipse-jee-2021-09-R-win32-x86_64.zip",
		"Git":       "https://github.com/git-for-windows/git/releases/download/v2.34.1.windows.1/Git-2.34.1-64-bit.exe",
		"Maven":     "https://dlcdn.apache.org/maven/maven-3/3.8.3/binaries/apache-maven-3.8.3-bin.zip",
		"OpenJdk11": "https://download.java.net/java/GA/jdk11/13/GPL/openjdk-11.0.1_windows-x64_bin.zip",
		"Lombok":    "https://projectlombok.org/downloads/lombok.jar",
		"Notepad":   "https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.1.9.1/npp.8.1.9.1.Installer.x64.exe",
		"ConEmu":    "https://download.fosshub.com/Protected/expiretime=1636888904;badurl=aHR0cHM6Ly93d3cuZm9zc2h1Yi5jb20vQ29uRW11Lmh0bWw=/7253d451ada51c2054be0702c1fe244f0b786c220ae58926a07cc3198d933f41/5b85860af9ee5a5c3e979f45/613e772663102e500262817b/ConEmuSetup.210912.exe",
	}

	var wg sync.WaitGroup

	total := len(programs)

	wg.Add(total)

	count := 0

	for program, url := range programs {
		go downloadFile(program, url, &wg, destinationFolder, &count, &total)
	}

	wg.Wait()

	fmt.Println()
	log.Println("Comprobando variables de entorno...")
	envKeys := [...]string{"JAVA_HOME", "MAVEN_HOME", "JRE_HOME"}
	for _, key := range envKeys {
		checkEnvVar(key)
	}
	fmt.Println()

	log.Println("Comprobando la variable PATH")
	pathEnvKeys := envKeys[:len(envKeys)-1]
	for _, key := range pathEnvKeys {
		checkEnvPath(key)
	}

	fmt.Println()
	log.Print("Proceso terminado, aprieta enter para salir...")
	fmt.Scanln()

}

func checkEnvPath(value string) {

	fullVar := os.Getenv(value) + string(os.PathSeparator) + "bin"

	if val, _ := os.LookupEnv("PATH"); strings.Contains(val, fullVar) {
		log.Println("El PATH contiene:", value)
	} else {
		log.Println("El PATH no contiene ", value)
	}
}

func createDestinationFolder() string {

	folder := os.Getenv("USERPROFILE") + string(os.PathSeparator) + "downloads"

	_, err := os.Stat(folder)

	fmt.Println()
	if os.IsNotExist(err) {
		log.Println("Creando carpeta:", folder)
		os.Mkdir(folder, 0777)
	}
	log.Println("Aqui se descargarán los programas:", folder)

	fmt.Println()

	return folder + string(os.PathSeparator)
}

func checkEnvVar(key string) {
	if val, exists := os.LookupEnv(key); exists {
		log.Printf("%v está seteada con valor: %v", key, val)
	} else {
		log.Printf("%v no encontrada \n", key)
	}
}

func downloadFile(program string, url string, wg *sync.WaitGroup, destinationFolder string, progress *int, total *int) {

	log.Printf("Descargando: %v...", program)
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	splittedUrl := strings.Split(url, "/")
	programNameIndex := len(splittedUrl) - 1

	programLocalUrl := destinationFolder + splittedUrl[programNameIndex]

	if exists(programLocalUrl) {
		*progress++
		log.Printf("[%v/%v] %v ya está descargado en la ruta %v", *progress, *total, program, programLocalUrl)
		return
	}

	out, err := os.Create(programLocalUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err == nil {
		*progress++
		log.Printf("[%v/%v] Se ha descargado: %v en: %v", *progress, *total, program, programLocalUrl)
	}

}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
