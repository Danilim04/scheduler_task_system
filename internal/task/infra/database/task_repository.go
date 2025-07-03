package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"scheduler_task_system/internal/task/entity"
)

const (
	tasksCollection = "tasks"
)

type TaskRepositoryMongo struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewTaskRepositoryMongo(db *mongo.Database) *TaskRepositoryMongo {
	return &TaskRepositoryMongo{
		database:   db,
		collection: db.Collection(tasksCollection),
	}
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
		"config":      task.Config,
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

	cursor, err := r.collection.Find(ctx, bson.M{"status": true}, opts)

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
		Config:      doc["config"].(map[string]interface{}),
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
