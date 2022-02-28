package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type Node map[interface{}]interface{}

type Structure struct {
	Root Node `json:"root"`
}

func main() {
	fp := flag.String("f", "./config.yaml", "Dirección relativa del archivo de configuración")
	op := flag.String("o", "./out/", "Directorio relativo donde escribir la estructura resultante, debe existir previamente")
	flag.Parse()

	if fp == nil || op == nil {
		log.Fatalln("Ha ocurrido un error: Flags no definidas")
	}

	log.Printf("Escribiendo estructura: %s en %s", *fp, *op)

	s, err := readStructure(*fp)
	if err != nil {
		log.Fatalln("No se pudo leer el archivo de estructura:", err)
	}

	for k, v := range s.Root {
		node, ok := v.(Node)
		if !ok {
			log.Printf("Error: %s is not a node: %v. Type: %s\n", k, node, reflect.TypeOf(v))
			continue
		}

		key, ok := k.(string)
		if !ok {
			log.Printf("Error: %s is not a valid string key. Type: %s\n", k, reflect.TypeOf(v))
			continue
		}

		err := walkNode(filepath.Join(*op, key), node)
		if err != nil {
			log.Fatalln("No se pudo caminar el archivo de estructura:", err)
		}
	}
}

func readStructure(fp string) (Structure, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return Structure{}, err
	}
	var s Structure
	if err := yaml.Unmarshal(b, &s); err != nil {
		return Structure{}, err
	}
	return s, nil
}

func walkNode(dir string, n Node) error {
	if err := createFolderIfDoesNotExist(dir); err != nil {
		return fmt.Errorf("dir: %s, error: %w", dir, err)
	}

	for k, v := range n {
		key, ok := k.(string)
		if !ok {
			return fmt.Errorf("%s is not a valid string key: type: %s", k, reflect.TypeOf(v))
		}

		if err := createFolderIfDoesNotExist(filepath.Join(dir, key)); err != nil {
			return fmt.Errorf("dir: %s, error: %w", dir, err)
		}

		subdir, ok := v.(Node)
		if !ok {
			return fmt.Errorf("cast failed: dir: %s, subdirs: %v", dir, subdir)
		}

		err := walkNode(filepath.Join(dir, key), subdir)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkDirExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}

func createFolderIfDoesNotExist(dir string) error {
	if !checkDirExists(dir) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return fmt.Errorf("dir: %s, err: %s", dir, err)
		}
	}
	return nil
}
