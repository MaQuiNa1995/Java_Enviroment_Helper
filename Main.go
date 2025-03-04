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

	programs := map[string]string{
		"Eclipse":   "https://www.eclipse.org/downloads/download.php?file=/oomph/epp/2024-12/R/eclipse-inst-jre-win64.exe&mirror_id=1319",
		"Putty":     "https://the.earth.li/~sgtatham/putty/latest/w64/putty-64bit-0.83-installer.msi",
		"Curl":      "https://curl.se/download/curl-8.12.1.zip",
		"Firefox":   "https://www.mozilla.org/firefox/download/thanks/",
		"Discord":   "https://stable.dl2.discordapp.net/distro/app/stable/win/x64/1.0.9184/DiscordSetup.exe",
		"Spotify":   "https://download.scdn.co/SpotifySetup.exe",
		"Git":       "https://github.com/git-for-windows/git/releases/download/v2.48.1.windows.1/Git-2.48.1-64-bit.exe",
		"Maven":     "https://dlcdn.apache.org/maven/maven-3/3.9.9/binaries/apache-maven-3.9.9-bin.zip",
		"OpenJdk23": "https://download.oracle.com/java/23/archive/jdk-23.0.1_windows-x64_bin.zip",
		"Lombok":    "https://projectlombok.org/downloads/lombok.jar",
		"Notepad":   "https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.1.9.1/npp.8.1.9.1.Installer.x64.exe",
		"ConEmu":    "https://github.com/Maximus5/ConEmu/releases/download/v23.07.24/ConEmuSetup.230724.exe",
	}

	destinationFolder := createDestinationFolder()

	log.Println("Descargando Programas...")

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

	fmt.Println()
	if _, err := os.Stat(folder); os.IsNotExist(err) {
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
		log.Printf("[%v/%v] %v ya estaba descargado en la ruta: %v", *progress, *total, program, programLocalUrl)
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
		log.Printf("[%v/%v] Se ha descargado %v en: %v", *progress, *total, program, programLocalUrl)
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
