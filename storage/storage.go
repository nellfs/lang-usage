package storage

import "github.com/nellfs/lang-usage/types"

type Storage interface {
  CreateCodeReport(*types.CodeReport) (error)
  CreateLanguage(*types.Language) (error) 
	GetCodeReport(int) (*types.CodeReport, error)
  GetLanguage(*string) ([]*types.Language, error)
	GetLanguageIDByName(string) (int, error)
  GetLastRequestID() (int, error)
}
