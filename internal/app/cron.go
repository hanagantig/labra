package app

import (
	"context"
	"github.com/hanagantig/cron"
	"github.com/hanagantig/gracy"
)

func (a App) RunCron() error {
	myCron := cron.New(cron.WithSeconds())
	_, err := myCron.AddFunc("*/10 * * * * *", "process_new_uploaded_files", func() {
		_ = a.c.GetUseCase().ProcessNewUploadedFiles(context.Background())
	})
	if err != nil {
		return err
	}

	_, err = myCron.AddFunc("*/10 * * * * *", "process_recognized_files", func() {
		_ = a.c.GetUseCase().ProcessRecognizedFiles(context.Background())
	})
	if err != nil {
		return err
	}

	myCron.Start()
	gracy.AddCallback(func() error {
		myCron.Stop()
		return nil
	})

	return nil
}
