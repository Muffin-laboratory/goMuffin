package customid

import "fmt"

const (
	PaginationContainerFirst   = "/muffin-pages/first"
	PaginationContainerPrev    = "/muffin-pages/prev"
	PaginationContainerPages   = "/muffin-pages/pages"
	PaginationContainerNext    = "/muffin-pages/next"
	PaginationContainerLast    = "/muffin-pages/last"
	PaginationContainerModal   = "/muffin-pages/modal"
	PaginationContainerSetPage = "/muffin-pages/modal/set"
)

func MakePaginationContainerPrev(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerPrev, id)
}

func MakePaginationContainerFirst(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerFirst, id)
}

func MakePaginationContainerPages(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerPages, id)
}

func MakePaginationContainerNext(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerNext, id)
}

func MakePaginationContainerLast(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerLast, id)
}

func MakePaginationContainerModal(id string) string {
	return fmt.Sprintf("%s/%s", PaginationContainerModal, id)
}
