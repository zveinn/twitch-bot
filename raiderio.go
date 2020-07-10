package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type RAIDERResponse struct {
	Name              string `json:"name"`
	Race              string `json:"race"`
	Class             string `json:"class"`
	ActiveSpecName    string `json:"active_spec_name"`
	ActiveSpecRole    string `json:"active_spec_role"`
	Gender            string `json:"gender"`
	Faction           string `json:"faction"`
	AchievementPoints int    `json:"achievement_points"`
	HonorableKills    int    `json:"honorable_kills"`
	ThumbnailURL      string `json:"thumbnail_url"`
	Region            string `json:"region"`
	Realm             string `json:"realm"`
	ProfileURL        string `json:"profile_url"`
	ProfileBanner     string `json:"profile_banner"`

	// route dependant
	Gear GEAR `json:"gear"`
	Score []RIOSCORE `json:"mythic_plus_scores_by_season"`
	
}

type Player struct {
	Base RAIDERResponse `json:"base"`
}
var PlayerCache = make(map[string]*Player)

func RaiderIOCharacter(region,realm,character string, fields []string) *Player {

if PlayerCache["region="+region+"&realm="+realm+"&name="+character] != nil {
	log.Println("PLAYER FROM CACHE..")
	return PlayerCache["region="+region+"&realm="+realm+"&name="+character]
}
	client := &http.Client{}
	parameters := "region="+region+"&realm="+realm+"&name="+character+"&fields="+strings.Join(fields, ",")
	
	req, err := http.NewRequest("GET", "https://raider.io/api/v1/characters/profile?"+parameters, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var RR RAIDERResponse
	err = json.Unmarshal(bodyBytes, &RR)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return nil
	}
if PlayerCache["region="+region+"&realm="+realm+"&name="+character] == nil {
	NewPlayer := Player{
		Base: RR,
	}
	PlayerCache["region="+region+"&realm="+realm+"&name="+character] = &NewPlayer
}


	GetScore(region,realm,character)
	defer resp.Body.Close()
	return PlayerCache["region="+region+"&realm="+realm+"&name="+character]
}
func GetScore(region,realm,character string) {
client := &http.Client{}
	parameters := "region="+region+"&realm="+realm+"&name="+character+"&fields=mythic_plus_scores_by_season:current"
	
	req, err := http.NewRequest("GET", "https://raider.io/api/v1/characters/profile?"+parameters, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var RR RAIDERResponse
	err = json.Unmarshal(bodyBytes, &RR)
	if err != nil {
		log.Println(err, string(debug.Stack()))
		return 
	}

	PlayerCache["region="+region+"&realm="+realm+"&name="+character].Base.Score = RR.Score

}
type RIOSCORE  struct {
	Season string `json:"season"`
	Scores struct {
		All    float64 `json:"all"`
		Dps    float64 `json:"dps"`
		Healer float64     `json:"healer"`
		Tank   float64     `json:"tank"`
		Spec0  float64     `json:"spec_0"`
		Spec1  float64     `json:"spec_1"`
		Spec2  float64 `json:"spec_2"`
		Spec3  float64     `json:"spec_3"`
	} `json:"scores"`
}

type GEAR              struct {
		ItemLevelEquipped int     `json:"item_level_equipped"`
		ItemLevelTotal    int     `json:"item_level_total"`
		ArtifactTraits    float64 `json:"artifact_traits"`
		Corruption        struct {
			Added     int `json:"added"`
			Resisted  int `json:"resisted"`
			Total     int `json:"total"`
			CloakRank int `json:"cloakRank"`
			Spells    []struct {
				ID     int         `json:"id"`
				Name   string      `json:"name"`
				Icon   string      `json:"icon"`
				School int         `json:"school"`
				Rank   interface{} `json:"rank"`
			} `json:"spells"`
		} `json:"corruption"`
		Items struct {
			Head struct {
				ItemID            int    `json:"item_id"`
				ItemLevel         int    `json:"item_level"`
				ItemQuality       int    `json:"item_quality"`
				Icon              string `json:"icon"`
				IsLegionLegendary bool   `json:"is_legion_legendary"`
				IsAzeriteArmor    bool   `json:"is_azerite_armor"`
				AzeritePowers     []struct {
					ID    int `json:"id"`
					Spell struct {
						ID     int         `json:"id"`
						Name   string      `json:"name"`
						Icon   string      `json:"icon"`
						School int         `json:"school"`
						Rank   interface{} `json:"rank"`
					} `json:"spell"`
					Tier int `json:"tier"`
				} `json:"azerite_powers"`
				Corruption struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"head"`
			Neck struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				HeartOfAzeroth    struct {
					Essences []struct {
						Slot  int `json:"slot"`
						ID    int `json:"id"`
						Rank  int `json:"rank"`
						Power struct {
							ID      int `json:"id"`
							Essence struct {
								ID          int    `json:"id"`
								Name        string `json:"name"`
								ShortName   string `json:"shortName"`
								Description string `json:"description"`
							} `json:"essence"`
							TierID          int `json:"tierId"`
							MajorPowerSpell struct {
								ID     int    `json:"id"`
								Name   string `json:"name"`
								Icon   string `json:"icon"`
								School int    `json:"school"`
								Rank   string `json:"rank"`
							} `json:"majorPowerSpell"`
							MinorPowerSpell struct {
								ID     int    `json:"id"`
								Name   string `json:"name"`
								Icon   string `json:"icon"`
								School int    `json:"school"`
								Rank   string `json:"rank"`
							} `json:"minorPowerSpell"`
						} `json:"power"`
					} `json:"essences"`
					Level    int     `json:"level"`
					Progress float64 `json:"progress"`
				} `json:"heart_of_azeroth"`
				Corruption struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"neck"`
			Shoulder struct {
				ItemID            int    `json:"item_id"`
				ItemLevel         int    `json:"item_level"`
				ItemQuality       int    `json:"item_quality"`
				Icon              string `json:"icon"`
				IsLegionLegendary bool   `json:"is_legion_legendary"`
				IsAzeriteArmor    bool   `json:"is_azerite_armor"`
				AzeritePowers     []struct {
					ID    int `json:"id"`
					Spell struct {
						ID     int         `json:"id"`
						Name   string      `json:"name"`
						Icon   string      `json:"icon"`
						School int         `json:"school"`
						Rank   interface{} `json:"rank"`
					} `json:"spell"`
					Tier int `json:"tier"`
				} `json:"azerite_powers"`
				Corruption struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"shoulder"`
			Back struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added     int `json:"added"`
					Resisted  int `json:"resisted"`
					Total     int `json:"total"`
					CloakRank int `json:"cloakRank"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"back"`
			Chest struct {
				ItemID            int    `json:"item_id"`
				ItemLevel         int    `json:"item_level"`
				ItemQuality       int    `json:"item_quality"`
				Icon              string `json:"icon"`
				IsLegionLegendary bool   `json:"is_legion_legendary"`
				IsAzeriteArmor    bool   `json:"is_azerite_armor"`
				AzeritePowers     []struct {
					ID    int `json:"id"`
					Spell struct {
						ID     int         `json:"id"`
						Name   string      `json:"name"`
						Icon   string      `json:"icon"`
						School int         `json:"school"`
						Rank   interface{} `json:"rank"`
					} `json:"spell"`
					Tier int `json:"tier"`
				} `json:"azerite_powers"`
				Corruption struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"chest"`
			Waist struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"waist"`
			Wrist struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"wrist"`
			Hands struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"hands"`
			Legs struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"legs"`
			Feet struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"feet"`
			Finger1 struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []int `json:"gems"`
				Bonuses []int `json:"bonuses"`
			} `json:"finger1"`
			Finger2 struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"finger2"`
			Trinket1 struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"trinket1"`
			Trinket2 struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []int `json:"gems"`
				Bonuses []int `json:"bonuses"`
			} `json:"trinket2"`
			Mainhand struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"mainhand"`
			Offhand struct {
				ItemID            int           `json:"item_id"`
				ItemLevel         int           `json:"item_level"`
				ItemQuality       int           `json:"item_quality"`
				Icon              string        `json:"icon"`
				IsLegionLegendary bool          `json:"is_legion_legendary"`
				IsAzeriteArmor    bool          `json:"is_azerite_armor"`
				AzeritePowers     []interface{} `json:"azerite_powers"`
				Corruption        struct {
					Added    int `json:"added"`
					Resisted int `json:"resisted"`
					Total    int `json:"total"`
				} `json:"corruption"`
				Gems    []interface{} `json:"gems"`
				Bonuses []int         `json:"bonuses"`
			} `json:"offhand"`
		} `json:"items"`
	} 