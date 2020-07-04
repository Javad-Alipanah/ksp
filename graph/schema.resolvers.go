package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"ksp/graph/generated"
	"ksp/graph/model"
	"ksp/internal/pkg/alg"
	database "ksp/internal/pkg/db/mysql"
	db_model "ksp/internal/pkg/model"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func (r *mutationResolver) CreateBoard(ctx context.Context, newBoard model.NewBoard) (*model.Board, error) {
	if newBoard.Size < 3 {
		return nil, errors.New("minimum size is 3")
	}

	if newBoard.Start.X >= newBoard.Size || newBoard.Start.Y >= newBoard.Size || newBoard.Start.X < 0 || newBoard.Start.Y < 0 {
		return nil, errors.New("invalid start point")
	}

	if newBoard.Target.X >= newBoard.Size || newBoard.Target.Y >= newBoard.Size || newBoard.Target.X < 0 || newBoard.Target.Y < 0 {
		return nil, errors.New("invalid target point")
	}

	start, _ := json.Marshal(newBoard.Start)
	target, _ := json.Marshal(newBoard.Target)

	board := db_model.Board{
		Size:   newBoard.Size,
		Start:  string(start),
		Target: string(target),
		Path:   "[]",
	}

	id := board.Save()

	// calculate path asynchronously
	go func() {
		source := model.Point{
			X: newBoard.Start.X,
			Y: newBoard.Start.Y,
		}

		sink := model.Point{
			X: newBoard.Target.X,
			Y: newBoard.Target.Y,
		}
		path := alg.CalcPath(source, sink, newBoard.Size)
		if path != nil {
			jsonVal, err := json.Marshal(path)
			if err != nil {
				log.Fatal(err)
			}

			board.ID = strconv.FormatInt(id, 10)
			board.Path = string(jsonVal)
			trx, err := database.Db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			_ = board.Update(trx)
			err = trx.Commit()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	return &model.Board{
		ID:     strconv.FormatInt(id, 10),
		Size:   newBoard.Size,
		Start:  &model.Point{X: newBoard.Start.X, Y: newBoard.Start.Y},
		Target: &model.Point{X: newBoard.Target.X, Y: newBoard.Target.Y},
		Path:   nil,
	}, nil
}

func (r *mutationResolver) UpdateBoard(ctx context.Context, board model.UpdateBoard) (*model.Board, error) {
	tx, err := database.Db.Begin()
	if err != nil {
		log.Warn(err)
		return nil, errors.New("internal server error")
	}

	id, err := strconv.ParseUint(board.ID, 10, 64)
	if err != nil {
		_ = tx.Rollback()
		return nil, errors.New("invalid id")
	}

	// get previous value and merge with new value
	prev := db_model.SelectForUpdate(id, tx)
	if prev == nil {
		_ = tx.Rollback()
		return nil, errors.New("not found")
	}

	if board.Size != nil && *board.Size < 3 {
		_ = tx.Rollback()
		return nil, errors.New("minimum size is 3")
	}

	var newSize *int
	var newStart *model.NewPoint
	var newTarget *model.NewPoint
	if newSize = board.Size; newSize != nil {
		prev.Size = *newSize
	}

	if newStart = board.Start; newStart != nil {
		jsonVal, _ := json.Marshal(*newStart)
		prev.Start = string(jsonVal)
	} else {
		newStart = new(model.NewPoint)
		err = json.Unmarshal([]byte(prev.Start), newStart)
		if err != nil {
			log.Fatal(err)
		}
	}

	if newTarget = board.Target; newTarget != nil {
		jsonVal, _ := json.Marshal(*newTarget)
		prev.Target = string(jsonVal)
	} else {
		newTarget = new(model.NewPoint)
		err = json.Unmarshal([]byte(prev.Target), newTarget)
		if err != nil {
			log.Fatal(err)
		}
	}

	size := prev.Size
	if newStart.X >= size || newStart.Y >= size || newStart.X < 0 || newStart.Y < 0 {
		_ = tx.Rollback()
		return nil, errors.New("incompatible input")
	}

	if newTarget.X >= size || newTarget.Y >= size || newTarget.X < 0 || newTarget.Y < 0 {
		_ = tx.Rollback()
		return nil, errors.New("incompatible input")
	}

	prev.ID = board.ID
	prev.Path = "[]"
	prev.Update(tx)
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return nil, errors.New("internal server error")
	}

	// calculate path asynchronously
	go func() {
		source := model.Point{
			X: newStart.X,
			Y: newStart.Y,
		}

		sink := model.Point{
			X: newTarget.X,
			Y: newTarget.Y,
		}
		path := alg.CalcPath(source, sink, prev.Size)
		if path != nil {
			jsonVal, err := json.Marshal(path)
			if err != nil {
				log.Fatal(err)
			}

			prev.Path = string(jsonVal)
			trx, err := database.Db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			_ = prev.Update(trx)
			err = trx.Commit()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	return &model.Board{
		ID:   prev.ID,
		Size: prev.Size,
		Start: &model.Point{
			X: newStart.X,
			Y: newStart.Y,
		},
		Target: &model.Point{
			X: newTarget.X,
			Y: newTarget.Y,
		},
		Path: nil,
	}, nil
}

func (r *mutationResolver) DeleteBoard(ctx context.Context, id string) (*bool, error) {
	n := db_model.Delete(id)
	ret := n > 0
	return &ret, nil
}

func (r *queryResolver) Boards(ctx context.Context) ([]*model.Board, error) {
	var result []*model.Board
	dbResult := *db_model.GetAll()
	for _, board := range dbResult {
		var start model.Point
		var target model.Point
		var path []*model.Point
		_ = json.Unmarshal([]byte(board.Start), &start)
		_ = json.Unmarshal([]byte(board.Target), &target)
		_ = json.Unmarshal([]byte(board.Path), &path)
		result = append(result, &model.Board{
			ID:   board.ID,
			Size: board.Size,
			Start: &model.Point{
				X: start.X,
				Y: start.Y,
			},
			Target: &model.Point{
				X: target.X,
				Y: target.Y,
			},
			Path: path,
		})
	}
	return result, nil
}

func (r *queryResolver) Board(ctx context.Context, id string) (*model.Board, error) {
	var result *model.Board
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	dbResult := db_model.Get(ID)
	if dbResult == nil {
		return nil, errors.New("not found")
	}
	var start model.Point
	var target model.Point
	var path []*model.Point
	_ = json.Unmarshal([]byte(dbResult.Start), &start)
	_ = json.Unmarshal([]byte(dbResult.Target), &target)
	_ = json.Unmarshal([]byte(dbResult.Path), &path)
	result = &model.Board{
		ID:   dbResult.ID,
		Size: dbResult.Size,
		Start: &model.Point{
			X: start.X,
			Y: start.Y,
		},
		Target: &model.Point{
			X: target.X,
			Y: target.Y,
		},
		Path: path,
	}
	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
