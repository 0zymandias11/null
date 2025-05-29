package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"example.com/Go_Land/internal/env/store"
)

var usernames = []string{
	"Jake", "Barry", "Alice", "Sophia", "Liam", "Noah", "Olivia", "Emma", "Mason", "Logan",
	"Ava", "Ethan", "Lucas", "Mia", "Harper", "Evelyn", "James", "Benjamin", "Amelia", "Ella",
	"Henry", "Sebastian", "Chloe", "Aiden", "Alexander", "Grace", "Nora", "Lily", "Daniel", "Michael",
	"Jackson", "Levi", "Zoe", "Victoria", "Elijah", "Samuel", "Aria", "Layla", "David", "Carter",
	"Wyatt", "Jayden", "Isabella", "Scarlett", "Julian", "Matthew", "Hazel", "Joseph", "Luke", "Gabriel",
}

var titles = []string{
	"Biking is fun",
	"10 VSCode extensions you must try",
	"Ram Ahuja randi ka pilla",
	"Anukriti chakko ka loda leti",
	"Here's how to learn GraphQL",
	"Why you should worship Slaanesh",
	"Top 10 Go tips for backend developers",
	"Mastering Kubernetes in 7 days",
	"Is AI really taking your job?",
	"How to optimize MongoDB queries",
	"The truth about startup burnout",
	"Why coffee is a developer’s best friend",
	"An introvert's guide to tech conferences",
	"How to build a game engine in Go",
	"Prasanna Madarchod !!",
	"Learning Rust: Day 1",
	"Why your team hates daily standups",
	"Productivity hacks no one talks about",
	"Is TypeScript worth the hype?",
	"Building a GraphQL server from scratch",
	"How to handle imposter syndrome",
	"My year without Google",
	"Deploying Go apps with Docker",
	"Why I love writing technical blogs",
	"The rise and fall of PHP",
	"What's wrong with agile today",
	"Backend devs need to learn frontend",
	"Learning Vim ruined my life (and I love it)",
	"Why Python is not always the answer",
	"Machine learning without a PhD",
	"Open source made me a better dev",
	"Is ChatGPT replacing software engineers?",
	"Tech interviews are broken. Here's why.",
	"100 days of code: what I learned",
	"How I built a SaaS in 30 days",
	"Is remote work here to stay?",
	"My favorite CS books of all time",
	"Serverless: pros, cons, and myths",
	"How I failed my first startup",
	"Speeding up Go with channels",
	"Building CLI tools with Cobra",
	"Dark mode is not just aesthetics",
	"Becoming a DevOps engineer",
	"Go vs Java: A backend comparison",
	"The art of writing clean APIs",
	"How to name things in code",
	"10 underrated Git commands",
	"The case for digital minimalism",
	"Why I switched to Arch Linux",
	"Making peace with legacy code",
	"Building hobby projects that last",
	"The productivity trap in tech",
	"Pair programming: blessing or curse?",
	"Why documentation is underrated",
	"How to run effective code reviews",
	"Contributing to open source the right way",
	"Writing readable code in Go",
	"Why design patterns still matter",
	"Life after the bootcamp",
	"What recruiters don’t tell you",
	"My favorite terminal tools",
	"How to manage side projects",
	"Debugging: from painful to powerful",
	"What I learned teaching code",
	"Why monoliths are underrated",
	"Understanding OAuth2 step-by-step",
	"Why I stopped using Redux",
	"Data structures I actually used at work",
	"How to avoid overengineering",
	"Building a blog with Hugo",
	"Using Go to build REST APIs",
	"How to make tech meetups less awkward",
	"Learning system design as a junior dev",
	"Why static typing saves lives",
	"Designing with empathy in mind",
	"My favorite dev podcasts",
	"How to think in SQL",
	"Why burnout happens silently",
	"Making your resume stand out",
	"Getting better at writing docs",
	"How to become a 10x debugger",
	"Why tests are your best friend",
	"Go modules: the good, bad, and ugly",
	"How to write better commit messages",
	"The power of community in tech",
	"Practicing deep work as a dev",
	"Building a RESTful API in Go",
	"Why coding interviews feel broken",
	"How I learned Go by building tools",
	"Handling failure in production",
	"Writing fast code in Go",
	"Tips for reading open source code",
	"Why I started self-hosting everything",
	"The hard parts of mentorship",
	"Learning GraphQL the hard way",
	"Managing imposter syndrome in tech",
	"Being a junior in a senior’s world",
	"What's next for web development?",
	"How to speak at your first tech event",
	"Using Go for automation scripts",
	"What I wish I knew before freelancing",
	"Making your CLI tools delightful",
	"GraphQL vs REST: what to choose?",
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	log.Printf("Starting seeding process...")
	rand.Seed(42) // Seed the random number generator for reproducibility
	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("Error creating user: %v", err)
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("Error creating Post: %v", err)
		}
	}

	comments := generateComments(500, posts, users)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("Error creating Comment: %v", err)
		}
	}

	log.Printf("Seeding completed successfully with %d users, %d posts, and %d comments", len(users), len(posts), len(comments))
	return nil
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "345132adsfg",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		posts[i] = &store.Post{
			UserID:  users[rand.Intn(len(users))].ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: fmt.Sprintf("This is the content of post %d", i+1),
			Tags: []string{
				titles[rand.Intn(len(titles))],
				titles[rand.Intn(len(titles))],
			},
		}
	}
	return posts
}

func generateComments(num int, posts []*store.Post, user []*store.User) []*store.Comment {
	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  user[rand.Intn(len(user))].ID,
			Content: fmt.Sprintf("This is comment %d", i+1),
		}
	}
	return comments
}
