package mysql

import (
	"context"

	v1 "github.com/ThinkHao/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/fields"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	gorm "gorm.io/gorm"

	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/internal/pkg/util/gormutil"
)

type ipgroups struct {
	db *gorm.DB
}

func newIPGroups(ds *datastore) *ipgroups {
	return &ipgroups{db: ds.db}
}

// Create creates a new ipgroup.
func (i *ipgroups) Create(ctx context.Context, ipgroup *v1.IPGroup, opts metav1.CreateOptions) error {
	return i.db.Create(&ipgroup).Error
}

// Update updates an user account information.
func (i *iproups) Update(ctx context.Context, ipgroup *v1.IPGroup, opts metav1.UpdateOptions) error {
	return i.db.Save(ipgroup).Error
}

// Delete deletes the user by the user identifier.
// func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
// 	// delete related policy first
// 	pol := newPolicies(&datastore{u.db})
// 	if err := pol.DeleteByUser(ctx, username, opts); err != nil {
// 		return err
// 	}

// 	if opts.Unscoped {
// 		u.db = u.db.Unscoped()
// 	}

// 	err := u.db.Where("name = ?", username).Delete(&v1.User{}).Error
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return errors.WithCode(code.ErrDatabase, err.Error())
// 	}

// 	return nil
// }

// // DeleteCollection batch deletes the users.
// func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
// 	// delete related policy first
// 	pol := newPolicies(&datastore{u.db})
// 	if err := pol.DeleteCollectionByUser(ctx, usernames, opts); err != nil {
// 		return err
// 	}

// 	if opts.Unscoped {
// 		u.db = u.db.Unscoped()
// 	}

// 	return u.db.Where("name in (?)", usernames).Delete(&v1.User{}).Error
// }

// Get return an user by the user identifier.
func (i *ipgroups) Get(ctx context.Context, ipgroup_name string, opts metav1.GetOptions) (*v1.IPGroup, error) {
	ipgroup := &v1.IPGroup{}
	err := i.db.Where("ipgroup_name = ?", ipgroup_name).First(&ipgroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return ipgroup, nil
}

// List return all users.
func (i *ipgroups) List(ctx context.Context, opts metav1.ListOptions) (*v1.IPGroupList, error) {
	ret := &v1.IPGroupList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	ipgroupname, _ := selector.RequiresExactMatch("ipgroup_name")
	d := u.db.Where("ipgroup_name like ?", "%"+ipgroupname+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

// ListOptional show a more graceful query method.
// func (u *users) ListOptional(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
// 	ret := &v1.UserList{}
// 	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

// 	where := v1.User{}
// 	whereNot := v1.User{
// 		IsAdmin: 0,
// 	}
// 	selector, _ := fields.ParseSelector(opts.FieldSelector)
// 	username, found := selector.RequiresExactMatch("name")
// 	if found {
// 		where.Name = username
// 	}

// 	d := u.db.Where(where).
// 		Not(whereNot).
// 		Offset(ol.Offset).
// 		Limit(ol.Limit).
// 		Order("id desc").
// 		Find(&ret.Items).
// 		Offset(-1).
// 		Limit(-1).
// 		Count(&ret.TotalCount)

// 	return ret, d.Error
// }