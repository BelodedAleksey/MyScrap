package autochess

import "fmt"

//Test func
func Test() {
	fmt.Println("TOP500")
	players, err := GetTop500()
	if err != nil {
		fmt.Printf("Error GetTop500: %s", err)
	}
	for _, p := range players {
		fmt.Println("Игрок: ", p.Name+"\t"+p.ID+"\t"+p.Rank+"\t"+p.MMR+"\t"+p.IconURL)
	}

	fmt.Println("SEARCH KIRITO")
	players, err = GetPlayersByName("kirito")
	if err != nil {
		fmt.Printf("Error GetPlayersByName: %s", err)
	}
	for _, p := range players {
		fmt.Println("Игрок: ", p.Name+"\t"+p.ID+"\t"+p.Rank+"\t"+p.IconURL)
	}
	fmt.Println("SEARCH BY ID")
	p, err := GetPlayerByID(players[0].ID)
	if err != nil {
		fmt.Printf("Error GetPlayersByID: %s", err)
	}
	fmt.Println("Игрок: ", p.Name)
	fmt.Println("ID: ", p.ID)
	fmt.Println("Ранг: ", p.Rank)
	fmt.Println("Очки ранга: ", p.RankPoint)
	fmt.Println("Макс. ранг: ", p.MaxRank)
	fmt.Println("Очки макс ранга: ", p.MaxRankPoint)
	fmt.Println("Конфеты: ", p.Candy)
	fmt.Println("Уровень: ", p.Level)
	fmt.Println("Очки уровня: ", p.LevelPoint)
	fmt.Println("Матчей сыграно: ", p.MatchesPlayed)
	fmt.Println("Среднее место: ", p.AveragePlace)
	fmt.Println("Топ-1: ", p.Top1)
	fmt.Println("Топ-3: ", p.Top3)
	fmt.Println("Раунды: ", p.Rounds)
	fmt.Println("ИГРЫ")
	for _, g := range p.Games {
		fmt.Println("Место: ", g.Place)
		fmt.Println("Дата: ", g.Data)
		fmt.Println("Лобби: ", g.Lobby)
		fmt.Println("ММР: ", g.MMR)
		fmt.Println("Время: ", g.Time)
		fmt.Println("Раунды: ", g.Rounds)
		fmt.Println("В/П/Н: ", g.VPN)
		fmt.Println("Золото: ", g.Gold)
		fmt.Println("Здоровье: ", g.Health)
		fmt.Println("Цена сборки: ", g.Cost)
	}
}
