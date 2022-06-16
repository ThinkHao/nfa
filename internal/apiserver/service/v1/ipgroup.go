package v1

import (
	"context"
	"sync"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/pkg/log"
)

// UserSrv defines functions used to handle user request.
type IPGroupSrv interface {
	Create(ctx context.Context, ipgroup *v1.IPGroup, opts metav1.CreateOptions) error
	Update(ctx context.Context, ipgroup *v1.IPGroup, opts metav1.UpdateOptions) error
	Get(ctx context.Context, ipgroup_name string, opts metav1.GetOptions) (*v1.IPGroup, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.IPGroupList, error)
}

type ipgroupService struct {
	store store.Factory
}

var _ IPGroupSrv = (*ipgroupService)(nil)

func newIPGroups(srv *service) *ipgroupService {
	return &ipgroupService{store: srv.store}
}

// List returns user list in the storage. This function has a good performance.
func (i *ipgroupService) List(ctx context.Context, opts metav1.ListOptions) (*v1.IPGroupList, error) {
	ipgroups, err := u.store.IPGroups().List(ctx, opts)
	if err != nil {
		log.L(ctx).Errorf("list ipgroups from storage failed: %s", err.Error())

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	var m sync.Map

	// Improve query efficiency in parallel
	for _, ipgroup := range ipgroups.Items {
		wg.Add(1)

		go func(ipgroup *v1.IPGroup) {
			defer wg.Done()

			// some cost time process
			policies, err := i.store.Policies().List(ctx, user.Name, metav1.ListOptions{})
			if err != nil {
				errChan <- errors.WithCode(code.ErrDatabase, err.Error())

				return
			}

			m.Store(ipgroup.ID, &v1.IPGroup{
				ObjectMeta: metav1.ObjectMeta{
					IPGroupID: ipgroup.ipgroup_id,
					IPGroupName:       ipgroup.ipgroup_name,
				},
			})
		}(ipgroup)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	// infos := make([]*v1.User, 0)
	infos := make([]*v1.IPGroup, 0, len(ipgroups.Items))
	for _, ipgroup := range ipgroups.Items {
		info, _ := m.Load(ipgroup.ipgroup_id)
		infos = append(infos, info.(*v1.IPGroup))
	}

	log.L(ctx).Debugf("get %d ipgroups from backend storage.", len(infos))

	return &v1.IPGroupList{ListMeta: ipgroups.ListMeta, Items: infos}, nil
}

// ListWithBadPerformance returns user list in the storage. This function has a bad performance.
// func (u *userService) ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
// 	users, err := u.store.Users().List(ctx, opts)
// 	if err != nil {
// 		return nil, errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	infos := make([]*v1.User, 0)
// 	for _, user := range users.Items {
// 		policies, err := u.store.Policies().List(ctx, user.Name, metav1.ListOptions{})
// 		if err != nil {
// 			return nil, errors.WithCode(code.ErrDatabase, err.Error())
// 		}

// 		infos = append(infos, &v1.User{
// 			ObjectMeta: metav1.ObjectMeta{
// 				ID:        user.ID,
// 				Name:      user.Name,
// 				CreatedAt: user.CreatedAt,
// 				UpdatedAt: user.UpdatedAt,
// 			},
// 			Nickname:    user.Nickname,
// 			Email:       user.Email,
// 			Phone:       user.Phone,
// 			TotalPolicy: policies.TotalCount,
// 		})
// 	}

// 	return &v1.UserList{ListMeta: users.ListMeta, Items: infos}, nil
// }

// func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
// 	if err := u.store.Users().Create(ctx, user, opts); err != nil {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	return nil
// }

// func (u *userService) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
// 	if err := u.store.Users().DeleteCollection(ctx, usernames, opts); err != nil {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	return nil
// }

// func (u *userService) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
// 	if err := u.store.Users().Delete(ctx, username, opts); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (i *ipgroupService) Get(ctx context.Context, ipgroup_name string, opts metav1.GetOptions) (*v1.IPGroup, error) {
	ipgroup, err := u.store.IPGroups().Get(ctx, ipgroup_name, opts)
	if err != nil {
		return nil, err
	}

	return ipgroup, nil
}

// func (u *userService) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
// 	if err := u.store.Users().Update(ctx, user, opts); err != nil {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	return nil
// }

// func (u *userService) ChangePassword(ctx context.Context, user *v1.User) error {
// 	// Save changed fields.
// 	if err := u.store.Users().Update(ctx, user, metav1.UpdateOptions{}); err != nil {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	return nil
// }
