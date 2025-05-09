package models

import "gorm.io/gorm"

type Task struct {
    gorm.Model
    Title     string
    Completed bool
    UserID    uint `gorm:"index"`
}
type TaskPreview struct {
    ID        uint   `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func (TaskPreview) TableName() string {
    return "tasks"
}