package script

import (
	"fmt"
	"io"

	"github.com/GiGurra/tais2/internal/game"
)

// WriteSnapshot writes a human-readable dump of the world state.
func WriteSnapshot(w io.Writer, s *game.Scenario) {
	fmt.Fprintf(w, "=== SNAPSHOT tick=%d alive=%d ===\n", s.Tick, s.World.AliveCount)
	s.World.Each(game.MaskPosition, func(idx int32) bool {
		e := &s.World.Entities[idx]
		pos := s.World.Position[idx]

		playerID := int32(-1)
		if e.Mask&game.MaskOwner != 0 {
			playerID = s.World.Owner[idx].PlayerID
		}

		ut := game.UnitNone
		if e.Mask&game.MaskUnitType != 0 {
			ut = s.World.UnitTypeComp[idx].Type
		}

		hp := ""
		if e.Mask&game.MaskHealth != 0 {
			h := s.World.Health[idx]
			hp = fmt.Sprintf(" hp=%d/%d", h.Current, h.Max)
		}

		fmt.Fprintf(w, "  [%d] player=%d type=%d pos=%d,%d%s\n",
			idx, playerID, ut, pos.X/1000, pos.Y/1000, hp)
		return true
	})
	fmt.Fprintln(w, "=== END ===")
}
