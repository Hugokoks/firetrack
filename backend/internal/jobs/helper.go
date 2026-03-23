package jobs

func applyJobUpdates(job *Job, input UpdateJobInput) {
	if input.Title != nil {
		job.Title = *input.Title
	}
	if input.CustomerName != nil {
		job.CustomerName = input.CustomerName
	}
	if input.Address != nil {
		job.Address = *input.Address
	}
	if input.City != nil {
		job.City = input.City
	}
	if input.Zip != nil {
		job.Zip = input.Zip
	}
	if input.Country != nil {
		job.Country = *input.Country
	}
	if input.Latitude != nil {
		job.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		job.Longitude = input.Longitude
	}
	if input.ScheduledStart != nil {
		job.ScheduledStart = *input.ScheduledStart
	}
	if input.ScheduledEnd != nil {
		job.ScheduledEnd = input.ScheduledEnd
	}
	if input.CompletedAt != nil {
		job.CompletedAt = input.CompletedAt
	}
	if input.Status != nil {
		job.Status = *input.Status
	}
	if input.Priority != nil {
		job.Priority = *input.Priority
	}
	if input.AssignedUserID != nil {
		job.AssignedUserID = input.AssignedUserID
	}
	if input.Description != nil {
		job.Description = input.Description
	}
}
