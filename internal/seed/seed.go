package seed

import (
	"context"
	"fmt"

	"github.com/mohdhujaifa/profile/internal/model"
	"github.com/mohdhujaifa/profile/internal/repository"
)

type Store interface {
	HasProfile(ctx context.Context) (bool, error)
	InsertProfile(ctx context.Context, p model.Profile) error
	InsertSkill(ctx context.Context, category, name string, order int) error
	InsertExperience(ctx context.Context, e model.Experience) (int64, error)
	InsertProject(ctx context.Context, p model.Project) error
	InsertEducation(ctx context.Context, e model.Education) error
}

func ResumeIfEmpty(ctx context.Context, store Store) error {
	ok, err := store.HasProfile(ctx)
	if err != nil {
		return fmt.Errorf("seed check: %w", err)
	}
	if ok {
		return nil
	}
	return Apply(ctx, store)
}

func Apply(ctx context.Context, store Store) error {
	profile := model.Profile{
		FullName:    "Mohd Hujaifa",
		Title:       "Software Engineer · Golang · Distributed Systems",
		Email:       "huzaifaanis40@gmail.com",
		Phone:       "8173846711",
		Location:    "Allahabad",
		LinkedInURL: "https://www.linkedin.com/hujaifademocracy/",
		GitHubURL:   "https://github.com/Sayyed98/",
		Summary:     "Software Engineer with 3+ years of experience specializing in building scalable backend systems with Golang. Optimized MySQL queries and implemented Redis caching, significantly improving API response times by up to 40% and reducing database load for an e-commerce platform. Consistently delivered robust microservices, designed RESTful APIs with stringent validation, and engineered asynchronous workflows using RabbitMQ. Focused on leveraging strong concurrent programming skills to develop high-performance, maintainable backend solutions.",
	}
	if err := store.InsertProfile(ctx, profile); err != nil {
		return err
	}

	skills := []struct {
		category string
		items    []string
	}{
		{"Programming & Core Technologies", []string{"Golang", "Goroutines", "Channels", "Worker Pools", "Data Structures & Algorithms", "Linux"}},
		{"Web Services & APIs", []string{"Gin", "RESTful API Design", "Microservices", "gRPC", "Webhooks"}},
		{"Databases & Caching", []string{"MySQL", "PostgreSQL", "SQL Optimization", "Redis"}},
		{"Cloud & Containerization", []string{"Docker", "Kubernetes", "AWS", "S3"}},
		{"Messaging & Version Control", []string{"RabbitMQ", "Git"}},
		{"Development Practices & Testing", []string{"Agile Development", "Testify", "Go Testing", "Performance Optimization", "Code Refactoring", "Code Reviews", "SMTP-based email delivery"}},
	}
	order := 0
	for _, g := range skills {
		for _, item := range g.items {
			if err := store.InsertSkill(ctx, g.category, item, order); err != nil {
				return err
			}
			order++
		}
	}

	experiences := []model.Experience{
		{
			Company:   "Infogain India Pvt Ltd",
			Client:    "HP Inc.",
			Role:      "Software Engineer",
			Location:  "Noida",
			StartDate: repository.Date(2025, 6, 1),
			EndDate:   repository.DatePtr(2025, 12, 1),
			Bullets: []string{
				"Developing and maintaining scalable back-end services using Golang for HP’s internal CPQ (Configure, Price, Quote) platform.",
				"Collaborated with cross-functional teams to design and deliver RESTful APIs and microservices for enterprise-level applications.",
				"Worked on the Integrated Quote (IQ) module, implementing functionalities such as Get Price, Accept Price, and Price Escalation to support additional discount workflows.",
				"Implemented logic to automatically attach matching products based on quote configurations.",
				"Performed performance optimization, code refactoring, and code reviews to improve system efficiency, reliability, and maintainability.",
			},
		},
		{
			Company:   "Xtrox Technology Pvt Ltd",
			Client:    "",
			Role:      "Software Engineer",
			Location:  "Lucknow",
			StartDate: repository.Date(2024, 7, 1),
			EndDate:   repository.DatePtr(2025, 2, 1),
			Bullets: []string{
				"Collaborated with US-based teams on the CleanClaims.com platform, contributing to backend development and system optimization.",
				"Optimized Golang services to improve performance, scalability, and maintainability.",
				"Implemented system action logging and automation macros to enhance auditability and operational workflows.",
				"Built backend tools for report generation, SMTP-based email delivery, and project history tracking.",
				"Designed and implemented admin tools for sensor management, stage unlocking, and Multi-Org franchise provisioning.",
				"Developed an Estimate Tool, enabling scope definition, cost estimation, and project conversion with robust version control and approval processes.",
				"Developed hybrid REST + gRPC microservices for low-latency inter-service communication.",
			},
		},
		{
			Company:   "Scalent Infotech Pvt Ltd",
			Client:    "",
			Role:      "Software Engineer",
			Location:  "Pune",
			StartDate: repository.Date(2022, 7, 1),
			EndDate:   repository.DatePtr(2024, 6, 1),
			Bullets: []string{
				"Collaborated with cross-functional teams to design, develop, and maintain microservices for an e-commerce website with seven microservices.",
				"Owned the Product Microservice, implementing core features and improving scalability.",
				"Optimized MySQL queries and database schema, improving API response time by up to 40%.",
				"Implemented Redis caching to reduce response time and database load.",
				"Built RabbitMQ-based asynchronous workflows for background jobs and event-driven processing.",
				"Designed and maintained RESTful APIs with proper validation and error handling.",
				"Collaborated with the front-end team to integrate APIs seamlessly into the user interface.",
				"Actively participated in agile development processes, including sprint planning, daily stand-ups, and retrospectives.",
			},
		},
	}
	for _, e := range experiences {
		if _, err := store.InsertExperience(ctx, e); err != nil {
			return err
		}
	}

	projects := []model.Project{
		{
			Name:    "CORE",
			Summary: "Recruitment, employee management, and project tracking platform modules that improved organizational efficiency.",
			Bullets: []string{
				"Implemented critical modules that enhanced recruitment, employee management, and project tracking processes.",
				"These improvements contributed to overall efficiency and productivity of the company.",
			},
		},
		{
			Name:    "CleanClaims.com",
			Summary: "Insurance claims operations platform with audit logging, admin tooling, and high-throughput APIs.",
			Bullets: []string{
				"Implemented system action logging to enhance project tracking.",
				"Developed and maintained project-specific system macros.",
				"Designed tools for report generation, SMTP-based email, and history tracking.",
				"Built a sensor management interface for superadmins, enabling cross-organization transfers.",
				"Enhanced the franchise setup system, streamlining Multi-Org Admin onboarding.",
				"Implemented goroutine-based concurrency and worker pools to handle high-throughput API requests.",
			},
		},
		{
			Name:    "E-commerce Website",
			Summary: "High-volume e-commerce platform backed by seven microservices, MySQL, Redis, and RabbitMQ.",
			Bullets: []string{
				"Enriched platform functionality while keeping the system robust and scalable.",
				"Delivered an e-commerce platform capable of handling a high volume of users and transactions.",
			},
		},
	}
	for _, p := range projects {
		if err := store.InsertProject(ctx, p); err != nil {
			return err
		}
	}

	education := []model.Education{
		{
			Institution: "Sam Higginbottom University of Agriculture, Technology & Sciences",
			Degree:      "Master of Computer Application (M.C.A)",
			Location:    "Naini, Allahabad",
			StartDate:   repository.Date(2018, 1, 1),
			EndDate:     repository.DatePtr(2021, 12, 1),
		},
		{
			Institution: "C.S.J.M. University Kanpur",
			Degree:      "Bachelor of Science",
			Location:    "Allahabad",
			StartDate:   repository.Date(2015, 1, 1),
			EndDate:     repository.DatePtr(2018, 12, 1),
		},
	}
	for _, e := range education {
		if err := store.InsertEducation(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
