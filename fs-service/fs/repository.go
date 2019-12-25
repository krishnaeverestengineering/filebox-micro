package fs

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type repo struct {
	db     driver.Database
	logger log.Logger
}

type Repository interface {
	CreateUserCollection(context context.Context, userID string) (bool, error)
	CreateFolder(ctx context.Context) error
	ListDirectoryFiles(ctx context.Context) error
}

func NewRepo(config Config, logger log.Logger) (Repository, error) {
	client, err := newConnection(config)
	if err != nil {
		return nil, err
	}
	db, err := client.Database(context.Background(), config.Db.DBName)
	if err != nil {
		return nil, err
	}
	return &repo{
		db:     db,
		logger: logger,
	}, nil
}

func (r *repo) CreateUserCollection(ctx context.Context, userID string) (bool, error) {
	ecopt := &driver.CreateCollectionOptions{
		Type: driver.CollectionType(3),
	}
	dcopt := &driver.CreateCollectionOptions{
		Type: driver.CollectionType(2),
	}
	dCol, err := r.db.CreateCollection(ctx, "doc_"+userID, dcopt)
	if err != nil {
		return false, err
	}
	index, _, err := dCol.EnsurePersistentIndex(nil, []string{"id"}, &driver.EnsurePersistentIndexOptions{
		Name:         "FileId",
		Unique:       true,
		Sparse:       false,
		InBackground: false,
	})
	if err != nil {
		return false, err
	}
	fmt.Println(index)
	eCol, err := r.db.CreateCollection(ctx, "edge_"+userID, ecopt)
	if err != nil {
		return false, err
	}
	gopt := &driver.CreateGraphOptions{
		EdgeDefinitions: []driver.EdgeDefinition{driver.EdgeDefinition{
			Collection: eCol.Name(),
			From:       []string{dCol.Name()},
			To:         []string{dCol.Name()},
		}},
	}
	_, e := r.db.CreateGraph(ctx, userID, gopt)
	if e != nil {
		return false, e
	}
	//creates root node
	// c.CreateNode(context, UserFile{
	// 	UserID:   userID,
	// 	FileID:   utils.NewHash(userID),
	// 	FileName: "Root Folder",
	// 	ParentId: "root",
	// 	Type:     FileType(1),
	// })
	return true, nil
}

func (r *repo) CreateFolder(ctx context.Context) error {
	return nil
}

func (r *repo) ListDirectoryFiles(ctx context.Context) error {
	return nil
}

func newConnection(config Config) (driver.Client, error) {
	con, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: config.Client.Endpoints,
	})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     con,
		Authentication: driver.BasicAuthentication(config.Db.User, config.Db.Password),
	})
	return client, err
}
