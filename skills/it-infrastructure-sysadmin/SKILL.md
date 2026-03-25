---
name: IT Infrastructure & Sysadmin Engine
description: Comprehensive IT infrastructure, system administration, networking, troubleshooting, Active Directory, DNS, DHCP, virtualization, cloud infrastructure, and IT automation expertise.
version: 1.0.0
author: Custom Meta-Skill
tags: [IT, sysadmin, networking, infrastructure, active-directory, DNS, DHCP, virtualization, cloud, troubleshooting]
---

# IT Infrastructure & Sysadmin Engine

## Purpose
Provide expert-level IT infrastructure knowledge covering networking, system administration, troubleshooting methodologies, and cloud/on-premises infrastructure management.

## Networking Fundamentals

### OSI Model (Troubleshooting Layer by Layer)
| Layer | Name | Protocols | Troubleshooting |
|-------|------|-----------|-----------------|
| 7 | Application | HTTP, DNS, SMTP, FTP | Check app logs, API responses |
| 6 | Presentation | SSL/TLS, JPEG, ASCII | Check certificates, encoding |
| 5 | Session | NetBIOS, RPC, PPTP | Check session state, auth |
| 4 | Transport | TCP, UDP | Check ports, firewalls, `netstat` |
| 3 | Network | IP, ICMP, OSPF, BGP | Check routing, `ping`, `traceroute` |
| 2 | Data Link | Ethernet, ARP, VLAN | Check switches, MAC tables, `arp` |
| 1 | Physical | Cables, WiFi, Fiber | Check cables, link lights, signal |

### DNS Deep Dive
- **Record Types**: A (IPv4), AAAA (IPv6), CNAME (alias), MX (mail), TXT (verification/SPF/DKIM), NS (nameserver), SOA (authority), SRV (service), PTR (reverse)
- **Resolution Flow**: Client cache → Local DNS → Root → TLD → Authoritative
- **Troubleshooting**: `nslookup`, `dig`, `host`, check TTL, check propagation
- **Common Issues**: Stale cache, misconfigured records, TTL too high/low, DNSSEC validation failures

### DHCP
- **DORA Process**: Discover → Offer → Request → Acknowledge
- **Scope Management**: Address pools, reservations, exclusions, lease duration
- **Troubleshooting**: Check scope exhaustion, relay agents, lease conflicts, `ipconfig /release && ipconfig /renew`

### Subnetting Quick Reference
| CIDR | Subnet Mask | Hosts | Use Case |
|------|-------------|-------|----------|
| /32 | 255.255.255.255 | 1 | Host route |
| /30 | 255.255.255.252 | 2 | Point-to-point |
| /28 | 255.255.255.240 | 14 | Small segment |
| /24 | 255.255.255.0 | 254 | Standard LAN |
| /16 | 255.255.0.0 | 65,534 | Large network |

## System Administration

### Windows Server
- **Active Directory**: Domain controllers, OUs, GPOs, FSMO roles, replication, trusts
- **Group Policy**: Computer vs User config, precedence (Local → Site → Domain → OU), `gpupdate /force`, `gpresult /r`
- **File Services**: NTFS permissions, share permissions, DFS, quotas
- **Certificate Services**: PKI, CA hierarchy, certificate templates, auto-enrollment
- **Hyper-V**: Virtual switches, checkpoints, live migration, replication

### Linux Server
- **Package Management**: apt (Debian/Ubuntu), yum/dnf (RHEL/CentOS), pacman (Arch)
- **Service Management**: `systemctl start/stop/enable/status service`
- **User Management**: `useradd`, `usermod`, `passwd`, `/etc/passwd`, `/etc/shadow`, `sudo`
- **Cron Jobs**: `crontab -e`, format: `min hour dom month dow command`
- **Log Management**: `/var/log/`, `journalctl`, `logrotate`, syslog/rsyslog
- **Performance**: `top`/`htop`, `vmstat`, `iostat`, `sar`, `free -h`, `df -h`

## Troubleshooting Methodology

