package genesis

import (
	"context"
	"genesis/db"
	"genesis/graphql"

	"github.com/gofrs/uuid"
	"github.com/ninja-software/terror"
	"github.com/volatiletech/null"
)

////////////////
//  Resolver  //
////////////////

// UserTask resolver
func (r *Resolver) UserTask() graphql.UserTaskResolver {
	return &userTaskResolver{r}
}

type userTaskResolver struct{ *Resolver }

func (r *userTaskResolver) Task(ctx context.Context, obj *db.UserTask) (*db.Task, error) {
	result, err := r.UserTaskStore.GetTask(obj.TaskID.String)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

func (r *userTaskResolver) User(ctx context.Context, obj *db.UserTask) (*db.User, error) {
	result, err := r.UserTaskStore.GetUser(obj.UserID)
	if err != nil {
		return nil, terror.New(err, "get sku")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserTask(ctx context.Context, id *string) (*db.UserTask, error) {
	userTaskUUID, err := uuid.FromString(*id)
	userTask, err := r.UserTaskStore.Get(userTaskUUID)
	if err != nil {
		return nil, terror.New(err, "get userTask")
	}
	return userTask, nil
}

func (r *queryResolver) UserTasks(ctx context.Context, search graphql.SearchFilter, limit int, offset int) (*graphql.UserTasksResult, error) {
	total, userTasks, err := r.UserTaskStore.SearchSelect(search, limit, offset)
	if err != nil {
		return nil, terror.New(err, "list userTask")
	}

	result := &graphql.UserTasksResult{
		UserTasks: userTasks,
		Total:     int(total),
	}

	return result, nil
}

///////////////
// Mutations //
///////////////

// UserTaskCreate creates an userTask
func (r *mutationResolver) UserTaskCreate(ctx context.Context, input graphql.UpdateUserTask) (*db.UserTask, error) {
	// Create UserTask
	ut := &db.UserTask{}

	userTaskID, _ := uuid.NewV4()
	ut.ID = userTaskID.String()

	if input.TaskID != "" {
		return nil, terror.New(terror.ErrParse, "create userTask: Task ID is required")
	}

	// get task
	task, err := r.UserTaskStore.GetTask(input.TaskID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create userTask: Task with given ID is not found")
	}

	// get subtasks
	subtasks, err := r.TaskStore.GetSubtasks(input.TaskID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create userTask: Error while fetching subtasks")
	}

	// get user
	userID, err := r.Auther.UserIDFromContext(ctx)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "create userTask: Error while fetching user")
	}

	// get user
	// user, err := r.UserTaskStore.GetUser(input.UserID)
	// if err != nil {
	// 	return nil, terror.New(terror.ErrParse, "create userTask: Task with given ID is not found")
	// }

	ut.TaskID = null.StringFrom(task.ID)
	ut.UserID = userID.String()
	ut.Status = "Incomplete"
	ut.IsComplete = false

	created, err := r.UserTaskStore.Insert(ut)
	if err != nil {
		return nil, terror.New(err, "create userTask")
	}

	// Add subuserTask
	if len(subtasks) >= 0 {
		for i := range subtasks {
			st := &db.UserSubtask{}
			id, _ := uuid.NewV4()
			st.ID = id.String()
			st.UserTaskID = null.StringFrom(created.ID)
			st.SubtaskID = null.StringFrom(subtasks[i].ID)
			_, err = r.UserTaskStore.InsertSubtask(st)
			if err != nil {
				return nil, terror.New(err, "create subuserTask")
			}
		}
	}

	return created, nil
}

// UserTaskUpdate updates a userTask
func (r *mutationResolver) UserTaskUpdate(ctx context.Context, id string, input graphql.UpdateUserTask) (*db.UserTask, error) {
	// Get UserTask
	userTaskUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	ut, err := r.UserTaskStore.Get(userTaskUUID)
	if err != nil {
		return nil, terror.New(err, "update userTask")
	}

	// Check archived state
	if ut.IsComplete {
		return nil, terror.New(ErrArchived, "update userTask: Task is already completed")
	}

	// Update UserTask
	ut.IsComplete = true

	updated, err := r.UserTaskStore.Update(ut)
	if err != nil {
		return nil, terror.New(err, "update userTask")
	}

	return updated, nil
}
