package fs

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type repo struct {
	db     driver.Database
	logger log.Logger
}

type FileType int

const (
	Folder   = FileType(1)
	TextFile = FileType(2)
)

type UserFile struct {
	UserID   string   `json:"userId"`
	FileID   string   `json:"id"`
	FileName string   `json:"filename"`
	ParentId string   `json:"parentId"`
	RootId   string   `json:"rootId"`
	Type     FileType `json:"type"`
	Path     string   `json:"path"`
}

type DirectoryEdge struct {
	*driver.EdgeDocument
	UserID string
}

const (
	EDGE_COLLECTION = "edge_"
	DOC_COLLECTION  = "doc_"
)

type Repository interface {
	GetFullPath(ctx context.Context, id string, root string, userID string) string
	CreateUser(context context.Context, userID string) (bool, error)
	CreateFolder(ctx context.Context, data UserFile) error
	ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error)
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

func createUserGraph(ctx context.Context, userID string, eColName string, dColName string, db driver.Database) (driver.Graph, error) {
	gopt := &driver.CreateGraphOptions{
		EdgeDefinitions: []driver.EdgeDefinition{driver.EdgeDefinition{
			Collection: eColName,
			From:       []string{dColName},
			To:         []string{dColName},
		}},
	}
	return db.CreateGraph(ctx, userID, gopt)
}

func createUserEdgeCollection(ctx context.Context, userID string, db driver.Database) (driver.Collection, error) {
	ecopt := &driver.CreateCollectionOptions{
		Type: driver.CollectionType(3),
	}
	eCol, err := db.CreateCollection(ctx, joinStrings(EDGE_COLLECTION, userID), ecopt)
	return eCol, err
}

func createUserDocumentCollection(ctx context.Context, userID string, db driver.Database) (driver.Collection, error) {
	dcopt := &driver.CreateCollectionOptions{
		Type: driver.CollectionType(2),
	}
	dCol, err := db.CreateCollection(ctx, joinStrings(DOC_COLLECTION, userID), dcopt)
	if err != nil {
		return nil, err
	}
	_, _, er := dCol.EnsurePersistentIndex(ctx, []string{"id"}, &driver.EnsurePersistentIndexOptions{
		Name:         "FileId",
		Unique:       true,
		Sparse:       false,
		InBackground: false,
	})
	if er != nil {
		return nil, er
	}
	return dCol, nil
}

func joinStrings(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func getDocumentKey(c context.Context, collection string, id string, db driver.Database) (string, error) {
	query := `FOR doc IN @@collection
				FILTER doc.path == @id
					RETURN doc`
	bindVars := map[string]interface{}{
		"@collection": collection,
		"id":          id,
	}
	cursor, err := db.Query(c, query, bindVars)
	if err != nil {
		return "", err
	}
	defer cursor.Close()
	for {
		var file UserFile
		meta, err := cursor.ReadDocument(nil, &file)
		if driver.IsNoMoreDocuments(err) || err != nil {
			break
		}
		return meta.Key, nil
	}
	return "", nil
}

func (r *repo) GetFullPath(ctx context.Context, id string, rootId string, userID string) string {
	query := `FOR v, e IN OUTBOUND SHORTEST_PATH @from TO @to GRAPH @graph
				RETURN v`
	bindVars := map[string]interface{}{
		"from":  joinStrings(DOC_COLLECTION, userID, "/", rootId),
		"to":    joinStrings(DOC_COLLECTION, userID, "/", id),
		"graph": userID,
	}
	cursor, err := r.db.Query(ctx, query, bindVars)
	if err != nil {
		return ""
	}
	var path string = ""
	for {
		var file UserFile
		_, err := cursor.ReadDocument(nil, &file)
		if driver.IsNoMoreDocuments(err) || err != nil {
			break
		}
		path = filepath.Join(path, file.FileName)
	}
	return path
}

func (r *repo) CreateUser(ctx context.Context, userID string) (bool, error) {
	exists, err := r.db.CollectionExists(ctx, joinStrings(DOC_COLLECTION, userID))
	if exists || err != nil {
		return true, nil
	}
	dCol, err := createUserDocumentCollection(ctx, userID, r.db)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	eCol, err := createUserEdgeCollection(ctx, userID, r.db)
	if err != nil {
		return false, err
	}
	graph, err := createUserGraph(ctx, userID, eCol.Name(), dCol.Name(), r.db)
	if err != nil {
		return false, err
	}
	fmt.Println(graph.Name())
	er := r.CreateFolder(ctx, UserFile{
		UserID:   userID,
		FileID:   NewHash(userID),
		FileName: "Root Folder",
		ParentId: "root",
		Type:     Folder,
	})
	//r.db.AbortTransaction()
	if er != nil {
		return false, er
	}
	return true, nil
}

func (r *repo) CreateFolder(ctx context.Context, data UserFile) error {
	graph, err := r.db.Graph(ctx, data.UserID)
	if err != nil {
		return err
	}
	eCol, _, err := graph.EdgeCollection(ctx, joinStrings(EDGE_COLLECTION, data.UserID))
	if err != nil {
		return err
	}
	uCol, err := r.db.Collection(ctx, joinStrings(DOC_COLLECTION, data.UserID))
	if err != nil {
		return err
	}
	meta, err := uCol.CreateDocument(ctx, data)
	if err != nil {
		return err
	}
	key, err := getDocumentKey(nil, uCol.Name(), data.ParentId, r.db)
	if err != nil {
		return err
	}
	from := joinStrings(uCol.Name(), "/", key)
	to := joinStrings(uCol.Name(), "/", meta.Key)
	eMeta, err := eCol.CreateDocument(nil, DirectoryEdge{
		UserID: data.UserID,
		EdgeDocument: &driver.EdgeDocument{
			From: driver.DocumentID(from),
			To:   driver.DocumentID(to),
		},
	})
	fmt.Println(eMeta, meta)
	return nil
}

func readUserFileDataCursor(c driver.Cursor) []UserFile {
	data := make([]UserFile, 0)
	for {
		var file UserFile
		_, err := c.ReadDocument(nil, &file)
		if driver.IsNoMoreDocuments(err) || err != nil {
			break
		}
		data = append(data, file)
	}
	return data
}

func (r *repo) ListDirectoryFiles(ctx context.Context, targetID string, userID string) ([]UserFile, error) {
	uCol, err := r.db.Collection(ctx, joinStrings(DOC_COLLECTION, userID))
	if err != nil {
		return []UserFile{}, err
	}

	key, err := getDocumentKey(nil, uCol.Name(), targetID, r.db)
	if err != nil {
		return []UserFile{}, err
	}

	dir := joinStrings(DOC_COLLECTION, userID)
	query := `FOR doc IN @@collection
				FILTER doc.path == @targetId
				FOR v, e, p IN 1..1 OUTBOUND @path GRAPH @graph
					RETURN v`
	bindVars := map[string]interface{}{
		"@collection": dir,
		"targetId":    targetID,
		"path":        joinStrings(dir, "/", key),
		"graph":       userID,
	}
	cursor, err := r.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	return readUserFileDataCursor(cursor), nil
}
