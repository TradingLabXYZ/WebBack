package main

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/copier"
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
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'all', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userd', 'all', current_timestamp, current_timestamp);`)
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
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'usera', 'all', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86B', 'userb', 'private', current_timestamp, current_timestamp),
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86C', 'userc', 'followers', current_timestamp, current_timestamp), 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86D', 'userd', 'subscribers', current_timestamp, current_timestamp);`)

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

func TestCheckVisibility(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)
	Db.Exec(`
		INSERT INTO coins (
			coinid, name, symbol, slug)
		VALUES
			(1, 'Bitcoin', 'BTC', 'Bitcoin'),
			(1000, 'USDC', 'USDC', 'usdc');`)
	Db.Exec(`
		INSERT INTO lastprices (
			updatedat, coinid, price)
		VALUES
			(current_timestamp, 1, 45000),
			(current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO trades(
			code, userwallet, createdat, updatedat,
			firstpair, secondpair)
		VALUES
		('useratr1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', current_timestamp, current_timestamp, 1000, 1);`)
	Db.Exec(`
		INSERT INTO subtrades(
			code, userwallet, tradecode, createdat, updatedat,
			type, quantity, avgprice, total, reason)
		VALUES ('userasub1', '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'useratr1', current_timestamp, current_timestamp, 'BUY', 1, 65000, 65000, 'TESTART');`)

	user, _ := SelectUser("wallet", "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A")
	initial_snapshot := user.GetSnapshot()

	// <test code>
	t.Run(fmt.Sprintf("Test totalcounttrades is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET totalcounttrades = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.CountTrades-temp_snapshot.CountTrades != 1 {
			t.Fatal("Failed totalcounttrades is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test totalportfolio is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET totalportfolio = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.TotalPortfolioUsd != "45,000" || temp_snapshot.TotalPortfolioUsd != "0" {
			t.Fatal("Failed totalportfolio is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test totalreturn is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET totalreturn = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.TotalReturnUsd != "-20,000" || temp_snapshot.TotalReturnUsd != "0" {
			t.Fatal("Failed totalreturn is equals to 0")
		}
		if initial_snapshot.TotalReturnBtc != "-0.44" || temp_snapshot.TotalReturnBtc != "0" {
			t.Fatal("Failed totalreturn is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test totalroi is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET totalroi = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Roi != -30.77 || temp_snapshot.Roi != 0 {
			t.Fatal("Failed totalreturn is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test tradeqtyavailable is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET tradeqtyavailable = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].QtyAvailable != "1" && temp_snapshot.Trades[0].QtyAvailable != "0" {
			t.Fatal("Failed tradeqtyavailable is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test tradevalue is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET tradevalue = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].TotalValueUsd != 45000 || temp_snapshot.Trades[0].TotalValueUsd != 0 {
			t.Fatal("Failed tradevalue is equals to 0")
		}
		if initial_snapshot.Trades[0].TotalValueUsdS != "45,000" || temp_snapshot.Trades[0].TotalValueUsdS != "0" {
			t.Fatal("Failed tradevalue is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test tradereturn is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET tradereturn = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].TotalReturnUsd != 20000 && temp_snapshot.Trades[0].TotalReturnUsd != 0 {
			t.Fatal("Failed tradereturn is equals to 0")
		}
		if initial_snapshot.Trades[0].TotalReturnBtc != -0.4444444444444444 || temp_snapshot.Trades[0].TotalReturnBtc != 0 {
			t.Fatal("Failed tradereturn is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test traderoi is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET traderoi = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].Roi != -30.8 || temp_snapshot.Trades[0].Roi != 0 {
			t.Fatal("Failed traderoi is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test subtradereasons is equals to null"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET subtradereasons = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].Subtrades[0].Reason != "TESTART" || temp_snapshot.Trades[0].Subtrades[0].Reason != "" {
			t.Fatal("Failed subtradereasons is equals to null")
		}
	})
	t.Run(fmt.Sprintf("Test subtradequantity is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET subtradequantity = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].Subtrades[0].Quantity != 1 || temp_snapshot.Trades[0].Subtrades[0].Quantity != 0 {
			t.Fatal("Failed subtradereasons is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test subtradeavgprice is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET subtradeavgprice = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].Subtrades[0].AvgPrice != 65000 || temp_snapshot.Trades[0].Subtrades[0].AvgPrice != 0 {
			t.Fatal("Failed subtradeavgprice is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test subtradetotal is equals to 0"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET subtradetotal = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if initial_snapshot.Trades[0].Subtrades[0].Total != 65000 || temp_snapshot.Trades[0].Subtrades[0].Total != 0 {
			t.Fatal("Failed subtradetotal is equals to 0")
		}
	})
	t.Run(fmt.Sprintf("Test subtradesall is empty"), func(t *testing.T) {
		Db.Exec(`
			UPDATE visibilities
			SET subtradesall = FALSE
			WHERE wallet = '0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A';`)
		temp_snapshot := TradesSnapshot{}
		copier.CopyWithOption(&temp_snapshot, &initial_snapshot, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		user.CheckVisibility(&temp_snapshot)
		if len(initial_snapshot.Trades[0].Subtrades) != 1 || len(temp_snapshot.Trades[0].Subtrades) != 0 {
			t.Fatal("Failed subtradesall is empty")
		}
	})
	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}

func TestGenerateApiToken(t *testing.T) {
	// <setup code>
	Db.Exec(
		`INSERT INTO users (wallet, username, privacy, createdat, updatedat) VALUES 
		('0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', 'userd', 'all', current_timestamp, current_timestamp);`)
	Db.Exec(
		`INSERT INTO visibilities (wallet, totalcounttrades, totalportfolio,
			totalreturn, totalroi, tradeqtyavailable, tradevalue, tradereturn,
			traderoi, subtradesall, subtradereasons, subtradequantity, subtradeavgprice, subtradetotal)
		VALUES (
			'0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A', TRUE, TRUE, TRUE, TRUE,
			TRUE, TRUE, TRUE ,TRUE, TRUE, TRUE, TRUE, TRUE, TRUE);`)

	// <test code>
	t.Run(fmt.Sprintf("Test wrong header"), func(t *testing.T) {
		req := httptest.NewRequest("GET", "/generate_api_token", nil)
		req.Header.Set("Authorization", "Bearer sessionId=")
		w := httptest.NewRecorder()
		GenerateApiToken(w, req)
		res := w.Result()
		if res.StatusCode != 401 {
			t.Fatal("Failed test generate APi token, wrong header")
		}
	})

	t.Run(fmt.Sprintf("Test wrong origin"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("api", "Europe|Berlin")
		req := httptest.NewRequest("GET", "/generate_api_token", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GenerateApiToken(w, req)
		res := w.Result()
		if res.StatusCode != 400 {
			t.Fatal("Failed test generate APi token, wrong origin")
		}
	})
	t.Run(fmt.Sprintf("Test deleting previous session"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web", "Europe|Berlin")
		_, _ = user.InsertSession("api", "Europe|Berlin")
		_, _ = user.InsertSession("api", "Europe|Berlin")
		req := httptest.NewRequest("GET", "/generate_api_token", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GenerateApiToken(w, req)
		var count_row int
		_ = Db.QueryRow(`
					SELECT
						COUNT(*)
					FROM sessions
					WHERE userwallet = $1
					AND origin = 'api';`, "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A").Scan(&count_row)
		if count_row != 1 {
			t.Fatal("Failed test generate APi token, deleting previous row")
		}
	})
	t.Run(fmt.Sprintf("Succesfully create session"), func(t *testing.T) {
		user := User{Wallet: "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A"}
		session, _ := user.InsertSession("web", "Europe|Berlin")
		_, _ = user.InsertSession("api", "Europe|Berlin")
		_, _ = user.InsertSession("api", "Europe|Berlin")
		req := httptest.NewRequest("GET", "/generate_api_token", nil)
		req.Header.Set("Authorization", "Bearer sessionId="+session.Code)
		w := httptest.NewRecorder()
		GenerateApiToken(w, req)
		session_received, _ := ioutil.ReadAll(w.Body)
		session_received_s := string(session_received)
		var session_code string
		_ = Db.QueryRow(`
					SELECT
						code
					FROM sessions
					WHERE userwallet = $1
					AND origin = 'api';`, "0x29D7d1dd5B6f9C864d9db560D72a247c178aE86A").Scan(&session_code)

		if session_received_s != session_code {
			t.Fatal("Failed test generate APi token, successfully create session")
		}
	})

	// <tear-down code>
	Db.Exec(`DELETE FROM users WHERE 1 = 1;`)
}
