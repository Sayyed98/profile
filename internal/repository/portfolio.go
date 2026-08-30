package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mohdhujaifa/profile/internal/model"
)

type PortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) GetPortfolio(ctx context.Context) (model.Portfolio, error) {
	var p model.Portfolio
	profile, err := r.getProfile(ctx)
	if err != nil {
		return p, err
	}
	skills, err := r.getSkills(ctx)
	if err != nil {
		return p, err
	}
	exps, err := r.getExperiences(ctx)
	if err != nil {
		return p, err
	}
	projects, err := r.getProjects(ctx)
	if err != nil {
		return p, err
	}
	edu, err := r.getEducation(ctx)
	if err != nil {
		return p, err
	}
	p.Profile = profile
	p.Skills = skills
	p.Experience = exps
	p.Projects = projects
	p.Education = edu
	return p, nil
}

func (r *PortfolioRepository) HasProfile(ctx context.Context) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM profiles`).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *PortfolioRepository) InsertProfile(ctx context.Context, p model.Profile) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO profiles (full_name, title, email, phone, location, linkedin_url, github_url, summary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.FullName, p.Title, p.Email, p.Phone, p.Location, p.LinkedInURL, p.GitHubURL, p.Summary)
	return err
}

func (r *PortfolioRepository) InsertSkill(ctx context.Context, category, name string, order int) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO skills (category, name, sort_order) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE sort_order = VALUES(sort_order)`, category, name, order)
	return err
}

func (r *PortfolioRepository) InsertExperience(ctx context.Context, e model.Experience) (int64, error) {
	var end any
	if e.EndDate != nil {
		end = e.EndDate.Format("2006-01-02")
	}
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO experiences (company, client, role, location, start_date, end_date, sort_order)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.Company, e.Client, e.Role, e.Location, e.StartDate.Format("2006-01-02"), end, 0)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for i, b := range e.Bullets {
		if _, err := r.db.ExecContext(ctx, `
			INSERT INTO experience_bullets (experience_id, content, sort_order) VALUES (?, ?, ?)`,
			id, b, i); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *PortfolioRepository) InsertProject(ctx context.Context, p model.Project) error {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO projects (name, summary, sort_order) VALUES (?, ?, ?)`, p.Name, p.Summary, 0)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for i, b := range p.Bullets {
		if _, err := r.db.ExecContext(ctx, `
			INSERT INTO project_bullets (project_id, content, sort_order) VALUES (?, ?, ?)`,
			id, b, i); err != nil {
			return err
		}
	}
	return nil
}

func (r *PortfolioRepository) InsertEducation(ctx context.Context, e model.Education) error {
	var end any
	if e.EndDate != nil {
		end = e.EndDate.Format("2006-01-02")
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO education (institution, degree, location, start_date, end_date, sort_order)
		VALUES (?, ?, ?, ?, ?, ?)`,
		e.Institution, e.Degree, e.Location, e.StartDate.Format("2006-01-02"), end, 0)
	return err
}

func (r *PortfolioRepository) InsertContact(ctx context.Context, m model.ContactMessage) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO contact_messages (name, email, subject, body, status) VALUES (?, ?, ?, ?, ?)`,
		m.Name, m.Email, m.Subject, m.Body, m.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PortfolioRepository) UpdateContactStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE contact_messages SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *PortfolioRepository) getProfile(ctx context.Context) (model.Profile, error) {
	var p model.Profile
	err := r.db.QueryRowContext(ctx, `
		SELECT id, full_name, title, email, IFNULL(phone,''), IFNULL(location,''),
		       IFNULL(linkedin_url,''), IFNULL(github_url,''), summary
		FROM profiles ORDER BY id LIMIT 1`).Scan(
		&p.ID, &p.FullName, &p.Title, &p.Email, &p.Phone, &p.Location, &p.LinkedInURL, &p.GitHubURL, &p.Summary)
	if err != nil {
		return p, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

func (r *PortfolioRepository) getSkills(ctx context.Context) ([]model.SkillGroup, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT category, name FROM skills ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	index := map[string]int{}
	var groups []model.SkillGroup
	for rows.Next() {
		var category, name string
		if err := rows.Scan(&category, &name); err != nil {
			return nil, err
		}
		i, ok := index[category]
		if !ok {
			groups = append(groups, model.SkillGroup{Category: category})
			i = len(groups) - 1
			index[category] = i
		}
		groups[i].Items = append(groups[i].Items, name)
	}
	return groups, rows.Err()
}

func (r *PortfolioRepository) getExperiences(ctx context.Context) ([]model.Experience, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, company, IFNULL(client,''), role, IFNULL(location,''), start_date, end_date
		FROM experiences ORDER BY start_date DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Experience
	for rows.Next() {
		var e model.Experience
		var end sql.NullTime
		if err := rows.Scan(&e.ID, &e.Company, &e.Client, &e.Role, &e.Location, &e.StartDate, &end); err != nil {
			return nil, err
		}
		if end.Valid {
			t := end.Time
			e.EndDate = &t
		}
		list = append(list, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range list {
		bullets, err := r.bullets(ctx, `SELECT content FROM experience_bullets WHERE experience_id = ? ORDER BY sort_order, id`, list[i].ID)
		if err != nil {
			return nil, err
		}
		list[i].Bullets = bullets
	}
	return list, nil
}

func (r *PortfolioRepository) getProjects(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, summary FROM projects ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Summary); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range list {
		bullets, err := r.bullets(ctx, `SELECT content FROM project_bullets WHERE project_id = ? ORDER BY sort_order, id`, list[i].ID)
		if err != nil {
			return nil, err
		}
		list[i].Bullets = bullets
	}
	return list, nil
}

func (r *PortfolioRepository) getEducation(ctx context.Context) ([]model.Education, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, institution, degree, IFNULL(location,''), start_date, end_date
		FROM education ORDER BY start_date DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.Education
	for rows.Next() {
		var e model.Education
		var end sql.NullTime
		if err := rows.Scan(&e.ID, &e.Institution, &e.Degree, &e.Location, &e.StartDate, &end); err != nil {
			return nil, err
		}
		if end.Valid {
			t := end.Time
			e.EndDate = &t
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *PortfolioRepository) bullets(ctx context.Context, q string, id int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func DatePtr(year int, month time.Month, day int) *time.Time {
	t := Date(year, month, day)
	return &t
}
