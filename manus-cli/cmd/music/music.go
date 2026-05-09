package music

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/onfire7777/manus-cli/internal/platform"
)

// Cmd is the top-level music command group.
var Cmd = &cobra.Command{
	Use:   "music",
	Short: "Music prompt crafting framework (genres, structure, multi-clip)",
	Long: `Music generation prompt crafting framework.
Covers prompt structure, genre syntax, multi-clip strategy,
and best practices for AI music generation.`,
}

var promptCmd = &cobra.Command{
	Use:   "prompt <description>",
	Short: "Generate a structured music prompt from a description",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		desc := strings.Join(args, " ")
		fmt.Printf("Music Prompt for: %s\n\n", desc)
		fmt.Println("Generated prompt structure:")
		fmt.Printf("  Genre: [inferred from description]\n")
		fmt.Printf("  Tempo: [matched to mood]\n")
		fmt.Printf("  Instruments: [genre-appropriate]\n")
		fmt.Printf("  Structure: [intro → verse → chorus → bridge → outro]\n")
		fmt.Printf("  Mood: %s\n", desc)
		fmt.Println()
		fmt.Println("Use 'manus music reference' to see the full prompt crafting guide.")
		return nil
	},
}

var referenceCmd = &cobra.Command{
	Use:   "reference",
	Short: "Show the full music prompt crafting reference guide",
	RunE: func(cmd *cobra.Command, args []string) error {
		refPath := filepath.Join(platform.SkillsDir(), "music-prompter", "references", "prompt_guide.md")
		if _, err := os.Stat(refPath); err != nil {
			refPath = filepath.Join(platform.RepoDir(), "music-prompter", "references", "prompt_guide.md")
		}
		if _, err := os.Stat(refPath); err != nil {
			// Inline fallback
			fmt.Println("Music Prompt Crafting Framework")
			fmt.Println("================================")
			fmt.Println()
			fmt.Println("Structure Syntax:")
			fmt.Println("  [Intro] [Verse 1] [Chorus] [Verse 2] [Chorus] [Bridge] [Chorus] [Outro]")
			fmt.Println()
			fmt.Println("Key Elements:")
			fmt.Println("  - Genre/Style tag (e.g., 'cinematic orchestral', 'lo-fi hip hop')")
			fmt.Println("  - Tempo (BPM or descriptor: slow/medium/fast)")
			fmt.Println("  - Key instruments")
			fmt.Println("  - Mood/Emotion keywords")
			fmt.Println("  - Production style (clean, distorted, reverb-heavy)")
			fmt.Println()
			fmt.Println("Multi-Clip Strategy:")
			fmt.Println("  - Clip 1: Intro + Verse (30s)")
			fmt.Println("  - Clip 2: Chorus + Bridge (30s)")
			fmt.Println("  - Clip 3: Final Chorus + Outro (30s)")
			return nil
		}
		data, _ := os.ReadFile(refPath)
		fmt.Print(string(data))
		return nil
	},
}

var genresCmd = &cobra.Command{
	Use:   "genres",
	Short: "List supported genres and their prompt patterns",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Supported Genres & Prompt Patterns:")
		fmt.Println()
		fmt.Printf("  %-25s %s\n", "GENRE", "TYPICAL PROMPT ELEMENTS")
		fmt.Printf("  %-25s %s\n", "-----", "----------------------")
		fmt.Printf("  %-25s %s\n", "Cinematic Orchestral", "strings, brass, percussion, epic, sweeping")
		fmt.Printf("  %-25s %s\n", "Lo-Fi Hip Hop", "vinyl crackle, mellow keys, boom-bap, chill")
		fmt.Printf("  %-25s %s\n", "Synthwave/Retrowave", "analog synths, arpeggios, 80s, neon")
		fmt.Printf("  %-25s %s\n", "Ambient/Atmospheric", "pads, drones, reverb, ethereal, space")
		fmt.Printf("  %-25s %s\n", "Jazz", "piano, upright bass, brushes, swing, smooth")
		fmt.Printf("  %-25s %s\n", "Electronic/EDM", "drops, builds, synth leads, four-on-floor")
		fmt.Printf("  %-25s %s\n", "Rock/Alternative", "electric guitar, drums, bass, distortion")
		fmt.Printf("  %-25s %s\n", "Classical", "piano, strings, woodwinds, dynamic, formal")
		fmt.Printf("  %-25s %s\n", "Folk/Acoustic", "acoustic guitar, mandolin, warm, intimate")
		fmt.Printf("  %-25s %s\n", "Metal", "heavy distortion, double kick, aggressive")
	},
}

func init() {
	Cmd.AddCommand(promptCmd)
	Cmd.AddCommand(referenceCmd)
	Cmd.AddCommand(genresCmd)
}
