package dump

import (
	"fmt"
	"math/rand"
)

var titleAdjectives = []string{
	"Dark", "Lost", "Broken", "Silent", "Golden", "Iron", "Crimson", "Forgotten",
	"Hidden", "Shattered", "Ancient", "Eternal", "Burning", "Rising", "Fallen",
	"Twisted", "Sacred", "Hollow", "Savage", "Endless", "Fading", "Blazing",
	"Frozen", "Haunted", "Cursed", "Blessed", "Ruined", "Stolen", "Buried",
	"Distant", "Neon", "Phantom", "Rogue", "Shining", "Sunken", "Wandering",
	"Wicked", "Restless", "Violent", "Radiant", "Murky", "Lonely", "Hungry",
	"Glowing", "Electric", "Desperate", "Crumbling", "Bitter", "Ashen", "Alive",
}

var titleNouns = []string{
	"Empire", "Shadow", "Legacy", "Dawn", "Storm", "Kingdom", "Frontier", "Echo",
	"Horizon", "Abyss", "Throne", "Cipher", "Wanderer", "Reckoning", "Covenant",
	"Requiem", "Tempest", "Verdict", "Vigil", "Vortex", "Wraith", "Zenith",
	"Asylum", "Beacon", "Citadel", "Dominion", "Ember", "Faction", "Gateway",
	"Harbinger", "Inferno", "Labyrinth", "Meridian", "Nexus", "Omen", "Pinnacle",
	"Relic", "Sanctum", "Specter", "Talisman", "Uprising", "Bastion", "Chronicle",
	"Crucible", "Descent", "Epoch", "Fortress", "Genesis", "Hollow", "Illusion",
}

var names = []string{
	"Axel", "Vera", "Caden", "Lyra", "Orion", "Nadia", "Ezra", "Mira",
	"Silas", "Cleo", "Dante", "Isla", "Jasper", "Freya", "Kai", "Lena",
	"Marcus", "Nova", "Owen", "Petra", "Quinn", "Rhea", "Seth", "Tara",
	"Ulric", "Veda", "Wade", "Xena", "York", "Zara", "Arlo", "Brin",
	"Cole", "Dex", "Eden", "Flynn", "Gray", "Hale", "Ines", "Jude",
}

var characters = []string{
	"a rogue detective", "a young soldier", "a disgraced scientist", "a lone warrior",
	"an exiled princess", "a hardened mercenary", "a brilliant hacker", "a retired assassin",
	"a rebellious heir", "a mysterious stranger", "a grieving father", "a resourceful spy",
	"a visionary engineer", "a forgotten hero", "a reckless pilot", "a seasoned general",
	"a naive farm boy", "a cunning thief", "a desperate mother", "a wandering monk",
	"a corrupt politician", "a gifted surgeon", "a fearless journalist", "a haunted veteran",
	"an unlikely prophet", "a street-smart orphan", "a disgraced knight", "a reluctant leader",
}

var settings = []string{
	"a post-apocalyptic wasteland", "the depths of outer space", "a medieval kingdom",
	"a dystopian megacity", "an underwater civilization", "a war-torn country",
	"a remote arctic outpost", "a sprawling desert empire", "a cyberpunk metropolis",
	"a mythical lost island", "a parallel universe", "a world overrun by machines",
	"a dying alien planet", "the criminal underworld", "a secretive government facility",
	"a time-fractured landscape", "a haunted mountain village", "a futuristic colony ship",
	"a forgotten jungle kingdom", "a virtual reality simulation", "a crumbling space station",
	"a plague-ravaged Europe", "a steampunk airship city", "a realm between life and death",
}

var conflicts = []string{
	"uncovers a conspiracy that threatens everything they hold dear",
	"must forge an unlikely alliance to survive an unstoppable threat",
	"is forced to confront the sins of their past",
	"races against time to prevent a catastrophic war",
	"discovers a dark secret that turns allies into enemies",
	"must choose between loyalty and justice",
	"fights to reclaim what was stolen from them",
	"navigates a web of betrayal and deception",
	"struggles to protect the ones they love from an ancient evil",
	"is hunted by a shadowy organization with limitless power",
	"must unlock a buried truth before it destroys the world",
	"faces an impossible choice that will alter the course of history",
	"battles inner demons while fighting an external war",
	"seeks revenge against those who wronged them",
	"tries to escape a system designed to destroy them",
	"must unite fractured factions before time runs out",
}

var themes = []string{
	"redemption and sacrifice", "the cost of power", "love across impossible odds",
	"survival against all odds", "the nature of identity", "loyalty and betrayal",
	"justice versus vengeance", "hope in the face of despair", "the burden of leadership",
	"family and belonging", "the corruption of ambition", "truth and deception",
	"freedom and oppression", "grief and healing", "fate versus free will",
	"the price of progress", "courage under fire", "trust and forgiveness",
}

var descTemplates = []string{
	"In %s, %s %s. A gripping tale of %s.",
	"Set against the backdrop of %s, %s finds themselves drawn into a conflict where they must %s. At its heart, a story about %s.",
	"When %s stumbles upon a hidden truth in %s, they %s. A powerful exploration of %s.",
	"%s navigates the dangers of %s and %s. An unflinching look at %s.",
	"Deep within %s, %s is thrust into chaos when they %s — forcing a confrontation with %s.",
	"Against all odds, %s survives %s only to discover they must %s. A meditation on %s.",
}

func generateMovies(count int) map[string]string {
	r := rand.New(rand.NewSource(42))
	movies := make(map[string]string, count)

	pick := func(slice []string) string {
		return slice[r.Intn(len(slice))]
	}

	titlePatterns := []func() string{
		func() string { return fmt.Sprintf("The %s %s", pick(titleAdjectives), pick(titleNouns)) },
		func() string { return fmt.Sprintf("%s %s", pick(titleAdjectives), pick(titleNouns)) },
		func() string { return fmt.Sprintf("The %s of %s", pick(titleNouns), pick(titleNouns)) },
		func() string { return fmt.Sprintf("%s: %s %s", pick(names), pick(titleAdjectives), pick(titleNouns)) },
		func() string {
			return fmt.Sprintf("%s and the %s %s", pick(names), pick(titleAdjectives), pick(titleNouns))
		},
		func() string { return fmt.Sprintf("Beyond the %s %s", pick(titleAdjectives), pick(titleNouns)) },
		func() string { return fmt.Sprintf("Rise of the %s %s", pick(titleAdjectives), pick(titleNouns)) },
		func() string { return fmt.Sprintf("Last %s of %s", pick(titleNouns), pick(names)) },
	}

	for len(movies) < count {
		pattern := titlePatterns[r.Intn(len(titlePatterns))]
		title := pattern()

		if _, exists := movies[title]; exists {
			title = fmt.Sprintf("%s (%d)", title, r.Intn(8999)+1000)
		}

		// templates 0, 1 and 4 take the setting first; the rest take the character first
		i := r.Intn(len(descTemplates))
		tmpl := descTemplates[i]
		var desc string
		if i == 0 || i == 1 || i == 4 {
			desc = fmt.Sprintf(tmpl, pick(settings), pick(characters), pick(conflicts), pick(themes))
		} else {
			desc = fmt.Sprintf(tmpl, pick(characters), pick(settings), pick(conflicts), pick(themes))
		}

		movies[title] = desc
	}

	return movies
}
