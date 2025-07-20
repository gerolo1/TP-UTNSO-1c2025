package dtos

type ConfigCPUDTO struct {
	NumberOfLevels int `json:"number_of_levels"`
	PageSize       int `json:"page_size"`
	EntriesPerPage int `json:"entries_per_page"`
}
