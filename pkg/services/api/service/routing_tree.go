package service

import (
	"errors"
	"lab_2/internal/base"
	"lab_2/pkg/services/api/models"
)

var (
	ErrPortNotFound = errors.New("service with specified port not found")
	ErrNameNotFound = errors.New("service with specified name not found")
)

type ServiceRoutingTree struct {
	base.RoutingTree[models.ProxyService]
}

func NewServiceRoutingTree() ServiceRoutingTree {
	return ServiceRoutingTree{
		RoutingTree: base.NewRoutingTree[models.ProxyService](),
	}
}

func (u *ServiceRoutingTree) FindByPort(port uint16) (value models.ProxyService, err error) {
	return findByPort(&u.RoutingTree, port)
}

func (u *ServiceRoutingTree) FindByName(name string) (value models.ProxyService, err error) {
	return findByName(&u.RoutingTree, name)
}

func (u *ServiceRoutingTree) Add(uri string, value models.ProxyService) (err error) {
	uriSlice := base.SplitURLPath(removeLeadingSlash(uri))
	return u.RoutingTree.Add(uriSlice, value)
}

func (u *ServiceRoutingTree) Find(uri string) (models.ProxyService, error) {
	uriSlice := base.SplitURLPath(removeLeadingSlash(uri))
	return u.RoutingTree.Find(uriSlice, true, func(node *base.RoutingTree[models.ProxyService], value *models.ProxyService, key string) (newNode *base.RoutingTree[models.ProxyService], found bool) {
		if newNode, found = node.Branches[key]; found {
			if newNode.Item != nil {
				*value = *newNode.Item
				return
			}
		}
		return
	})
}

func findByPort(u *base.RoutingTree[models.ProxyService], port uint16) (value models.ProxyService, err error) {
	for _, node := range u.Branches {
		if node.Item != nil && node.Item.Port == port {
			value = *node.Item
			return
		}
		value, err = findByPort(node, port)
		if err == nil {
			return
		}
	}
	err = ErrPortNotFound
	return
}

func findByName(u *base.RoutingTree[models.ProxyService], name string) (value models.ProxyService, err error) {
	for _, node := range u.Branches {
		if node.Item != nil && node.Item.Name == name {
			value = *node.Item
			return
		}
		value, err = findByName(node, name)
		if err == nil {
			return
		}
	}
	err = ErrNameNotFound
	return
}

func removeLeadingSlash(in string) string {
	if len(in) != 0 && in[0] == '/' {
		return in[1:]
	}
	return in
}
