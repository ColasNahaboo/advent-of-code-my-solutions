package main

import (
	"testing"
)

func fight(hh, hm, bh, bd, hard int) int {
	u := Universe{[SpellCount]Spell{
		Spell{"MagicMissile", 53, 4, 0, 0, 0, 0, 0, 0}, // 0
		Spell{"Drain", 73, 2, 0, 2, 0, 0, 0, 0},        // 1
		Spell{"Shield", 113, 0, 0, 0, 6, 0, 7, 0},      // 2
		Spell{"Poison", 173, 0, 0, 0, 6, 3, 0, 0},      // 3
		Spell{"Recharge", 229, 0, 0, 0, 5, 0, 0, 101},  // 4
	},
		false,
	}
	w := World{id: NewID(), effects: [SpellCount]Effect{
		Effect{MagicMissile, 0, 0, 0, 0}, // 0
		Effect{Drain, 0, 0, 0, 0},        // 1
		Effect{Shield, 0, 0, 7, 0},       // 2
		Effect{Poison, 0, 3, 0, 0},       // 3
		Effect{Recharge, 0, 0, 0, 101},   // 4
	},
	}
	w.hero.hp = hh
	w.hero.mana = hm
	w.boss.hp = bh
	w.boss.dmg = bd
	if hard != 0 {
		u.hard = true
	}
	curmana := maxint
	w.curmana = &curmana
	return fightRound(&w, u)
}

func Test_fightRound(t *testing.T) {
	t.Run("10 10", func(t *testing.T) {
		got := fight(50, 500, 10, 10, 0)
		expected := 159
		if got != expected {
			t.Errorf("expected '%#v' but got '%#v'", expected, got)
		}
	})
	t.Run("50 10", func(t *testing.T) {
		got := fight(50, 500, 50, 10, 0)
		expected := 900
		if got != expected {
			t.Errorf("expected '%#v' but got '%#v'", expected, got)
		}
	})
	t.Run("50 10 hard", func(t *testing.T) {
		got := fight(50, 500, 50, 10, 1)
		expected := 1309
		if got != expected {
			t.Errorf("expected '%#v' but got '%#v'", expected, got)
		}
	})
	t.Run("80 10", func(t *testing.T) {
		got := fight(50, 500, 80, 10, 0)
		expected := 1990
		if got != expected {
			t.Errorf("expected '%#v' but got '%#v'", expected, got)
		}
	})
}
