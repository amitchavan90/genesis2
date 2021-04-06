package genesis

import (
	"context"
	"fmt"
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
		return nil, terror.New(err, "get task")
	}
	return result, nil
}

func (r *userTaskResolver) User(ctx context.Context, obj *db.UserTask) (*db.User, error) {
	result, err := r.UserTaskStore.GetUser(obj.UserID)
	if err != nil {
		return nil, terror.New(err, "get user")
	}
	return result, nil
}

func (r *userTaskResolver) UserSubtasks(ctx context.Context, obj *db.UserTask) ([]*db.UserSubtask, error) {
	result, err := r.UserTaskStore.GetSubtasks(obj.ID)
	if err != nil {
		return nil, terror.New(err, "get user")
	}
	return result, nil
}

///////////////
//   Query   //
///////////////

func (r *queryResolver) UserTask(ctx context.Context, code *string) (*db.UserTask, error) {
	userTask, err := r.UserTaskStore.GetByCode(*code)
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
	// Get Task count (for Task Code)
	count, err := r.UserTaskStore.Count()
	if err != nil {
		return nil, terror.New(err, "create user task: Error while fetching user task count from db")
	}

	// Create UserTask
	ut := &db.UserTask{
		Code: fmt.Sprintf("UT%05d", count+1),
	}

	userTaskID, _ := uuid.NewV4()
	ut.ID = userTaskID.String()

	if input.TaskID == "" {
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

	// Verify if user already opted for task --- TODO

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
			ust := &db.UserSubtask{}
			id, _ := uuid.NewV4()
			ust.ID = id.String()
			ust.UserTaskID = null.StringFrom(created.ID)
			ust.SubtaskID = null.StringFrom(subtasks[i].ID)
			ust.Status = "Incomplete"
			ust.IsComplete = false
			_, err = r.UserTaskStore.InsertSubtask(ust)
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

// UserTaskApprove archives a userTask
func (r *mutationResolver) UserTaskApprove(ctx context.Context, id string) (*db.UserTask, error) {
	// Get UserTask
	userTaskUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	ut, err := r.UserTaskStore.Get(userTaskUUID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}

	if ut.IsComplete {
		return nil, terror.New(terror.ErrParse, "User task is already complete and approved")
	}

	ut.IsComplete = true
	ut.Status = "Complete"

	// Get user
	userUUID, err := uuid.FromString(ut.UserID)
	if err != nil {
		return nil, terror.New(terror.ErrParse, "")
	}
	u, err := r.UserStore.Get(userUUID)

	u.WalletPoints += ut.R.Task.LoyaltyPoints

	// Update UserTask
	updated, err := r.UserTaskStore.Update(ut)
	if err != nil {
		return nil, terror.New(err, "update user task")
	}

	// Update WalletTransactions
	wtID, _ := uuid.NewV4()
	wt := &db.WalletTransaction{
		ID:            wtID.String(),
		UserID:        u.ID,
		LoyaltyPoints: ut.R.Task.LoyaltyPoints,
		IsCredit:      true,
		Message:       "Loyalty points awarded by completing the task",
	}

	_, err = r.UserStore.InsertWalletTransaction(wt)
	if err != nil {
		return nil, terror.New(err, "create wallet transaction")
	}

	// Update User
	_, err = r.UserStore.Update(u)
	if err != nil {
		return nil, terror.New(err, "update user")
	}

	return updated, nil
}
