package main

import (
	"debug/elf"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// defined in add_amd64.s
func callCAdd(addr uintptr, a, b int32) int32

func main() {
	file := "example/libadd.so"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	f, err := elf.Open(file)
	if err != nil {
		log.Fatalf("Failed to open ELF: %v", err)
	}
	defer f.Close()

	reqSize := uint64(0)
	for i, prog := range f.Progs {
		if prog.Type != elf.PT_LOAD {
			continue
		}
		if i == 0 && prog.Vaddr != 0 {
			log.Fatalf("First PT_LOAD segment must have Vaddr 0, got 0x%x", prog.Vaddr)
		}
		reqSize = max(
			reqSize,
			prog.Vaddr+prog.Memsz,
		)
	}

	pageSize := os.Getpagesize()
	numPages := (reqSize + uint64(pageSize) - 1) / uint64(pageSize)
	memSlice, err := syscall.Mmap(
		-1,                                   // fd
		0,                                    // offset
		int(numPages)*pageSize,               // length
		syscall.PROT_READ|syscall.PROT_WRITE, // RW for now, will change to RX later (W^X prevents RWx)
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS, // private mapping, not backed by a file
	)
	if err != nil {
		log.Fatalf("Failed to mmap memory: %v", err)
	}

	for _, prog := range f.Progs {
		if prog.Type != elf.PT_LOAD {
			continue
		}
		if prog.Filesz == 0 {
			continue
		}
		start := prog.Vaddr
		end := prog.Vaddr + prog.Filesz
		io.ReadFull(prog.Open(), memSlice[start:end])
	}

	err = syscall.Mprotect(memSlice, syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		log.Fatalf("Failed to mark memory as executable: %v", err)
	}

	// 1. Get the symbol table from the parsed ELF file
	symbols, err := f.Symbols()
	if err != nil {
		log.Fatalf("Failed to read symbols: %v", err)
	}

	var funcOffset uint64
	found := false

	// 2. Loop through the symbols to find the one named "add"
	for _, sym := range symbols {
		if sym.Name == "add" {
			// sym.Value is the relative offset from the start of the file
			funcOffset = sym.Value
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Symbol 'add' not found in ELF file")
	}

	// 3. Calculate the absolute pointer inside our allocated memory slice
	funcAddress := uintptr(unsafe.Pointer(&memSlice[funcOffset]))

	// Call it using our assembly bridge!
	result := callCAdd(funcAddress, 10, 32)
	fmt.Printf("Result from C: %d\n", result)
}
