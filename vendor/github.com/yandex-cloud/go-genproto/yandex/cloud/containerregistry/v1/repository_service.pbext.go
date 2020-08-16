// Code generated by protoc-gen-goext. DO NOT EDIT.

package containerregistry

func (m *GetRepositoryRequest) SetRepositoryId(v string) {
	m.RepositoryId = v
}

func (m *GetRepositoryByNameRequest) SetRepositoryName(v string) {
	m.RepositoryName = v
}

func (m *ListRepositoriesRequest) SetRegistryId(v string) {
	m.RegistryId = v
}

func (m *ListRepositoriesRequest) SetFolderId(v string) {
	m.FolderId = v
}

func (m *ListRepositoriesRequest) SetPageSize(v int64) {
	m.PageSize = v
}

func (m *ListRepositoriesRequest) SetPageToken(v string) {
	m.PageToken = v
}

func (m *ListRepositoriesRequest) SetFilter(v string) {
	m.Filter = v
}

func (m *ListRepositoriesRequest) SetOrderBy(v string) {
	m.OrderBy = v
}

func (m *ListRepositoriesResponse) SetRepositories(v []*Repository) {
	m.Repositories = v
}

func (m *ListRepositoriesResponse) SetNextPageToken(v string) {
	m.NextPageToken = v
}
