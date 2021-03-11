package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
)

////////////////
//  Resolver  //
////////////////

// Task resolver
func (r *Resolver) Task() graphql.TaskResolver {
	return &taskResolver{r}
}

type taskResolver struct{ *Resolver }

func (r *taskResolver) FinishDate(ctx context.Context, obj *db.Task) (*time.Time, error) {
	taskUUID, err := uuid.FromString(obj.ID)
	t, err := r.TaskStore.Get(taskUUID)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return &t.FinishDate.Time, nil
}

func (r *taskResolver) Sku(ctx context.Context, obj *db.Task) (*db.StockKeepingUnit, error) {
	result, err := r.TaskStore.GetSku(obj.SkuID.String)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

func (r *taskResolver) Subtasks(ctx context.Context, obj *db.Task) ([]*db.Subtask, error) {
	result, err := r.TaskStore.GetSubtasks(obj.ID)
	if err != nil {
		return nil, terror.New(err, "get subtasks")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) Task(ctx context.Context, id *string) (*db.Task, error) {
	taskUUID, err := uuid.FromString(*id)
	task, err := r.TaskStore.Get(taskUUID)
	if err != nil {
		return nil, terror.New(err, "get task")
	}
	return task, nil
}

func (r *queryResolver) Tasks(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.TasksResult, error) {
	total, tasks, err := r.TaskStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list task")
	}

	result := &graphql.TasksResult{
		Tasks: tasks,
		Total: int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// TaskCreate creates an task
func (r *mutationResolver) TaskCreate(ctx context.Context, input graphql.UpdateTask) (*db.Task, error) {
	// Create Task
	t := &db.Task{}

	taskID, _ := uuid.NewV4()
	t.ID = taskID.String()

	if input.Title != "" {
		t.Title = input.Title
	}
	if input.Description != "" {
		t.Description = input.Description
	}

	t.LoyaltyPoints = input.LoyaltyPoints
	t.MaximumPeople = input.MaximumPeople
	t.IsTimeBound = input.IsTimeBound
	t.IsPeopleBound = input.IsPeopleBound
	t.IsProductRelevant = input.IsProductRelevant
	t.IsFinal = false

	if input.IsTimeBound {
		if input.FinishDate == nil {
			return nil, terror.New(terror.ErrParse, "create task: finish date is not provided")
		}
		t.FinishDate = null.TimeFrom(*input.FinishDate)
	}

	if input.IsProductRelevant {
		if input.SkuID == nil {
			return nil, terror.New(terror.ErrParse, "create task: sku ID is not provided")
		}
		skuUUID, err := uuid.FromString(input.SkuID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		_, err = r.SKUStore.Get(skuUUID)
		if err != nil {
			return nil, terror.New(err, "create task")
		}
		t.SkuID = null.StringFrom(input.SkuID.String)
	}

	created, err := r.TaskStore.Insert(t)
	if err != nil {
		return nil, terror.New(err, "create task")
	}

	// Add subtask
	if len(input.Subtasks) >= 0 {
		for i := range input.Subtasks {
			st := &db.Subtask{}
			id, _ := uuid.NewV4()
			st.ID = id.String()
			st.TaskID = null.StringFrom(created.ID)
			st.Title = input.Subtasks[i].Title
			st.Description = input.Subtasks[i].Description
			_, err = r.TaskStore.InsertSubtask(st)
			if err != nil {
				return nil, terror.New(err, "create subtask")
			}
		}
	}

	return created, nil
}

// TaskUpdate updates a task
func (r *mutationResolver) TaskUpdate(ctx context.Context, id string, input graphql.UpdateTask) (*db.Task, error) {
	// Get Task
	taskUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	t, err := r.TaskStore.Get(taskUUID)
	if err != nil {
		return nil, terror.New(err, "update task")
	}

	// Check archived state
	if t.Archived {
		return nil, terror.New(ErrArchived, "cannot update archived task")
	}

	// Update Task
	if input.Title != "" {
		t.Title = input.Title
	}
	if input.Description != "" {
		t.Description = input.Description
	}

	t.LoyaltyPoints = input.LoyaltyPoints
	t.MaximumPeople = input.MaximumPeople
	t.IsTimeBound = input.IsTimeBound
	t.IsPeopleBound = input.IsPeopleBound
	t.IsProductRelevant = input.IsProductRelevant

	if input.IsTimeBound {
		if input.FinishDate == nil {
			return nil, terror.New(terror.ErrParse, "create task")
		}
		t.FinishDate = null.TimeFrom(*input.FinishDate)
	}

	if input.IsProductRelevant {
		if input.SkuID == nil {
			return nil, terror.New(terror.ErrParse, "create task: sku ID is not provided")
		}
		skuUUID, err := uuid.FromString(input.SkuID.String)
		if err != nil {
			return nil, terror.New(terror.ErrParse, "")
		}
		_, err = r.SKUStore.Get(skuUUID)
		if err != nil {
			return nil, terror.New(err, "create task")
		}
		t.SkuID = null.StringFrom(input.SkuID.String)
	}

	updated, err := r.TaskStore.Update(t)
	if err != nil {
		return nil, terror.New(err, "update task")
	}

	return updated, nil
}

// TaskArchive archives an task
func (r *mutationResolver) TaskArchive(ctx context.Context, id string) (*db.Task, error) {
	// Get Task
	taskUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Archive Task
	updated, err := r.TaskStore.Archive(taskUUID)
	if err != nil {
		return nil, terror.New(err, "update task")
	}

	return updated, nil
}

// TaskUnarchive unarchives an task
func (r *mutationResolver) TaskUnarchive(ctx context.Context, id string) (*db.Task, error) {
	// Get Task
	taskUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	// Unarchive Task
	updated, err := r.TaskStore.Unarchive(taskUUID)
	if err != nil {
		return nil, terror.New(err, "update task")
	}

	return updated, nil
}
