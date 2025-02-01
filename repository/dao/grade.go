package dao

import (
	"context"
	"github.com/asynccnu/be-grade/repository/model"
	"gorm.io/gorm"
)

// GradeDAO 数据库操作的集合
type GradeDAO interface {
	FirstOrCreate(ctx context.Context, grade *model.Grade) error
	FindGrades(ctx context.Context, studentId string, Xnm int64, Xqm int64) ([]model.Grade, error)
	BatchInsertOrUpdate(ctx context.Context, grades []model.Grade) (updateGrade []model.Grade, err error)
}

type gradeDAO struct {
	db *gorm.DB
}

// NewDatabaseStruct  构建数据库操作实例
func NewGradeDAO(db *gorm.DB) GradeDAO {
	return &gradeDAO{db: db}
}

// FirstOrCreate 会自动查找是否存在记录,如果不存在则会存储
func (d *gradeDAO) FirstOrCreate(ctx context.Context, grade *model.Grade) error {
	return d.db.WithContext(ctx).Where("student_id = ? AND jxb_id = ?", grade.Studentid, grade.JxbId).FirstOrCreate(grade).Error
}

// FindAllGradesByStudentId 搜索成绩,xnm(学年名),xqm(学期名)条件为可选
func (d *gradeDAO) FindGrades(ctx context.Context, studentId string, Xnm int64, Xqm int64) ([]model.Grade, error) {
	// 定义查询结果的容器
	var grades []model.Grade

	// 构建查询
	query := d.db.WithContext(ctx).Where("student_id = ?", studentId)
	if Xnm != 0 { // 如果 Xnm 有值，拼接学年条件
		query = query.Where("xnm = ?", Xnm)
	}

	if Xqm != 0 { // 如果 Xqm 有值，拼接学期条件
		query = query.Where("xqm = ?", Xqm)
	}

	// 执行查询
	err := query.Find(&grades).Error
	if err != nil {
		return nil, err
	}

	return grades, nil
}

// 批量更新(被AI战胜了,我想不出这么优秀的思路)
func (d *gradeDAO) BatchInsertOrUpdate(ctx context.Context, grades []model.Grade) (updateGrade []model.Grade, err error) {

	// 提取所有 student_id 和 jxb_id
	ids := make([]string, len(grades))
	for i, grade := range grades {
		ids[i] = grade.Studentid + grade.JxbId
	}

	// 查询数据库中已有的记录
	var existingGrades []model.Grade
	if err := d.db.WithContext(ctx).
		Where("CONCAT(student_id, jxb_id) IN ?", ids).
		Find(&existingGrades).Error; err != nil {
		return nil, err
	}

	// 用 map 比对现有数据和插入数据
	existingMap := make(map[string]model.Grade)
	for _, grade := range existingGrades {
		existingMap[grade.Studentid+grade.JxbId] = grade
	}

	// 找出需要插入或更新的数据
	var toInsert []model.Grade
	for _, grade := range grades {
		key := grade.Studentid + grade.JxbId
		if _, exists := existingMap[key]; !exists {
			toInsert = append(toInsert, grade)
		}
	}

	// 批量插入需要新增的记录
	if len(toInsert) > 0 {
		if err := d.db.WithContext(ctx).Create(&toInsert).Error; err != nil {
			return nil, err
		}
	}

	return toInsert, nil
}
