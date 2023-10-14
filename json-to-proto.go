package main

import (
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func HandleJSONToProto() {
	err := filepath.Walk("commands/", walk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking through commands/ directory: %v\n", err)
		os.Exit(1)
	}

}

func walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	jsonData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	cmd := &Command{}
	if err := protojson.Unmarshal(jsonData, cmd); err != nil {
    fmt.Println(path)
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	protoData, err := proto.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal proto: %v", err)
	}

	outPath := "pb/" + path
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(outPath, protoData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
