package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/TaushifReza/go-social/internal/model"
	"github.com/TaushifReza/go-social/internal/store"
)

var username = []string{
	"james.wilson", "emma_thompson", "mgarcia88", "liam.smith.92", "sophia.chen",
	"noah_walker", "olivia_perez", "william_brown", "ava.martinez", "lucas.jones",
	"isabella_davis", "mason_moore", "mia_clark", "ethan_lewis", "chloe.young",
	"alex.rodriguez", "sofia_hill", "jacob.scott", "hannah_b", "logan.green",
	"zoe_adams", "jackson_white", "lily.king", "ryan_harris", "madison_nelson",
	"c_campbell", "aiden_mitchell", "evelyn_roberts", "matthew_hall", "abigail.lee",
	"dbaker_77", "sam_turner", "mila_phillips", "david_evans", "scarlett_wright",
	"joseph.torres", "aria_parker", "samuel_collins", "penelope_edwards", "sebastian_stewart", "nora_flores", "carter_morris", "hazel_nguyen", "wyatt_murphy", "aubrey_rivera",
	"julian_cook", "stella_rogers", "isaac_morgan", "natalie_peterson", "henry_gray",
}

var titles = []string{
	"The Future of AI", "Mastering the Terminal", "Digital Nomad Life",
	"Cooking Made Simple", "Morning Rituals", "Minimalist Living",
	"Remote Work Tips", "The Power of Habit", "Beginner's Guide to Go",
	"Travel on a Budget", "Home Office Setup", "Healthy Eating 101",
	"Understanding Crypto", "Modern Web Design", "Productivity Hacks",
	"Mental Health Matters", "Sustainable Fashion", "Fitness for Busy People",
	"Investing for Beginners", "The Joy of Reading",
}

var contents = []string{
	"Exploring how machine learning is reshaping our daily lives and industries.",
	"A deep dive into essential commands to boost your command-line efficiency.",
	"How to maintain a career while traveling the world on your own terms.",
	"Easy recipes and time-saving techniques for the busy home cook.",
	"Small changes to your morning that lead to a more productive day.",
	"The art of decluttering your space and mind for a focused life.",
	"Creating a healthy work-life balance while working from your living room.",
	"How tiny, consistent actions lead to massive long-term transformations.",
	"Everything you need to know to write your first program in Go.",
	"Pro tips for seeing the world without breaking your bank account.",
	"Designing an ergonomic and inspiring workspace in any corner of your home.",
	"Practical advice for transitioning to a whole-foods-based lifestyle.",
	"Decoding the basics of blockchain and what it means for the future.",
	"The latest trends in UI/UX that make websites both beautiful and functional.",
	"Tools and techniques to get more done in less time with less stress.",
	"Simple strategies to prioritize your emotional well-being every day.",
	"Why choosing quality over quantity is better for you and the planet.",
	"Quick workout routines designed for those with a packed schedule.",
	"A beginner-friendly guide to building wealth through smart stock choices.",
	"Rediscovering the mental and emotional benefits of getting lost in a book.",
}

var tags = []string{
	"AI", "CLI", "Lifestyle", "Cooking", "Wellness",
	"Minimalism", "WorkFromHome", "PersonalGrowth", "Golang", "Travel",
	"OfficeSetup", "Nutrition", "Finance", "WebDesign", "Productivity",
	"MentalHealth", "EcoFriendly", "Fitness", "Investing", "Books",
}

var comments = []string{
	"Really eye-opening look at where AI is headed!",
	"I finally understand how to use grep properly now.",
	"The dream! Which country are you visiting next?",
	"Tried the 15-minute pasta, it was actually delicious.",
	"Waking up at 5 AM changed my whole week, thanks!",
	"I just threw out three bags of clutter. Feeling lighter.",
	"The tip about noise-canceling headphones is a lifesaver.",
	"Atomic habits are the only way I've stayed consistent.",
	"Great tutorial, the syntax for slices finally clicked.",
	"I never thought about using credit card points that way.",
	"Where did you get that wooden desk? It looks amazing.",
	"Meal prepping on Sundays is the only way I eat healthy.",
	"Can you explain the difference between coins and tokens?",
	"Dark mode is definitely the way to go for 2026 design.",
	"The Pomodoro technique never fails me during finals.",
	"It's so important to talk about burnout, thank you.",
	"Thrifting is such a fun way to find unique pieces.",
	"Do you have a version of this workout for beginners?",
	"Index funds are definitely the safest bet for me.",
	"Just added three of these to my Kindle reading list!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	// 1. Generate and Create Users
	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			fmt.Println("ERROR creating user: ", err)
			return
		}
	}

	// 2. Generate and Create Posts
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("ERROR creating post: ", err)
			return
		}
	}

	// 3. Generate and Create Comments
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("ERROR creating post: ", err)
			return
		}
	}

	fmt.Println("Seed Completed")
}

func generateUsers(num int) []*model.User {
	users := make([]*model.User, num)

	for i := 0; i < num; i++ {
		baseName := username[i%len(username)]
		users[i] = &model.User{
			UserName: fmt.Sprintf("%s_%d", baseName, i),
			Email:    fmt.Sprintf("%s_%d@example.com", baseName, i),
			Password: "Admin123@",
		}
	}

	return users
}

func generatePosts(num int, users []*model.User) []*model.Posts {
	posts := make([]*model.Posts, num)
	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]

		posts[i] = &model.Posts{
			UserID:  user.ID,
			Title:   titles[rand.IntN(len(titles))],
			Content: contents[rand.IntN(len(contents))],
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*model.User, posts []*model.Posts) []*model.Comment {
	cms := make([]*model.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &model.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}

	return cms
}
