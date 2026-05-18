package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

type MemberRepo struct{}

func (r *MemberRepo) FindByID(id uint) (*model.Member, error) {
	var m model.Member
	err := pkg.DB.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MemberRepo) List() ([]model.Member, error) {
	var list []model.Member
	err := pkg.DB.Order("created_at").Find(&list).Error
	return list, err
}

func (r *MemberRepo) Count() (int64, error) {
	var count int64
	err := pkg.DB.Model(&model.Member{}).Count(&count).Error
	return count, err
}

func (r *MemberRepo) CountAdmins() (int64, error) {
	var count int64
	err := pkg.DB.Model(&model.Member{}).Where("role = ?", "admin").Count(&count).Error
	return count, err
}

func (r *MemberRepo) Create(m *model.Member) error {
	return pkg.DB.Create(m).Error
}

func (r *MemberRepo) Update(m *model.Member) error {
	return pkg.DB.Save(m).Error
}

func (r *MemberRepo) SoftDelete(id uint) error {
	return pkg.DB.Delete(&model.Member{}, id).Error
}

func (r *MemberRepo) FindByUsername(username string) (*model.Member, error) {
	var m model.Member
	err := pkg.DB.Where("username = ?", username).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// CountActivityRecords counts all activity records for a member across all modules
func (r *MemberRepo) CountActivityRecords(memberID uint) (int64, error) {
	var total int64

	// Count ledgers created
	var lc int64
	pkg.DB.Model(&model.Ledger{}).Where("creator_id = ?", memberID).Count(&lc)
	total += lc

	// Count ledger_members (shared ledger entries)
	var lm int64
	pkg.DB.Table("ledger_members").Where("member_id = ?", memberID).Count(&lm)
	total += lm

	// Count todos created or assigned
	var tc int64
	pkg.DB.Model(&model.Todo{}).Where("creator_id = ? OR assignee_id = ?", memberID, memberID).Count(&tc)
	total += tc

	// Count todo_logs
	var tlc int64
	pkg.DB.Model(&model.TodoLog{}).Where("operator_id = ?", memberID).Count(&tlc)
	total += tlc

	// Count wishes created
	var wc int64
	pkg.DB.Model(&model.Wish{}).Where("creator_id = ?", memberID).Count(&wc)
	total += wc

	// Count wish_votes
	var wv int64
	pkg.DB.Model(&model.WishVote{}).Where("member_id = ?", memberID).Count(&wv)
	total += wv

	// Count posts created
	var pc int64
	pkg.DB.Model(&model.Post{}).Where("creator_id = ?", memberID).Count(&pc)
	total += pc

	// Count topics created
	var topc int64
	pkg.DB.Model(&model.Topic{}).Where("creator_id = ?", memberID).Count(&topc)
	total += topc

	// Count comments created
	var cc int64
	pkg.DB.Model(&model.Comment{}).Where("creator_id = ?", memberID).Count(&cc)
	total += cc

	// Count votes created
	var vc int64
	pkg.DB.Model(&model.Vote{}).Where("creator_id = ?", memberID).Count(&vc)
	total += vc

	// Count vote_records
	var vr int64
	pkg.DB.Model(&model.VoteRecord{}).Where("member_id = ?", memberID).Count(&vr)
	total += vr

	// Count likes
	var likec int64
	pkg.DB.Model(&model.Like{}).Where("member_id = ?", memberID).Count(&likec)
	total += likec

	return total, nil
}
