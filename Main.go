package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const downloadFolder = "downloads"

func main() {

	fmt.Println()
	log.Println("Comprobando variables de entorno...")
	envKeys := [...]string{"JAVA_HOME", "JDK_HOME", "JRE_HOME", "MAVEN_HOME"}
	for _, key := range envKeys {
		checkEnvVar(key)
	}
	fmt.Println()

	log.Println("Comprobando la variable PATH")
	pathEnvKeys := [...]string{"JAVA_HOME", "JDK_HOME", "MAVEN_HOME"}
	for _, key := range pathEnvKeys {
		checkEnvPath(key)
	}
	createFolder(downloadFolder)

	programs := map[string]string{
		"eclipse": "https://rhlx01.hs-esslingen.de/pub/Mirrors/eclipse/technology/epp/downloads/release/2021-09/R/eclipse-jee-2021-09-R-win32-x86_64.zip",
		"maven":   "https://dlcdn.apache.org/maven/maven-3/3.8.3/binaries/apache-maven-3.8.3-bin.zip",
		"openJdk": "https://download.java.net/java/GA/jdk11/13/GPL/openjdk-11.0.1_windows-x64_bin.zip",
		"lombok":  "https://projectlombok.org/downloads/lombok.jar",
		"notepad": "https://github.com/notepad-plus-plus/notepad-plus-plus/releases/download/v8.1.9.1/npp.8.1.9.1.Installer.x64.exe",
		"conEmu":  "https://download.fosshub.com/Protected/expiretime=1636888904;badurl=aHR0cHM6Ly93d3cuZm9zc2h1Yi5jb20vQ29uRW11Lmh0bWw=/7253d451ada51c2054be0702c1fe244f0b786c220ae58926a07cc3198d933f41/5b85860af9ee5a5c3e979f45/613e772663102e500262817b/ConEmuSetup.210912.exe",
	}

	var wg sync.WaitGroup
	wg.Add(len(programs))

	for program, url := range programs {
		go downloadFile(program, url, &wg)
		wg.Done()
	}
	wg.Wait()
	log.Print("Acabado")
}

func checkEnvPath(value string) {

	fullVar := os.Getenv(value) + "\\bin"

	if val, _ := os.LookupEnv("PATH"); strings.Contains(val, fullVar) {
		log.Println("El PATH contiene:", value)
	} else {
		log.Println("El PATH no contiene ", value)
	}
}

func createFolder(url string) {
	fmt.Println()
	log.Println("Creating a folder with url:", url)
	os.Mkdir(url, 0777)
	fmt.Println()
}

func checkEnvVar(key string) {
	if val, exists := os.LookupEnv(key); exists {
		log.Printf("Variable de entorno: %v  está seteada con valor: %v", key, val)
	} else {
		log.Printf("Variable de entorno: %v no encontrada \n", key)
	}
}

func downloadFile(program string, url string, wg *sync.WaitGroup) {
	log.Println("Se empieza a descargar:", program)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	splittedUrl := strings.Split(url, "/")
	programNameIndex := len(splittedUrl) - 1

	out, err := os.Create("." + string(os.PathSeparator) + downloadFolder + string(os.PathSeparator) + splittedUrl[programNameIndex])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err == nil {
		log.Println("Se ha descargado:", program)
	}
}
