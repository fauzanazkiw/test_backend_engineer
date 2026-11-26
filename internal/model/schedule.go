package model

import "time"

// Schedule merepresentasikan jadwal tayang sebuah film pada cabang dan studio tertentu.
type Schedule struct {
	ID          int64     `json:"id"`
	MovieID     int64     `json:"movie_id"`
	BranchID    int64     `json:"branch_id"`
	StudioID    int64     `json:"studio_id"`
	ShowDate    string    `json:"show_date"`
	ShowTime    string    `json:"show_time"`
	TicketPrice int64     `json:"ticket_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ScheduleInput adalah data yang sudah tervalidasi dan siap disimpan.
type ScheduleInput struct {
	MovieID     int64
	BranchID    int64
	StudioID    int64
	ShowDate    string
	ShowTime    string
	TicketPrice int64
}

// ScheduleFilter dipakai untuk memfilter daftar jadwal tayang.
type ScheduleFilter struct {
	MovieID  *int64
	BranchID *int64
	StudioID *int64
	ShowDate string
}
