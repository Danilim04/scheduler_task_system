package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"scheduler_task_system/internal/core/entity"
)

const (
	tasksCollection = "tasks"
	database        = "azapfy"
)

type TaskRepositoryMongo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewTaskRepositoryMongo(client *mongo.Client) (*TaskRepositoryMongo, error) {
	return &TaskRepositoryMongo{
		database:   client.Database(database),
		collection: client.Database(database).Collection(tasksCollection),
	}, nil
}

func (r *TaskRepositoryMongo) ExistsByID(ctx context.Context, id entity.TaskID) (bool, error) {

	filter := bson.M{"task_id": id}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil || count < 0 {
		return false, err
	}

	return true, nil
}

func (r *TaskRepositoryMongo) Save(ctx context.Context, task *entity.Task) error {

	doc := bson.M{
		"task_id":     string(task.TaskId),
		"name":        task.Name,
		"description": task.Description,
		"payload":     task.Payload,
		"schedule":    task.Schedule,
		"status":      string(task.Status),
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepositoryMongo) Update(ctx context.Context, task *entity.Task) error {
	if err := task.IsValid(); err != nil {
		return err
	}
	filter := bson.M{"task_id": task.TaskId}
	_, err := r.collection.UpdateOne(ctx, filter, task)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryMongo) FindByID(ctx context.Context, id entity.TaskID) (*entity.Task, error) {

	filter := bson.M{"task_id": id}

	var result bson.M
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("nenhuma task encontrada com o id informado")
	}
	if err != nil {
		return nil, errors.New("erro ao buscar task: " + err.Error())
	}

	task, err := r.bsonToTask(result)
	if err != nil {
		return nil, errors.New("erro ao converter task do banco: " + err.Error())
	}

	return task, nil
}

func (r *TaskRepositoryMongo) FindAll(ctx context.Context) ([]*entity.Task, error) {

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"status": string(entity.TaskStatusActive)}, opts)

	if err != nil {
		return nil, err
	}

	var tasks []*entity.Task

	for cursor.Next(ctx) {
		var result bson.M
		task, err := r.bsonToTask(result)
		if err != nil {
			return nil, errors.New("erro ao converter task do banco: " + err.Error())
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepositoryMongo) DeleteByID(ctx context.Context, id entity.TaskID) error {

	filter := bson.M{"task_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepositoryMongo) bsonToTask(doc bson.M) (*entity.Task, error) {
	task := &entity.Task{
		TaskId:      doc["task_id"].(entity.TaskID),
		Name:        doc["name"].(string),
		Description: doc["description"].(string),
		Payload:     doc["Payload"].([]byte),
		Schedule:    doc["schedule"].(entity.Schedule),
		Status:      entity.TaskStatus(doc["status"].(string)),
		CreatedAt:   doc["created_at"].(time.Time),
		UpdatedAt:   doc["updated_at"].(time.Time),
	}
	if err := task.IsValid(); err != nil {
		return nil, err
	}
	return task, nil
}

func ConnectMongodb() (*mongo.Client, error) {
	uri := "mongodb://dev:dev123@mongodb:27017/"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func DisconnectMongodb(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}
