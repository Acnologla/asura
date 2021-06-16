package test

import (
	"asura/src/utils/rinha"
	"testing"
	"time"
)

func TestBackground(t *testing.T) {
	t.Run("TestGetBackground", func(t *testing.T) {
		val := rinha.GetBackground(rinha.Galo{})
		if val != "https://i.imgur.com/F64ybgg.jpg" {
			t.Errorf("This must return the default background")
		}
	})
	t.Run("TestGetBackgrounds", func(t *testing.T) {
		galo := rinha.Galo{
			Cosmetics: []int{1, 2},
		}
		backgrounds, indexes := rinha.GetBackgrounds(galo.Cosmetics)
		if len(backgrounds) != 2 || len(indexes) != 2 {
			t.Errorf("The length must be 2")
		}
		if indexes[0] != 0 || indexes[1] != 1 {
			t.Errorf("This must be true")
		}
	})
}

func TestGetName(t *testing.T) {
	name := rinha.GetName("acno", rinha.Galo{
		Name: "acno o galo",
	})
	if name != "acno o galo" {
		t.Errorf("This must be 'acno o galo'")
	}
}

func TestIsIntInList(t *testing.T) {
	arr := []int{1, 2}
	if !rinha.IsIntInList(1, arr) {
		t.Errorf("1 is in the intList")
	}
	if rinha.IsIntInList(10, arr) {
		t.Errorf("10 isnt in the list")
	}
}

func TestLootbox(t *testing.T) {
	galo := rinha.Galo{
		CommonLootbox: 1,
	}
	t.Run("TestHaveLootbox", func(t *testing.T) {
		if !rinha.HaveLootbox(galo, "comum") {
			t.Errorf("this must return true")
		}
	})
	t.Run("TestGetNewLb", func(t *testing.T) {
		newGalo := rinha.GetNewLb("comum", galo, false)
		if newGalo.CommonLootbox != 0 {
			t.Errorf("this sould have 0 lootboxs")
		}
	})
}

func TestClan(t *testing.T) {
	member := rinha.ClanMember{
		ID: 0,
		Xp: 10,
	}
	clan := rinha.Clan{
		Members: []rinha.ClanMember{member},
	}
	t.Run("TestIsInClan", func(t *testing.T) {
		if !rinha.IsInClan(clan, 0) {
			t.Errorf("This must be true")
		}
		if rinha.IsInClan(clan, 20) {
			t.Errorf("This must be false")
		}
	})
	t.Run("TestGetMember", func(t *testing.T) {
		m := rinha.GetMember(clan, 0)
		if m.Xp != 10 {
			t.Errorf("This must be 10")
		}
	})
	t.Run("TestFindMemberIndex", func(t *testing.T) {
		if rinha.FindMemberIndex(clan, 0) != 0 {
			t.Errorf("This sould be 0 ")
		}
	})
	t.Run("TestRemoveMember", func(t *testing.T) {
		clan.Members = rinha.RemoveMember(clan, 0)
		if len(clan.Members) != 0 {
			t.Errorf("This must be true")
		}
	})
}

func TestMission(t *testing.T) {
	missions := []rinha.Mission{rinha.CreateMission(), rinha.CreateMission()}
	t.Run("TestRemoveMission", func(t *testing.T) {
		missions = rinha.RemoveMission(missions, 0)
		if len(missions) != 1 {
			t.Errorf("Length must be 1")
		}
	})
	t.Run("TestPopulateMissions", func(t *testing.T) {
		galo := rinha.Galo{}
		newMissions := rinha.PopulateMissions(galo)
		if len(newMissions) != 3 {
			t.Errorf("Length must be 3")
		}
		oneMission := rinha.PopulateMissions(rinha.Galo{
			LastMission: uint64(time.Now().Unix()) - 60*60*24,
		})
		if len(oneMission) != 1 {
			t.Errorf("Length must be 1")
		}
	})
}

func TestItems(t *testing.T) {
	t.Run("TestAddItem", func(t *testing.T) {
		items, _ := rinha.AddItem(1, []int{})
		if len(items) != 1 {
			t.Errorf("Length must be 1")
		}
		if rinha.Items[items[0]].Level != 1 {
			t.Errorf("Level must be 0")
		}
	})
}

func TestIsVip(t *testing.T) {
	galo := rinha.Galo{
		Vip: uint64(time.Now().Unix()) + (60 * 60 * 24 * 30),
	}
	if !rinha.IsVip(galo) {
		t.Errorf("Galo must be vip")
	}
	galo.Vip = galo.Vip - (60 * 60 * 24 * 30) - 1
	if rinha.IsVip(galo) {
		t.Errorf("Galo must not be vip")
	}
}

func TestGetGaloDB(t *testing.T) {
	_, err := rinha.GetGaloDB(365948625676795904)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestHaveGalo(t *testing.T) {
	galo := rinha.Galo{
		Galos: []rinha.SubGalo{{
			Type: 1,
		}},
	}
	if !rinha.HaveGalo(1, galo.Galos) {
		t.Errorf("This must return true")
	}
	if rinha.HaveGalo(2, galo.Galos) {
		t.Errorf("This must return false")
	}
}

func TestCalcDamage(t *testing.T) {
	min, max := rinha.CalcDamage(&rinha.Skill{
		Damage: [2]int{10, 20},
	}, rinha.Galo{
		Type: 1,
	})
	if min != 10 || max != 20 {
		t.Errorf("Max damage must be 20 and min damage must be 10")
	}
	min, max = rinha.CalcDamage(&rinha.Skill{
		Damage: [2]int{10, 20},
	}, rinha.Galo{
		Type:      1,
		GaloReset: 1,
	})

	if min != 12 || max != 24 {
		t.Errorf("Max damage must be 24 and min damage must be 102")
	}
}

func TestUpgrades(t *testing.T) {
	galo := rinha.Galo{
		Upgrades: []int{0, 1},
	}
	t.Run("TestHasUpgrades", func(t *testing.T) {
		if !rinha.HasUpgrade(galo.Upgrades, 0, 1) {
			t.Errorf("This must return true")
		}
		if !rinha.HasUpgrade(galo.Upgrades, 0) {
			t.Errorf("This must return true")
		}
		if rinha.HasUpgrade(galo.Upgrades, 1) {
			t.Errorf("This must return false")
		}
	})
	t.Run("TestGetCurrentUpgrade", func(t *testing.T) {
		upgrade := rinha.GetCurrentUpgrade(galo)
		if upgrade.Name != "Mais bonus" {
			t.Errorf("Upgrade name must be 'Mais bonus'")
		}
	})
	t.Run("TestHavePoint", func(t *testing.T) {
		galo = rinha.Galo{
			UserXp: 100,
		}
		if !rinha.HavePoint(galo) {
			t.Errorf("This must return true")
		}
		galo.UserXp = 0
		if rinha.HavePoint(galo) {
			t.Errorf("This must return false")
		}
	})
}

func TestBattle(t *testing.T) {
	t.Run("TestCheckItem", func(t *testing.T) {
		galo := rinha.Galo{
			Items: []int{1},
		}
		effect, payload := rinha.CheckItem(&galo)
		if effect != 2 || payload != 8 {
			t.Errorf("Effect must be 2 and payload must be 8")
		}
	})
}