### The 7-Step IT Troubleshooting Framework
1. **Identify the Problem**: What exactly is broken? Who is affected? When did it start? What changed?
2. **Establish a Theory**: Based on symptoms, what are the most likely causes? (Start with the simplest)
3. **Test the Theory**: Can you reproduce? Does the evidence support your theory?
4. **Establish a Plan**: What's the fix? What's the rollback plan? Who needs to be notified?
5. **Implement the Fix**: Apply the change, monitor for side effects
6. **Verify Full Functionality**: Test the fix AND test that nothing else broke
7. **Document**: Root cause, fix applied, prevention measures

### The "5 Whys" for Root Cause Analysis
- Problem: Website is down
- Why? The web server crashed
- Why? It ran out of memory
- Why? A memory leak in the application
- Why? An unclosed database connection in a loop
- Why? No connection pooling was implemented
- **Root Cause**: Missing connection pooling → Fix: Implement connection pooling

## Cloud Infrastructure

### AWS Core Services
| Category | Service | Purpose |
|----------|---------|---------|
| Compute | EC2, Lambda, ECS, EKS | Run applications |
| Storage | S3, EBS, EFS, Glacier | Store data |
| Database | RDS, DynamoDB, ElastiCache | Manage data |
| Networking | VPC, Route 53, CloudFront, ELB | Connect and deliver |
| Security | IAM, KMS, WAF, GuardDuty | Protect resources |
| Monitoring | CloudWatch, CloudTrail, X-Ray | Observe and audit |

### Azure Equivalent Services
| AWS | Azure | Purpose |
|-----|-------|---------|
| EC2 | Virtual Machines | Compute |
| S3 | Blob Storage | Object storage |
| RDS | Azure SQL | Managed database |
| VPC | Virtual Network | Networking |
| IAM | Azure AD / Entra ID | Identity |
| Lambda | Azure Functions | Serverless |
| CloudWatch | Azure Monitor | Monitoring |

### Infrastructure as Code
- **Terraform**: Multi-cloud, declarative, state management, modules
- **CloudFormation**: AWS-native, JSON/YAML templates
- **Ansible**: Agentless, SSH-based, playbooks for configuration management
- **Pulumi**: Infrastructure as real code (TypeScript, Python, Go)

## Virtualization & Containers

### Virtualization
- **Type 1 Hypervisors**: VMware ESXi, Hyper-V, KVM (bare metal)
- **Type 2 Hypervisors**: VirtualBox, VMware Workstation (hosted)
- **Key Concepts**: vCPU allocation, memory ballooning, thin provisioning, snapshots, live migration

### Containers (Docker/Kubernetes)
- **Docker**: Images, containers, Dockerfile, docker-compose, volumes, networks
- **Kubernetes**: Pods, Deployments, Services, Ingress, ConfigMaps, Secrets, PVCs
- **Best Practices**: Multi-stage builds, non-root users, health checks, resource limits, image scanning

## Monitoring & Alerting

### The 4 Golden Signals (Google SRE)
1. **Latency**: How long requests take (track p50, p95, p99)
2. **Traffic**: How much demand (requests/sec, concurrent users)
3. **Errors**: Rate of failed requests (5xx, timeouts, exceptions)
4. **Saturation**: How full is the system (CPU, memory, disk, connections)

### Monitoring Stack Options
- **Metrics**: Prometheus + Grafana, Datadog, CloudWatch
- **Logs**: ELK Stack (Elasticsearch, Logstash, Kibana), Loki, Splunk
- **Traces**: Jaeger, Zipkin, AWS X-Ray, Datadog APM
- **Alerting**: PagerDuty, OpsGenie, Prometheus Alertmanager

## Backup & Disaster Recovery

### 3-2-1 Backup Rule
- **3** copies of data
- **2** different storage media
- **1** offsite/cloud copy

### Recovery Objectives
- **RPO** (Recovery Point Objective): How much data loss is acceptable?
- **RTO** (Recovery Time Objective): How quickly must service be restored?
- Design backup frequency and DR strategy around these objectives
