package main

import "github.com/nellfs/lang-usage/storage"

type Server struct {
	listenAddr string
	store     storage.Storage 
}
