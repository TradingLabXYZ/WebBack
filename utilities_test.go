package main

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// <setup code>
	// <test-code>
	t.Run(fmt.Sprintf("Test correct encryption"), func(t *testing.T) {
		have, _ := Encrypt("r")
		want := "4dc7c9ec434ed06502767136789763ec11d2c4b7"
		if have != want {
			t.Error("Failed correct encryption")
		}
	})

	t.Run(fmt.Sprintf("Test encrypting empty string"), func(t *testing.T) {
		_, err := Encrypt("")
		if err.Error() != "Empty string" {
			t.Error("TestEncrypt: empty string")
		}
	})

	// <tear-down code>
}

func TestCheckRelation(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, plan, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', 'basic', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userd', 'all', 'basic', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO followers (followfrom, followto, createdat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', current_timestamp);`)
	Db.Exec(
		`INSERT INTO subscribers (subscribefrom, subscribeto, createdat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp);`)
	// <test code>
	t.Run(fmt.Sprintf("Test user not follow when userid is null"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: ""}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		if user_connection.IsFollower || user_connection.IsSubscriber {
			t.Fatal("Failed user not follow when userid is null")
		}
	})
	t.Run(fmt.Sprintf("Test usera follows userb"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		if !user_connection.IsFollower {
			t.Fatal("Failed test usera follows userb")
		}
	})
	t.Run(fmt.Sprintf("Test userc does not follows userb"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		if user_connection.IsFollower {
			t.Fatal("Failed test userc does not follows userb")
		}
	})
	t.Run(fmt.Sprintf("Test userb is subscribed to userc"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		if !user_connection.IsSubscriber {
			t.Fatal("Failed test userb is subscibed to userc")
		}
	})
	t.Run(fmt.Sprintf("Test usera is not subscribed to userc"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		if user_connection.IsSubscriber {
			t.Fatal("Failed test usera is not subscribed to userc")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestCheckPrivacy(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, plan, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'private', 'basic', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc', 'followers', 'basic', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd', 'subscribers', 'basic', current_timestamp, current_timestamp);`)

	// <test code>
	t.Run(fmt.Sprintf("Test user with privacy ALL is fully visibile"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed test user with privacy ALL is fully visibile")
		}
	})

	t.Run(fmt.Sprintf("Test user not authenticated try to access not ALL users"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: ""}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "KO" {
			t.Fatal("Failed user not authenticated try to access not ALL users")
		}
	})

	t.Run(fmt.Sprintf("Test user PRIVATE always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed user PRIVATE always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user FOLLOWERS always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed user FOLLOWERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user SUBSCRIBERS always able to see its profile if authenticated"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed user SUBSCRIBERS always able to see its profile if authenticated")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when PRIVATE"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Reason != "private" {
			t.Fatal("Failed user cannot access other user when PRIVATE")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when FOLLOWERS and not following"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Reason != "user is not follower" {
			t.Fatal("Failed user cannot access other user when FOLLOWERS and not following")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when FOLLOWERS and yes following"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C")
		Db.Exec(`
				INSERT INTO followers (followfrom, followto, createdat)
				VALUES ('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', current_timestamp);`)
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed user can access other user when FOLLOWERS and yes following")
		}
	})

	t.Run(fmt.Sprintf("Test user cannot access other user when SUBSCRIBERS and not subscribers"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Reason != "user is not subscriber" {
			t.Fatal("Failed user cannot access other user when SUBSCRIBERS and not subscribers")
		}
	})

	t.Run(fmt.Sprintf("Test user can access other user when SUBSCRIBERS and yes subscriber"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D")
		Db.Exec(`
				INSERT INTO subscribers (subscribefrom, subscribeto, createdat)
				VALUES ('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', current_timestamp);`)
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Status != "OK" {
			t.Fatal("Failed user can access other user when SUBSCRIBERS and yes subscriber")
		}
	})

	t.Run(fmt.Sprintf("Test user with privacy not legit"), func(t *testing.T) {
		observed, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A")
		observer := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B"}
		observed.Privacy = "random"
		user_connection := Connection{
			Observer: observer,
			Observed: observed,
		}
		user_connection.CheckConnection()
		user_connection.CheckPrivacy()
		if user_connection.Privacy.Reason != "unknown reason" {
			t.Fatal("Failed user with privacy not legit")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
