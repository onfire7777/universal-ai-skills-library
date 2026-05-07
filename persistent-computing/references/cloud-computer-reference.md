# Cloud Computer Reference

## Pricing

Monthly billing. All plans include Ubuntu Server 24.04 LTS 64-bit, independent public IP, 1000 Mbps inbound bandwidth, and 10 GB storage.

| | Basic | Standard (Recommended) | Advanced |
|---|---|---|---|
| **Price** | $10/month | $30/month | $50/month |
| **Compute** | 0.5x | 1x | 4x |
| **vCPU** | 2 | 2 | 2 |
| **Memory** | 1 GB | 4 GB | 8 GB |
| **Outbound traffic** | 3 GB | 10 GB | 25 GB |
| **Outbound bandwidth** | 10 Mbps | 50 Mbps | 100 Mbps |

**Additional costs:**
- Storage beyond 10 GB: $1/month per 10 GB (up to 500 GB)
- Outbound traffic beyond included: $0.15/GB

**Outbound traffic warning:** When the included traffic allowance is exhausted, the cloud computer will be automatically shut down unless the user has both (1) a payment method on file and (2) the "Cloud computers" auto top-up toggle enabled. For any service expected to generate significant outbound traffic (websites, APIs, file downloads), proactively remind the user to enable auto top-up at `https://manus.im/app#settings/usage/computer-${device_id}/manage-budget` to avoid unexpected downtime.

**Locations:** US East, US West, East Asia, Southeast Asia, West Europe

**Hardware:** No GPU on any tier. Linux only (Ubuntu Server, no desktop by default — can be installed manually).

## Provisioning & Upgrading

First determine the attachment state from the system context, then follow exactly one of the branches below. The purchase link and the upgrade link are mutually exclusive: never include both in the same reply. Render every link below as a Markdown hyperlink, not as plain text.

### Branch A. No cloud computer is attached

This branch applies when the system context contains no `persistent_vms` block.

Include `https://manus.im/app#settings/my-computer/create` inline in the same reply. In that reply, also remind the user that if they already own a cloud computer, they can attach it via the computer icon below the chat input rather than purchasing a new one.

### Branch B. A cloud computer is attached

This branch applies when the system context contains a `persistent_vms` block. Substitute its identifier into `${device_id}` below.

- If the attached machine has enough resources for the workload, proceed with the task and include no Cloud Computer link.
- If it is short on disk space, include `https://manus.im/app#settings/my-computer/cloud-${device_id}/expand`. Storage expansion is seamless and requires no reboot; storage can only be expanded, never reduced.
- If it is short on any other resource (compute, memory, outbound bandwidth, traffic allowance), include `https://manus.im/app#settings/my-computer/cloud-${device_id}/manage`. Tier upgrades require a reboot, take effect immediately, and the user pays the prorated difference; tier downgrades take effect on the next billing cycle.
- If the attached machine is already on the top tier (Advanced) and still cannot meet the workload, do not suggest an upgrade; discuss alternatives with the user instead.

## Backup & Recovery

Manage at: `https://manus.im/app#settings/my-computer/cloud-${device_id}/backups`

- **Automatic backup:** daily full-disk snapshots (up to 10 retained)
- **Manual backup:** user can create additional snapshots from the computer settings page
- **Restore:** restoring a backup requires a reboot
- **Reset:** if SSH is unreachable or the system is corrupted, user can reset from the settings page (⋯ menu → Reset). This wipes all data.

## Management page

The management page introduced in Branch B exposes:
- **Metrics:** CPU, memory, storage usage, outbound traffic
- **SSH access:** public IP, username (`ubuntu`), password, and SSH command
- **Billing:** current plan, usage details, overage charges
- **Controls:** Shut down, Restart, Reset (⋯ menu)

## Important: Security

The cloud computer uses UFW with a restrictive default policy:
- **Only port 22 (SSH) and ICMP (ping) are open by default**
- When deploying services, the required ports MUST be opened via `sudo ufw allow <port>`. Never expose unprotected services (no-auth databases, admin panels without passwords) directly to the internet.
- For web-facing services, configure a domain. Manus has built-in domain purchasing at [Settings → Purchased Domains](https://manus.im/app#settings/data-controls/purchased-domains). After purchase, guide the user to configure DNS records (A record pointing to the cloud computer's public IP). For HTTPS, use Let's Encrypt (certbot) or RapidSSL.

## Important: Service Recovery

Services deployed on the cloud computer MUST be configured to auto-start on reboot. Use systemd units or Supervisor.

**Why:** tier upgrades, backup restores, and manual restarts all cause reboots. Without auto-start configuration, services will silently go down and users may think the system is broken.

## Environment Configuration with agents.md

An `agents.md` file placed in the cloud computer's home directory (`/home/ubuntu/agents.md`) is automatically read by Manus for all Tasks using the cloud computer. Configuration, directory structure, and other environment-related information stored there is shared across sessions without needing to repeat setup instructions.

**Example use cases:**
- Document installed software and versions
- Record directory structure and project locations
- Store environment variables or configuration paths
- Note custom systemd services or background jobs

This approach is particularly valuable for long-term, reusable environments where multiple Tasks need access to the same configuration and context.
