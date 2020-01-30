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
	CreateFolder(ctx context.Context, data UserFile) (interface{}, error)
	ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error)
	DeleteFileOrFolder(ctx context.Context, targetID string, userID string)
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
				FILTER doc.id == @id
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
	col := joinStrings(DOC_COLLECTION, userID)
	rootKey, err := getDocumentKey(ctx, col, rootId, r.db)
	targetKey, err := getDocumentKey(ctx, col, id, r.db)
	query := `FOR v, e IN OUTBOUND SHORTEST_PATH @from TO @to GRAPH @graph
				RETURN v`
	bindVars := map[string]interface{}{
		"from":  joinStrings(DOC_COLLECTION, userID, "/", rootKey),
		"to":    joinStrings(DOC_COLLECTION, userID, "/", targetKey),
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
	_, er := r.CreateFolder(ctx, UserFile{
		UserID:   userID,
		FileID:   "/",
		FileName: "/",
		ParentId: "/",
		Type:     Folder,
	})
	//r.db.AbortTransaction()
	if er != nil {
		return false, er
	}
	return true, nil
}

func (r *repo) CreateFolder(ctx context.Context, data UserFile) (interface{}, error) {
	graph, err := r.db.Graph(ctx, data.UserID)
	if err != nil {
		return nil, err
	}
	eCol, _, err := graph.EdgeCollection(ctx, joinStrings(EDGE_COLLECTION, data.UserID))
	if err != nil {
		return nil, err
	}
	uCol, err := r.db.Collection(context.Background(), joinStrings(DOC_COLLECTION, data.UserID))
	if err != nil {
		return nil, err
	}
	meta, err := uCol.CreateDocument(context.Background(), UserFile{
		FileName: data.FileName,
		FileID:   data.FileID,
		Path:     data.FileName,
		Type:     Folder,
		ParentId: data.ParentId,
		UserID:   data.UserID,
	})
	if err != nil {
		return nil, err
	}
	key, err := getDocumentKey(context.Background(), uCol.Name(), data.ParentId, r.db)
	if err != nil {
		return nil, err
	}
	from := joinStrings(uCol.Name(), "/", key)
	to := joinStrings(uCol.Name(), "/", meta.Key)
	eMeta, err := eCol.CreateDocument(context.Background(), DirectoryEdge{
		UserID: data.UserID,
		EdgeDocument: &driver.EdgeDocument{
			From: driver.DocumentID(from),
			To:   driver.DocumentID(to),
		},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(eMeta, meta)
	return r.ListDirectoryFiles(ctx, data.ParentId, data.UserID)
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
	dir := joinStrings(DOC_COLLECTION, userID)
	query := `FOR doc IN @@collection
				FILTER doc.id == @targetId
				FOR v, e, p IN 1..1 OUTBOUND doc._id GRAPH @gh
					FILTER v.id != "/"
						RETURN v`
	bindVars := map[string]interface{}{
		"@collection": dir,
		"targetId":    targetID,
		"gh":          userID,
	}
	cursor, err := r.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	return readUserFileDataCursor(cursor), nil
}

func (r *repo) DeleteFileOrFolder(ctx context.Context, targetID string, userID string) {
	// dir := joinStrings(DOC_COLLECTION, userID)
	// query := `FOR doc in @@collection
	// 			FILTER doc.id == @targetId
	// action := `function (params) {
	// 	var db = require('@arangodb').db;
	// 	db._query('FOR i IN doc_113176837686976104031 RETURN i');
	// }`

	// d, err := r.db.Transaction(nil, action, &driver.TransactionOptions{
	// 	ReadCollections: []string{
	// 		"doc_113176837686976104031",
	// 	},
	// 	WaitForSync: true,
	// })
	// fmt.Println(d, err)
	// key, err := getDocumentKey(ctx, joinStrings(DOC_COLLECTION, userID), targetID, r.db)
	// if err != nil {
	// 	return
	// }
	// eCol, err := r.db.Collection(ctx, joinStrings(EDGE_COLLECTION, userID))
	// if err != nil {
	// 	return
	// }
	// removeEdgeQuery := `FOR doc IN @@collection
	// 						FILTER doc.id == @targetID
	// 							`

}

func deleteEdge() {

}
