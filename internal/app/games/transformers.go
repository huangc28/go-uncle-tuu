package games

import "huangc28/go-ios-iap-vendor/internal/app/models"

//prod_name,
//prod_id,
//game_bundle_id

type TrfmedGames struct {
	GameBundleID string `json:"game_bundle_id"`
	ReadableName string `json:"readable_name"`
}

func transformProducts(prods []*models.Game) interface{} {
	outputs := make([]TrfmedGames, 0, len(prods))

	for _, prod := range prods {
		outputs = append(outputs, TrfmedGames{
			GameBundleID: prod.GameBundleID,
			ReadableName: prod.ReadableName,
		})
	}

	return struct {
		Games []TrfmedGames `json:"games"`
	}{outputs}
}
