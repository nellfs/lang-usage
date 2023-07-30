package storage

import "github.com/nellfs/lang-usage/types"

type Storage interface {
  CreateCodeReport(*types.CodeReport) (error)
  CreateLanguage(*types.Language) (error) 
	GetCodeReport(int) (*types.CodeReport, error)
	GetLanguageIDByName(string) (*types.Language, error)
  GetLastRequestID() (int, error)
}
