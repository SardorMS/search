package search

import (
	"context"
	"log"
	"testing"
)

func TestAll_success(t *testing.T) {

	root := context.Background()
	files := []string{
		"../../data/file1.txt",
		"../../data/file2.txt",
		"../../data/file3.txt",
	}
	channel := All(root, "pipeline", files)
	result, err := <-channel
	if !err {
		t.Errorf("Error_All(): %v", err)
		return
	}
	log.Printf("result: %v\n", result)

	result, err = <-channel
	if !err {
		t.Errorf("Error_All(): %v", err)
		return
	}
	log.Printf("result: %v\n", result)

	result, err = <-channel
	if !err {
		t.Errorf("Error_All(): %v", err)
		return
	}
	log.Printf("result: %v \n", result)

}

func TestAll_notSuccess(t *testing.T) {

	root := context.Background()
	files := []string{""}
	channel := All(root, "pipeline", files)
	result, err := <-channel
	if err {
		t.Errorf("Error_All(): %v", err)
		return
	}
	log.Println("result:", result)

}
