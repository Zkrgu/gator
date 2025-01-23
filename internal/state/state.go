package state

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/zkrgu/gator/internal/config"
	"github.com/zkrgu/gator/internal/database"
)

type State struct {
	Config       *config.Config
	DBConnection *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No arguments supplied\n")
	}
	user, err := s.DBConnection.GetUserByName(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Set user to %s\n", s.Config.CurrentUserName)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No arguments supplied\n")
	}
	user, err := s.DBConnection.CreateUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Created user %s with id %s\n", user.Name, user.ID)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DBConnection.Reset(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DBConnection.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Not enough arguments supplied\n")
	}
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Parsing at interval %v", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough arguments supplied\n")
	}
	user, err := s.DBConnection.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	_, err = s.DBConnection.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.DBConnection.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s\t%s\t%s\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func HandlerFollowing(s *State, cmd Command) error {
	user, err := s.DBConnection.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	feeds, err := s.DBConnection.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feeds) == 0 {
		return nil
	}
	fmt.Printf("Feeds for %s:\n", feeds[0].Name)
	for _, feed := range feeds {
		fmt.Printf("* %s\t%s\n", feed.Name_2, feed.Url)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command) error {
	user, err := s.DBConnection.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.DBConnection.GetFeedByUrl(context.Background(), cmd.Args[0])
	_, err = s.DBConnection.CreateFeedFollowsForUser(context.Background(), database.CreateFeedFollowsForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command) error {
	user, err := s.DBConnection.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	err = s.DBConnection.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return err
	}
	return nil
}

func HandlerBrowse(s *State, cmd Command) error {
	var limit int32 = 2
	if len(cmd.Args) > 0 {
		tmp, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(tmp)
	}
	user, err := s.DBConnection.GetUserByName(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	posts, err := s.DBConnection.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		fmt.Println("No posts yet")
	}

	for _, item := range posts {
		fmt.Println(item.PublishedAt)
		fmt.Println(item.Title)
		fmt.Println(item.Description)
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	fun, ok := c.Commands[cmd.Name]

	if !ok {
		return fmt.Errorf("No Command with that name found\n")
	}

	return fun(s, cmd)
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var feed RSSFeed
	err = xml.NewDecoder(res.Body).Decode(&feed)

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, err
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return nil
}

func scrapeFeeds(s *State) error {
	next_feed, err := s.DBConnection.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	s.DBConnection.MarkFeedFetched(context.Background(), next_feed.ID)
	feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			fmt.Printf("failed to parse %v.\n%v\n", item.PubDate, err)
			continue
		}
		s.DBConnection.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: t,
			FeedID:      next_feed.ID,
		})
		fmt.Println(item.Title)
	}

	return nil
}
