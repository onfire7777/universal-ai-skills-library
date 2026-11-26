# Skill source status

The November 25, 2026 review validated the complete canonical catalog and compared
950 recorded GitHub imports against the current trees of their 11 upstream repositories.
It compared entire bundle paths, including supporting files, from the recorded
import-date snapshot to the upstream tree. Refreshed bundles record that tree and
upstream folder in their `.openskills.json` files.

| Result | Bundles |
| --- | ---: |
| Refreshed from upstream | 61 |
| Upstream files unchanged since import or already present locally | 847 |
| Deliberately retained while upstream changes need a compatibility review | 25 |
| Recorded source lacked a matching current skill folder | 17 |
| No recorded GitHub import to compare | 845 |

A source freshness comparison is separate from catalog validation. The catalog
accepts all 1,795 canonical bundles. It checks paths, script inventory, and
duplicate instruction bodies. Bundles without a matching upstream snapshot
remain maintained in this repository.

## Retained variants

Open Design now publishes some formerly generic templates as specific,
example-driven decks. The generic versions remain available for existing users.
Three generated Printing Press fixture bundles retain their checked-in contract.
Impeccable's upstream update introduces a standalone launcher and replaces
multiple runtime assets. Its existing bundle remains available while those
runtime changes need separate verification.

| Source | Refreshed | Retained changed | No matching folder |
| --- | ---: | ---: | ---: |
| [ComposioHQ/awesome-claude-skills](https://github.com/ComposioHQ/awesome-claude-skills) | 0 | 0 | 0 |
| [JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman) | 6 | 0 | 0 |
| [VoltAgent/awesome-design-md](https://github.com/VoltAgent/awesome-design-md) | 0 | 0 | 1 |
| [anthropics/skills](https://github.com/anthropics/skills) | 5 | 0 | 0 |
| [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills) | 0 | 0 | 0 |
| [lightpanda-io/browser](https://github.com/lightpanda-io/browser) | 0 | 0 | 1 |
| [mattpocock/skills](https://github.com/mattpocock/skills) | 9 | 0 | 12 |
| [mvanhorn/cli-printing-press](https://github.com/mvanhorn/cli-printing-press) | 8 | 3 | 1 |
| [nextlevelbuilder/ui-ux-pro-max-skill](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill) | 6 | 0 | 0 |
| [nexu-io/open-design](https://github.com/nexu-io/open-design) | 27 | 21 | 2 |
| [pbakaus/impeccable](https://github.com/pbakaus/impeccable) | 0 | 1 | 0 |

The [source integration guide](SOURCE_INTEGRATIONS.md) covers optional
runtime tools and provider authentication.
