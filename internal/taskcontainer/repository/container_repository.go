package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/model"
)

const dbTimeout = time.Second * 5

type ContainerRepository interface {
	AllTaskContainers() ([]*model.TaskContainer, error)
	GetById(id string) (*model.TaskContainer, error)
	GetContainersByGroupId(groupId int) ([]model.TaskContainer, error)
	CreateContainer(container model.TaskContainer) error
	DeleteContainer(id string) error
	RemoveContainerByUsergroupId(groupId int) error
}

type ContainerRepo struct {
	DB *sql.DB
}

func NewContainerRepository(db *sql.DB) *ContainerRepo {
	return &ContainerRepo{
		DB: db,
	}
}

func (m *ContainerRepo) AllTaskContainers() ([]*model.TaskContainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetAllContainers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []*model.TaskContainer
	for rows.Next() {
		container, err := scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}

		containers = append(containers, container)
	}
	return containers, nil
}

func (m *ContainerRepo) GetById(id string) (*model.TaskContainer, error) {
	rows, err := m.DB.Query(sqlGetById, id)
	if err != nil {
		return nil, err
	}

	container := new(model.TaskContainer)
	for rows.Next() {
		container, err = scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}
	}
	return container, err
}

func (m *ContainerRepo) GetContainersByGroupId(groupId int) ([]model.TaskContainer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetContainersByGroupId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	containers := []model.TaskContainer{}
	for rows.Next() {
		container, err := scanRowsIntoContainer(rows)
		if err != nil {
			return nil, err
		}

		containers = append(containers, *container)
	}
	return containers, nil
}

func (m *ContainerRepo) CreateContainer(c model.TaskContainer) error {
	_, err := m.DB.Exec(sqlCreateContainer, c.Id, c.Name, c.Description, c.IsActive, c.Activity_level, c.Type, c.UsergroupId)
	if err != nil {
		return fmt.Errorf("unable to insert into taskcontainer table : %w", err)
	}
	return nil
}

func (m *ContainerRepo) DeleteContainer(id string) error {
	_, err := m.DB.Exec(sqlDeleteContainer, id)
	if err != nil {
		return fmt.Errorf("unable to remove task container : %w", err)
	}
	return nil
}

func (m *ContainerRepo) RemoveContainerByUsergroupId(groupId int) error {
	_, err := m.DB.Exec(sqlDeleteContainerByUsergroupId, groupId)
	if err != nil {
		return fmt.Errorf("unable to remove container by usergroup Id : %w", err)
	}
	return nil
}

func scanRowsIntoContainer(rows *sql.Rows) (*model.TaskContainer, error) {
	container := new(model.TaskContainer)
	err := rows.Scan(
		&container.Id,
		&container.Name,
		&container.Description,
		&container.IsActive,
		&container.UsergroupId,
	)
	if err != nil {
		return nil, err
	}

	return container, nil
}
