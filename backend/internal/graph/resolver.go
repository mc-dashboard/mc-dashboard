package graph

import "github.com/rohanvsuri/minecraft-dashboard/internal/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	todo  *model.Todo
	todos []*model.Todo
}
